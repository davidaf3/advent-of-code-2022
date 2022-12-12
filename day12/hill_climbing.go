package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"log"
	"os"
)

type MinHeap struct {
	array    []*Node
	start    *Node
	distance []int
}

func (h MinHeap) Len() int {
	return len(h.array)
}

func (h MinHeap) Less(i, j int) bool {
	iDist := h.distance[h.array[i].index]
	jDist := h.distance[h.array[j].index]

	if iDist == -1 || jDist == -1 {
		return iDist > jDist
	}

	return iDist < jDist
}

func (h MinHeap) Swap(i, j int) {
	h.array[i], h.array[j] = h.array[j], h.array[i]
}

func (h *MinHeap) Push(x any) {
	h.array = append(h.array, x.(*Node))
}

func (h *MinHeap) Pop() any {
	n := len(h.array) - 1
	x := h.array[n]
	h.array = h.array[:n]
	return x
}

type Node struct {
	i      int
	j      int
	index  int
	height byte
}

func (n *Node) getNeighbours(heightMap [][]byte, testHeight func(byte, byte) bool) []*Node {
	var neighbours []*Node

	for _, delta := range []int{-1, 1} {
		if n.i+delta >= 0 && n.i+delta < len(heightMap) &&
			testHeight(n.height, heightMap[n.i+delta][n.j]) {
			neighbours = append(neighbours, newNode(heightMap, n.i+delta, n.j))
		}

		if n.j+delta >= 0 && n.j+delta < len(heightMap[n.i]) &&
			testHeight(n.height, heightMap[n.i][n.j+delta]) {
			neighbours = append(neighbours, newNode(heightMap, n.i, n.j+delta))
		}
	}

	return neighbours
}

func (n *Node) equals(other *Node) bool {
	return n.i == other.i && n.j == other.j
}

func newNode(heightMap [][]byte, i, j int) *Node {
	return &Node{
		i:      i,
		j:      j,
		index:  j + i*len(heightMap[0]),
		height: heightMap[i][j],
	}
}

func shortestPath(heightMap [][]byte, nodes []*Node, start *Node,
	testEnd func(*Node) bool, testHeight func(byte, byte) bool) ([]int, *Node) {
	visited := make([]bool, len(nodes))
	distance := make([]int, len(nodes))

	for i := range distance {
		distance[i] = -1
	}
	distance[start.index] = 0

	frontier := &MinHeap{
		start:    start,
		distance: distance,
	}
	heap.Init(frontier)
	heap.Push(frontier, start)

	for frontier.Len() != 0 {
		closest := heap.Pop(frontier).(*Node)

		if visited[closest.index] {
			continue
		}

		if testEnd(closest) {
			return distance, closest
		}

		if distance[closest.index] == -1 {
			return nil, nil
		}

		closestDist := distance[closest.index]
		for _, neighbour := range closest.getNeighbours(heightMap, testHeight) {
			neighbourDist := distance[neighbour.index]
			if !visited[neighbour.index] &&
				(neighbourDist == -1 || closestDist+1 < neighbourDist) {
				distance[neighbour.index] = closestDist + 1
				frontier.Push(neighbour)
			}
		}

		visited[closest.index] = true
	}

	return nil, nil
}

func parseHeightMap(scanner *bufio.Scanner) ([][]byte, []*Node, *Node, *Node) {
	var heightMap [][]byte
	var nodes []*Node
	var start *Node
	var end *Node

	for scanner.Scan() {
		row := make([]byte, len(scanner.Bytes()))
		copy(row, scanner.Bytes())
		heightMap = append(heightMap, row)

		i := len(heightMap) - 1
		for j, height := range row {
			var node *Node

			switch height {
			case 'S':
				heightMap[i][j] = 'a'
				node = newNode(heightMap, i, j)
				start = node
			case 'E':
				heightMap[i][j] = 'z'
				node = newNode(heightMap, i, j)
				end = node
			default:
				node = newNode(heightMap, i, j)
			}

			nodes = append(nodes, node)
		}
	}

	return heightMap, nodes, start, end
}

func main() {
	f, err := os.Open("input.txt")
	errorHandler(err)
	defer f.Close()

	scanner := bufio.NewScanner(f)
	heightMap, nodes, start, end := parseHeightMap(scanner)
	testEnd := func(n *Node) bool { return n.equals(end) }
	testHeight := func(hSrc byte, hDst byte) bool { return hDst <= hSrc+1 }
	distances, _ := shortestPath(heightMap, nodes, start, testEnd, testHeight)

	fmt.Println(distances[end.index])

	testEnd = func(n *Node) bool { return heightMap[n.i][n.j] == 'a' }
	testHeight = func(hSrc byte, hDst byte) bool { return hDst+1 >= hSrc }
	distances, scenicStart := shortestPath(heightMap, nodes, end, testEnd, testHeight)

	fmt.Println(distances[scenicStart.index])
}

func errorHandler(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
