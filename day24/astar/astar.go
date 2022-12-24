package astar

import (
	"container/heap"
)

type MinHeap []State

func (h MinHeap) Len() int {
	return len(h)
}

func (h MinHeap) Less(i, j int) bool {
	return h[i].GetCost()+h[i].GetHeuristicValue() <
		h[j].GetCost()+h[j].GetHeuristicValue()
}

func (h MinHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h *MinHeap) Push(x any) {
	*h = append(*h, x.(State))
}

func (h *MinHeap) Pop() any {
	n := len(*h)
	x := (*h)[n-1]
	*h = (*h)[:n-1]
	return x
}

type State interface {
	GetCost() int
	GetHeuristicValue() int
	GetHash() string
}

type Problem[S State] interface {
	GetInitialState(func(S, Problem[S]) int) S
	IsFinal(S) bool
	Expand(S, func(S, Problem[S]) int) []S
}

func AStar[S State](problem Problem[S], h func(S, Problem[S]) int) State {
	frontier := &MinHeap{}
	heap.Init(frontier)
	heap.Push(frontier, problem.GetInitialState(h))

	visited := make(map[string]bool)

	for frontier.Len() > 0 {
		state := heap.Pop(frontier).(S)
		if _, stateVisited := visited[state.GetHash()]; stateVisited {
			continue
		}

		if problem.IsFinal(state) {
			return state
		}

		visited[state.GetHash()] = true

		for _, child := range problem.Expand(state, h) {
			heap.Push(frontier, child)
		}
	}

	return nil
}
