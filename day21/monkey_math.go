package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
)

type Expression interface {
	getValue(map[string]Expression) int
	containsHuman(map[string]Expression) bool
	getHumanNumber(int, map[string]Expression) int
}

type BinaryOp struct {
	op                 byte
	left               string
	right              string
	containsHumanCache *bool
	valueCache         *int
}

func (o *BinaryOp) getValue(expressions map[string]Expression) int {
	if o.valueCache == nil {
		leftValue := expressions[o.left].getValue(expressions)
		rightValue := expressions[o.right].getValue(expressions)
		var value int

		switch o.op {
		case '+':
			value = leftValue + rightValue
		case '*':
			value = leftValue * rightValue
		case '-':
			value = leftValue - rightValue
		default:
			value = leftValue / rightValue
		}

		o.valueCache = &value
	}

	return *o.valueCache
}

func (o *BinaryOp) containsHuman(expressions map[string]Expression) bool {
	if o.containsHumanCache == nil {
		contains :=
			expressions[o.left].containsHuman(expressions) ||
				expressions[o.right].containsHuman(expressions)
		o.containsHumanCache = &contains
	}

	return *o.containsHumanCache
}

func (o *BinaryOp) getHumanNumber(expected int, expressions map[string]Expression) int {
	left, right := expressions[o.left], expressions[o.right]
	humanBranch, monkeyBranch := getHumanMonkeyBranch(left, right, expressions)
	other := monkeyBranch.getValue(expressions)

	switch o.op {
	case '+':
		return humanBranch.getHumanNumber(expected-other, expressions)
	case '*':
		return humanBranch.getHumanNumber(expected/other, expressions)
	case '-':
		if monkeyBranch == left {
			expected = -expected
		}

		return humanBranch.getHumanNumber(expected+other, expressions)
	default:
		if monkeyBranch == left {
			return humanBranch.getHumanNumber(other/expected, expressions)
		}

		return humanBranch.getHumanNumber(expected*other, expressions)
	}
}

type Number struct {
	value   int
	isHuman bool
}

func (n *Number) getValue(expressions map[string]Expression) int {
	return n.value
}

func (n *Number) containsHuman(expressions map[string]Expression) bool {
	return n.isHuman
}

func (n *Number) getHumanNumber(expected int, expressions map[string]Expression) int {
	return expected
}

func parseExpression(scanner *bufio.Scanner) (Expression, string) {
	line := scanner.Text()
	name := line[:4]

	if line[6] >= '0' && line[6] <= '9' {
		number, err := strconv.Atoi(line[6:])
		errorHandler(err)
		return &Number{number, name == "humn"}, name
	}

	return &BinaryOp{
		op:    line[11],
		left:  line[6:10],
		right: line[13:],
	}, name
}

func getHumanMonkeyBranch(left, right Expression, expressions map[string]Expression) (Expression, Expression) {
	if left.containsHuman(expressions) {
		return left, right
	}

	return right, left
}

func getHumanNumber(root *BinaryOp, expressions map[string]Expression) int {
	rootLeft, rootRight := expressions[root.left], expressions[root.right]
	humanBranch, monkeyBranch := getHumanMonkeyBranch(rootLeft, rootRight, expressions)
	expected := monkeyBranch.getValue(expressions)
	return humanBranch.getHumanNumber(expected, expressions)
}

func main() {
	f, err := os.Open("input.txt")
	errorHandler(err)
	defer f.Close()

	var root *BinaryOp
	expressions := map[string]Expression{}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		expression, name := parseExpression(scanner)
		expressions[name] = expression
		if name == "root" {
			root = expression.(*BinaryOp)
		}
	}

	fmt.Println(root.getValue(expressions))
	fmt.Println(getHumanNumber(root, expressions))
}

func errorHandler(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
