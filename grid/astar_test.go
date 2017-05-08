package grid

import "testing"

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


func (n TestNode) GetNeighbors() []Neighbor {
	results := make([]Neighbor, 0, len(n.Neighbors))
	for _, neighbor := range n.Neighbors {
		distance := n.Position.DistanceTo(neighbor.Position)
		results = append(results, TestPath{TestNode: neighbor, Distance: distance})
	}
	return results
}

var Node1 = &TestNode{ Position: Position {0, 0}, Neighbors: make([]*TestNode, 0)}
var Node2 = &TestNode{ Position: Position {1, 0}, Neighbors: make([]*TestNode, 0)}
var Node3 = &TestNode{ Position: Position {2, 0}, Neighbors: make([]*TestNode, 0)}
var Node4 = &TestNode{ Position: Position {3, 0}, Neighbors: make([]*TestNode, 0)}

func init() {
	Node1.Neighbors = append(Node1.Neighbors, Node2)
	Node2.Neighbors = append(Node2.Neighbors, Node1, Node3)
	Node3.Neighbors = append(Node3.Neighbors, Node2, Node4)
	Node4.Neighbors = append(Node4.Neighbors, Node3)
}

func distance(n1, n2 Node) float64 {
	return n1.(*TestNode).Position.DistanceTo(n2.(*TestNode).Position)
}

func TestSimplePath(t *testing.T) {
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
