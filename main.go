package main

import (
	"fmt"
	"os"
	"strconv"

	paths "stations/findingPaths"
	helpers "stations/helpers"
	moves "stations/movingTrains"
)

func main() {
	args := os.Args[1:]
	if len(args) != 4 {
		fmt.Fprintln(os.Stderr, "Error: incorrect number of command line arguments")
		os.Exit(1)
	}

	filepath := args[0]
	startName := args[1]
	endName := args[2]
	numTrainStr := args[3]

	// Parse number of trains
	trains, err := strconv.Atoi(numTrainStr)
	if err != nil || trains <= 0 {
		fmt.Fprintln(os.Stderr, "Error: number of trains is not a valid positive integer")
		os.Exit(1)
	}

	// fmt.Printf(
	// 	"%s,  %s,  %d\n",
	// 	startName, endName, trains,
	// )

	// Parse map
	stationsGraph, idToStations, stationsToId, err := helpers.ParseMap(filepath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing map: %v\n", err)
		os.Exit(1)
	}

	// Validate start & end stations
	start, ok := stationsToId[startName]
	if !ok {
		fmt.Fprintf(os.Stderr, "Error: start station %q does not exist\n", startName)
		os.Exit(1)
	}

	end, ok := stationsToId[endName]
	if !ok {
		fmt.Fprintf(os.Stderr, "Error: end station %q does not exist\n", endName)
		os.Exit(1)
	}

	// Check if start and end are the same
	if start == end {
		fmt.Fprintf(os.Stderr, "Error: start and end station are the same")
		os.Exit(1)
	}

	// Find paths
	allPaths := paths.PathFinder(stationsGraph, start, end, trains)
	if len(allPaths) == 0 {
		fmt.Fprintln(os.Stderr, "Error: no path between the start and end stations")
		os.Exit(1)
	}

	// Move trains
	moves.MoveTrains(allPaths, idToStations, trains)
}
