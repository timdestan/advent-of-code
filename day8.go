package main

import (
	"fmt"
	"slices"
	"strings"
)

func day8() {
	var points []vec3
	for _, line := range readDataLines("day8.txt") {
		parts := strings.Split(line, ",")
		assert(len(parts) == 3, "bad line")
		points = append(points, vec3{
			x: float64(parseInt(parts[0])),
			y: float64(parseInt(parts[1])),
			z: float64(parseInt(parts[2])),
		})
	}

	n := len(points)

	type link struct {
		i, j int
		dist float64
	}

	var distances []link
	for i := range n {
		for j := i + 1; j < n; j++ {
			distances = append(distances, link{
				i:    i,
				j:    j,
				dist: euclideanDistance(&points[i], &points[j]),
			})
		}
	}

	slices.SortFunc(distances, func(a, b link) int {
		// Assuming no funny business with NaNs...
		if a.dist < b.dist {
			return -1
		} else if a.dist > b.dist {
			return 1
		} else {
			return 0
		}
	})

	set := NewDisjointSet(n)

	areAllConnected := func() bool {
		x := set.Find(0)
		for i := range set.Size {
			if set.Find(i) != x {
				return false
			}
		}
		return true
	}

	for _, l := range distances {
		set.Union(l.i, l.j)
		if areAllConnected() {
			fmt.Printf("Product of last connection: %d.\n", int64(points[l.i].x)*int64(points[l.j].x))
			break
		}
	}

	// clusters := make(map[int]int)
	// for i := range n {
	// 	clusters[set.Find(i)]++
	// }

	// counts := slices.Collect(maps.Values(clusters))
	// slices.Sort(counts)
	// slices.Reverse(counts)

	// product := 1
	// for _, x := range counts[:3] {
	// 	product *= x
	// }

	// fmt.Printf("Product of first 3 cluster: %d.\n", product)
}

func init() {
	functionRegistry[8] = day8
}
