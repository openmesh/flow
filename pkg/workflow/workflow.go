package workflow

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/openmesh/flow/pkg/orderedmap"
	"sync"
)

// Workflow implements a workflow. A workflow makes use of a directed acyclic graph data structure
type Workflow struct {
	ID          uuid.UUID
	Name        string
	Description string
	mu          sync.Mutex
	nodes       orderedmap.OrderedMap
}

func NewWorkflow(name string, description string) *Workflow {
	return &Workflow{
		ID:          uuid.New(),
		Name:        name,
		Description: description,
		mu:          sync.Mutex{},
		nodes:       orderedmap.OrderedMap{},
	}
}

// AddNode adds a node to the workflow
func (w *Workflow) AddNode(n *Node) error {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.nodes.Put(n.ID, n)
	return nil
}

// DeleteNode deletes a node and all the edges referencing it from the graph.
func (w *Workflow) DeleteNode(n *Node) error {
	w.mu.Lock()
	w.mu.Unlock()

	// Check that node exists
	if !w.nodes.Contains(n) {
		return fmt.Errorf("node with ID %v could not be found", n.ID)
	}

	w.nodes.Remove(n.ID)
	return nil
}

// AddEdge adds a directed edge between two existing nodes of the graph.
func (w *Workflow) AddEdge(tail *Node, head *Node) error {
	if !w.nodes.Contains(tail) || !w.nodes.Contains(head) {
		return fmt.Errorf("node with ID %v not found", tail)
	}

	if tail.Children.Contains(head) {
		return fmt.Errorf("edge (%v,%v) already exists", tail.ID, head.ID)
	}

	// Add edge
	tail.Children.Add(head)
	head.Parents.Add(tail)

	return nil
}

// DeleteEdge deletes a directed edge between two existing vertices from the
// graph.
func (w *Workflow) DeleteEdge(tail *Node, head *Node) error {
	for _, child := range tail.Children.Values() {
		if child == head {
			tail.Children.Remove(child)
		}
	}
	return nil
}

// GetNode returns a node from the graph with a given ID.
func (w *Workflow) GetNode(id interface{}) (*Node, error) {
	var node *Node

	n, found := w.nodes.Get(id)
	if !found {
		return node, fmt.Errorf("node %s not found in the graph", id)
	}

	node = n.(*Node)

	return node, nil
}

// Order returns the number of nodes in the graph.
func (w *Workflow) Order() int {
	return w.nodes.Size()
}

// Size returns the number of edges in the graph.
func (w *Workflow) Size() int {
	count := 0
	for _, n := range w.nodes.Values() {
		count = count + n.(*Node).Children.Size()
	}
	return count
}

// SinkNodes returns nodes with no children defined by the graph edges.
func (w *Workflow) SinkNodes() []*Node {
	var sinkNodes []*Node

	for _, node := range w.nodes.Values() {
		if node.(*Node).Children.Size() == 0 {
			sinkNodes = append(sinkNodes, node.(*Node))
		}
	}

	return sinkNodes
}

// SourceNodes return vertices with no parent defined by the graph edges.
func (w *Workflow) SourceNodes() []*Node {
	var sourceVertices []*Node

	for _, node := range w.nodes.Values() {
		if node.(*Node).Parents.Size() == 0 {
			sourceVertices = append(sourceVertices, node.(*Node))
		}
	}

	return sourceVertices
}

// Successors return nodes that are children of a given node.
func (w *Workflow) Successors(node *Node) ([]*Node, error) {
	var successors []*Node

	_, found := w.GetNode(node.ID)
	if found != nil {
		return successors, fmt.Errorf("node %s not found in the graph", node.ID)
	}

	for _, n := range node.Children.Values() {
		successors = append(successors, n.(*Node))
	}

	return successors, nil
}

// Predecessors return nodes that are a parent of a given node.
func (w *Workflow) Predecessors(node *Node) ([]*Node, error) {
	var predecessors []*Node

	_, found := w.GetNode(node.ID)
	if found != nil {
		return predecessors, fmt.Errorf("node %s not found in the graph", node.ID)
	}

	for _, v := range node.Parents.Values() {
		predecessors = append(predecessors, v.(*Node))
	}

	return predecessors, nil
}

// String implements stringer interface.
//
// Prints an string representation of this instance.
func (w *Workflow) String() string {
	result := fmt.Sprintf("DAG Nodes: %d - Edges: %d\n", w.Order(), w.Size())
	result += fmt.Sprintf("Vertices:\n")
	for _, vertex := range w.nodes.Values() {
		vertex = vertex.(*Node)
		result += fmt.Sprintf("%s", vertex)
	}

	return result
}