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

type WorryLevel struct {
	modulos map[int]int
}

func (w *WorryLevel) Add(n int) {
	for divisor, modulo := range w.modulos {
		w.modulos[divisor] = (modulo + n) % divisor
	}
}

func (w *WorryLevel) Mul(n int) {
	for divisor, modulo := range w.modulos {
		w.modulos[divisor] = (modulo * n) % divisor
	}
}

func (w *WorryLevel) Pow2() {
	for divisor, modulo := range w.modulos {
		w.modulos[divisor] = (modulo * modulo) % divisor
	}
}

type Monkey struct {
	startingItems  []int
	items          []*WorryLevel
	operation      func(*WorryLevel)
	divisor        int
	ifTrue         int
	ifFalse        int
	timesInspected int
}

func (m *Monkey) turn(monkeys []*Monkey) {
	for _, item := range m.items {
		m.operation(item)
		if m.test(item) {
			monkeys[m.ifTrue].catch(item)
		} else {
			monkeys[m.ifFalse].catch(item)
		}
	}

	m.timesInspected += len(m.items)
	m.items = m.items[:0]
}

func (m *Monkey) test(level *WorryLevel) bool {
	return level.modulos[m.divisor] == 0
}

func (m *Monkey) catch(item *WorryLevel) {
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
		startingItems:  startingItems,
		items:          []*WorryLevel{},
		operation:      makeOperation(attributes[1][1], attributes[1][2]),
		divisor:        divisor,
		ifTrue:         ifTrue,
		ifFalse:        ifFalse,
		timesInspected: 0,
	}
}

func setItems(monkeys []*Monkey) {
	divisors := make([]int, len(monkeys))
	for i, monkey := range monkeys {
		divisors[i] = monkey.divisor
	}

	for _, monkey := range monkeys {
		for _, item := range monkey.startingItems {
			modulos := make(map[int]int, len(divisors))
			for _, divisor := range divisors {
				modulos[divisor] = item % divisor
			}

			monkey.items = append(monkey.items, &WorryLevel{modulos})
		}
	}
}

func makeOperation(operator, valueStr string) func(*WorryLevel) {
	if valueStr == "old" {
		return func(old *WorryLevel) { old.Pow2() }
	}

	value, err := strconv.Atoi(valueStr)
	errorHandler(err)

	switch operator {
	case "+":
		return func(old *WorryLevel) { old.Add(value) }
	default:
		return func(old *WorryLevel) { old.Mul(value) }
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

	setItems(monkeys)
	for i := 0; i < 10000; i++ {
		round(monkeys)
	}

	fmt.Println(getMonkeyBusiness(monkeys))
}

func errorHandler(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
