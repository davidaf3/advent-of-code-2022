package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

const (
	rock     = 'R'
	paper    = 'P'
	scissors = 'S'
)

const (
	draw = 'D'
	win  = 'W'
	lose = 'L'
)

var scores = map[byte]int{
	rock:     1,
	paper:    2,
	scissors: 3,
}

var encodings = map[string]byte{
	"A": rock,
	"B": paper,
	"C": scissors,
	"X": rock,
	"Y": paper,
	"Z": scissors,
}

var outcomeEncodigns = map[string]byte{
	"X": lose,
	"Y": draw,
	"Z": win,
}

var wins = map[byte]byte{
	rock:     scissors,
	paper:    rock,
	scissors: paper,
}

var loses = map[byte]byte{
	scissors: rock,
	rock:     paper,
	paper:    scissors,
}

type Round struct {
	mine     byte
	opponent byte
}

func newRoundFromChosen(strat []string) *Round {
	return &Round{
		mine:     encodings[strat[1]],
		opponent: encodings[strat[0]],
	}
}

func newRoundFromOutcome(strat []string) *Round {
	opponent := encodings[strat[0]]
	outcome := outcomeEncodigns[strat[1]]

	var mine byte
	switch outcome {
	case win:
		mine = loses[opponent]
	case lose:
		mine = wins[opponent]
	default:
		mine = opponent
	}

	return &Round{
		mine:     mine,
		opponent: opponent,
	}
}

func (r *Round) outcome() int {
	if r.mine == r.opponent {
		return 3
	}

	if wins[r.mine] == r.opponent {
		return 6
	}

	return 0
}

func (r *Round) score() int {
	return r.outcome() + scores[r.mine]
}

func main() {
	f, err := os.Open("input.txt")
	errorHandler(err)
	defer f.Close()

	partOneScore := 0
	partTwoScore := 0
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		strat := strings.Split(scanner.Text(), " ")
		partOneScore += newRoundFromChosen(strat).score()
		partTwoScore += newRoundFromOutcome(strat).score()
	}

	fmt.Println(partOneScore)
	fmt.Println(partTwoScore)
}

func errorHandler(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
