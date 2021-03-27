package flow

import (
	"context"
	"github.com/google/uuid"
	"time"
)

type Node struct {
	ID          uuid.UUID    `json:"id" db:"id,omitempty"`
	CreatedAt   time.Time    `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at" db:"updated_at"`
	WorkflowID  uuid.UUID    `json:"workflow_id" db:"workflow_id"`
	Integration string       `json:"integration" db:"integration"`
	Action      string       `json:"action" db:"action"`
	Params      []*Param     `json:"params" db:"-"`
	ParentIDs   []*uuid.UUID `json:"parent_ids" db:"-"`
	ChildrenIDs []*uuid.UUID `json:"children_ids" db:"-"`
}

type Edge struct {
	HeadID uuid.UUID `json:"head_id" db:"head_id"`
	Head   *Node     `json:"head" db:"-"`
	TailID uuid.UUID `json:"tail_id" db:"tail_id"`
	Tail   *Node     `json:"tail" db:"-"`
}

type ParamType string

const (
	ParamTypeValue     ParamType = "value"
	ParamTypeReference           = "reference"
)

type Param struct {
	ID        uuid.UUID `json:"id" db:"id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	Key       string    `json:"key" db:"key"`
	Value     string    `json:"value" db:"value"`
	Type      ParamType `json:"type" db:"type"`
}

type NodeService interface {
	GetNodeByID(ctx context.Context, id uuid.UUID) (*Node, error)
	GetNodes(ctx context.Context, filter NodeFilter) ([]*Node, int, error)
	CreateNode(ctx context.Context, node *Node) (*Node, error)
	UpdateNode(ctx context.Context, upd NodeUpdate) (*Node, error)
	DeleteNode(ctx context.Context, id uuid.UUID) error
}

type NodeUpdate struct {
	ID          uuid.UUID `json:"id"`
	WorkflowID  uuid.UUID `json:"workflow_id"`
	Integration string    `json:"integration"`
	Action      string    `json:"action"`
	// TODO look at allowing users to update params
}

type DeleteNodeRequest struct {
	ID uuid.UUID `json:"id"`
}

type NodeFilter struct {
	WorkflowID *uuid.UUID `json:"workflow_id"`
	Offset     int        `json:"offset"`
	Limit      int        `json:"limit"`
}
