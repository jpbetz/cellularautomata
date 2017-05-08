package grid

import (
	"container/heap"
	"math"
)

type Node interface {
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

type nodeToPriorityQueueNodeMap map[Node]*priorityQueueNode

func (d nodeToPriorityQueueNodeMap) GetFromStartOrInf(key Node) (result float64) {
	if v, ok := d[key]; ok {
		return v.fromStartScore
	} else {
		return math.Inf(1)
	}
}

type HuristicCostEstimateFunc func (p1, p2 Node) float64

func FindPath(start, goal Node, estimateCost HuristicCostEstimateFunc) (*Path, bool) {
	closedSet := make(map[Node]bool)

	openSet := make(nodeToPriorityQueueNodeMap)
	open := make(priorityQueue, 0)
	heap.Init(&open)

	startCandidate := &priorityQueueNode{node: start, toGoalScoreViaCell: estimateCost(start, goal), fromStartScore: 0}
	heap.Push(&open, startCandidate)
	openSet[start] = startCandidate

	cameFrom := make(map[Node]Node)

	for open.Len() > 0 {
		current := heap.Pop(&open).(*priorityQueueNode)
		if current.node == goal {
			// build the Nodes
			path := make([]Node, 1)
			path[0] = current.node
			var from, ok = cameFrom[current.node]
			for ok {
				path = append(path, from)
				from, ok = cameFrom[from]
			}
			return &Path{path}, true
		}
		delete(openSet, current.node)
		closedSet[current.node] = true
		for _, neighbor := range current.node.GetNeighbors() {
			neighborNode := neighbor.GetNode()
			if closedSet[neighborNode] != true {
				tentativeFromStartScore := openSet.GetFromStartOrInf(current.node) + neighbor.GetDistance()

				var neighborCandidate, ok = openSet[neighborNode]
				if !ok {
					neighborCandidate = &priorityQueueNode{
						node:               neighborNode,
						toGoalScoreViaCell: math.Inf(1),
						fromStartScore:     math.Inf(1),
					}
					openSet[neighborNode] = neighborCandidate
					heap.Push(&open, neighborCandidate)
				} else if tentativeFromStartScore >= neighborCandidate.fromStartScore {
					// not a better Nodes
					continue
				}
				cameFrom[neighborNode] = current.node
				neighborCandidate.fromStartScore = tentativeFromStartScore
				neighborCandidate.toGoalScoreViaCell = tentativeFromStartScore + estimateCost(neighborNode, goal)
				heap.Fix(&open, neighborCandidate.index)
			}
		}
	}
	return nil, false
}