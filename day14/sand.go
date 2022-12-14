package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func signOrZero(n int) int {
	if n == 0 {
		return 0
	}

	if n < 0 {
		return -1
	}

	return 1
}

type Point struct {
	x, y int
}

func parsePoint(text string) Point {
	splitted := strings.Split(text, ",")
	x, err := strconv.Atoi(splitted[0])
	errorHandler(err)
	y, err := strconv.Atoi(splitted[1])
	errorHandler(err)
	return Point{x, y}
}

type Rock struct {
	points []Point
}

func (r *Rock) fill(grid *Grid) {
	prev := r.points[0]

	for _, cur := range r.points[1:] {
		dX, dY := signOrZero(cur.x-prev.x), signOrZero(cur.y-prev.y)
		for i, j := prev.x, prev.y; i != cur.x || j != cur.y; i, j = i+dX, j+dY {
			grid.blocked[j-grid.tl.y][i-grid.tl.x] = true
		}
		prev = cur
	}

	grid.blocked[prev.y-grid.tl.y][prev.x-grid.tl.x] = true
}

func parseRock(text string) (*Rock, Point, Point) {
	var rock Rock
	splitted := strings.Split(text, " -> ")
	rock.points = make([]Point, len(splitted))
	first := parsePoint(splitted[0])
	tL, bR := first, first
	rock.points[0] = first

	for i, pointStr := range splitted[1:] {
		point := parsePoint(pointStr)
		updateCorners(&tL, &bR, point, point)
		rock.points[i+1] = point
	}

	return &rock, tL, bR
}

func updateCorners(tL, bR *Point, newTL, newBR Point) {
	if newTL.x < tL.x {
		tL.x = newTL.x
	}

	if newTL.y < tL.y {
		tL.y = newTL.y
	}

	if newBR.x > bR.x {
		bR.x = newBR.x
	}

	if newBR.y > bR.y {
		bR.y = newBR.y
	}
}

type Grid struct {
	blocked [][]bool
	tl      Point
	w, h    int
}

func (g *Grid) simulateUntilOverflow(i int, stopOnOverflow bool) int {
	var overflows bool
	for ; ; i++ {
		sand := Point{500 - g.tl.x, 0}
		for stop := false; !stop; {
			stop, overflows = g.simulateStep(&sand, stopOnOverflow)
			if overflows {
				return i
			}
		}

		if sand.y == 0 {
			return i + 1
		}
	}
}

func (g *Grid) simulateStep(sand *Point, stopOnOverflow bool) (bool, bool) {
	for _, move := range []Point{{0, 1}, {-1, 1}, {1, 1}} {
		nextX, nextY := sand.x+move.x, sand.y+move.y
		if stopOnOverflow && nextY >= g.h-2 {
			return true, true
		}

		if !g.blocked[nextY][nextX] {
			sand.x, sand.y = nextX, nextY
			return false, false
		}
	}

	g.blocked[sand.y][sand.x] = true
	return true, false
}

func newGrid(rocks []*Rock, rocksTL, rocksBR Point) *Grid {
	h := rocksBR.y + 3
	w := 2*h + 1

	blocked := make([][]bool, h)
	for i := 0; i < h; i++ {
		blocked[i] = make([]bool, w)
	}

	for j := 0; j < w; j++ {
		blocked[h-1][j] = true
	}

	grid := &Grid{blocked, Point{500 - (w-1)/2, 0}, w, h}
	for _, rock := range rocks {
		rock.fill(grid)
	}

	return grid
}

func main() {
	f, err := os.Open("input.txt")
	errorHandler(err)
	defer f.Close()

	var rocks []*Rock
	gridTL, gridBR := Point{500, 0}, Point{500, 0}

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		rock, tL, bR := parseRock(scanner.Text())
		updateCorners(&gridTL, &gridBR, tL, bR)
		rocks = append(rocks, rock)
	}

	grid := newGrid(rocks, gridTL, gridBR)
	firstOverflow := grid.simulateUntilOverflow(0, true)
	startBlocked := grid.simulateUntilOverflow(firstOverflow, false)

	fmt.Println(firstOverflow)
	fmt.Println(startBlocked)
}

func errorHandler(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
