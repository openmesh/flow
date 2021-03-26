package flow

import (
	"context"
	"github.com/google/uuid"
)

type Node struct {
	ID          uuid.UUID    `json:"id"`
	WorkflowID  uuid.UUID    `json:"workflow_id"`
	Integration string       `json:"integration"`
	Action      string       `json:"action"`
	Params      []*Param     `json:"params"`
	ParentIDs   []*uuid.UUID `json:"parent_ids"`
	ChildrenIDs []*uuid.UUID `json:"children_ids"`
}

type Edge struct {
	HeadID uuid.UUID `json:"head_id"`
	Head   *Node     `json:"head"`
	TailID uuid.UUID `json:"tail_id"`
	Tail   *Node     `json:"tail"`
}

type ParamType string

const (
	ParamTypeValue     ParamType = "value"
	ParamTypeReference           = "reference"
)

type Param struct {
	ID    uuid.UUID `json:"id"`
	Key   string    `json:"key"`
	Value string    `json:"value"`
	Type  ParamType `json:"type"`
}

type NodeService interface {
	CreateNode(ctx context.Context, req CreateNodeRequest) (*Node, error)
	UpdateNode(ctx context.Context, req UpdateNodeRequest) (*Node, error)
	DeleteNode(ctx context.Context, req DeleteNodeRequest) error
}

type CreateNodeRequest struct {
	WorkflowID  uuid.UUID    `json:"workflow_id"`
	Integration string       `json:"integration"`
	Action      string       `json:"action"`
	Params      []*ParamDTO  `json:"params"`
	ParentIDs   []*uuid.UUID `json:"parent_ids"`
	ChildrenIDs []*uuid.UUID `json:"children_ids"`
}

type ParamDTO struct {
	Key   string    `json:"key"`
	Value string    `json:"value"`
	Type  ParamType `json:"type"`
}

type UpdateNodeRequest struct {
	ID          uuid.UUID   `json:"id"`
	WorkflowID  uuid.UUID   `json:"workflow_id"`
	Integration string      `json:"integration"`
	Action      string      `json:"action"`
	Params      []*ParamDTO `json:"params"`
}

type DeleteNodeRequest struct {
	ID uuid.UUID `json:"id"`
}
