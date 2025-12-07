package main

// Embed FS to get access to the input file
import (
	"embed"
	"path/filepath"
	"strconv"
	"strings"
)

//go:embed data/*.txt
var dataFS embed.FS

func assert(cond bool, msg string) {
	if !cond {
		panic(msg)
	}
}

func must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}

func readDataFile(filename string) string {
	return string(must(dataFS.ReadFile(filepath.Join("data", filename))))
}

func readDataLines(filename string) []string {
	var rv []string
	for _, line := range strings.Split(readDataFile(filename), "\n") {
		if line == "" {
			continue
		}
		rv = append(rv, line)
	}
	return rv
}

// readBlankSeparatedDataLines reads chunks of lines, separated by blank lines.
func readBlankSeparatedDataLines(filename string) [][]string {
	var rv [][]string
	for chunk := range strings.SplitSeq(readDataFile(filename), "\n\n") {
		var chunkLines []string
		for line := range strings.SplitSeq(chunk, "\n") {
			if line == "" {
				continue
			}
			chunkLines = append(chunkLines, line)
		}
		rv = append(rv, chunkLines)
	}
	return rv
}

func parseInt(s string) int64 {
	return must(strconv.ParseInt(s, 10, 64))
}

type gridCoord struct {
	i, j int
}

// adjacentGridCells finds the adjacent neighbors (up to 8) in a 2D grid.
func adjacentGridCells[T any](grid [][]T, i, j int) []gridCoord {
	var res []gridCoord
	for di := -1; di <= 1; di++ {
		for dj := -1; dj <= 1; dj++ {
			if 0 == di && 0 == dj {
				continue
			}
			ii := i + di
			jj := j + dj
			if ii < 0 || ii >= len(grid) {
				continue
			}
			if jj < 0 || jj >= len(grid[i]) {
				continue
			}
			res = append(res, gridCoord{i: ii, j: jj})
		}
	}
	return res
}

// transposeGrid transposes rows and columns
func transposeGrid[T any](grid [][]T) [][]T {
	var transposed [][]T
	nrows := len(grid)
	if nrows == 0 {
		return transposed
	}
	maxcols := 0
	for i := range grid {
		maxcols = max(maxcols, len(grid[i]))
	}
	for j := range maxcols {
		var trow []T
		for i := range nrows {
			if j >= len(grid[i]) {
				continue
			}
			trow = append(trow, grid[i][j])
		}
		transposed = append(transposed, trow)
	}
	return transposed
}
