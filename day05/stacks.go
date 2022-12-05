package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Stack[T any] struct {
	items []T
}

func newStack[T any]() *Stack[T] {
	var items []T
	return &Stack[T]{items}
}

func (s *Stack[T]) push(item T) {
	s.items = append(s.items, item)
}

func (s *Stack[T]) pop() T {
	lastIdx := len(s.items) - 1
	last := s.items[lastIdx]
	s.items = s.items[:lastIdx]
	return last
}

func (s *Stack[T]) peek() T {
	return s.items[len(s.items)-1]
}

type Command struct {
	moves int
	from  int
	to    int
}

func newCommandFromText(moves, from, to string) *Command {
	movesN, err := strconv.Atoi(moves)
	errorHandler(err)
	fromN, err := strconv.Atoi(from)
	errorHandler(err)
	toN, err := strconv.Atoi(to)
	errorHandler(err)
	return &Command{movesN, fromN - 1, toN - 1}
}

type CrateMover9000Command struct {
	Command
}

func (c *CrateMover9000Command) run(stacks []*Stack[byte]) {
	for i := 0; i < c.moves; i++ {
		crate := stacks[c.from].pop()
		stacks[c.to].push(crate)
	}
}

type CrateMover9001Command struct {
	Command
}

func (c *CrateMover9001Command) run(stacks []*Stack[byte]) {
	var crates []byte
	for i := 0; i < c.moves; i++ {
		crates = append(crates, stacks[c.from].pop())
	}
	for i := len(crates) - 1; i >= 0; i-- {
		stacks[c.to].push(crates[i])
	}
}

func getStacksFromText(lines []string) []*Stack[byte] {
	lastLine := len(lines) - 1
	stackIndexes := strings.Split(strings.Trim(lines[lastLine], " "), " ")
	nStacks, err := strconv.Atoi(stackIndexes[len(stackIndexes)-1])
	errorHandler(err)

	var stacks []*Stack[byte]
	for i := 0; i < nStacks; i++ {
		stacks = append(stacks, newStack[byte]())
	}

	for i := lastLine - 1; i >= 0; i-- {
		fillStacks(lines[i], stacks)
	}

	return stacks
}

func fillStacks(line string, stacks []*Stack[byte]) {
	for stack := 0; stack*4 < len(line); stack++ {
		linePos := stack * 4
		if line[linePos] == '[' {
			stacks[stack].push(line[linePos+1])
		}
	}
}

func main() {
	f, err := os.Open("input.txt")
	errorHandler(err)
	defer f.Close()

	var stackLines []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan(); scanner.Text() != ""; scanner.Scan() {
		stackLines = append(stackLines, scanner.Text())
	}

	firstStacks := getStacksFromText(stackLines)
	secondStacks := getStacksFromText(stackLines)
	commandRegex, err := regexp.Compile("move ([0-9]+) from ([0-9]+) to ([0-9]+)")
	errorHandler(err)
	for scanner.Scan() {
		args := commandRegex.FindStringSubmatch(scanner.Text())
		command := newCommandFromText(args[1], args[2], args[3])
		(&CrateMover9000Command{*command}).run(firstStacks)
		(&CrateMover9001Command{*command}).run(secondStacks)
	}

	var firstTop []byte
	var secondTop []byte
	for i := range firstStacks {
		firstTop = append(firstTop, firstStacks[i].peek())
		secondTop = append(secondTop, secondStacks[i].peek())
	}

	fmt.Println(string(firstTop))
	fmt.Println(string(secondTop))
}

func errorHandler(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
