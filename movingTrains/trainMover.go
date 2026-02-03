package moves

import (
	"sort"
)

type Move struct {
	Train  int
	Target string
}

// MoveTrains generates the sequence of moves for nTrains
func MoveTrains(paths [][]int, idToStations map[int]string, nTrains int) [][]Move {
	if nTrains <= 0 || len(paths) == 0 {
		return nil
	}

	// Sort paths so shorter ones are preferred for assignment
	sort.Slice(paths, func(i, j int) bool { return len(paths[i]) < len(paths[j]) })

	// Precompute the number of edges in each path
	pathEdges := make([]int, len(paths))
	for i, p := range paths {
		pathEdges[i] = len(p) - 1
	}

	// Assign each train to a path using a simple load‑balancing heuristic:
	// cost = path length + number of trains already assigned to that path
	assigned := make([]int, len(paths))
	trainPath := make([]int, nTrains)

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

	// Track each train's current position along its assigned path
	pos := make([]int, nTrains)

	// Normalize an edge so (A→B) and (B→A) are treated the same
	edgeKey := func(a, b int) [2]int {
		if a < b {
			return [2]int{a, b}
		}
		return [2]int{b, a}
	}

	var allTurns [][]Move

	// Simulate movement until all trains reach their destination
	for {
		var turnMoves []Move
		done := 0

		// Intermediate stations occupied this turn
		occupiedStation := make(map[int]bool)

		// Tracks used this turn (to prevent collisions)
		usedEdge := make(map[[2]int]bool)

		for ti := 0; ti < nTrains; ti++ {
			p := paths[trainPath[ti]]
			if pos[ti] >= len(p)-1 {
				done++
				continue
			}
			curStation := p[pos[ti]]
			nextPos := pos[ti] + 1
			nextStation := p[nextPos]

			// Prevent two trains from using the same track in the same turn
			ek := edgeKey(curStation, nextStation)
			if usedEdge[ek] {
				continue
			}

			// Prevent two trains from entering the same intermediate station
			if nextPos != len(p)-1 {
				if occupiedStation[nextStation] {
					continue
				}
			}

			// Train already at destination
			pos[ti] = nextPos
			usedEdge[ek] = true

			if nextPos != len(p)-1 {
				occupiedStation[nextStation] = true
			}

			turnMoves = append(turnMoves, Move{
				Train:  ti + 1,
				Target: idToStations[nextStation],
			})
		}

		if done == nTrains {
			break
		}
		allTurns = append(allTurns, turnMoves)
	}
	return allTurns
}
