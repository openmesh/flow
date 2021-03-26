package workflow

import (
	"fmt"
	"github.com/google/uuid"
	"sync"
)

type Graph struct {
	mu    sync.Mutex
	nodes map[uuid.UUID]*Node
}

func NewGraph() *Graph {
	return &Graph{
		mu:    sync.Mutex{},
		nodes: make(map[uuid.UUID]*Node),
	}
}

// AddNode adds a node to the workflow
func (g *Graph) AddNode(n *Node) error {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.nodes[n.ID] = n
	return nil
}

// DeleteNode deletes a node and all the edges referencing it from the graph.
func (g *Graph) DeleteNode(n *Node) error {
	g.mu.Lock()
	g.mu.Unlock()

	// Check that node exists
	if !g.containsNode(n) {
		return fmt.Errorf("node with ID %v could not be found", n.ID)
	}

	delete(g.nodes, n.ID)

	return nil
}

// AddEdge adds a directed edge between two existing nodes of the graph.
func (g *Graph) AddEdge(tail *Node, head *Node) error {
	g.mu.Lock()
	g.mu.Unlock()

	if !g.containsNode(tail) || !g.containsNode(head) {
		return fmt.Errorf("node with ID %v not found", tail)
	}

	if tail.hasChild(head) {
		return fmt.Errorf("edge (%v,%v) already exists", tail.ID, head.ID)
	}

	// Add edge
	tail.Children[head.ID] = head
	head.Parents[tail.ID] = tail

	return nil
}

// DeleteEdge deletes a directed edge between two existing nodes from the
// graph.
func (g *Graph) DeleteEdge(tail *Node, head *Node) error {
	for _, child := range tail.Children {
		if child.ID == head.ID {
			delete(tail.Children, child.ID)
		}
	}
	return nil
}

// GetNode returns a node from the graph with a given ID.
func (g *Graph) GetNode(id uuid.UUID) (*Node, error) {
	n, found := g.nodes[id]
	if !found {
		return n, fmt.Errorf("node %s not found in the graph", id)
	}

	return n, nil
}

// Order returns the number of nodes in the graph.
func (g *Graph) Order() int {
	return len(g.nodes)
}

// Size returns the number of edges in the graph.
func (g *Graph) Size() int {
	count := 0
	for key := range g.nodes {
		count = count + len(g.nodes[key].Children)
	}
	return count
}

// SinkNodes returns nodes with no children defined by the graph edges.
func (g *Graph) SinkNodes() []*Node {
	var sinkNodes []*Node

	for key := range g.nodes {
		if len(g.nodes[key].Children) == 0 {
			sinkNodes = append(sinkNodes, g.nodes[key])
		}
	}

	return sinkNodes
}

// SourceNodes return vertices with no parent defined by the graph edges.
func (g *Graph) SourceNodes() []*Node {
	var sourceVertices []*Node

	for key := range g.nodes {
		if len(g.nodes[key].Parents) == 0 {
			sourceVertices = append(sourceVertices, g.nodes[key])
		}
	}

	return sourceVertices
}

// Successors return nodes that are children of a given node.
func (g *Graph) Successors(node *Node) ([]*Node, error) {
	var successors []*Node

	_, found := g.GetNode(node.ID)
	if found != nil {
		return successors, fmt.Errorf("node %s not found in the graph", node.ID)
	}

	for _, n := range node.Children {
		successors = append(successors, n)
	}

	return successors, nil
}

// Predecessors return nodes that are a parent of a given node.
func (g *Graph) Predecessors(node *Node) ([]*Node, error) {
	var predecessors []*Node

	_, found := g.GetNode(node.ID)
	if found != nil {
		return predecessors, fmt.Errorf("node %s not found in the graph", node.ID)
	}

	for _, n := range node.Parents {
		predecessors = append(predecessors, n)
	}

	return predecessors, nil
}

// String implements stringer interface.
//
// Prints an string representation of this instance.
func (g *Graph) String() string {
	result := fmt.Sprintf("DAG Nodes: %d - Edges: %d\n", g.Order(), g.Size())
	result += fmt.Sprintf("Vertices:\n")
	for _, node := range g.nodes {
		result += fmt.Sprintf("%s", node)
	}

	return result
}

// containsNode returns a bool representing whether or not the graph contains a node
// with an ID equal to the given node's ID.
func (g *Graph) containsNode(n *Node) bool {
	for key := range g.nodes {
		if key == n.ID {
			return true
		}
	}
	return false
}
