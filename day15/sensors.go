package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
)

func abs(n int) int {
	if n < 0 {
		return -n
	}

	return n
}

type Point struct {
	x, y int
}

type SensorArea struct {
	t, r, b, l Point
}

func (s *SensorArea) getSegmentAt(y int) *Segment {
	if s.t.y > y || s.b.y < y {
		return nil
	}

	var delta int
	if y <= s.l.y {
		delta = y - s.t.y
	} else {
		delta = s.b.y - y
	}

	start := s.t.x - delta
	return newSegment(start, start+2*delta)
}

type Segment struct {
	start, end int
	excluded   []*Segment
}

func (s *Segment) overlapsWithoutExcluding(other *Segment) bool {
	return other.start >= s.start && other.start <= s.end ||
		other.end >= s.start && other.end <= s.end ||
		other.start < s.start && other.end > s.end
}

func (s *Segment) exclude(other *Segment) {
	for _, area := range s.excluded {
		if area.overlapsWithoutExcluding(other) {
			area.exclude(other)
		}
	}

	overlapping := s.getOverlappingSegment(other)
	s.excluded = append(s.excluded, overlapping)
}

func (s *Segment) count() int {
	excludedSum := 0
	for _, excludedSegment := range s.excluded {
		excludedSum += excludedSegment.count()
	}

	return s.end - s.start + 1 - excludedSum
}

func (s *Segment) countInRange(rangeSegment *Segment) int {
	if !s.overlapsWithoutExcluding(rangeSegment) {
		return 0
	}

	excludedSum := 0
	for _, excludedSegment := range s.excluded {
		excludedSum += excludedSegment.countInRange(rangeSegment)
	}

	overlapping := s.getOverlappingSegment(rangeSegment)
	return overlapping.end - overlapping.start + 1 - excludedSum
}

func (s *Segment) getOverlappingSegment(other *Segment) *Segment {
	var start, end int

	if s.start > other.start {
		start = s.start
	} else {
		start = other.start
	}

	if s.end < other.end {
		end = s.end
	} else {
		end = other.end
	}

	return newSegment(start, end)
}

func newSegment(start, end int) *Segment {
	return &Segment{start, end, []*Segment{}}
}

func getEmptyPositionsAt(y int, sensorAreas []*SensorArea, beacons []Point) int {
	var segments []*Segment
	for _, sensorArea := range sensorAreas {
		if segment := sensorArea.getSegmentAt(y); segment != nil {
			for _, otherSegment := range segments {
				if segment.overlapsWithoutExcluding(otherSegment) {
					otherSegment.exclude(segment)
				}
			}

			segments = append(segments, segment)
		}
	}

	emptyPosCount := 0
	for _, segment := range segments {
		emptyPosCount += segment.count()
	}

	for _, beacon := range beacons {
		if beacon.y == y {
			emptyPosCount--
		}
	}

	return emptyPosCount
}

func getBeaconAt(startX, endX, y int, sensorAreas []*SensorArea, beacons []Point) *Point {
	constraints := newSegment(startX, endX)
	var segments []*Segment
	for _, sensorArea := range sensorAreas {
		if segment := sensorArea.getSegmentAt(y); segment != nil {
			if !segment.overlapsWithoutExcluding(constraints) {
				continue
			}

			for _, otherSegment := range segments {
				if segment.overlapsWithoutExcluding(otherSegment) {
					otherSegment.exclude(segment)
				}
			}

			segments = append(segments, segment)
		}
	}

	emptyPosCount := 0
	for _, segment := range segments {
		emptyPosCount += segment.countInRange(constraints)
	}

	if emptyPosCount > endX-startX {
		return nil
	}

	for x := startX; x <= endX; x++ {
		point := newSegment(x, x)
		overlaps := false
		for _, segment := range segments {
			if segment.overlapsWithoutExcluding(point) {
				overlaps = true
				break
			}
		}

		if !overlaps {
			return &Point{x, y}
		}
	}

	return nil
}

var numberRegex *regexp.Regexp = regexp.MustCompile("-?[0-9]+")

func parseSensorAndBeacon(scanner *bufio.Scanner) (*SensorArea, Point) {
	numberStrs := numberRegex.FindAllString(scanner.Text(), -1)
	var numbers [4]int
	for i := 0; i < len(numbers); i++ {
		number, err := strconv.Atoi(numberStrs[i])
		errorHandler(err)
		numbers[i] = number
	}

	sensor := Point{numbers[0], numbers[1]}
	beacon := Point{numbers[2], numbers[3]}

	d := abs(sensor.x-beacon.x) + abs(sensor.y-beacon.y)
	return &SensorArea{
		t: Point{sensor.x, sensor.y - d},
		l: Point{sensor.x - d, sensor.y},
		b: Point{sensor.x, sensor.y + d},
		r: Point{sensor.x + d, sensor.y},
	}, beacon
}

func main() {
	f, err := os.Open("input.txt")
	errorHandler(err)
	defer f.Close()

	var areas []*SensorArea
	var beacons []Point
	visitedBeacons := make(map[string]bool)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		sensor, beacon := parseSensorAndBeacon(scanner)
		areas = append(areas, sensor)
		if _, ok := visitedBeacons[fmt.Sprint(beacon)]; !ok {
			beacons = append(beacons, beacon)
			visitedBeacons[fmt.Sprint(beacon)] = true
		}
	}

	fmt.Println(getEmptyPositionsAt(2000000, areas, beacons))

	pointFound := make(chan Point)
	start, end := 0, 4000000
	step := (end - start) / 40

	for y := start; y < end; y += step {
		go func(startY int) {
			for y := startY; y < startY+step; y++ {
				if point := getBeaconAt(start, end, y, areas, beacons); point != nil {
					pointFound <- *point
				}
			}
		}(y)
	}

	if point := getBeaconAt(start, end, end, areas, beacons); point != nil {
		pointFound <- *point
	}

	point := <-pointFound
	fmt.Println(point.x*4000000 + point.y)
}

func errorHandler(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
