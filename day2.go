package main

import (
	"fmt"
	"strings"
)

func day2() {
	isSilly := func(s string) bool {
	outerLoop:
		for patternSize := 1; patternSize <= len(s)/2; patternSize++ {
			if len(s)%patternSize != 0 {
				continue
			}
			pattern := s[:patternSize]
			for i := patternSize; i < len(s); i += patternSize {
				if s[i:i+patternSize] != pattern {
					continue outerLoop
				}
			}
			return true
		}
		return false
	}
	var sumOfSilly int64

	for _, group := range strings.Split(readDataFile("day2.txt"), ",") {
		group = strings.TrimSpace(group)
		if group == "" {
			continue
		}
		dashIndex := strings.Index(group, "-")
		if dashIndex < 0 {
			panic("bad range:" + group)
		}
		start := parseInt(group[:dashIndex])
		end := parseInt(group[dashIndex+1:])
		for i := start; i <= end; i++ {
			if isSilly(fmt.Sprintf("%d", i)) {
				// fmt.Printf("%d is silly.\n", i)
				sumOfSilly += i
			}
		}
	}

	fmt.Printf("Sum of silly: %d\n", sumOfSilly)
}
