package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"sort"
	"strconv"
)

type Monkey struct {
	items          []int
	operation      func(int) int
	divisor        int
	ifTrue         int
	ifFalse        int
	timesInspected int
}

func (m *Monkey) turn(monkeys []*Monkey) {
	for _, item := range m.items {
		item := m.operation(item)
		item /= 3
		if m.test(item) {
			monkeys[m.ifTrue].catch(item)
		} else {
			monkeys[m.ifFalse].catch(item)
		}
	}

	m.timesInspected += len(m.items)
	m.items = m.items[:0]
}

func (m *Monkey) test(item int) bool {
	return item%m.divisor == 0
}

func (m *Monkey) catch(item int) {
	m.items = append(m.items, item)
}

var numberRegex *regexp.Regexp = regexp.MustCompile("[0-9]+")
var opRegex *regexp.Regexp = regexp.MustCompile("[*+]|[0-9]+|old")

var regexps [5]*regexp.Regexp = [5]*regexp.Regexp{
	numberRegex,
	opRegex,
	numberRegex,
	numberRegex,
	numberRegex,
}

func parseMonkey(scanner *bufio.Scanner) *Monkey {
	var attributes [5][]string
	for i, regex := range regexps {
		scanner.Scan()
		attributes[i] = regex.FindAllString(scanner.Text(), -1)
	}
	scanner.Scan()

	startingItems := make([]int, len(attributes[0]))
	for i, itemStr := range attributes[0] {
		item, err := strconv.Atoi(itemStr)
		errorHandler(err)
		startingItems[i] = item
	}

	divisor, err := strconv.Atoi(attributes[2][0])
	errorHandler(err)
	ifTrue, err := strconv.Atoi(attributes[3][0])
	errorHandler(err)
	ifFalse, err := strconv.Atoi(attributes[4][0])
	errorHandler(err)

	return &Monkey{
		items:          startingItems,
		operation:      makeOperation(attributes[1][1], attributes[1][2]),
		divisor:        divisor,
		ifTrue:         ifTrue,
		ifFalse:        ifFalse,
		timesInspected: 0,
	}
}

func makeOperation(operator, valueStr string) func(int) int {
	if valueStr == "old" {
		return func(old int) int { return old * old }
	}

	value, err := strconv.Atoi(valueStr)
	errorHandler(err)

	switch operator {
	case "+":
		return func(old int) int { return old + value }
	default:
		return func(old int) int { return old * value }
	}
}

func round(monkeys []*Monkey) {
	for _, monkey := range monkeys {
		monkey.turn(monkeys)
	}
}

func getMonkeyBusiness(monkeys []*Monkey) int64 {
	sort.Slice(monkeys, func(i, j int) bool {
		return monkeys[i].timesInspected > monkeys[j].timesInspected
	})

	var monkeyBusiness int64 = 1
	for _, monkey := range monkeys[:2] {
		monkeyBusiness *= int64(monkey.timesInspected)
	}

	return monkeyBusiness
}

func main() {
	f, err := os.Open("input.txt")
	errorHandler(err)
	defer f.Close()

	var monkeys []*Monkey
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		monkeys = append(monkeys, parseMonkey(scanner))
	}

	for i := 0; i < 20; i++ {
		round(monkeys)
	}

	fmt.Println(getMonkeyBusiness(monkeys))
}

func errorHandler(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
