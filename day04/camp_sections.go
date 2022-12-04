package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type Assignment struct {
	start int
	end   int
}

func (a *Assignment) contains(other *Assignment) bool {
	return a.start <= other.start && a.end >= other.end
}

func (a *Assignment) overlaps(other *Assignment) bool {
	return a.start <= other.start && a.end >= other.start ||
		a.start <= other.end && a.end >= other.end ||
		a.contains(other) ||
		other.contains(a)
}

func newAssignmentFromText(text string) *Assignment {
	splitted := strings.Split(text, "-")
	start, err := strconv.Atoi(splitted[0])
	errorHandler(err)
	end, err := strconv.Atoi(splitted[1])
	errorHandler(err)
	return &Assignment{start, end}
}

type Pair struct {
	first  *Assignment
	second *Assignment
}

func (p *Pair) oneContainsOther() bool {
	return p.first.contains(p.second) || p.second.contains(p.first)
}

func (p *Pair) assingmentsOverlap() bool {
	return p.first.overlaps(p.second)
}

func newPairFromText(text string) *Pair {
	splitted := strings.Split(text, ",")
	return &Pair{
		first:  newAssignmentFromText(splitted[0]),
		second: newAssignmentFromText(splitted[1]),
	}
}

func main() {
	f, err := os.Open("input.txt")
	errorHandler(err)
	defer f.Close()

	contained := 0
	overlapping := 0
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		pair := newPairFromText(scanner.Text())
		if pair.oneContainsOther() {
			contained += 1
			overlapping += 1
		} else if pair.assingmentsOverlap() {
			overlapping += 1
		}
	}

	fmt.Println(contained)
	fmt.Println(overlapping)
}

func errorHandler(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
