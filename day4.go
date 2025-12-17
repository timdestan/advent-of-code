package main

import "fmt"

func day4() {
	initGrid := func() [][]rune {
		var g [][]rune
		for _, row := range readDataLines("day4.txt") {
			var rowSlice []rune
			for _, c := range row {
				rowSlice = append(rowSlice, c)
			}
			g = append(g, rowSlice)
		}
		return g
	}

	grid := initGrid()

	countAdjacent := func(i, j int) int64 {
		var total int64
		for _, neighbor := range adjacentGridCells(grid, i, j) {
			if grid[neighbor.i][neighbor.j] == '@' {
				total += 1
			}
		}
		// fmt.Printf("countAdjacent(%d, %d) = %d\n", i, j, total)
		return total
	}

	total := 0

	for {
		accessibleRolls := 0
		for i, row := range grid {
			for j, cell := range row {
				if cell == '@' && countAdjacent(i, j) < 4 {
					// We are removing these as we go, which is different
					// from the example but should result in the same outcome.
					grid[i][j] = 'X'
					accessibleRolls += 1
				}
			}
		}

		// fmt.Printf("There are %d accessible rolls.\n", accessibleRolls)
		total += accessibleRolls
		if accessibleRolls == 0 {
			break
		}
	}
	fmt.Printf("Total: %d\n", total)
}

func init() {
	functionRegistry[4] = day4
}
