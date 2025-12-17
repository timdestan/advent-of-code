package main

import "fmt"

func day7() {
	var grid [][]rune
	for _, line := range readDataLines("day7.txt") {
		grid = append(grid, []rune(line))
	}

	var start gridCoord

	// printGrid := func() {
	// 	for _, line := range grid {
	// 		fmt.Printf("  %s\n", string(line))
	// 	}
	// }

	for i, line := range grid {
		for j, cell := range line {
			if cell == 'S' {
				start.i = i
				start.j = j
			}
		}
	}
	// fmt.Printf("Start = %v\n", start)

	// printGrid()

	// Values are the number of timelines in this state.
	rayPositions := map[int]int64{start.j: 1}
	var numSplits int64
	for i := start.i + 1; i < len(grid); i++ {
		nextpos := make(map[int]int64)

		for j, numTimelines := range rayPositions {
			// fmt.Printf("check [%d][%d]\n", i, j)
			switch grid[i][j] {
			case '.':
				nextpos[j] += numTimelines
			case '^':
				numSplits++
				if j-1 >= 0 {
					nextpos[j-1] += numTimelines
				}
				if j+1 < len(grid[i]) {
					nextpos[j+1] += numTimelines
				}
			default:
				panic("unexpected char: " + string(grid[i][j]))
			}
		}

		for j := range nextpos {
			grid[i][j] = '|'
		}
		// fmt.Printf("nextpos: %v\n", nextpos)
		rayPositions = nextpos
	}

	// printGrid()

	fmt.Printf("numSplits: %d\n", numSplits)
	var totalTimelines int64
	for _, num := range rayPositions {
		totalTimelines += num
	}
	fmt.Printf("totalTimelines: %d\n", totalTimelines)

}

func init() {
	functionRegistry[7] = day7
}
