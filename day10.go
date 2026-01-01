package main

// This one requires lp_solve to be installed and available
// in cgo's include path,
// see https://pkg.go.dev/github.com/draffensperger/golp#section-readme

import (
	"fmt"
	"math"
	"slices"
	"strings"

	"github.com/draffensperger/golp"
)

func day10() {
	type instruction struct {
		indices []int
	}

	type machine struct {
		targetState   []bool
		targetJoltage []int
		instructions  []instruction
	}

	var machines []machine
	for _, line := range readDataLines("day10.txt") {
		var m machine
		assert(line[0] == '[')

		end := strings.IndexRune(line, ']')
		for i := 1; i < end; i++ {
			switch line[i] {
			case '#':
				m.targetState = append(m.targetState, true)
			case '.':
				m.targetState = append(m.targetState, false)
			default:
				panic("bad line: " + line)
			}
		}

		fields := strings.Fields(strings.TrimSpace(line[end+1:]))
		for _, f := range fields {
			// fmt.Printf("field: %s\n", f)
			switch f[0] {
			case '(':
				inst := instruction{indices: parseIntsWithSep(f[1:len(f)-1], ",")}
				m.instructions = append(m.instructions, inst)
			case '{':
				// These should always be last.
				joltage := parseIntsWithSep(f[1:len(f)-1], ",")
				assert(len(joltage) == len(m.targetState))
				m.targetJoltage = joltage
			default:
				panic("bad field: " + f)
			}
		}
		machines = append(machines, m)
	}
	total := 0

	for _, m := range machines {
		// fmt.Printf("machine: %v\ntotalJoltage: %d\n", m, sum(m.targetJoltage))
		numRows := len(m.targetJoltage)
		numCols := len(m.instructions)

		lp := golp.NewLP(numRows, numCols)
		// Set to minimize by default.
		for i, joltage := range m.targetJoltage {
			row := make([]float64, numCols)
			for j, inst := range m.instructions {
				if slices.Contains(inst.indices, i) {
					row[j] = 1.0
				}
			}
			lp.AddConstraint(row, golp.EQ, float64(joltage))
		}
		{
			objFn := make([]float64, numCols)
			for j := range objFn {
				objFn[j] = 1.0
				lp.SetInt(j, true)
			}
			lp.SetObjFn(objFn)
		}
		// fmt.Printf("LP: %v", lp.WriteToString())
		assert(lp.Solve() == golp.OPTIMAL, "failed lp_solve")
		cost := 0
		for _, v := range lp.Variables() {
			assert(v == math.Trunc(v), "non-integer solution")
			cost += int(v)
		}
		assert(cost > 0)
		// fmt.Printf("Cost is %d\n", cost)
		total += cost
	}

	fmt.Printf("Total min cost across all machines is %d\n", total)
}

func init() {
	functionRegistry[10] = day10
}
