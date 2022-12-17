package part1

import (
	"fmt"
	"sort"

	"github.com/davidaf3/advent-of-code-2022/day16/astar"
	v "github.com/davidaf3/advent-of-code-2022/day16/valve"
)

type ValvesState struct {
	minute         int
	valve          string
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

func newValveState(minute int, valve string, openValves map[string]int,
	cost int, h func(*ValvesState, astar.Problem[*ValvesState]) int, problem *ValvesProblem) *ValvesState {
	hash := fmt.Sprintf("%d%s", minute, valve)

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
		minute:     minute,
		valve:      valve,
		openValves: openValves,
		cost:       cost,
		hash:       hash,
	}

	state.heuristicValue = h(state, problem)
	return state
}

type ValvesProblem struct {
	valves map[string]*v.Valve
	sorted []*v.Valve
}

func (p *ValvesProblem) GetInitialState(h func(*ValvesState, astar.Problem[*ValvesState]) int) *ValvesState {
	return newValveState(0, "AA", make(map[string]int), 0, h, p)
}

func (p *ValvesProblem) IsFinal(state *ValvesState) bool {
	return state.minute == 30
}

func (p *ValvesProblem) Expand(state *ValvesState, h func(*ValvesState, astar.Problem[*ValvesState]) int) []*ValvesState {
	var children []*ValvesState

	nextCost := state.cost
	for valve := range state.openValves {
		nextCost += p.valves[valve].FlowRate
	}

	if _, open := state.openValves[state.valve]; !open && p.valves[state.valve].FlowRate > 0 {
		openValves := astar.CopyMap(state.openValves)
		openValves[state.valve] = state.minute + 1
		children = append(children,
			newValveState(state.minute+1, state.valve, openValves, nextCost, h, p))
	}

	neighbours := p.valves[state.valve].Neighbours
	for _, neighbour := range neighbours {
		openValves := astar.CopyMap(state.openValves)
		children = append(children,
			newValveState(state.minute+1, neighbour, openValves, nextCost, h, p))
	}

	return children
}

func NewValvesProblem(valves map[string]*v.Valve) *ValvesProblem {
	var sorted []*v.Valve

	for _, valve := range valves {
		sorted = append(sorted, valve)
	}

	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].FlowRate > sorted[j].FlowRate
	})

	return &ValvesProblem{valves, sorted}
}

func H(state *ValvesState, problem astar.Problem[*ValvesState]) int {
	estimatedCost := 0
	var remaining []*v.Valve

	for _, valve := range problem.(*ValvesProblem).sorted {
		if _, open := state.openValves[valve.Name]; open {
			estimatedCost += valve.FlowRate * (30 - state.minute)
		} else {
			remaining = append(remaining, valve)
		}
	}

	minute := state.minute
	for _, valve := range remaining {
		estimatedCost += valve.FlowRate * (30 - minute)
		minute += 2
		if minute > 30 {
			break
		}
	}

	return estimatedCost
}
