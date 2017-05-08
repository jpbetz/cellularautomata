package grid

import (
	"testing"
	"fmt"
)

type TestPath struct {
	TestNode Node
	Distance float64
}

func (p TestPath) GetNode() Node {
	return p.TestNode
}

func (p TestPath) GetDistance() float64 {
	return p.Distance
}

type TestNode struct {
	Position Position
	Neighbors []*TestNode
}

func (n TestNode) Id() NodeId {
	return n.Position
}

func (n TestNode) GetNeighbors() []Neighbor {
	results := make([]Neighbor, 0, len(n.Neighbors))
	for _, neighbor := range n.Neighbors {
		distance := n.Position.DistanceTo(neighbor.Position)
		results = append(results, TestPath{TestNode: neighbor, Distance: distance})
	}
	return results
}

func testNode(x, y int) *TestNode {
	return &TestNode{ Position: Position {x, y}, Neighbors: make([]*TestNode, 0)}
}

func link(n1, n2 *TestNode) {
	n1.Neighbors = append(n1.Neighbors, n2)
	n2.Neighbors = append(n2.Neighbors, n1)
}

type NodeType int

const (
	O NodeType = iota
	B
	S
	G
)

func toNodes(types [][]NodeType) (*TestNode, *TestNode) {
	var start *TestNode
	var goal *TestNode

	var results = make([][]*TestNode, len(types))
	// create the nodes, unconnected
	for i := 0; i < len(types); i++ {
		row := make([]*TestNode, len(types[i]))
		for j := 0; j < len(types[i]); j++ {
			switch types[i][j] {
			case O:
				row[j] = testNode(i, j)
			case B:
				row[j] = nil
			case S:
				start = testNode(i, j)
				row[j] = start
			case G:
				goal = testNode(i, j)
				row[j] = goal
			}
		}
		results[i] = row
	}

	// connect the nodes
	for i := 0; i < len(results); i++ {
		for j := 0; j < len(results[i]); j++ {
			for k := i-1; k <= i+1; k++ {
				for l := j-1; l <= j+1; l++ {
					validK := k >= 0 && k < len(results)
					validL := l >= 0 && l < len(results[i])
					if validK && validL && results[k][l] != nil && results[i][j] != nil && !(k == i && l == j) {
						results[i][j].Neighbors = append(results[i][j].Neighbors, results[k][l])
					}
				}
			}
		}
	}
	return start, goal
}

func distance(n1, n2 Node) float64 {
	return n1.(*TestNode).Position.DistanceTo(n2.(*TestNode).Position)
}

func TestSimplePath(t *testing.T) {
	Node1 := testNode(0, 0)
	Node2 := testNode(1, 0)
	Node3 := testNode(2, 0)
	Node4 := testNode(3, 0)

	link(Node1, Node2)
	link(Node2, Node3)
	link(Node3, Node4)

	path, ok := FindPath(Node1, Node4, distance)
	if !ok {
		t.Error("Expected FindPath to return ok=true, but got ok=false.")
	}
	route := path.Nodes
	t.Logf("%v %v %v %v", route[0], route[1], route[2], route[3])
	if len(route) != 4 {
		t.Error("Expected Nodes to be length 4")
	}
	if route[0].(*TestNode) != Node4 {
		t.Error("Expected Nodes to end at Node4")
	}

	if route[3].(*TestNode) != Node1 {
		t.Error("Expected Nodes to start at Node1")
	}
}

func TestGrid(t *testing.T) {
	start, goal := toNodes([][]NodeType{
		{S, O, O, O},
		{B, B, B, O},
		{O, O, O, O},
		{O, B, B, B},
		{O, O, O, G},
	})

	path, ok := FindPath(start, goal, distance)
	if !ok {
		t.Error("Expected FindPath to return ok=true, but got ok=false.")
	} else {
		expected := []Position {
			{4, 3},
			{4, 2},
			{4, 1},
			{3, 0},
			{2, 1},
			{2, 2},
			{1, 3},
			{0, 2},
			{0, 1},
			{0, 0},
		}
		route := path.Nodes
		for i, n := range route {
			if n.(*TestNode).Position != expected[i] {
				t.Error(fmt.Sprintf("Expected path[%d] to be %v but found %v", i, expected[i], n.(*TestNode).Position))
			}
		}
	}
}

func TestShortLong(t *testing.T) {
	start, goal := toNodes([][]NodeType{
		{S, O, O, O},
		{O, B, B, O},
		{G, B, B, O},
		{O, B, B, O},
		{O, O, O, O},
	})

	path, ok := FindPath(start, goal, distance)
	if !ok {
		t.Error("Expected FindPath to return ok=true, but got ok=false.")
	} else {
		expected := []Position {
			{2, 0},
			{1, 0},
			{0, 0},
		}
		route := path.Nodes
		for i, n := range route {
			if n.(*TestNode).Position != expected[i] {
				t.Error(fmt.Sprintf("Expected path[%d] to be %v but found %v. Actual path: %#v", i, expected[i], n.(*TestNode).Position, route))
			}
		}
	}
}

func TestAngle(t *testing.T) {
	start, goal := toNodes([][]NodeType{
		{S, O, O},
		{O, O, O},
		{O, O, G},
	})

	path, ok := FindPath(start, goal, distance)
	if !ok {
		t.Error("Expected FindPath to return ok=true, but got ok=false.")
	} else {
		expected := []Position {
			{2, 2},
			{1, 1},
			{0, 0},
		}
		route := path.Nodes
		for i, n := range route {
			if n.(*TestNode).Position != expected[i] {
				t.Error(fmt.Sprintf("Expected path[%d] to be %v but found %v. Actual path: %#v", i, expected[i], n.(*TestNode).Position, route))
			}
		}
	}
}