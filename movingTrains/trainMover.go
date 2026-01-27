package moves

import (
	"fmt"
	"os"
	"sort"
	"strings"
)

type move struct {
	train  int
	target string
}

func MoveTrains(paths [][]int, idToStations map[int]string, nTrains int) {
	if nTrains <= 0 || len(paths) == 0 {
		return
	}

	// Sort paths by length so assignment is stable.
	sort.Slice(paths, func(i, j int) bool { return len(paths[i]) < len(paths[j]) })

	// Cost-based assignment: pick path minimizing (pathLenEdges + alreadyAssigned)
	pathEdges := make([]int, len(paths))
	for i, p := range paths {
		pathEdges[i] = len(p) - 1
	}

	assigned := make([]int, len(paths))
	trainPath := make([]int, nTrains) // train i uses paths[trainPath[i]]

	for t := 0; t < nTrains; t++ {
		bestP := 0
		bestCost := pathEdges[0] + assigned[0]
		for p := 1; p < len(paths); p++ {
			c := pathEdges[p] + assigned[p]
			if c < bestCost {
				bestCost = c
				bestP = p
			}
		}
		trainPath[t] = bestP
		assigned[bestP]++
	}

	// Position = index in the path slice. Start at 0 (at the start station).
	pos := make([]int, nTrains)

	edgeKey := func(a, b int) [2]int {
		if a < b {
			return [2]int{a, b}
		}
		return [2]int{b, a}
	}

	for {
		var turnMoves []move
		done := 0

		occupiedStation := make(map[int]bool) // intermediate stations only
		usedEdge := make(map[[2]int]bool)     // each track can be used once per turn

		for ti := 0; ti < nTrains; ti++ {
			p := paths[trainPath[ti]]

			// Train already at destination
			if pos[ti] >= len(p)-1 {
				done++
				continue
			}

			curStation := p[pos[ti]]
			nextPos := pos[ti] + 1
			nextStation := p[nextPos]

			// Rule: each track can only be used once per turn (undirected edge)
			ek := edgeKey(curStation, nextStation)
			if usedEdge[ek] {
				continue
			}

			// Rule: only one train per intermediate station per turn (end station unlimited)
			if nextPos != len(p)-1 {
				if occupiedStation[nextStation] {
					continue
				}
			}

			// Apply move
			pos[ti] = nextPos
			usedEdge[ek] = true
			if nextPos != len(p)-1 {
				occupiedStation[nextStation] = true
			}

			turnMoves = append(turnMoves, move{
				train:  ti + 1,
				target: idToStations[nextStation],
			})
		}

		if done == nTrains {
			break
		}

		// Stable output order: T1, T2, ...
		sort.SliceStable(paths, func(i, j int) bool {
			if len(paths[i]) != len(paths[j]) {
				return len(paths[i]) < len(paths[j])
			}
			// tie-break: compare the full path
			a, b := paths[i], paths[j]
			m := len(a)
			if len(b) < m {
				m = len(b)
			}
			for k := 0; k < m; k++ {
				if a[k] != b[k] {
					return a[k] < b[k]
				}
			}
			return len(a) < len(b)
		})

		out := make([]string, 0, len(turnMoves))
		for _, m := range turnMoves {
			out = append(out, fmt.Sprintf("T%d-%s", m.train, m.target))
		}
		if len(out) > 0 {
			fmt.Fprintln(os.Stdout, strings.Join(out, " "))
		}
	}
}
