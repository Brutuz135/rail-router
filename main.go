package main

import (
	"fmt"
	paths "here/findingPaths"
	moves "here/movingTrains"
	"os"
	"strconv"
)

func main() {


	stationsGraph, idToStations, stationsToId := parseMap(os.Args[1]) 

	fmt.Println("graph:", stationsGraph)
	fmt.Println("idToStationsi:", idToStations)
	fmt.Println("stationsToId", stationsToId)
	return

	start, _ := strconv.Atoi(os.Args[2]) // change to int!!!
	end, _ := strconv.Atoi(os.Args[3])   //change to int!!!

	trains := 9

	paths := paths.PathFinder(stationsGraph, start, end)
	moves.MoveTrains(paths, idToStations, trains)


	//read arguments and check for errors


	//call pathFinder(), networkfile, start, end, (number of trains???)


	//call trainCender() send and track trains trough optimal routes.

	//end function

}
