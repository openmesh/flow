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
	// Triggers []Trigger
	// Name     string
	// Steps    []Step
}

type Node struct {
	ID          uuid.UUID `json:"id"`
	Integration string    `json:"integration"`
	Action      string    `json:"action"`
	IsSource    bool      `json:"is_source"`
	// TODO decide if this is a good name
	Parameters []*Parameter `json:"parameters"`
	// inputs
	// outputs
}

type Edge struct {
	HeadID uuid.UUID `json:"head_id"`
	Head   *Node     `json:"head"`
	TailID uuid.UUID `json:"tail_id"`
	Tail   *Node     `json:"tail"`
}

type Parameter struct {
	ID   uuid.UUID `json:"id"`
	Path string    `json:"path"`
}

//type Step struct {
//	Name         string
//	Dependencies []string
//	Action       Action
//}

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
