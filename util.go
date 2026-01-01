package main

// Embed FS to get access to the input file
import (
	"cmp"
	"embed"
	"fmt"
	"maps"
	"math"
	"slices"
	"strconv"
	"strings"

	"golang.org/x/exp/constraints"
)

func assert(cond bool, msgs ...string) {
	if !cond {
		if len(msgs) == 0 {
			panic("assertion failed!")
		} else {
			panic("assertion failed: " + strings.Join(msgs, ""))
		}
	}
}

func assertOk(err error) {
	if err != nil {
		panic(err)
	}
}

func must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}

func slicesEqual[T comparable](a, b []T) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func sum[T Number](xs []T) T {
	var acc T
	for _, x := range xs {
		acc += x
	}
	return acc
}

//go:embed data/*.txt
var dataFS embed.FS

func readDataFile(filename string) string {
	return string(must(dataFS.ReadFile("data/" + filename)))
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

func parseInt(s string) int {
	return int(must(strconv.ParseInt(s, 10, 64)))
}

func parseIntsWithSep(s string, sep string) []int {
	var ints []int
	for piece := range strings.SplitSeq(s, sep) {
		ints = append(ints, parseInt(piece))
	}
	return ints
}

func parseFloat(s string) float64 {
	return must(strconv.ParseFloat(s, 64))
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

type DisjointSet struct {
	Size    int
	Parents []int
}

func NewDisjointSet(size int) *DisjointSet {
	s := &DisjointSet{
		Size:    size,
		Parents: make([]int, size),
	}
	for i := range size {
		s.Parents[i] = i
	}
	return s
}

func (s *DisjointSet) Find(x int) int {
	if s.Parents[x] != x {
		s.Parents[x] = s.Find(s.Parents[x])
		return s.Parents[x]
	} else {
		return x
	}
}

func (s *DisjointSet) Union(x, y int) {
	// Replace nodes by roots
	x = s.Find(x)
	y = s.Find(y)
	if x == y {
		return
	}
	s.Parents[y] = x
}

func SortedKeys[Map ~map[K]V, K cmp.Ordered, V any](m Map) []K {
	return slices.Sorted(maps.Keys(m))
}

type vec2 struct {
	x, y float64
}

type vec3 struct {
	x, y, z float64
}

func euclideanDistance(a, b *vec3) float64 {
	dx := a.x - b.x
	dy := a.y - b.y
	dz := a.z - b.z
	return math.Sqrt(dx*dx + dy*dy + dz*dz)
}

type Number interface {
	constraints.Integer | constraints.Float
}

func abs[T Number](x T) T {
	if x < 0 {
		return -x
	}
	return x
}

func sign[T Number](x T) int {
	if x < 0 {
		return -1
	}
	if x > 0 {
		return 1
	}
	return 0
}

func min[T cmp.Ordered](a, b T) T {
	if a < b {
		return a
	}
	return b
}

func min3[T cmp.Ordered](a, b, c T) T {
	return min(min(a, b), c)
}

type interval[T constraints.Integer] struct {
	min T // inclusive
	max T // exclusive
}

func (ivl *interval[T]) Contains(x T) bool {
	return x >= ivl.min && x < ivl.max
}

func (ivl *interval[T]) Intersect(other *interval[T]) *interval[T] {
	return &interval[T]{
		min: min(ivl.min, other.min),
		max: max(ivl.max, other.max),
	}
}

func (ivl *interval[T]) Length() T {
	return ivl.max - ivl.min
}

func (ivl *interval[T]) IsValid() bool {
	return ivl.min < ivl.max
}

func (ivl *interval[T]) String() string {
	return fmt.Sprintf("[%v, %v)", ivl.min, ivl.max)
}

type intervalSet[T constraints.Integer] struct {
	ivls      []*interval[T]
	compacted bool
}

func (set *intervalSet[T]) Add(ivl interval[T]) {
	set.ivls = append(set.ivls, &ivl)
	set.compacted = false
}

func (set *intervalSet[T]) AddPoint(point T) {
	set.ivls = append(set.ivls, &interval[T]{min: point, max: point + 1})
	set.compacted = false
}

func (set *intervalSet[T]) Contains(v T) bool {
	for _, ivl := range set.ivls {
		if ivl.Contains(v) {
			return true
		}
	}
	return false
}

// FindMatchingInterval finds the first interval that contains the point.
func (set *intervalSet[T]) FindMatchingInterval(v T) *interval[T] {
	for _, ivl := range set.ivls {
		if ivl.Contains(v) {
			return ivl
		}
	}
	return nil
}

func (set *intervalSet[T]) ContainsInterval(ivl interval[T]) bool {
	assert(ivl.IsValid(), "ContainsInterval() called with invalid interval.")
	assert(set.compacted, "ContainsInterval() called on non-compacted set.")
	// When compacted we know the set ivls are in increasing order and
	// do not overlap.
	for _, setIvl := range set.ivls {
		if setIvl.Contains(ivl.min) && setIvl.Contains(ivl.max-1) {
			return true
		}
	}
	return false
}

func (set *intervalSet[T]) Compact() {
	if len(set.ivls) == 0 {
		return
	}
	slices.SortFunc(set.ivls, func(a, b *interval[T]) int {
		return int(a.min - b.min)
	})
	var compacted []*interval[T]
	var curr *interval[T]
	for _, ivl := range set.ivls {
		if curr == nil {
			curr = ivl
		} else if ivl.min <= curr.max {
			curr = &interval[T]{
				min: min(curr.min, ivl.min),
				max: max(curr.max, ivl.max),
			}
		} else {
			compacted = append(compacted, curr)
			curr = ivl
		}
	}
	compacted = append(compacted, curr)
	set.ivls = compacted
	set.compacted = true
}

func (set *intervalSet[T]) String() string {
	var sb strings.Builder
	for i, ivl := range set.ivls {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(fmt.Sprintf("%v", ivl))
	}
	return sb.String()
}

//
// Heap
//

type Heap[T any] struct {
	Cmp          func(T, T) bool
	Items        []T
	Pushes, Pops int
}

func (h Heap[T]) Len() int           { return len(h.Items) }
func (h Heap[T]) Less(i, j int) bool { return h.Cmp(h.Items[i], h.Items[j]) }
func (h Heap[T]) Swap(i, j int)      { h.Items[i], h.Items[j] = h.Items[j], h.Items[i] }

func (h *Heap[T]) Push(x any) {
	h.Pushes++
	h.Items = append(h.Items, x.(T))
}

func (h *Heap[T]) Pop() any {
	h.Pops++
	old := h.Items
	n := len(old)
	x := old[n-1]
	h.Items = old[0 : n-1]
	return x
}
