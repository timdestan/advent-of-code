package main

import (
	"fmt"
	"slices"
	"strings"
)

func day5() {
	type dataRange struct {
		min, max int64
	}
	var numbers []int64
	var ranges []dataRange

	for _, line := range readDataLines("day5.txt") {
		if i := strings.IndexRune(line, '-'); i >= 0 {
			ranges = append(ranges, dataRange{
				min: parseInt(line[:i]),
				max: parseInt(line[i+1:]),
			})
		} else {
			numbers = append(numbers, parseInt(line))
		}
	}

	type endpoint struct {
		index     int64
		dataRange dataRange
		isStart   bool
	}
	var endpoints []endpoint
	for _, dr := range ranges {
		endpoints = append(endpoints, endpoint{
			index:     dr.min,
			dataRange: dr,
			isStart:   true,
		})
		endpoints = append(endpoints, endpoint{
			index:     dr.max + 1, // make exclusive
			dataRange: dr,
			isStart:   false,
		})
	}
	slices.SortStableFunc(endpoints, func(a, b endpoint) int {
		return int(a.index - b.index)
	})

	var validNumbers int64

	var lastIdx int64 = -1
	numOpen := 0
	for _, ep := range endpoints {
		// fmt.Printf("ep %v\n", ep)
		if numOpen > 0 {
			// fmt.Printf("total += %d\n", (ep.index - lastStart))
			validNumbers += (ep.index - lastIdx)
		}
		lastIdx = ep.index
		if ep.isStart {
			numOpen++
		} else {
			numOpen--
		}
	}

	fmt.Printf("%d possible valid.\n", validNumbers)
}
