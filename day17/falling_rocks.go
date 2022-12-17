package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

type Point struct {
	x, y int
}

type Rock interface {
	getTL() Point
	moveLeft(*Chamber)
	moveRight(*Chamber)
	moveDown()
	canMoveDown(*Chamber) bool
	updateResting(*Chamber) bool
}

type BaseRock struct {
	tL   Point
	w, h int
}

func (r *BaseRock) getTL() Point {
	return r.tL
}

func (r *BaseRock) moveLeft(c *Chamber) {
	if r.tL.x > 0 {
		for y := r.tL.y; y > r.tL.y-r.h; y-- {
			if c.cells[y-c.bottomH][r.tL.x-1] {
				return
			}
		}

		r.tL.x -= 1
	}
}

func (r *BaseRock) moveRight(c *Chamber) {
	if r.tL.x+r.w < c.w {
		for y := r.tL.y; y > r.tL.y-r.h; y-- {
			if c.cells[y-c.bottomH][r.tL.x+r.w] {
				return
			}
		}

		r.tL.x += 1
	}
}

func (r *BaseRock) moveDown() {
	r.tL.y -= 1
}

func (r *BaseRock) canMoveDown(c *Chamber) bool {
	rBottomY := r.tL.y - r.h + 1
	if rBottomY-c.bottomH == 0 {
		return false
	}

	for x := r.tL.x; x < r.tL.x+r.w; x++ {
		if c.cells[rBottomY-c.bottomH-1][x] {
			return false
		}
	}

	return true
}

func (r *BaseRock) setHighest(c *Chamber) {
	if r.tL.y+1 > c.highest {
		c.highest = r.tL.y + 1
	}
}

func (r *BaseRock) checkFullRows(c *Chamber) bool {
	for y := r.tL.y; y > r.tL.y-r.h; y-- {
		rowFull := true
		for x := 0; x < c.w && rowFull; x++ {
			rowFull = rowFull && c.cells[y-c.bottomH][x]
		}

		if rowFull {
			c.cells = c.cells[y-c.bottomH+1:]
			c.bottomH = y + 1
			return y == r.tL.y
		}
	}

	return false
}

func (r *BaseRock) updateResting(c *Chamber) bool {
	r.setHighest(c)

	for y := r.tL.y; y > r.tL.y-r.h; y-- {
		for x := r.tL.x; x < r.tL.x+r.w; x++ {
			c.cells[y-c.bottomH][x] = true
		}
	}

	return r.checkFullRows(c)
}

type LRock struct {
	BaseRock
}

func (r *LRock) moveLeft(c *Chamber) {
	if r.tL.x > 0 && !c.cells[r.tL.y-r.h+1-c.bottomH][r.tL.x-1] {
		for y := r.tL.y; y > r.tL.y-r.h+1; y-- {
			if c.cells[y-c.bottomH][r.tL.x+r.w-2] {
				return
			}
		}

		r.tL.x -= 1
	}
}

func (r *LRock) updateResting(c *Chamber) bool {
	r.setHighest(c)

	for y := r.tL.y; y > r.tL.y-r.h; y-- {
		c.cells[y-c.bottomH][r.tL.x+r.w-1] = true
	}

	for x := r.tL.x; x < r.tL.x+r.w-1; x++ {
		c.cells[r.tL.y-r.h+1-c.bottomH][x] = true
	}

	return r.checkFullRows(c)
}

type PlusRock struct {
	BaseRock
}

func (r *PlusRock) moveLeft(c *Chamber) {
	if r.tL.x > 0 && !c.cells[r.tL.y-1-c.bottomH][r.tL.x-1] &&
		!c.cells[r.tL.y-c.bottomH][r.tL.x] &&
		!c.cells[r.tL.y-r.h+1-c.bottomH][r.tL.x] {
		r.tL.x -= 1
	}
}

func (r *PlusRock) moveRight(c *Chamber) {
	if r.tL.x+r.w < c.w && !c.cells[r.tL.y-1-c.bottomH][r.tL.x+r.w] &&
		!c.cells[r.tL.y-c.bottomH][r.tL.x+r.w-1] &&
		!c.cells[r.tL.y-r.h+1-c.bottomH][r.tL.x+r.w-1] {
		r.tL.x += 1
	}
}

func (r *PlusRock) canMoveDown(c *Chamber) bool {
	rBottomY := r.tL.y - r.h + 1
	return rBottomY-c.bottomH != 0 &&
		!c.cells[rBottomY-c.bottomH-1][r.tL.x+1] &&
		!c.cells[rBottomY-c.bottomH][r.tL.x] &&
		!c.cells[rBottomY-c.bottomH][r.tL.x+r.w-1]
}

func (r *PlusRock) updateResting(c *Chamber) bool {
	r.setHighest(c)

	for y := r.tL.y; y > r.tL.y-r.h; y-- {
		c.cells[y-c.bottomH][r.tL.x+1] = true
	}

	for x := r.tL.x; x < r.tL.x+r.w; x++ {
		c.cells[r.tL.y-1-c.bottomH][x] = true
	}

	return r.checkFullRows(c)
}

type SimulationResult struct {
	curMove         int
	curRock         int
	cells           [][]bool
	newRocks        int
	heightAdded     int
	bottomAdded     int
	nextStateString string
}

