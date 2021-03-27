package flow

import (
	"context"
	"github.com/google/uuid"
	"time"
)

type Workflow struct {
	ID          uuid.UUID `json:"id" db:"id,omitempty"`
	UserID      uuid.UUID `json:"-" db:"user_id"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	Nodes       []*Node   `json:"nodes" db:"-"`
}

type WorkflowService interface {
	GetWorkflowByID(ctx context.Context, id uuid.UUID) (*Workflow, error)
	GetWorkflows(ctx context.Context, filter WorkflowFilter) ([]*Workflow, int, error)
	CreateWorkflow(ctx context.Context, workflow *Workflow) error
	UpdateWorkflow(ctx context.Context, id uuid.UUID, upd WorkflowUpdate) (*Workflow, error)
	DeleteWorkflow(ctx context.Context, uuid uuid.UUID) error
}

type WorkflowUpdate struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
	// TODO maybe add nodes to the update
}

type WorkflowFilter struct {
	Offset      int     `json:"page"`
	Limit       int     `json:"limit"`
	Name        *string `json:"name"`
	Description *string `json:"description"`
}
