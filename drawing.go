package main

import (
	"fmt"
	"math"
	"math/rand"

	"github.com/fogleman/gg"
)

type drawingContext struct {
	dc            *gg.Context
	maxX, maxY    float64
	width, height int
}

func newContextFromPoints(points []vec2) *drawingContext {
	c := &drawingContext{width: 1200, height: 1200}
	for _, p := range points {
		c.maxX = max(p.x, c.maxY)
		c.maxY = max(p.y, c.maxY)
	}
	c.dc = gg.NewContext(c.width, c.height)
	return c
}

func (c *drawingContext) scaleX(x float64) float64 {
	return 0.1*float64(c.width) + (0.8*float64(c.width)/c.maxX)*x
}

func (c *drawingContext) scaleY(y float64) float64 {
	return 0.1*float64(c.height) + (0.8*float64(c.height)/c.maxY)*y
}

func (c *drawingContext) drawConnectedPoints(points []vec2) {
	dc := c.dc
	for i, a := range points {
		b := points[(i+1)%len(points)]
		dc.Push()
		dc.DrawLine(c.scaleX(a.x), c.scaleY(a.y), c.scaleX(b.x), c.scaleY(b.y))
		dc.SetLineWidth(1)
		dc.Stroke()
		dc.Pop()
	}
}

// drawConnectedPoints draws lines between adjacent points
//
// We assume all the values are positive for now.
func drawConnectedPoints(points []vec2, pngFileName string) error {
	c := newContextFromPoints(points)

	c.dc.SetRGB(1, 1, 1)
	c.dc.Clear()

	c.dc.SetRGB(1, 0, 0)
	c.drawConnectedPoints(points)

	return c.dc.SavePNG(pngFileName)
}

type rectangle struct {
	x1, x2, y1, y2 float64
}

func drawConnectedPointsAndRectangle(points []vec2, rect rectangle, pngFileName string) error {
	return drawConnectedPointsAndRectangles(points, []rectangle{rect}, pngFileName)
}

func drawConnectedPointsAndRectangles(points []vec2, rects []rectangle, pngFileName string) error {
	c := newContextFromPoints(points)

	c.dc.SetRGB(1, 1, 1)
	c.dc.Clear()

	c.dc.SetRGB(1, 0, 0)
	c.drawConnectedPoints(points)

	for _, rect := range rects {
		// Random-ish blue color
		h := 200.0 + rand.Float64()*80.0
		s := 1.0 - rand.Float64()*0.3
		v := 1.0 - rand.Float64()*0.4
		c.dc.SetRGB(hsvToRGB(h, s, v))

		rpoints := make([]vec2, 4)
		rpoints[0] = vec2{rect.x1, rect.y1}
		rpoints[1] = vec2{rect.x2, rect.y1}
		rpoints[2] = vec2{rect.x2, rect.y2}
		rpoints[3] = vec2{rect.x1, rect.y2}
		c.drawConnectedPoints(rpoints)
	}
	return c.dc.SavePNG(pngFileName)
}

func drawConnectedPointsWithRowColExtents(points []vec2, rowExtents map[int]*intervalSet[int], colExtents map[int]*intervalSet[int], pngFileName string) error {
	c := newContextFromPoints(points)

	c.dc.SetRGB(1, 1, 1)
	c.dc.Clear()

	for _, k := range SortedKeys(rowExtents) {
		if k%500 != 0 {
			continue
		}
		dc := c.dc
		for _, ivl := range rowExtents[k].ivls {
			dc.Push()
			dc.SetRGB(0, 1, 0)
			dc.DrawLine(c.scaleX(float64(k)), c.scaleY(float64(ivl.min)), c.scaleX(float64(k)), c.scaleY(float64(ivl.max-1)))
			dc.SetLineWidth(1)
			dc.Stroke()
			dc.Pop()
		}
	}

	for _, k := range SortedKeys(colExtents) {
		if k%500 != 0 {
			continue
		}
		dc := c.dc
		for _, ivl := range colExtents[k].ivls {
			dc.Push()
			dc.SetRGB(0, 0, 1)
			dc.DrawLine(c.scaleX(float64(ivl.min)), c.scaleY(float64(k)), c.scaleX(float64(ivl.max-1)), c.scaleY(float64(k)))
			dc.SetLineWidth(1)
			dc.Stroke()
			dc.Pop()
		}
	}

	c.dc.SetRGB(1, 0, 0)
	c.drawConnectedPoints(points)

	return c.dc.SavePNG(pngFileName)
}

type asciiGrid struct {
	grid             [][]rune
	numRows, numCols int
}

func newAsciiGrid(rows, cols int, initValue rune) *asciiGrid {
	g := &asciiGrid{
		grid:    make([][]rune, rows),
		numRows: rows,
		numCols: cols,
	}
	for i := range g.grid {
		g.grid[i] = make([]rune, cols)
		for j := range g.grid[i] {
			g.grid[i][j] = initValue
		}
	}
	return g
}

func (g *asciiGrid) setCell(r, c int, val rune) {
	g.grid[r][c] = val
}

func (g *asciiGrid) drawWithHeaders() {
	// Draw header
	fmt.Printf("  ")
	for i := range g.numCols {
		fmt.Printf("%2d", i)
	}
	fmt.Println()
	for i, row := range g.grid {
		fmt.Printf("%2d", i)
		for _, cell := range row {
			fmt.Printf(" %c", cell)
		}
		fmt.Println()
	}
	fmt.Println()
}

// hsvToRGB converts HSV values (h: 0-360, s: 0-1, v: 0-1) to RGB values (r, g, b: 0-1).
func hsvToRGB(h, s, v float64) (r, g, b float64) {
	kr := math.Mod(5+h*6, 6)
	kg := math.Mod(3+h*6, 6)
	kb := math.Mod(1+h*6, 6)

	r = 1 - math.Max(min3(kr, 4-kr, 1), 0)
	g = 1 - math.Max(min3(kg, 4-kg, 1), 0)
	b = 1 - math.Max(min3(kb, 4-kb, 1), 0)

	return
}
