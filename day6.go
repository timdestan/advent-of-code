package main

import (
	"fmt"
	"strings"
)

func day6() {
	var grid [][]rune
	for _, line := range readDataLines("day6.txt") {
		grid = append(grid, []rune(line))
	}
	grid = transposeGrid(grid)

	// Chunk into problems
	var problems [][]string
	var currentProblem []string
	for _, row := range grid {
		row := string(row)
		if strings.TrimSpace(row) == "" {
			problems = append(problems, currentProblem)
			currentProblem = nil
		} else {
			currentProblem = append(currentProblem, row)
		}
	}
	if currentProblem != nil {
		problems = append(problems, currentProblem)
	}

	ops := map[byte]func(int64, int64) int64{
		'*': func(a, b int64) int64 {
			return a * b
		},
		'+': func(a, b int64) int64 {
			return a + b
		},
	}

	var total int64
	for _, p := range problems {
		// fmt.Printf("\nProblem\n")
		// for _, row := range p {
		// 	fmt.Printf("  %q\n", row)
		// }

		// The first row is always max length, with operator in last column.
		ncols := len(p[0])
		op := p[0][ncols-1]
		assert(ops[op] != nil, "unsupposed op")

		parseRow := func(i int) int64 {
			var x int64
			for j := range min(len(p[i]), ncols-1) {
				if p[i][j] == ' ' {
					continue
				}
				x *= 10
				x += parseInt(string(p[i][j]))
			}
			return x
		}

		acc := parseRow(0)
		for i := 1; i < len(p); i++ {
			acc = ops[op](acc, parseRow(i))
		}

		// fmt.Printf("  acc = %d\n", acc)

		total += acc
	}

	fmt.Printf("Total: %d\n", total)
}
