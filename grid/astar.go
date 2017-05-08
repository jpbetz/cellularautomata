package grid

import (
	"container/heap"
	"math"
)

type NodeId interface{}

type Node interface {
	Id() NodeId
	GetNeighbors() []Neighbor
}

type Neighbor interface {
	GetNode() Node
	GetDistance() float64
}

type Path struct {
	Nodes []Node
}

type priorityQueueNode struct {
	node               Node
	toGoalScoreViaCell float64
	fromStartScore     float64
	index              int
}

type priorityQueue []*priorityQueueNode

func (pq priorityQueue) Len() int { return len(pq) }

func (pq priorityQueue) Less(i, j int) bool {
	return pq[i].toGoalScoreViaCell < pq[j].toGoalScoreViaCell
}

func (pq priorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *priorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*priorityQueueNode)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *priorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	item.index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}

type HeuristicCostEstimateFunc func(p1, p2 Node) float64

func FindPath(start, goal Node, estimateCost HeuristicCostEstimateFunc) (path *Path, ok bool) {
	// https://en.wikipedia.org/wiki/A*_search_algorithm

	closedSet := make(map[NodeId]bool)
	openSet := make(map[NodeId]*priorityQueueNode)
	openQueue := make(priorityQueue, 0)
	heap.Init(&openQueue)
	startCandidate := &priorityQueueNode{
		node:               start,
		toGoalScoreViaCell: estimateCost(start, goal),
		fromStartScore:     0,
	}
	heap.Push(&openQueue, startCandidate)
	openSet[start.Id()] = startCandidate
	cameFrom := make(map[NodeId]Node)

	for openQueue.Len() > 0 {
		current := heap.Pop(&openQueue).(*priorityQueueNode)
		if current.node == goal {
			return buildPath(cameFrom, current.node), true
		}
		delete(openSet, current.node.Id())
		closedSet[current.node.Id()] = true
		for _, neighbor := range current.node.GetNeighbors() {
			neighborNode := neighbor.GetNode()
			if closedSet[neighborNode.Id()] != true {
				tentativeFromStartScore := current.fromStartScore + neighbor.GetDistance()

				var neighborCandidate, ok = openSet[neighborNode.Id()]
				if !ok {
					neighborCandidate = &priorityQueueNode{
						node:               neighborNode,
						toGoalScoreViaCell: math.Inf(1),
						fromStartScore:     math.Inf(1),
					}
					openSet[neighborNode.Id()] = neighborCandidate
					heap.Push(&openQueue, neighborCandidate)
				} else if tentativeFromStartScore >= neighborCandidate.fromStartScore {
					// not a better Node
					continue
				}
				cameFrom[neighborNode.Id()] = current.node
				neighborCandidate.fromStartScore = tentativeFromStartScore
				neighborCandidate.toGoalScoreViaCell = tentativeFromStartScore + estimateCost(neighborNode, goal)
				heap.Fix(&openQueue, neighborCandidate.index)
			}
		}
	}
	return nil, false
}

func buildPath(cameFrom map[NodeId]Node, current Node) *Path {
	path := make([]Node, 1)
	path[0] = current
	var from, ok = cameFrom[current.Id()]
	for ok {
		path = append(path, from)
		from, ok = cameFrom[from.Id()]
	}
	return &Path{path}
}
