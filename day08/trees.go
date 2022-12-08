package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func checkTreeVisibility(i, j int, height byte, tallest *byte, visible [][]bool) {
	if height > *tallest {
		visible[i][j] = true
		*tallest = height
	}
}

func visibleTrees(trees [][]byte) int {
	w := len(trees[0])
	h := len(trees)
	visible := make([][]bool, h)
	for i := range visible {
		visible[i] = make([]bool, w)
		visible[i][0] = true
		visible[i][w-1] = true
		if i == 0 || i == h-1 {
			for j := 1; j < w-1; j++ {
				visible[i][j] = true
			}
		}
	}

	for j := 1; j < w-1; j++ {
		tallest := trees[0][j]
		for i := 1; i < h-1; i++ {
			checkTreeVisibility(i, j, trees[i][j], &tallest, visible)
		}
	}

	for j := 1; j < w-1; j++ {
		tallest := trees[h-1][j]
		for i := h - 2; i > 0; i-- {
			checkTreeVisibility(i, j, trees[i][j], &tallest, visible)
		}
	}

	for i := 1; i < h-1; i++ {
		tallest := trees[i][0]
		for j := 1; j < w-1; j++ {
			checkTreeVisibility(i, j, trees[i][j], &tallest, visible)
		}
	}

	totalVisble := w*2 + (h-2)*2
	for i := 1; i < h-1; i++ {
		tallest := trees[i][w-1]
		for j := w - 2; j > 0; j-- {
			checkTreeVisibility(i, j, trees[i][j], &tallest, visible)
			if visible[i][j] {
				totalVisble += 1
			}
		}
	}

	return totalVisble
}

func getScenicScore(trees [][]byte, i, j int) int {
	w := len(trees[0])
	h := len(trees)
	treeHeight := trees[i][j]
	left := 0
	right := 0
	up := 0
	down := 0

	for x := i - 1; x >= 0; x-- {
		left += 1
		if trees[x][j] >= treeHeight {
			break
		}
	}

	for x := i + 1; x < w; x++ {
		right += 1
		if trees[x][j] >= treeHeight {
			break
		}
	}

	for y := j - 1; y >= 0; y-- {
		up += 1
		if trees[i][y] >= treeHeight {
			break
		}
	}

	for y := j + 1; y < h; y++ {
		down += 1
		if trees[i][y] >= treeHeight {
			break
		}
	}

	return left * right * up * down
}

func maxScenicScore(trees [][]byte) int {
	max := 0
	for i := 0; i < len(trees); i++ {
		for j := 0; j < len(trees[i]); j++ {
			score := getScenicScore(trees, i, j)
			if score > max {
				max = score
			}
		}
	}

	return max
}

func main() {
	f, err := os.Open("input.txt")
	errorHandler(err)
	defer f.Close()

	var trees [][]byte
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := make([]byte, len(scanner.Bytes()))
		copy(line, scanner.Bytes())
		trees = append(trees, line)
	}

	fmt.Println(visibleTrees(trees))
	fmt.Println(maxScenicScore(trees))
}

func errorHandler(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
