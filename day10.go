package main

import (
	"fmt"
	"math"
	"slices"
	"strings"
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
		fmt.Printf("machine: %v\ntotalJoltage: %d\n", m, sum(m.targetJoltage))

		type soln struct {
			state []int
			cost  int
		}

		isAllZeroes := func(xs []int) bool {
			for _, x := range xs {
				if x != 0 {
					return false
				}
			}
			return true
		}

		key := func(xs []int) string {
			var sb strings.Builder
			for i, x := range xs {
				if i > 0 {
					sb.WriteString(",")
				}
				sb.WriteString(fmt.Sprintf("%d", x))
			}
			return sb.String()
		}

		applyInstructionBackwards := func(st []int, inst instruction) ([]int, bool) {
			newst := slices.Clone(st)
			for _, i := range inst.indices {
				if st[i] == 0 {
					return nil, false
				}
				newst[i] -= 1
			}
			return newst, true
		}

		costByState := make(map[string]int)

		var computeCost func(st []int) (int, bool)
		computeCost = func(st []int) (int, bool) {
			k := key(st)
			if v, ok := costByState[k]; ok {
				return v, true
			}
			if isAllZeroes(st) {
				return 0, true
			}
			mincost := math.MaxInt64
			for _, inst := range m.instructions {
				newst, ok := applyInstructionBackwards(st, inst)
				if !ok {
					continue
				}
				cost, ok := computeCost(newst)
				if ok {
					mincost = min(mincost, cost+1)
				}
			}
			ok := mincost != math.MaxInt64
			if ok {
				costByState[k] = mincost
			}
			return mincost, ok
		}
		cost, ok := computeCost(m.targetJoltage)
		assert(ok)
		fmt.Printf("Cost is %d\n", cost)
		total += cost
	}

	fmt.Printf("Total min cost across all machines is %d\n", total)
}

func init() {
	functionRegistry[10] = day10
}
