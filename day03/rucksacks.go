package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

type Rucksack struct {
	firstCompartment  []byte
	secondCompartment []byte
	typesSet          []bool
	errorType         byte
}

func newRucksack(contents []byte) *Rucksack {
	compartmentSize := len(contents) / 2
	firstCompartment := contents[:compartmentSize]
	secondCompartment := contents[compartmentSize:]

	typesSet := make([]bool, 'z'-'A'+1)
	var errorType byte

	for _, item := range firstCompartment {
		typesSet[item-'A'] = true
	}

	for _, item := range secondCompartment {
		if typesSet[item-'A'] {
			errorType = item
			continue
		}

		typesSet[item-'A'] = true
	}

	return &Rucksack{
		firstCompartment:  firstCompartment,
		secondCompartment: secondCompartment,
		typesSet:          typesSet,
		errorType:         errorType,
	}
}

func getPriority(itemType byte) int {
	if itemType >= 'a' {
		return int(itemType - 'a' + 1)
	}

	return int(itemType - 'A' + 27)
}

func getGroupType(group []*Rucksack) byte {
	for itemType := byte('A'); itemType <= 'Z'; itemType += 1 {
		if isGroupType(itemType, group) {
			return itemType
		}
	}

	for itemType := byte('a'); itemType <= 'z'; itemType += 1 {
		if isGroupType(itemType, group) {
			return itemType
		}
	}

	return '0'
}

func isGroupType(itemType byte, group []*Rucksack) bool {
	for _, rucksack := range group {
		if !rucksack.typesSet[itemType-'A'] {
			return false
		}
	}

	return true
}

func main() {
	f, err := os.Open("input.txt")
	errorHandler(err)
	defer f.Close()

	group := make([]*Rucksack, 3)
	groupCount := 0
	errorPriorities := 0
	groupPriorities := 0

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		rucksack := newRucksack(scanner.Bytes())
		errorPriorities += getPriority(rucksack.errorType)
		group[groupCount] = rucksack
		groupCount += 1
		if groupCount == 3 {
			groupType := getGroupType(group)
			groupPriorities += getPriority(groupType)
			groupCount = 0
		}
	}

	fmt.Println(errorPriorities)
	fmt.Println(groupPriorities)
}

func errorHandler(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
