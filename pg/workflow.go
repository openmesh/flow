package pg

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/openmesh/flow"
	"strings"
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

func (s workflowService) UpdateWorkflow(ctx context.Context, id uuid.UUID, req flow.WorkflowUpdate) (*flow.Workflow, error) {
	tx, err := s.db.beginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	workflow, err := updateWorkflow(ctx, tx, id, req)
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
	err = deleteRow(ctx, tx, id, "workflows")
	if err != nil {
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
	panic("implement me")
}

func createWorkflow(ctx context.Context, tx *Tx, w *flow.Workflow) error {
	err := insertRow(ctx, tx, w, "workflows")
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			return flow.Errorf(flow.ECONFLICT, err.(*pq.Error).Detail)
		}
		return flow.Errorf(flow.EINTERNAL, "Failed to insert workflow into database")
	}

	return nil
}

func updateWorkflow(ctx context.Context, tx *Tx, id uuid.UUID, upd flow.WorkflowUpdate) (*flow.Workflow, error) {
	// Fetch current entity state.
	workflow, err := getWorkflowByID(ctx, tx, id)
	if err != nil {
		return nil, err
	}

	if upd.Name != nil {
		workflow.Name = *upd.Name
	}
	if upd.Description != nil {
		workflow.Description = *upd.Description
	}

	if err := updateRow(ctx, tx, workflow, "workflows"); err != nil {
		return nil, err
	}

	return workflow, nil
}

func getWorkflowByID(ctx context.Context, tx *Tx, id uuid.UUID) (*flow.Workflow, error) {
	workflow := flow.Workflow{}
	err := getRowByID(ctx, tx, &workflow, id, "workflows")
	if err != nil {
		return nil, err
	}
	err = attachWorkflowNodes(ctx, tx, &workflow)
	if err != nil {
		return nil, err
	}

	return &workflow, nil
}

func attachWorkflowNodes(ctx context.Context, tx *Tx, workflow *flow.Workflow) error {
	var err error
	workflow.Nodes, _, err = getNodes(ctx, tx, flow.NodeFilter{WorkflowID: &workflow.ID})
	return err
}
