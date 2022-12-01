package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"log"
	"os"
	"strconv"
)

type MaxHeap []int

func (h MaxHeap) Len() int {
	return len(h)
}

func (h MaxHeap) Less(i, j int) bool {
	return h[i] > h[j]
}

func (h MaxHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h *MaxHeap) Push(x any) {
	*h = append(*h, x.(int))
}

func (h *MaxHeap) Pop() any {
	n := len(*h)
	x := (*h)[n-1]
	*h = (*h)[:n-1]
	return x
}

func main() {
	f, err := os.Open("input.txt")
	errorHandler(err)
	defer f.Close()

	calories := &MaxHeap{}
	heap.Init(calories)
	currentSum := 0

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			heap.Push(calories, currentSum)
			currentSum = 0
			continue
		}

		number, err := strconv.Atoi(line)
		errorHandler(err)
		currentSum += number
	}

	first := heap.Pop(calories).(int)
	topThree := first
	for i := 0; i < 2; i += 1 {
		topThree += heap.Pop(calories).(int)
	}

	fmt.Println(first)
	fmt.Println(topThree)
}

func errorHandler(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
