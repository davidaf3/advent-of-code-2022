package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type CPU struct {
	x                 int
	nCycles           int
	signalStrenghtSum int
	screen            [][]byte
}

func newCPU() *CPU {
	screen := make([][]byte, 6)
	for i := range screen {
		for j := 0; j < 40; j++ {
			screen[i] = append(screen[i], '.')
		}
	}

	return &CPU{1, 0, 0, screen}
}

func (c *CPU) cycle() {
	c.nCycles++
	if c.nCycles%20 == 0 && c.nCycles%40 != 0 {
		c.signalStrenghtSum += c.x * c.nCycles
	}

	col := (c.nCycles - 1) % 40
	if c.x == col || c.x == col+1 || c.x == col-1 {
		row := (c.nCycles - 1) / 40
		c.screen[row][col] = '#'
	}
}

func (c *CPU) addx(n int) {
	c.cycle()
	c.cycle()
	c.x += n
}

func (c *CPU) noop() {
	c.cycle()
}

func (c *CPU) run(inst string) {
	splitted := strings.Split(inst, " ")

	switch splitted[0] {
	case "addx":
		n, err := strconv.Atoi(splitted[1])
		errorHandler(err)
		c.addx(n)
	case "noop":
		c.noop()
	}
}

func main() {
	f, err := os.Open("input.txt")
	errorHandler(err)
	defer f.Close()

	cpu := newCPU()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		cpu.run(scanner.Text())
	}

	fmt.Println(cpu.signalStrenghtSum)
	for _, row := range cpu.screen {
		fmt.Println(string(row))
	}
}

func errorHandler(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
