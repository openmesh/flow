package pg

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"github.com/openmesh/flow"
)

type workflowService struct {
	db *DB
}

func NewWorkflowService(db *DB) flow.WorkflowService {
	return workflowService{
		db,
	}
}

func (s workflowService) CreateWorkflow(ctx context.Context, workflow *flow.Workflow) error {
	tx, err := s.db.beginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback()

	err = createWorkflow(ctx, tx, workflow)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (s workflowService) UpdateWorkflow(ctx context.Context, id uuid.UUID, upd flow.WorkflowUpdate) (*flow.Workflow, error) {
	tx, err := s.db.beginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	workflow, err := updateWorkflow(ctx, tx, id, upd)
	if err != nil {
		return nil, err
	}
	err = attachWorkflowNodes(ctx, tx, workflow)
	if err != nil {
		return nil, err
	}

	return workflow, tx.Commit()
}

func (s workflowService) DeleteWorkflow(ctx context.Context, id uuid.UUID) error {
	tx, err := s.db.beginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := deleteWorkflow(ctx, tx, id); err != nil {
		return err
	}
	return tx.Commit()
}

func (s workflowService) GetWorkflowByID(ctx context.Context, id uuid.UUID) (*flow.Workflow, error) {
	tx, err := s.db.beginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return nil, err
	}

	workflow, err := getWorkflowByID(ctx, tx, id)
	if err != nil {
		return nil, err
	}

	return workflow, nil
}

func (s workflowService) GetWorkflows(ctx context.Context, filter flow.WorkflowFilter) ([]*flow.Workflow, int, error) {
	tx, err := s.db.beginTx(ctx, nil)
	if err != nil {
		return nil, 0, err
	}
	defer tx.Rollback()

	// Fetch list of matching workflows.
	workflows, n, err := getWorkflows(ctx, tx, filter)
	if err != nil {
		return workflows, n, err
	}

	// Iterate over workflows and attach associated nodes.
	// TODO batch query.
	for _, workflow := range workflows {
		if err := attachWorkflowNodes(ctx, tx, workflow); err != nil {
			return workflows, n, err
		}
	}
	return workflows, n, nil
}

func createWorkflow(ctx context.Context, tx *Tx, w *flow.Workflow) error {
	// Get user ID from context and return unauthorized error if no value is set.
	userID := flow.UserIDFromContext(ctx)
	if userID == uuid.Nil {
		return flow.Errorf(flow.EUNAUTHORIZED, "User not authenticated.")
	}
	// Assign user to workflow
	w.UserID = userID

	//Prepare named statement
	stmt, err := tx.PrepareNamed(`
		INSERT INTO
			workflows 
		    	(
			    	user_id,
			    	name,
			    	description
				)
		VALUES 
		    (
				:user_id,
			    :name,
			    :description
			)
		RETURNING 
			*
	`)
	if err != nil {
		return err
	}
	// Execute query and assign
	var res flow.Workflow
	if err := stmt.Get(&res, w); err != nil {
		return err
	}
	w.ID = res.ID
	w.CreatedAt = res.CreatedAt
	w.UpdatedAt = res.UpdatedAt

	return nil
}

func updateWorkflow(ctx context.Context, tx *Tx, id uuid.UUID, upd flow.WorkflowUpdate) (*flow.Workflow, error) {
	// Fetch current entity state.
	workflow, err := getWorkflowByID(ctx, tx, id)
	if err != nil {
		return nil, err
	}
	// Verify that user is authorized to edit the workflow.
	if !flow.CanEditWorkflow(ctx, workflow) {
		return workflow, flow.Errorf(flow.EUNAUTHORIZED, "Only the workflow owner can edit it.")
	}
	// Update workflow properties
	if upd.Name != nil {
		workflow.Name = *upd.Name
	}
	if upd.Description != nil {
		workflow.Description = *upd.Description
	}

	workflow.UpdatedAt = tx.now

	// Execute update query.
	if _, err := tx.ExecContext(ctx, `
		UPDATE
			workflows
		SET
			name = $1,
		    description = $2,
		    updated_at = $3
		WHERE 
			id = $4
	`,
		workflow.Name,
		workflow.Description,
		workflow.UpdatedAt,
		workflow.ID,
	); err != nil {
		return workflow, err
	}

	return workflow, nil
}

func deleteWorkflow(ctx context.Context, tx *Tx, id uuid.UUID) error {
	// Verify that workflow exists and that the current user is the owner.
	workflow, err := getWorkflowByID(ctx, tx, id)
	if err != nil {
		return err
	} else if !flow.CanEditWorkflow(ctx, workflow) {
		return flow.Errorf(flow.EUNAUTHORIZED, "Only the workflow owner can delete it.")
	}

	// Delete row from the database.
	if _, err := tx.ExecContext(ctx, `DELETE FROM workflows WHERE id = $1`, id); err != nil {
		return err
	}
	return nil
}

func getWorkflowByID(ctx context.Context, tx *Tx, id uuid.UUID) (*flow.Workflow, error) {
	workflows, _, err := getWorkflows(ctx, tx, flow.WorkflowFilter{ID: &id})
	if err != nil {
		return nil, err
	} else if len(workflows) == 0 {
		return nil, &flow.Error{Code: flow.ENOTFOUND, Message: "Workflow not found."}
	}
	err = attachWorkflowNodes(ctx, tx, workflows[0])
	return workflows[0], nil
}

func getWorkflows(ctx context.Context, tx *Tx, filter flow.WorkflowFilter) ([]*flow.Workflow, int, error) {
	userID := flow.UserIDFromContext(ctx)
	where := []string{"user_id = $1"}
	args := []interface{}{userID}

	if v := filter.ID; v != nil {
		where, args = append(where, fmt.Sprintf("id = $%d", len(where)+1)), append(args, *v)
	}
	if v := filter.Description; v != nil {
		where, args = append(where, fmt.Sprintf("description = $%d", len(where)+1)), append(args, *v)
	}
	if v := filter.Name; v != nil {
		where, args = append(where, fmt.Sprintf("name = $%d", len(where)+1)), append(args, *v)
	}

	baseQuery := fmt.Sprintf("SELECT * FROM workflows %s", buildWhereClause(where))

	var n int
	err := tx.Get(
		&n,
		fmt.Sprintf("SELECT COUNT(*) FROM (%s) AS count;", baseQuery),
		args...,
	)
	if err != nil {
		return nil, n, err
	}

	query := baseQuery + `
		ORDER BY created_at ASC
	` + formatLimitOffset(filter.Limit, filter.Page)

	workflows := make([]*flow.Workflow, 0)
	if err := tx.Select(&workflows, query, args...); err != nil {
		return workflows, n, err
	}

	return workflows, n, nil
}

func attachWorkflowNodes(ctx context.Context, tx *Tx, workflow *flow.Workflow) error {
	var err error
	workflow.Nodes, _, err = getNodes(ctx, tx, flow.NodeFilter{WorkflowID: &workflow.ID})
	return err
}
