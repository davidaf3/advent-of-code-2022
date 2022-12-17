package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"

	"github.com/davidaf3/advent-of-code-2022/day16/astar"
	"github.com/davidaf3/advent-of-code-2022/day16/part1"
	"github.com/davidaf3/advent-of-code-2022/day16/part2"
	v "github.com/davidaf3/advent-of-code-2022/day16/valve"
)

var numberRegex *regexp.Regexp = regexp.MustCompile("[0-9]+")
var valveRegex *regexp.Regexp = regexp.MustCompile("[A-Z][A-Z]")

func parseValves(scanner *bufio.Scanner) map[string]*v.Valve {
	valves := make(map[string]*v.Valve)

	for scanner.Scan() {
		valvesInLine := valveRegex.FindAllString(scanner.Text(), -1)
		flowRateStr := numberRegex.FindAllString(scanner.Text(), -1)
		flowRate, err := strconv.Atoi(flowRateStr[0])
		errorHandler(err)
		valves[valvesInLine[0]] = &v.Valve{
			Name:       valvesInLine[0],
			FlowRate:   flowRate,
			Neighbours: valvesInLine[1:],
		}
	}

	return valves
}

func main() {
	f, err := os.Open("input.txt")
	errorHandler(err)
	defer f.Close()

	scanner := bufio.NewScanner(f)
	valves := parseValves(scanner)

	result := astar.AStar[*part1.ValvesState](part1.NewValvesProblem(valves), part1.H)
	fmt.Println(result.GetCost())

	result = astar.AStar[*part2.ValvesState](part2.NewValvesProblem(valves), part2.H)
	fmt.Println(result.GetCost())
}

func errorHandler(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