func copyCells(cells [][]bool) [][]bool {
	cellsCopy := make([][]bool, len(cells))
	for i, row := range cells {
		dst := make([]bool, len(row))
		copy(dst, row)
		cellsCopy[i] = dst
	}
	return cellsCopy
}

type Chamber struct {
	curMove      int
	moves        []byte
	curRock      int
	rockCreators [](func(int) Rock)
	cells        [][]bool
	restingRocks int
	w, bottomH   int
	highest      int
	cache        map[string]*SimulationResult
}

func (c *Chamber) getStateString() string {
	stateString := fmt.Sprintf("%d-%d-", c.curMove, c.curRock)
	for _, row := range c.cells {
		for _, cell := range row {
			if cell {
				stateString += "#"
			} else {
				stateString += "."
			}
		}
	}
	return stateString
}

func (c *Chamber) nextSideMove(rock Rock) {
	switch c.moves[c.curMove] {
	case '<':
		rock.moveLeft(c)
	case '>':
		rock.moveRight(c)
	}

	c.curMove++
	if c.curMove == len(c.moves) {
		c.curMove = 0
	}
}

func (c *Chamber) simulateRock(rock Rock) bool {
	for {
		c.nextSideMove(rock)
		if rock.canMoveDown(c) {
			rock.moveDown()
		} else {
			return rock.updateResting(c)
		}
	}
}

func (c *Chamber) simulateUntilOrRowFull(maxResting int) {
	initResting := c.restingRocks
	initHighest := c.highest
	initBottomH := c.bottomH
	stateString := c.getStateString()

	if state, ok := c.cache[stateString]; ok &&
		c.restingRocks+state.newRocks <= maxResting {
		c.curMove = state.curMove
		c.curRock = state.curRock
		c.cells = copyCells(state.cells)
		c.highest += state.heightAdded
		c.bottomH += state.bottomAdded
		c.restingRocks += state.newRocks
		return
	}

	for {
		rock := c.rockCreators[c.curRock](c.highest + 3)
		c.curRock++
		if c.curRock == len(c.rockCreators) {
			c.curRock = 0
		}

		for i := len(c.cells) + c.bottomH; i <= rock.getTL().y; i++ {
			c.cells = append(c.cells, make([]bool, c.w))
		}

		rowFull := c.simulateRock(rock)
		c.restingRocks++

		if rowFull {
			result := &SimulationResult{
				curMove:         c.curMove,
				curRock:         c.curRock,
				cells:           copyCells(c.cells),
				newRocks:        c.restingRocks - initResting,
				heightAdded:     c.highest - initHighest,
				bottomAdded:     c.bottomH - initBottomH,
				nextStateString: c.getStateString(),
			}

			for _, prevState := range c.cache {
				if prevState.nextStateString == stateString {
					prevState.curMove = result.curMove
					prevState.curRock = result.curRock
					prevState.cells = result.cells
					prevState.newRocks += result.newRocks
					prevState.heightAdded += result.heightAdded
					prevState.bottomAdded += result.bottomAdded
					prevState.nextStateString = result.nextStateString
				}
			}

			c.cache[stateString] = result
		}

		if rowFull || c.restingRocks == maxResting {
			return
		}
	}
}

func (c *Chamber) simulateUntil(maxResting int) int {
	for {
		c.simulateUntilOrRowFull(maxResting)

		for stateString, state := range c.cache {
			if stateString == state.nextStateString {
				loopingState := state
				initState := c.cache["0-0-"]
				loops := (maxResting - initState.newRocks) / loopingState.newRocks

				c.curMove = loopingState.curMove
				c.curRock = loopingState.curRock
				c.cells = copyCells(loopingState.cells)
				c.highest = initState.heightAdded + loopingState.heightAdded*loops
				c.bottomH = initState.bottomAdded + loopingState.bottomAdded*loops
				c.restingRocks = initState.newRocks + loopingState.newRocks*loops

				c.simulateUntilOrRowFull(maxResting)
				return c.highest
			}
		}

		if c.restingRocks == maxResting {
			return c.highest
		}
	}
}

func newChamber(moves []byte, w int) *Chamber {
	rockCreators := [](func(int) Rock){
		func(y int) Rock {
			return &BaseRock{Point{2, y}, 4, 1}
		},
		func(y int) Rock {
			return &PlusRock{BaseRock{Point{2, y + 2}, 3, 3}}
		},
		func(y int) Rock {
			return &LRock{BaseRock{Point{2, y + 2}, 3, 3}}
		},
		func(y int) Rock {
			return &BaseRock{Point{2, y + 3}, 1, 4}
		},
		func(y int) Rock {
			return &BaseRock{Point{2, y + 1}, 2, 2}
		},
	}

	return &Chamber{
		curMove:      0,
		moves:        moves,
		curRock:      0,
		rockCreators: rockCreators,
		cells:        [][]bool{},
		restingRocks: 0,
		w:            w,
		bottomH:      0,
		highest:      0,
		cache:        make(map[string]*SimulationResult),
	}
}

func main() {
	f, err := os.Open("input.txt")
	errorHandler(err)
	defer f.Close()

	scanner := bufio.NewScanner(f)
	scanner.Scan()

	c := newChamber(scanner.Bytes(), 7)
	fmt.Println(c.simulateUntil(2022))

	c = newChamber(scanner.Bytes(), 7)
	fmt.Println(c.simulateUntil(1000000000000))
}

func errorHandler(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
