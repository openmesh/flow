package workflow

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/openmesh/flow/pkg/orderedset"
)

// Node type implements a node of a directed acyclic graph (DAG)
type Node struct {
	ID       uuid.UUID
	Value    interface{}
	Parents  *orderedset.OrderedSet
	Children *orderedset.OrderedSet
}

// NewNode creates a new node.
func NewNode(value interface{}) *Node {
	n := &Node{
		ID:       uuid.New(),
		Parents:  orderedset.NewOrderedSet(),
		Children: orderedset.NewOrderedSet(),
		Value:    value,
	}

	return n
}

// Degree returns the number of parents and children of the node.
func (n *Node) Degree() int {
	return n.Parents.Size() + n.Children.Size()
}

// InDegree return the number of parents of the vertex or the number of edges
// entering on it.
func (n *Node) InDegree() int {
	return n.Parents.Size()
}

// OutDegree return the number of children of the vertex or the number of edges
// leaving it.
func (n *Node) OutDegree() int {
	return n.Children.Size()
}

// String implements stringer interface and prints an string representation
// of this instance.
func (n *Node) String() string {
	result := fmt.Sprintf("ID: %s - Parents: %d - Children: %d - Value: %v\n", n.ID, n.Parents.Size(), n.Children.Size(), n.Value)

	return result
}