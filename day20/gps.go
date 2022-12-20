package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
)

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

type Node[T any] struct {
	value T
	prev  *Node[T]
	next  *Node[T]
}

type CircularList[T any] struct {
	start  *Node[T]
	length int
}

func (l *CircularList[T]) append(node *Node[T]) {
	l.length++
	if l.start == nil {
		l.start = node
		node.next = node
		node.prev = node
		return
	}

	l.start.prev.next = node
	node.prev = l.start.prev
	l.start.prev = node
	node.next = l.start
}

func (l *CircularList[T]) get(i int) *Node[T] {
	current := l.start
	counter := 0

	for counter < i%l.length {
		counter++
		current = current.next
	}

	return current
}

func (l *CircularList[T]) find(pred func(T) bool) int {
	current := l.start
	i := 0

	for !pred(current.value) {
		current = current.next
		i++
	}

	return i
}

func (l *CircularList[T]) moveBack(node *Node[T]) {
	prev := node.prev
	next := node.next
	prev.prev.next = node
	node.prev = prev.prev
	node.next = prev
	prev.prev = node
	prev.next = next
	next.prev = prev
}

func (l *CircularList[T]) moveForward(node *Node[T]) {
	prev := node.prev
	next := node.next
	next.next.prev = node
	node.next = next.next
	node.prev = next
	next.next = node
	next.prev = prev
	prev.next = next
}

func newCircularList[T any](values []*Node[T]) *CircularList[T] {
	list := &CircularList[T]{}

	for _, value := range values {
		list.append(value)
	}

	return list
}

func decrypt(list *CircularList[int], nodes []*Node[int], key int) {
	for _, node := range nodes {
		var moveFn func(*Node[int])
		if node.value > 0 {
			moveFn = list.moveForward
		} else {
			moveFn = list.moveBack
		}

		for i := 0; i < abs(node.value*key)%(list.length-1); i++ {
			moveFn(node)
		}
	}
}

func getCoordinatesSum(list *CircularList[int], key int) int {
	zero := list.find(func(value int) bool { return value == 0 })
	sum := 0

	for i := 1000; i <= 3000; i += 1000 {
		sum += list.get(zero+i).value * key
	}

	return sum
}

func main() {
	f, err := os.Open("input.txt")
	errorHandler(err)
	defer f.Close()

	firstNumbers := []*Node[int]{}
	secondNumbers := []*Node[int]{}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		number, err := strconv.Atoi(scanner.Text())
		errorHandler(err)
		firstNumbers = append(firstNumbers, &Node[int]{value: number})
		secondNumbers = append(secondNumbers, &Node[int]{value: number})
	}

	list := newCircularList(firstNumbers)
	decrypt(list, firstNumbers, 1)
	fmt.Println(getCoordinatesSum(list, 1))

	list = newCircularList(secondNumbers)
	for i := 0; i < 10; i++ {
		decrypt(list, secondNumbers, 811589153)
	}
	fmt.Println(getCoordinatesSum(list, 811589153))
}

func errorHandler(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
