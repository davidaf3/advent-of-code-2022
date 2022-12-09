package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type Point struct {
	x int
	y int
}

type Rope struct {
	knots   []Point
	visited map[string]int
}

func newRope(knots int) *Rope {
	return &Rope{
		knots:   make([]Point, knots),
		visited: map[string]int{"0,0": 1},
	}
}

func (r *Rope) updateTail() {
	for i := 1; i < len(r.knots); i++ {
		prev := r.knots[i-1]
		current := r.knots[i]

		if prev.x > current.x+1 ||
			prev.x < current.x-1 ||
			prev.y > current.y+1 ||
			prev.y < current.y-1 {

			if current.x < prev.x {
				r.knots[i].x += 1
			} else if current.x > prev.x {
				r.knots[i].x -= 1
			}

			if current.y < prev.y {
				r.knots[i].y += 1
			} else if current.y > prev.y {
				r.knots[i].y -= 1
			}
		}
	}
}

var ropeMoves = map[byte](func(*Rope)){
	'R': func(r *Rope) { r.knots[0].x += 1 },
	'L': func(r *Rope) { r.knots[0].x -= 1 },
	'U': func(r *Rope) { r.knots[0].y += 1 },
	'D': func(r *Rope) { r.knots[0].y -= 1 },
}

func (r *Rope) move(direction byte, times int) {
	moveFn := ropeMoves[direction]
	for i := 0; i < times; i++ {
		moveFn(r)
		r.updateTail()

		last := r.knots[len(r.knots)-1]
		tailCoords := fmt.Sprintf("%d,%d", last.x, last.y)
		if timesVisited, ok := r.visited[tailCoords]; !ok {
			r.visited[tailCoords] = 1
		} else {
			r.visited[tailCoords] = timesVisited + 1
		}
	}
}

func main() {
	f, err := os.Open("input.txt")
	errorHandler(err)
	defer f.Close()

	shortRope := newRope(2)
	longRope := newRope(10)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		splitted := strings.Split(scanner.Text(), " ")
		times, err := strconv.Atoi(splitted[1])
		errorHandler(err)
		shortRope.move(splitted[0][0], times)
		longRope.move(splitted[0][0], times)
	}

	fmt.Println(len(shortRope.visited))
	fmt.Println(len(longRope.visited))
}

func errorHandler(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
