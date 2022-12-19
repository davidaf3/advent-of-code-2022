package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"

	"github.com/davidaf3/advent-of-code-2022/day19/astar"
)

func triangularNumber(n int) int {
	return n * (n + 1) / 2
}

type Blueprint struct {
	id    int
	costs [4][3]int
}

type RobotsState struct {
	minute            int
	robots, materials [4]int
	hash              string
	heuristic         int
}

func (s *RobotsState) GetCost() int {
	return s.materials[3]
}

func (s *RobotsState) GetHeuristicValue() int {
	return s.heuristic
}

func (s *RobotsState) GetHash() string {
	return s.hash
}

func newRobotState(minute int, robots, materials [4]int,
	h func(*RobotsState, astar.Problem[*RobotsState]) int,
	problem *RobotsProblem) *RobotsState {
	materialsCopy := [4]int{}
	robotsCopy := [4]int{}
	copy(materialsCopy[:], materials[:])
	copy(robotsCopy[:], robots[:])

	hash := strconv.Itoa(minute)
	for i := range materials {
		hash += fmt.Sprintf("-%d,%d", materials[i], robots[i])
	}

	state := &RobotsState{
		minute:    minute,
		robots:    robotsCopy,
		materials: materialsCopy,
		hash:      hash,
	}

	state.heuristic = h(state, problem)
	return state
}

type RobotsProblem struct {
	maxMinutes int
	blueprint  *Blueprint
}

func (p *RobotsProblem) GetInitialState(
	h func(*RobotsState, astar.Problem[*RobotsState]) int) *RobotsState {
	return newRobotState(0, [4]int{1}, [4]int{}, h, p)
}

func (p *RobotsProblem) IsFinal(s *RobotsState) bool {
	return s.minute == p.maxMinutes
}

func (p *RobotsProblem) Expand(s *RobotsState,
	h func(*RobotsState, astar.Problem[*RobotsState]) int) []*RobotsState {
	children := []*RobotsState{}
	nextMaterials := [4]int{}
	for i, n := range s.robots {
		nextMaterials[i] = s.materials[i] + n
	}

	children = append(children,
		newRobotState(s.minute+1, s.robots, nextMaterials, h, p))

	for i, cost := range p.blueprint.costs {
		canCreate := true
		for j, material := range cost {
			canCreate = canCreate && s.materials[j] >= material
		}

		if canCreate {
			s.robots[i] += 1
			for j, material := range cost {
				nextMaterials[j] -= material
			}

			children = append(children,
				newRobotState(s.minute+1, s.robots, nextMaterials, h, p))

			s.robots[i] -= 1
			for j, material := range cost {
				nextMaterials[j] += material
			}
		}
	}

	return children
}

func h(s *RobotsState, p astar.Problem[*RobotsState]) int {
	remainingMins := p.(*RobotsProblem).maxMinutes - s.minute
	return s.robots[3]*remainingMins + triangularNumber(remainingMins)
}

var numberRegex *regexp.Regexp = regexp.MustCompile("[0-9]+")

func parseBlueprint(scanner *bufio.Scanner) *Blueprint {
	input := numberRegex.FindAllString(scanner.Text(), -1)
	numbers := [7]int{}
	for i, numberStr := range input {
		number, err := strconv.Atoi(numberStr)
		errorHandler(err)
		numbers[i] = number
	}

	costs := [4][3]int{
		{numbers[1]},
		{numbers[2]},
		{numbers[3], numbers[4]},
		{numbers[5], 0, numbers[6]},
	}

	return &Blueprint{numbers[0], costs}
}

func getQualityLevelSum(blueprints []*Blueprint) int {
	sum := 0
	for _, blueprint := range blueprints {
		problem := &RobotsProblem{24, blueprint}
		state := astar.AStar[*RobotsState](problem, h)
		sum += blueprint.id * state.GetCost()
	}

	return sum
}

func getMaxGeodesProduct(blueprints []*Blueprint) int {
	product := 1
	for _, blueprint := range blueprints {
		problem := &RobotsProblem{32, blueprint}
		state := astar.AStar[*RobotsState](problem, h)
		product *= state.GetCost()
	}

	return product
}

func main() {
	f, err := os.Open("input.txt")
	errorHandler(err)
	defer f.Close()

	blueprints := []*Blueprint{}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		blueprints = append(blueprints, parseBlueprint(scanner))
	}

	fmt.Println(getQualityLevelSum(blueprints))
	fmt.Println(getMaxGeodesProduct(blueprints[:3]))
}

func errorHandler(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
