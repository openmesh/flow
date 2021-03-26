package flow

import (
	"context"
	"github.com/google/uuid"
)

type Workflow struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Nodes       []*Node   `json:"nodes"`
}

type WorkflowService interface {
	CreateWorkflow(ctx context.Context, req CreateWorkflowRequest) (*Workflow, error)
	UpdateWorkflow(ctx context.Context, req UpdateWorkflowRequest) (*Workflow, error)
	DeleteWorkflow(ctx context.Context, req DeleteWorkflowRequest) error
	GetWorkflowByID(ctx context.Context, req GetWorkflowByIDRequest) (*Workflow, error)
	GetWorkflows(ctx context.Context, req GetWorkflowsRequest) ([]*Workflow, int, error)
}

type CreateWorkflowRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type NodeDTO struct {
	Action      string `json:"action"`
	Integration string `json:"integration"`
}

type UpdateWorkflowRequest struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
}

type DeleteWorkflowRequest struct {
	ID uuid.UUID `json:"id"`
}

type GetWorkflowByIDRequest struct {
	ID uuid.UUID `json:"id"`
}

type GetWorkflowsRequest struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
}
