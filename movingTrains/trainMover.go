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

	// Assign trains to paths round-robin.
	trainPath := make([]int, nTrains) // train i uses paths[trainPath[i]]
	for i := 0; i < nTrains; i++ {
		trainPath[i] = i % len(paths)
	}

	// Position = index in the path slice. Start at 0 (at the start station).
	pos := make([]int, nTrains)

	// Run turns until all trains reach the end (pos == len(path)-1).
	for {
		var turnMoves []move
		done := 0

		// Track which stations are occupied THIS turn
		occupied := make(map[int]bool)

		for ti := 0; ti < nTrains; ti++ {
			p := paths[trainPath[ti]]
			// Train already at destination
			if pos[ti] >= len(p)-1 {
				done++
				continue
			}

			nextPos := pos[ti] + 1
			nextStation := p[nextPos]

			// Only one train per intermediate station per turn
			if nextPos != len(p)-1 { // not end station
				if occupied[nextStation] {
					continue // wait this turn
				}
			}

			// move train
			pos[ti] = nextPos
			occupied[nextStation] = true

			turnMoves = append(turnMoves, move{
				train:  ti + 1,
				target: idToStations[nextStation],
			})
		}

		// All trains finished
		if done == nTrains {
			break
		}

		// Stable output order: T1, T2, ...
		sort.Slice(turnMoves, func(i, j int) bool {
			return turnMoves[i].train < turnMoves[j].train
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
