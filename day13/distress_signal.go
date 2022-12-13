package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
)

type Packet interface {
	compareTo(other Packet) int
}

type PacketList struct {
	children []Packet
}

func (l *PacketList) compareTo(other Packet) int {
	switch other.(type) {
	case *PacketList:
		otherChildren := other.(*PacketList).children
		for i, child := range l.children {
			if i >= len(otherChildren) {
				return len(l.children) - len(otherChildren)
			}

			diff := child.compareTo(otherChildren[i])
			if diff != 0 {
				return diff
			}
		}

		return len(l.children) - len(otherChildren)
	case *PacketInt:
		return l.compareTo(&PacketList{[]Packet{other}})
	}

	return 0
}

type PacketInt struct {
	value int
}

func (i *PacketInt) compareTo(other Packet) int {
	switch other.(type) {
	case *PacketList:
		return (&PacketList{[]Packet{i}}).compareTo(other)
	case *PacketInt:
		return i.value - other.(*PacketInt).value
	}

	return 0
}

func parsePacket(text string, pos int) (Packet, int) {
	if text[pos] == '[' {
		return parseList(text, pos)
	}

	return parseInt(text, pos)
}

func parseInt(text string, pos int) (*PacketInt, int) {
	var valueStr []byte

	for ; text[pos] >= '0' && text[pos] <= '9'; pos++ {
		valueStr = append(valueStr, text[pos])
	}

	value, err := strconv.Atoi(string(valueStr))
	errorHandler(err)
	return &PacketInt{value}, pos
}

func parseList(text string, pos int) (*PacketList, int) {
	var children []Packet
	var child Packet

	if text[pos:pos+2] == "[]" {
		return &PacketList{[]Packet{}}, pos + 2
	}

	for text[pos] != ']' {
		pos++
		child, pos = parsePacket(text, pos)
		children = append(children, child)
	}

	pos++
	return &PacketList{children}, pos
}

func getDividerPacket(n int) Packet {
	return &PacketList{[]Packet{&PacketList{[]Packet{&PacketInt{n}}}}}
}

func getDecoderKey(packets []Packet) int {
	firstDivider := getDividerPacket(2)
	secondDivider := getDividerPacket(6)
	var firstIdx int

	for i := 0; i < len(packets); i++ {
		if packets[i].compareTo(firstDivider) == 0 {
			firstIdx = i + 1
			break
		}
	}

	for i := firstIdx + 1; i < len(packets); i++ {
		if packets[i].compareTo(secondDivider) == 0 {
			return firstIdx * (i + 1)
		}
	}

	return -1
}

func main() {
	f, err := os.Open("input.txt")
	errorHandler(err)
	defer f.Close()

	scanner := bufio.NewScanner(f)
	var packets []Packet
	var pair [2]Packet
	pairIdx := 1
	correctPairIdxSum := 0
	i := 0
	for scanner.Scan() {
		if scanner.Text() == "" {
			i = 0
			pairIdx++
			continue
		}

		pair[i], _ = parsePacket(scanner.Text(), 0)
		packets = append(packets, pair[i])
		i++

		if i == 2 && pair[0].compareTo(pair[1]) < 0 {
			correctPairIdxSum += pairIdx
		}
	}

	fmt.Println(correctPairIdxSum)

	packets = append(packets, getDividerPacket(2))
	packets = append(packets, getDividerPacket(6))
	sort.Slice(packets, func(i, j int) bool {
		return packets[i].compareTo(packets[j]) < 0
	})

	fmt.Println(getDecoderKey(packets))
}

func errorHandler(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
