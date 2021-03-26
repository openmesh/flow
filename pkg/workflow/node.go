package workflow

import (
	"fmt"
	"github.com/google/uuid"
)

// Node type implements a node of a directed acyclic graph (DAG)
type Node struct {
	ID       uuid.UUID
	Value    interface{}
	Parents  map[uuid.UUID]*Node
	Children map[uuid.UUID]*Node
}

// NewNode creates a new node.
func NewNode(value interface{}) *Node {
	n := &Node{
		ID:       uuid.New(),
		Parents:  make(map[uuid.UUID]*Node),
		Children: make(map[uuid.UUID]*Node),
		Value:    value,
	}

	return n
}

// Degree returns the number of parents and children of the node.
func (n *Node) Degree() int {
	return len(n.Parents) + len(n.Children)
}

// InDegree return the number of parents of the vertex or the number of edges
// entering on it.
func (n *Node) InDegree() int {
	return len(n.Parents)
}

// OutDegree return the number of children of the vertex or the number of edges
// leaving it.
func (n *Node) OutDegree() int {
	return len(n.Children)
}

// String implements stringer interface and prints an string representation
// of this instance.
func (n *Node) String() string {
	result := fmt.Sprintf("ID: %s - Parents: %d - Children: %d - Value: %v\n", n.ID, n.InDegree(), n.OutDegree(), n.Value)

	return result
}

func (n *Node) hasChild(c *Node) bool {
	for i := range n.Children {
		if n.Children[i].ID == c.ID {
			return true
		}
	}
	return false
}

func (n *Node) hasParent(c *Node) bool {
	for i := range n.Parents {
		if n.Parents[i].ID == c.ID {
			return true
		}
	}
	return false
}
