package workflow

import (
	"github.com/google/uuid"
)

// Workflow implements a workflow. A workflow makes use of a directed acyclic graph data structure
type Workflow struct {
	ID          uuid.UUID
	Name        string
	Description string
	Graph       Graph
}
