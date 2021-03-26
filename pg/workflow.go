package pg

import (
	"context"
	"database/sql"
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

func (s workflowService) CreateWorkflow(ctx context.Context, req flow.CreateWorkflowRequest) (*flow.Workflow, error) {
	tx, err := s.db.beginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	w := flow.Workflow{
		Name:        req.Name,
		Description: req.Description,
	}
	err = createWorkflow(ctx, tx, &w)
	if err != nil {
		return nil, err
	}

	return &w, tx.Commit()
}

func (s workflowService) UpdateWorkflow(ctx context.Context, req flow.UpdateWorkflowRequest) (*flow.Workflow, error) {
	panic("implement me")
}

func (s workflowService) DeleteWorkflow(ctx context.Context, req flow.DeleteWorkflowRequest) error {
	panic("implement me")
}

func (s workflowService) GetWorkflowByID(ctx context.Context, req flow.GetWorkflowByIDRequest) (*flow.Workflow, error) {
	panic("implement me")
}

func (s workflowService) GetWorkflows(ctx context.Context, req flow.GetWorkflowsRequest) ([]*flow.Workflow, int, error) {
	panic("implement me")
}

func createWorkflow(ctx context.Context, tx *Tx, w *flow.Workflow) error {
	// TODO assign workflow to current user
	err := insert(ctx, tx, &w, "workflows")
	if err != nil {
		return flow.Errorf(flow.EINTERNAL, "Failed to insert workflow into database")
	}

	return nil
}
