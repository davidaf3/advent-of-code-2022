package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func getDistinctPosition(datastream []byte, nDistinct int) int {
	for i := 0; i < len(datastream); {
		distinct, skip := isDistinct(datastream[i : i+nDistinct])
		if distinct {
			return i + nDistinct
		}
		i += skip
	}

	return -1
}

func isDistinct(data []byte) (bool, int) {
	visited := make([]int, 'z'-'a'+1)

	for i, datum := range data {
		if visited[datum-'a'] != 0 {
			return false, visited[datum-'a']
		}
		visited[datum-'a'] = i + 1
	}

	return true, 0
}

func main() {
	f, err := os.Open("input.txt")
	errorHandler(err)
	defer f.Close()

	scanner := bufio.NewScanner(f)
	scanner.Scan()
	datastream := scanner.Bytes()
	sop := getDistinctPosition(datastream, 4)
	som := getDistinctPosition(datastream[sop-4:], 14) + sop - 4

	fmt.Println(sop)
	fmt.Println(som)
}

func errorHandler(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
