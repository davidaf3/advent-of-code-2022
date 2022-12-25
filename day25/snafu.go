package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
)

func max(a, b int) int {
	if a > b {
		return a
	}

	return b
}

type SNAFU []int8

func add(a, b SNAFU) SNAFU {
	maxLen := max(len(a), len(b))
	res := SNAFU{}
	var carry int8 = 0

	for i := 0; i < maxLen; i++ {
		if i >= len(a) {
			a = append(a, 0)
		}

		if i >= len(b) {
			b = append(b, 0)
		}

		res = append(res, a[i]+b[i]+carry)

		if res[i] > 2 {
			res[i] -= 5
			carry = 1
		} else if res[i] < -2 {
			res[i] += 5
			carry = -1
		} else {
			carry = 0
		}
	}

	if carry != 0 {
		res = append(res, carry)
	}

	return res
}

func parseNumber(scanner *bufio.Scanner) SNAFU {
	number := SNAFU{}
	for i := len(scanner.Bytes()) - 1; i >= 0; i-- {
		switch scanner.Bytes()[i] {
		case '2':
			number = append(number, 2)
		case '1':
			number = append(number, 1)
		case '0':
			number = append(number, 0)
		case '-':
			number = append(number, -1)
		case '=':
			number = append(number, -2)
		}
	}

	return number
}

func Int(number SNAFU) int {
	res := 0
	for i, digit := range number {
		res += int(digit) * int(math.Pow(5, float64(i)))
	}

	return res
}

func String(snafu SNAFU) string {
	res := ""
	for _, digit := range snafu {
		switch digit {
		case 2:
			res = "2" + res
		case 1:
			res = "1" + res
		case 0:
			res = "0" + res
		case -1:
			res = "-" + res
		case -2:
			res = "=" + res
		}
	}

	leadingZeroes := 0
	for _, char := range res {
		if char != '0' {
			break
		}

		leadingZeroes++
	}

	return res[leadingZeroes:]
}

func main() {
	f, err := os.Open("input.txt")
	errorHandler(err)
	defer f.Close()

	numbers := []SNAFU{}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		numbers = append(numbers, parseNumber(scanner))
	}

	sum := SNAFU{}
	for _, number := range numbers {
		sum = add(number, sum)
	}

	fmt.Println(String(sum))
}

func errorHandler(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
