package main

import (
	"fmt"
	"strconv"
)

func day1() {
	curr := 50
	zeroCount := 0

	for _, line := range readDataLines("day1.txt") {
		if len(line) < 2 {
			panic("bad line:" + line)
		}
		delta := must(strconv.Atoi(line[1:]))
		var dir int
		switch line[0] {
		case 'L':
			dir = -1
		case 'R':
			dir = 1
		default:
			panic("bad direction:" + line)
		}
		for range delta {
			curr = (curr + dir) % 100
			if curr == 0 {
				zeroCount++
			}
		}
	}

	fmt.Printf("Zero count: %d\n", zeroCount)
}

func init() {
	functionRegistry[1] = day1
}
