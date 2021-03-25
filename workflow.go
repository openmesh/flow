package flow

import "context"

type Workflow struct {
	Triggers []Trigger
	Name     string
	Steps    []Step
}

type Step struct {
	Name         string
	Dependencies []string
	Action       Action
}

type WorkflowService interface {
	CreateWorkflow(ctx context.Context, req CreateWorkflowRequest) (*Workflow, error)
	UpdateWorkflow(ctx context.Context, req UpdateWorkflowRequest) (*Workflow, error)
	DeleteWorkflow(ctx context.Context, req DeleteWorkflowRequest) error
	GetWorkflowByID(ctx context.Context, req GetWorkflowByIDRequest) (*Workflow, error)
	GetWorkflows(ctx context.Context, req GetWorkflowsRequest) ([]*Workflow, int, error)
}

type CreateWorkflowRequest struct{}

type UpdateWorkflowRequest struct{}

type DeleteWorkflowRequest struct{}

type GetWorkflowByIDRequest struct{}

type GetWorkflowsRequest struct{}
