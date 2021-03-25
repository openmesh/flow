package pg

import (
	"context"
	"github.com/openmesh/events"
)

type workflowService struct {
	db *DB
}

func NewWorkflowService(db *DB) flow.WorkflowService {
	return workflowService{
		db,
	}
}

func (w workflowService) CreateWorkflow(ctx context.Context, req flow.CreateWorkflowRequest) (*flow.Workflow, error) {
	panic("implement me")
}

func (w workflowService) UpdateWorkflow(ctx context.Context, req flow.UpdateWorkflowRequest) (*flow.Workflow, error) {
	panic("implement me")
}

func (w workflowService) DeleteWorkflow(ctx context.Context, req flow.DeleteWorkflowRequest) error {
	panic("implement me")
}

func (w workflowService) GetWorkflowByID(ctx context.Context, req flow.GetWorkflowByIDRequest) (*flow.Workflow, error) {
	panic("implement me")
}

func (w workflowService) GetWorkflows(ctx context.Context, req flow.GetWorkflowsRequest) ([]*flow.Workflow, error) {
	panic("implement me")
}
