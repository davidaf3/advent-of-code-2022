package part2

import (
	"fmt"
	"sort"

	"github.com/davidaf3/advent-of-code-2022/day16/astar"
	v "github.com/davidaf3/advent-of-code-2022/day16/valve"
)

type ValvesState struct {
	minute         int
	valves         [2]string
	parentValves   [2]string
	openValves     map[string]int
	cost           int
	heuristicValue int
	hash           string
}

func (s *ValvesState) GetCost() int {
	return s.cost
}

func (s *ValvesState) GetHeuristicValue() int {
	return s.heuristicValue
}

func (s *ValvesState) GetHash() string {
	return s.hash
}

func newValveState(minute int, valves, parentValves [2]string, openValves map[string]int,
	cost int, h func(*ValvesState, astar.Problem[*ValvesState]) int, problem *ValvesProblem) *ValvesState {
	hash := fmt.Sprintf("%d%s%s", minute, valves[0], valves[1])

	var sortedValves []string
	for openValve := range openValves {
		sortedValves = append(sortedValves, openValve)
	}

	sort.Slice(sortedValves, func(i, j int) bool {
		return sortedValves[i] < sortedValves[j]
	})

	for _, openValve := range sortedValves {
		hash += fmt.Sprintf("%d%s", openValves[openValve], openValve)
	}

	state := &ValvesState{
		minute:       minute,
		valves:       valves,
		parentValves: parentValves,
		openValves:   openValves,
		cost:         cost,
		hash:         hash,
	}

	state.heuristicValue = h(state, problem)
	return state
}

type ValvesProblem struct {
	valves            map[string]*v.Valve
	sorted            []*v.Valve
	nonZeroFlowValves int
}

func (p *ValvesProblem) GetInitialState(h func(*ValvesState, astar.Problem[*ValvesState]) int) *ValvesState {
	return newValveState(0, [2]string{"AA", "AA"}, [2]string{"AA", "AA"}, make(map[string]int), 0, h, p)
}

func (p *ValvesProblem) IsFinal(state *ValvesState) bool {
	return state.minute == 26
}

func (p *ValvesProblem) Expand(state *ValvesState, h func(*ValvesState, astar.Problem[*ValvesState]) int) []*ValvesState {
	var children []*ValvesState

	nextCost := state.cost
	for valve := range state.openValves {
		nextCost += p.valves[valve].FlowRate
	}

	if len(state.openValves) == p.nonZeroFlowValves {
		return []*ValvesState{
			newValveState(26, state.valves, state.valves, state.openValves, state.cost+(nextCost-state.cost)*(26-state.minute), h, p)}
	}

	_, firstOpen := state.openValves[state.valves[0]]
	_, secondOpen := state.openValves[state.valves[1]]

	if !firstOpen && p.valves[state.valves[0]].FlowRate > 0 {
		for _, neighbour := range p.valves[state.valves[1]].Neighbours {
			if !secondOpen || state.parentValves[1] != neighbour {
				openValves := astar.CopyMap(state.openValves)
				openValves[state.valves[0]] = state.minute + 1
				children = append(children,
					newValveState(state.minute+1, [2]string{state.valves[0], neighbour}, state.valves, openValves, nextCost, h, p))
			}
		}
	}

	if !secondOpen && p.valves[state.valves[1]].FlowRate > 0 {
		for _, neighbour := range p.valves[state.valves[0]].Neighbours {
			if !firstOpen || state.parentValves[0] != neighbour {
				openValves := astar.CopyMap(state.openValves)
				openValves[state.valves[1]] = state.minute + 1
				children = append(children,
					newValveState(state.minute+1, [2]string{neighbour, state.valves[1]}, state.valves, openValves, nextCost, h, p))
			}
		}
	}

	if !firstOpen && !secondOpen && p.valves[state.valves[0]].FlowRate > 0 && p.valves[state.valves[1]].FlowRate > 0 {
		openValves := astar.CopyMap(state.openValves)
		openValves[state.valves[0]] = state.minute + 1
		openValves[state.valves[1]] = state.minute + 1
		children = append(children,
			newValveState(state.minute+1, [2]string{state.valves[0], state.valves[1]}, state.valves, openValves, nextCost, h, p))
	}

	firstNeighbours := p.valves[state.valves[0]].Neighbours
	for _, firstNeighbour := range firstNeighbours {
		secondNeighbours := p.valves[state.valves[1]].Neighbours
		for _, secondNeighbour := range secondNeighbours {
			if (!firstOpen || state.parentValves[0] != firstNeighbour) && (!secondOpen || state.parentValves[1] != secondNeighbour) {
				openValves := astar.CopyMap(state.openValves)
				children = append(children,
					newValveState(state.minute+1, [2]string{firstNeighbour, secondNeighbour}, state.valves, openValves, nextCost, h, p))
			}
		}
	}

	return children
}

func NewValvesProblem(valves map[string]*v.Valve) *ValvesProblem {
	var sorted []*v.Valve
	nonZeroFlowValves := 0

	for _, valve := range valves {
		sorted = append(sorted, valve)
		if valve.FlowRate > 0 {
			nonZeroFlowValves++
		}
	}

	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].FlowRate > sorted[j].FlowRate
	})

	return &ValvesProblem{valves, sorted, nonZeroFlowValves}
}

func H(state *ValvesState, problem astar.Problem[*ValvesState]) int {
	estimatedCost := 0
	var remaining []*v.Valve

	for _, valve := range problem.(*ValvesProblem).sorted {
		if _, open := state.openValves[valve.Name]; open {
			estimatedCost += valve.FlowRate * (26 - state.minute)
		} else {
			remaining = append(remaining, valve)
		}
	}

	minute := state.minute
	steps := 0
	for _, valve := range remaining {
		steps++
		estimatedCost += valve.FlowRate * (26 - minute)
		if steps == len(state.valves) {
			minute += 2
			if minute > 26 {
				break
			}
			steps = 0
		}
	}

	return estimatedCost
}
