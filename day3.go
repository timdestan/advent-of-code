package main

import (
	"fmt"
)

func day3() {
	parseLine := func(s string) []int {
		var res []int
		for i := range s {
			res = append(res, parseInt(s[i:i+1]))
		}
		return res
	}
	maxJoltage := func(nums []int) int {
		i := 0
		numBatteries := 12
		var total int
		for numBatteries > 0 {
			// Find the best battery to use. We need to leave space
			// for the remaining batteries.
			limit := len(nums) - (numBatteries - 1)
			maxJ := i
			for j := i; j < limit; j++ {
				if nums[j] > nums[maxJ] {
					maxJ = j
				}
			}
			// fmt.Printf("best x[%d] = %d\n", maxJ, nums[maxJ])
			total *= 10
			total += nums[maxJ]
			i = maxJ + 1
			numBatteries--
		}
		// fmt.Printf("max joltage is %d\n", total)
		return total
	}

	var total int
	for _, line := range readDataLines("day3.txt") {
		nums := parseLine(line)
		total += maxJoltage(nums)
	}
	fmt.Printf("total: %d\n", total)
}

func init() {
	functionRegistry[3] = day3
}
