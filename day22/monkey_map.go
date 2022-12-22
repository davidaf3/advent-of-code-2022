package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
)

func isDigit(char byte) bool {
	return char >= '0' && char <= '9'
}

func mod(a, b int) int {
	res := a % b
	if res < 0 {
		return b + res
	}

	return res
}

func getQuadrant(point [2]int, quadrantSize int) [2]int {
	quadrant := [2]int{point[0] / quadrantSize, point[1] / quadrantSize}
	for i := range quadrant {
		if point[i] < 0 && mod(point[i], quadrantSize) != 0 {
			quadrant[i] -= 1
		}
	}

	return quadrant
}

type Row struct {
	start   int
	blocked []bool
}

type Map struct {
	current   [2]int
	direction [2]int
	rows      []*Row
	wrapFn    func(*Map) bool
}

func (m *Map) turn(direction byte) {
	switch direction {
	case 'R':
		if m.direction[0] != 0 {
			m.direction[1] = -m.direction[0]
			m.direction[0] = 0
		} else {
			m.direction[0] = m.direction[1]
			m.direction[1] = 0
		}
	case 'L':
		if m.direction[0] != 0 {
			m.direction[1] = m.direction[0]
			m.direction[0] = 0
		} else {
			m.direction[0] = -m.direction[1]
			m.direction[1] = 0
		}
	}
}

func (m *Map) move(times int) {
	next := [2]int{}
	for i := 0; i < times; i++ {
		next[0] = m.current[0] + m.direction[0]
		next[1] = m.current[1] + m.direction[1]

		if next[0] < 0 || next[0] >= len(m.rows) {
			if m.wrapFn(m) {
				return
			}
			continue
		}

		row := m.rows[next[0]]
		if next[1] < row.start || next[1] >= row.start+len(row.blocked) {
			if m.wrapFn(m) {
				return
			}
			continue
		}

		if row.blocked[next[1]-row.start] {
			break
		}

		m.current = next
	}
}

func (m *Map) wrapAround() bool {
	reverse := [2]int{}
	for i := range m.direction {
		if m.direction[i] != 0 {
			reverse[i] = m.direction[i] * -1
		}
	}

	prev := [2]int{m.current[0], m.current[1]}
	next := [2]int{}
	for {
		next[0] = prev[0] + reverse[0]
		next[1] = prev[1] + reverse[1]

		if next[0] < 0 || next[0] >= len(m.rows) {
			break
		}

		row := m.rows[next[0]]
		if next[1] < row.start || next[1] >= row.start+len(row.blocked) {
			break
		}

		prev = next
	}

	row := m.rows[prev[0]]
	if row.blocked[prev[1]-row.start] {
		return true
	}

	m.current = prev
	return false
}

func (m *Map) wrapCube() bool {
	oldDirection := [2]int{m.direction[0], m.direction[1]}
	prev := [2]int{m.current[0], m.current[1]}
	next := [2]int{}
	inside := false

	for !inside {
		next[0] = prev[0] + m.direction[0]
		next[1] = prev[1] + m.direction[1]

		if next[0] >= 0 && next[0] < len(m.rows) {
			row := m.rows[next[0]]
			if next[1] >= row.start && next[1] < row.start+len(row.blocked) {
				inside = true
			}
		}

		// Probably only works on my input
		if mod(next[0], 50) == mod(next[1], 50) {
			quadrant := getQuadrant(next, 50)
			if (quadrant[0] == -1 && (quadrant[1] == -1 || quadrant[1] == 2)) ||
				(quadrant[0] == 0 && quadrant[1] == -2) ||
				(quadrant[0] == 1 && (quadrant[1] == 0 || quadrant[1] == 2)) ||
				(quadrant[0] == 2 && quadrant[1] == 4) ||
				(quadrant[0] == 3 && quadrant[1] == 1) ||
				(quadrant[0] == 4 && quadrant[1] == 3) {
				m.direction[0], m.direction[1] = -m.direction[1], -m.direction[0]
			}
		}

		// Same as above
		if mod(next[0], 50)+mod(next[1], 50) == 49 {
			quadrant := getQuadrant(next, 50)
			if (quadrant[0] == -1 && (quadrant[1] == 1 || quadrant[1] == 3)) ||
				(quadrant[0] == 0 && quadrant[1] == 4) ||
				(quadrant[0] == 2 && quadrant[1] == -2) ||
				(quadrant[0] == 3 && quadrant[1] == -1) ||
				(quadrant[0] == 4 && quadrant[1] == 0) {
				m.direction[0], m.direction[1] = m.direction[1], m.direction[0]
			}
		}

		prev = next
	}

	row := m.rows[prev[0]]
	if row.blocked[prev[1]-row.start] {
		m.direction = oldDirection
		return true
	}

	m.current = prev
	return false
}

func (m *Map) getFacing() int {
	switch m.direction {
	case [2]int{0, 1}:
		return 0
	case [2]int{1, 0}:
		return 1
	case [2]int{0, -1}:
		return 2
	default:
		return 3
	}
}

func (m *Map) getPassword() int {
	return 1000*(m.current[0]+1) + 4*(m.current[1]+1) + m.getFacing()
}

func newMap(rows []*Row, wrapFn func(*Map) bool) *Map {
	return &Map{[2]int{0, rows[0].start}, [2]int{0, 1}, rows, wrapFn}
}

func parseRow(scanner *bufio.Scanner) *Row {
	blocked := []bool{}
	start := -1

	for i, char := range scanner.Text() {
		if char != ' ' {
			if start == -1 {
				start = i
			}

			blocked = append(blocked, char == '#')
		}
	}

	return &Row{start, blocked}
}

func runInstructions(instructions string, m *Map) {
	for i := 0; i < len(instructions); i++ {
		if isDigit(instructions[i]) {
			numberStr := []byte{instructions[i]}
			j := i + 1

			for ; j < len(instructions) && isDigit(instructions[j]); j++ {
				numberStr = append(numberStr, instructions[j])
			}

			i = j - 1
			number, err := strconv.Atoi(string(numberStr))
			errorHandler(err)
			m.move(number)
		} else {
			m.turn(instructions[i])
		}
	}
}

func main() {
	f, err := os.Open("input.txt")
	errorHandler(err)
	defer f.Close()

	rows := []*Row{}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() && len(scanner.Bytes()) > 0 {
		rows = append(rows, parseRow(scanner))
	}

	scanner.Scan()
	instructions := scanner.Text()

	firstMap := newMap(rows, (*Map).wrapAround)
	runInstructions(instructions, firstMap)
	fmt.Println(firstMap.getPassword())

	secondMap := newMap(rows, (*Map).wrapCube)
	runInstructions(instructions, secondMap)
	fmt.Println(secondMap.getPassword())
}

func errorHandler(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
