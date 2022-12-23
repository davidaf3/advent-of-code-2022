package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

type Elf struct {
	current [2]int
	next    *[2]int
}

func (e *Elf) isAlone(grove *Grove) bool {
	for _, moves := range grove.movesToConsider {
		for _, move := range moves {
			if grove.elvesMap[e.current[0]+move[0]][e.current[1]+move[1]] != nil {
				return false
			}
		}
	}

	return true
}

func (e *Elf) considerMoves(moves [3][2]int, grove *Grove) {
	for _, move := range moves {
		if grove.elvesMap[e.current[0]+move[0]][e.current[1]+move[1]] != nil {
			return
		}
	}

	e.next = &[2]int{e.current[0] + moves[1][0], e.current[1] + moves[1][1]}
	if _, ok := grove.nextPositions[*e.next]; ok {
		grove.nextPositions[*e.next] += 1
	} else {
		grove.nextPositions[*e.next] = 1
	}
}

func (e *Elf) move(grove *Grove) {
	if e.next != nil && grove.nextPositions[*e.next] < 2 {
		grove.elvesMap[e.current[0]][e.current[1]] = nil
		grove.elvesMap[e.next[0]][e.next[1]] = e
		e.current = *e.next
	}

	e.next = nil
}

type Grove struct {
	round           int
	elves           []*Elf
	elvesMap        [][]*Elf
	nextPositions   map[[2]int]int
	movesToConsider [4][3][2]int
	firstConsidered int
}

func (g *Grove) simulateUntil(rounds int) {
	for i := 0; i < rounds; i++ {
		g.round++
		g.considerMoves()
		g.move()
	}
}

func (g *Grove) simulateUntilNoMoves() {
	for {
		g.round++
		g.considerMoves()
		if g.move() == 0 {
			return
		}
	}
}

func (g *Grove) considerMoves() {
	g.checkBorders()
	g.nextPositions = make(map[[2]int]int)

	for _, elf := range g.elves {
		if !elf.isAlone(g) {
			for i := 0; i < 4 && elf.next == nil; i++ {
				elf.considerMoves(g.movesToConsider[(g.firstConsidered+i)%4], g)
			}
		}
	}

	g.firstConsidered++
}

func (g *Grove) move() int {
	moved := 0
	for _, elf := range g.elves {
		last := elf.current
		elf.move(g)
		if last != elf.current {
			moved += 1
		}
	}

	return moved
}

func (g *Grove) getEmptyTiles() int {
	topLeft := [2]int{len(g.elvesMap), len(g.elvesMap[0])}
	bottomRight := [2]int{0, 0}

	for _, elf := range g.elves {
		for i := range topLeft {
			if elf.current[i] < topLeft[i] {
				topLeft[i] = elf.current[i]
			}
		}

		for i := range bottomRight {
			if elf.current[i] > bottomRight[i] {
				bottomRight[i] = elf.current[i]
			}
		}
	}

	h := bottomRight[0] - topLeft[0] + 1
	w := bottomRight[1] - topLeft[1] + 1
	return h*w - len(g.elves)
}

func (g *Grove) checkBorders() {
	for i, row := range g.elvesMap {
		if i == 0 || i == len(g.elvesMap)-1 {
			for _, elf := range row {
				if elf != nil {
					g.expandMap()
					return
				}
			}
		} else {
			if row[0] != nil || row[len(row)-1] != nil {
				g.expandMap()
				return
			}
		}
	}
}

func (g *Grove) expandMap() {
	newH := len(g.elvesMap) + 20
	newW := len(g.elvesMap[0]) + 20
	newMap := make([][]*Elf, newH)
	copy(newMap[10:], g.elvesMap)

	for i, row := range newMap {
		if len(row) == 0 {
			newMap[i] = make([]*Elf, newW)
		} else {
			newRow := make([]*Elf, newW)
			copy(newRow[10:], row)
			newMap[i] = newRow
		}
	}

	g.elvesMap = newMap
	for _, elf := range g.elves {
		elf.current[0] += 10
		elf.current[1] += 10
	}
}

func parseGrove(scanner *bufio.Scanner) *Grove {
	elves := []*Elf{}
	elvesMap := [][]*Elf{}
	i := 0

	for ; scanner.Scan(); i++ {
		row := make([]*Elf, len(scanner.Bytes()))

		for j, char := range scanner.Bytes() {
			if char == '#' {
				elf := &Elf{[2]int{i, j}, nil}
				elves = append(elves, elf)
				row[j] = elf
			} else {
				row[j] = nil
			}
		}

		elvesMap = append(elvesMap, row)
	}

	return &Grove{
		elves:    elves,
		elvesMap: elvesMap,
		movesToConsider: [4][3][2]int{
			{{-1, -1}, {-1, 0}, {-1, 1}},
			{{1, -1}, {1, 0}, {1, 1}},
			{{-1, -1}, {0, -1}, {1, -1}},
			{{-1, 1}, {0, 1}, {1, 1}},
		},
	}
}

func main() {
	f, err := os.Open("input.txt")
	errorHandler(err)
	defer f.Close()

	scanner := bufio.NewScanner(f)
	grove := parseGrove(scanner)

	grove.simulateUntil(10)
	fmt.Println(grove.getEmptyTiles())

	grove.simulateUntilNoMoves()
	fmt.Println(grove.round)
}

func errorHandler(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
