package main

import (
	"fmt"
	"slices"
	"strings"
)

func day11() {
	var nodes []string
	nodeSet := make(map[string]bool)

	edges := make(map[string][]string)

	addNode := func(n string) {
		if !nodeSet[n] {
			nodes = append(nodes, n)
			nodeSet[n] = true
		}
	}

	addEdge := func(src, sink string) {
		addNode(src)
		addNode(sink)
		edges[src] = append(edges[src], sink)
	}

	for _, line := range readDataLines("day11.txt") {
		colon := strings.IndexRune(line, ':')
		assert(colon >= 0)
		in := strings.TrimSpace(line[:colon])
		for _, out := range strings.Fields(line[colon+1:]) {
			addEdge(in, strings.TrimSpace(out))
		}
	}

	computeTopologicalSort := func() []string {
		// Not detecting cycles

		var order []string
		seen := make(map[string]bool)

		var visit func(n string)
		visit = func(n string) {
			if seen[n] {
				return
			}
			for _, m := range edges[n] {
				visit(m)
			}
			seen[n] = true
			order = append(order, n)
		}

		remaining := slices.Clone(nodes)
		for len(remaining) > 0 {
			visit(remaining[0])
			remaining = remaining[1:]
		}

		slices.Reverse(order)
		return order
	}

	for k, v := range edges {
		fmt.Printf(" %v = %v\n", k, v)
	}

	type pathInfo struct {
		numComplete, numAwaitingFft, numAwaitingDac, numAwaitingBoth int
	}

	accPathInfo := func(acc, x *pathInfo, n string) {
		acc.numComplete += x.numComplete
		switch n {
		case "fft":
			acc.numComplete += x.numAwaitingFft
			acc.numAwaitingDac += x.numAwaitingDac
			acc.numAwaitingDac += x.numAwaitingBoth
		case "dac":
			acc.numAwaitingFft += x.numAwaitingFft
			acc.numComplete += x.numAwaitingDac
			acc.numAwaitingFft += x.numAwaitingBoth
		default:
			acc.numAwaitingFft += x.numAwaitingFft
			acc.numAwaitingDac += x.numAwaitingDac
			acc.numAwaitingBoth += x.numAwaitingBoth
		}
	}

	paths := make(map[string]*pathInfo)
	paths["svr"] = &pathInfo{numAwaitingBoth: 1}

	for _, x := range computeTopologicalSort() {
		for _, outgoing := range edges[x] {
			if paths[outgoing] == nil {
				paths[outgoing] = &pathInfo{}
			}
			accPathInfo(paths[outgoing], paths[x], outgoing)
		}
	}
	fmt.Printf("Num paths to out: %v\n", paths["out"].numComplete)
}

func init() {
	functionRegistry[11] = day11
}
