package main

import (
	"fmt"
	"strings"
)

func day9() {
	type point struct {
		x, y int
	}

	var points []point

	rowExtent := make(map[int]*intervalSet[int])
	colExtent := make(map[int]*intervalSet[int])

	updateExtents := func(p point) {
		if rowExtent[p.x] == nil {
			rowExtent[p.x] = &intervalSet[int]{}
		}
		rowExtent[p.x].AddPoint(p.y)
		if colExtent[p.y] == nil {
			colExtent[p.y] = &intervalSet[int]{}
		}
		colExtent[p.y].AddPoint(p.x)
	}

	var maxrow, maxcol int
	for _, line := range readDataLines("day9.txt") {
		parts := strings.Split(line, ",")
		assert(len(parts) == 2, "bad line")
		p := point{
			x: parseInt(parts[0]),
			y: parseInt(parts[1]),
		}
		updateExtents(p)
		maxrow = max(maxrow, p.x)
		maxcol = max(maxcol, p.y)
		points = append(points, p)
	}
	// grid := newAsciiGrid(maxrow+1, maxcol+1, '.')
	for i, p1 := range points {
		// grid.setCell(p1.x, p1.y, '#')
		j := (i + 1) % len(points)
		p2 := points[j]
		if p1.x == p2.x {
			ydir := sign(p2.y - p1.y)
			q := p1
			for q.y != p2.y {
				updateExtents(q)
				q.y += ydir
			}
		} else if p1.y == p2.y {
			xdir := sign(p2.x - p1.x)
			q := p1
			for q.x != p2.x {
				updateExtents(q)
				q.x += xdir
			}
		} else {
			panic("expect either same row or same col")
		}
	}
	// grid.drawWithHeaders()

	customCompact := func(set *intervalSet[int]) {
		if len(set.ivls) == 0 {
			return
		}
		set.Compact()
		inside := false
		var compacted []*interval[int]
		var curr *interval[int]
		for _, ivl := range set.ivls {
			assert(curr == nil || ivl.min >= curr.min, "invervalSet invariant")
			if curr == nil {
				curr = ivl
			} else if inside {
				// fmt.Printf("..Merging %v and %v.\n", curr, ivl)
				curr = &interval[int]{
					min: curr.min,
					max: ivl.max,
				}
			} else {
				// fmt.Printf("..Adding %v.\n", curr)
				compacted = append(compacted, curr)
				curr = ivl
			}
			inside = !inside
		}
		// fmt.Printf("..Adding %v.\n", curr)
		compacted = append(compacted, curr)
		set.ivls = compacted
	}
	for _, k := range SortedKeys(rowExtent) {
		customCompact(rowExtent[k])
		// fmt.Printf("rowExtent[%d]: %v\n", k, rowExtent[k])
	}
	for _, k := range SortedKeys(colExtent) {
		customCompact(colExtent[k])
		// fmt.Printf("colExtent[%d]: %v\n", k, colExtent[k])
	}

	ordered := func(a, b int) (int, int) {
		if a > b {
			return b, a
		}
		return a, b
	}

	var best rectangle
	var rectangles []rectangle
	var bestArea int

	for i, p1 := range points {
		for j := i + 1; j < len(points); j++ {
			p2 := points[j]

			x1, x2 := ordered(p1.x, p2.x)
			y1, y2 := ordered(p1.y, p2.y)

			xivl := &interval[int]{x1, x2 + 1}
			yivl := &interval[int]{y1, y2 + 1}

			area := xivl.Length() * yivl.Length()

			// fmt.Printf(". area = %d", area)
			if !colExtent[yivl.min].ContainsInterval(*xivl) ||
				!colExtent[yivl.max-1].ContainsInterval(*xivl) ||
				!rowExtent[xivl.min].ContainsInterval(*yivl) ||
				!rowExtent[xivl.max-1].ContainsInterval(*yivl) {
				// fmt.Printf(" ... FAIL\n")
				continue
			} else {
				// fmt.Printf(" ... OK\n")
			}
			var rect rectangle
			rect.x1 = float64(xivl.min)
			rect.x2 = float64(xivl.max - 1)
			rect.y1 = float64(yivl.min)
			rect.y2 = float64(yivl.max - 1)

			if area > bestArea {
				bestArea = area
				best = rect
				rectangles = append(rectangles, rect)
			}
		}
	}

	fmt.Printf("Largest area: %d (%v)\n", bestArea, best)
	// 1377278750 is too low

	convertPoints := func(points []point) []vec2 {
		vs := make([]vec2, len(points))
		for i, p := range points {
			vs[i] = vec2{
				x: float64(p.x),
				y: float64(p.y),
			}
		}
		return vs
	}

	assertOk(drawConnectedPointsAndRectangle(convertPoints(points), best, "out.png"))
	assertOk(drawConnectedPointsAndRectangles(convertPoints(points), rectangles, "out_all_rects.png"))
	assertOk(drawConnectedPointsWithRowColExtents(convertPoints(points), rowExtent, colExtent, "out_row_col_extents.png"))
}

func init() {
	functionRegistry[9] = day9
}
