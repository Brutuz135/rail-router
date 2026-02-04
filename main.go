package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	paths "stations/findingPaths"
	helpers "stations/helpers"
	moves "stations/movingTrains"
)

// computes the train moves given a map file, start/end stations, and number of trains.
func FindTrainTurns(mapPath, startName, endName string, trains int) ([][]moves.Move, error) {
	// Parse the map
	graph, idToStations, stationsToId, err := helpers.ParseMap(mapPath)
	if err != nil {
		return nil, fmt.Errorf("parsing map: %w", err)
	}

	// Validate start station
	start, ok := stationsToId[startName]
	if !ok {
		return nil, fmt.Errorf("start station %q does not exist", startName)
	}

	// Validate end station
	end, ok := stationsToId[endName]
	if !ok {
		return nil, fmt.Errorf("end station %q does not exist", endName)
	}

	// Validate start ≠ end
	if start == end {
		return nil, fmt.Errorf("start and end stations are the same")
	}	

	// Compute all possible paths through the graph.
	allPaths := paths.PathFinder(graph, start, end, trains)
	if len(allPaths) == 0 {
		return nil, fmt.Errorf("no path between the start and end stations")
	}

	// Compute train moves based on the paths.
	return moves.MoveTrains(allPaths, idToStations, trains), nil
}

func main() {
	args := os.Args[1:]
	if len(args) != 4 {
		fmt.Fprintf(os.Stderr, "Error: expected 4 arguments, got %d\n", len(os.Args)-1)
		fmt.Fprintln(os.Stderr, "Usage: go run . [map file] [start station] [end station] [number of trains]")
		os.Exit(1)
	}

	// Parse number of trains
	trains, err := strconv.Atoi(args[3])
	if err != nil || trains <= 0 {
		fmt.Fprintln(os.Stderr, "Error: number of trains is not a valid positive integer")
		os.Exit(1)
	}

	// Compute the train movements
	turns, err := FindTrainTurns(args[0], args[1], args[2], trains)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Print each turn in the format: 
	// T1-station T2-station ...
	for _, turn := range turns {
		if len(turn) == 0 {
			continue
		}
		var parts []string

		// Format: Tn-station
		for _, m := range turn {
			parts = append(parts, fmt.Sprintf("T%d-%s", m.Train, m.Target))
		}
		fmt.Println(strings.Join(parts, " "))
	}
}
