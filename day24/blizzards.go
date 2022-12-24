package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/davidaf3/advent-of-code-2022/day23/astar"
)

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

type BlizzardsState struct {
	current   [2]int
	cost      int
	heuristic int
	hash      string
}

func (s *BlizzardsState) GetCost() int {
	return s.cost
}

func (s *BlizzardsState) GetHeuristicValue() int {
	return s.heuristic
}

func (s *BlizzardsState) GetHash() string {
	return s.hash
}

func newBlizzardsState(cost int, current [2]int,
	h func(*BlizzardsState, astar.Problem[*BlizzardsState]) int,
	problem *BlizzardsProblem) *BlizzardsState {
	hash := strconv.Itoa(cost)
	hash += fmt.Sprintf("-%d,%d", current[0], current[1])

	state := &BlizzardsState{
		cost:    cost,
		current: current,
		hash:    hash,
	}

	state.heuristic = h(state, problem)
	return state
}

type BlizzardsProblem struct {
	blizzardDirections [][2]int
	blizzardPositions  [][][2]int
	blizzardMaps       []map[[2]int]bool
	dimensions         [2]int
	initialMinute      int
	start, goal        [2]int
}

func (p *BlizzardsProblem) GetInitialState(
	h func(*BlizzardsState, astar.Problem[*BlizzardsState]) int) *BlizzardsState {
	p.getBlizzardsAt(p.initialMinute)
	return newBlizzardsState(p.initialMinute, p.start, h, p)
}

func (p *BlizzardsProblem) IsFinal(s *BlizzardsState) bool {
	return s.current == p.goal
}

func (p *BlizzardsProblem) getBlizzardsAt(minute int) map[[2]int]bool {
	for minute >= len(p.blizzardMaps) {
		last := p.blizzardPositions[len(p.blizzardPositions)-1]
		nextMap := make(map[[2]int]bool, len(last))
		next := make([][2]int, len(last))
		copy(next, last)

		for i, position := range last {
			direction := p.blizzardDirections[i]
			nextPosition := [2]int{}
			for j := range direction {
				nextPosition[j] = position[j] + direction[j]
				if nextPosition[j] < 1 || nextPosition[j] >= p.dimensions[j]-1 {
					nextPosition[j] = position[j] - (direction[j] * (p.dimensions[j] - 3))
				}
			}

			next[i] = nextPosition
			nextMap[nextPosition] = true
		}

		p.blizzardPositions = append(p.blizzardPositions, next)
		p.blizzardMaps = append(p.blizzardMaps, nextMap)
	}

	return p.blizzardMaps[minute]
}

func (p *BlizzardsProblem) Expand(s *BlizzardsState,
	h func(*BlizzardsState, astar.Problem[*BlizzardsState]) int) []*BlizzardsState {
	children := []*BlizzardsState{}
	blizzards := p.getBlizzardsAt(s.cost + 1)
	moves := [5][2]int{{0, 0}, {1, 0}, {-1, 0}, {0, 1}, {0, -1}}

	for _, move := range moves {
		next := [2]int{}
		for i := range move {
			next[i] = s.current[i] + move[i]
		}

		inBounds := (next[0] > 0 && next[0] < p.dimensions[0]-1 &&
			next[1] > 0 && next[1] < p.dimensions[1]-1) ||
			next == p.start || next == p.goal

		if blizzard, ok := blizzards[next]; inBounds && (!ok || !blizzard) {
			children = append(children, newBlizzardsState(s.cost+1, next, h, p))
		}
	}

	return children
}

func h(s *BlizzardsState, p astar.Problem[*BlizzardsState]) int {
	goal := p.(*BlizzardsProblem).goal
	return abs(s.current[0]-goal[0]) + abs(s.current[1]-goal[1])
}

func parseProblem(scanner *bufio.Scanner) *BlizzardsProblem {
	blizzardDirections := [][2]int{}
	blizzardPositions := make([][][2]int, 1)
	blizzardMaps := make([]map[[2]int]bool, 1)
	blizzardMaps[0] = make(map[[2]int]bool)
	dimensions := [2]int{}
	scanner.Scan()
	dimensions[1] = len(scanner.Bytes())

	i := 0
	for ; scanner.Scan() && scanner.Bytes()[1] != '#'; i++ {
		for j, cell := range scanner.Bytes()[1 : len(scanner.Bytes())-1] {
			if cell != '.' {
				blizzardPositions[0] = append(blizzardPositions[0], [2]int{i + 1, j + 1})
				blizzardMaps[0][[2]int{i + 1, j + 1}] = true

				var direction [2]int
				switch cell {
				case '>':
					direction = [2]int{0, 1}
				case '<':
					direction = [2]int{0, -1}
				case 'v':
					direction = [2]int{1, 0}
				case '^':
					direction = [2]int{-1, 0}
				}

				blizzardDirections = append(blizzardDirections, direction)
			}
		}
	}

	dimensions[0] = i + 2
	return &BlizzardsProblem{
		blizzardDirections: blizzardDirections,
		blizzardPositions:  blizzardPositions,
		blizzardMaps:       blizzardMaps,
		dimensions:         dimensions,
	}
}

func main() {
	f, err := os.Open("input.txt")
	errorHandler(err)
	defer f.Close()

	scanner := bufio.NewScanner(f)
	problem := parseProblem(scanner)
	start := [2]int{0, 1}
	goal := [2]int{problem.dimensions[0] - 1, problem.dimensions[1] - 2}

	problem.initialMinute = 0
	problem.start = start
	problem.goal = goal
	state := astar.AStar[*BlizzardsState](problem, h)
	fmt.Println(state.GetCost())

	problem.initialMinute = state.GetCost()
	problem.start = goal
	problem.goal = start
	state = astar.AStar[*BlizzardsState](problem, h)

	problem.initialMinute = state.GetCost()
	problem.start = start
	problem.goal = goal
	state = astar.AStar[*BlizzardsState](problem, h)
	fmt.Println(state.GetCost())
}

func errorHandler(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
