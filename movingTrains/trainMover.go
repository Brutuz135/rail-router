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

		for ti := 0; ti < nTrains; ti++ {
			p := paths[trainPath[ti]]
			if pos[ti] >= len(p)-1 {
				done++
				continue
			}

			// Move one step forward.
			pos[ti]++
			stationID := p[pos[ti]]
			turnMoves = append(turnMoves, move{
				train:  ti + 1,
				target: idToStations[stationID],
			})
		}

		if done == nTrains {
			break
		}

		// Stable required output order: T1, T2, T3...
		sort.Slice(turnMoves, func(i, j int) bool {
			return turnMoves[i].train < turnMoves[j].train
		}) // sort.Slice sorts using a less function. 

		out := make([]string, 0, len(turnMoves))
		for _, m := range turnMoves {
			out = append(out, fmt.Sprintf("T%d-%s", m.train, m.target))
		}
		if len(out) > 0 {
			fmt.Fprintln(os.Stdout, strings.Join(out, " "))
		}
	}
}

