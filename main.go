package main

import (
	paths "main/findingPaths"
	moves "main/movingTrains"
	"os"
	"strconv"
)

func main() {
	//Incomplete, needs to handle errors gracefully. But the main structure works.


	// get infromation sortet from the given map.
	//stationsGraph slice of ints containing all the connections of stationsGraph
	//idToStations a map to convert from int id to string sations, there is also stationsToId if needed.
	stationsGraph, idToStations, _ := parseMap(os.Args[1]) 

	//Handle arguments
	start, _ := strconv.Atoi(os.Args[2]) // now works is station is given like a number, but needs to handle names
	end, _ := strconv.Atoi(os.Args[3])    // now works is station is given like a number, but needs to handle names
	trains, _ := strconv.Atoi(os.Args[4])

	//finds all possible paths, with totally different in between stations and sorts them from shortest to longest
	paths := paths.PathFinder(stationsGraph, start, end)
	
	//writes out paths to a file acording to how many trains
	//Incomplete, Should write in columns and acording when train is sent. row represents a turn
	moves.MoveTrains(paths, idToStations, trains)
	
	//maybe here print out the train_turns file?

}
