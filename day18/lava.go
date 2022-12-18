package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type Cube struct {
	coords       [3]int
	blockedSides int
	exposedSides int
}

func (c *Cube) isNextTo(other *Cube, coord int) bool {
	for i := 0; i < 3; i++ {
		if i != coord && c.coords[i] != other.coords[i] {
			return false
		}
	}

	diff := c.coords[coord] - other.coords[coord]
	return diff == 1 || diff == -1
}

func (c *Cube) updateBlocked(other *Cube) {
	for i := 0; i < 3; i++ {
		if c.isNextTo(other, i) {
			c.blockedSides++
			other.blockedSides++
		}
	}
}

func getSurfaceArea(cubes []*Cube) int {
	area := 0
	for _, cube := range cubes {
		area += 6 - cube.blockedSides
	}

	return area
}

func getExteriorSurfaceArea(cubes []*Cube, min, max [3]int) int {
	ranges := [3]int{max[0] - min[0] + 2, max[1] - min[1] + 2, max[2] - min[2] + 2}

	cubeMap := make(map[string]*Cube)
	for _, cube := range cubes {
		i := cube.coords[0] - min[0] + 1
		j := cube.coords[1] - min[1] + 1
		k := cube.coords[2] - min[2] + 1
		coords := fmt.Sprintf("%d,%d,%d", i, j, k)
		cubeMap[coords] = cube
	}

	visited := make(map[string]bool)
	frontier := [][3]int{{0, 0, 0}}

	deltas := [6][3]int{
		{1, 0, 0},
		{-1, 0, 0},
		{0, 1, 0},
		{0, -1, 0},
		{0, 0, 1},
		{0, 0, -1},
	}

	for len(frontier) > 0 {
		air := frontier[len(frontier)-1]
		frontier = frontier[:len(frontier)-1]
		i, j, k := air[0], air[1], air[2]

		airStr := fmt.Sprintf("%d,%d,%d", i, j, k)
		if _, alreadyVisited := visited[airStr]; alreadyVisited {
			continue
		}

		visited[airStr] = true

		for _, delta := range deltas {
			iN, jN, kN := i+delta[0], j+delta[1], k+delta[2]
			nextAirStr := fmt.Sprintf("%d,%d,%d", iN, jN, kN)

			if iN >= 0 && iN <= ranges[0] && jN >= 0 && jN <= ranges[1] &&
				kN >= 0 && kN <= ranges[2] {
				if cube, ok := cubeMap[nextAirStr]; ok {
					cube.exposedSides += 1
				} else {
					frontier = append(frontier, [3]int{iN, jN, kN})
				}
			}
		}
	}

	area := 0
	for _, cube := range cubes {
		area += cube.exposedSides
	}

	return area
}

func parseCube(scanner *bufio.Scanner) *Cube {
	splitted := strings.Split(scanner.Text(), ",")
	var coords [3]int

	for i := 0; i < 3; i++ {
		coord, err := strconv.Atoi(splitted[i])
		errorHandler(err)
		coords[i] = coord
	}

	return &Cube{coords, 0, 0}
}

func main() {
	f, err := os.Open("input.txt")
	errorHandler(err)
	defer f.Close()

	var cubes []*Cube
	minCoords := [3]int{99, 99, 99}
	maxCoords := [3]int{0, 0, 0}
	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		cube := parseCube(scanner)

		for i := 0; i < 3; i++ {
			if cube.coords[i] < minCoords[i] {
				minCoords[i] = cube.coords[i]
			}

			if cube.coords[i] > maxCoords[i] {
				maxCoords[i] = cube.coords[i]
			}
		}

		for _, other := range cubes {
			cube.updateBlocked(other)
		}

		cubes = append(cubes, cube)
	}

	fmt.Println(getSurfaceArea(cubes))
	fmt.Println(getExteriorSurfaceArea(cubes, minCoords, maxCoords))
}

func errorHandler(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
