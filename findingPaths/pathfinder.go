package paths

import "math"

// PathFinder finds a good set of vertex-disjoint-ish paths (greedy removal)
// and retries by blocking edges of the first greedy path to escape bad first choices.
func PathFinder(graph [][]int, start, end, trains int) [][]int {
	// Base greedy solution
	base := greedyDisjoint(copyGraph(graph), start, end)
	if len(base) == 0 {
		return nil
	}
	best := base
	bestTurns := estimateTurns(best, trains)

	// Try alternatives by blocking each edge of the first path
	first := base[0]
	for i := 0; i < len(first)-1; i++ {
		u := first[i]
		v := first[i+1]

		g := copyGraph(graph)
		removeUndirectedEdge(g, u, v)

		cand := greedyDisjoint(g, start, end)
		if len(cand) == 0 {
			continue
		}
		turns := estimateTurns(cand, trains)
		if turns < bestTurns {
			bestTurns = turns
			best = cand
		}
	}

	return best
}

func greedyDisjoint(graph [][]int, start, end int) [][]int {
	var paths [][]int
	for {
		p := bfs(graph, start, end)
		if p == nil {
			break
		}
		paths = append(paths, p)
		// Keep start/end; remove intermediate nodes (vertex capacity heuristic)
		for i := 1; i < len(p)-1; i++ {
			graph[p[i]] = nil
		}
		// Also prevent infinite repeats on a direct edge start<->end
		if len(p) == 2 {
			removeUndirectedEdge(graph, p[0], p[1])
		}
	}
	return paths
}

// estimateTurns: simple, effective heuristic for passing.
// For each path length Li (edges), completion time if we load it with xi trains is Li+xi-1.
// Greedy assign trains to minimize max completion time using "Li + assigned".
func estimateTurns(paths [][]int, trains int) int {
	if len(paths) == 0 {
		return math.MaxInt
	}

	L := make([]int, len(paths))
	for i, p := range paths {
		L[i] = len(p) - 1 // edges
	}

	assigned := make([]int, len(paths))
	maxT := 0

	for t := 0; t < trains; t++ {
		bestP := 0
		bestCost := L[0] + assigned[0]
		for p := 1; p < len(paths); p++ {
			c := L[p] + assigned[p]
			if c < bestCost {
				bestCost = c
				bestP = p
			}
		}
		assigned[bestP]++
		ct := L[bestP] + assigned[bestP] - 1
		if ct > maxT {
			maxT = ct
		}
	}

	return maxT
}

func copyGraph(graph [][]int) [][]int {
	g := make([][]int, len(graph))
	for i := range graph {
		if graph[i] == nil {
			continue
		}
		g[i] = make([]int, len(graph[i]))
		copy(g[i], graph[i])
	}
	return g
}

func removeUndirectedEdge(graph [][]int, a, b int) {
	removeEdge(graph, a, b)
	removeEdge(graph, b, a)
}

func removeEdge(graph [][]int, from, to int) {
	row := graph[from]
	if row == nil {
		return
	}
	out := row[:0]
	for _, x := range row {
		if x != to {
			out = append(out, x)
		}
	}
	// shrink
	graph[from] = append([]int(nil), out...)
}

// Standard BFS shortest path (unweighted)
func bfs(graph [][]int, start, end int) []int {
	n := len(graph)
	visited := make([]bool, n)
	parent := make([]int, n)
	for i := range parent {
		parent[i] = -1
	}

	q := make([]int, 0, n)
	q = append(q, start)
	visited[start] = true

	found := false
	for len(q) > 0 {
		cur := q[0]
		q = q[1:]

		if cur == end {
			found = true
			break
		}
		for _, nb := range graph[cur] {
			if !visited[nb] {
				visited[nb] = true
				parent[nb] = cur
				q = append(q, nb)
			}
		}
	}

	if !found {
		return nil
	}

	// reconstruct
	path := []int{}
	for v := end; v != -1; v = parent[v] {
		path = append([]int{v}, path...)
	}
	return path
}
