package main

import (
	"fmt"
	paths "here/findingPaths"
	"os"
)

func main() {


	stationsGraph, idToStations, _ := parseMap(os.Args[1]) 
	//start := os.Args[2] // change to int!!!
	//end := os.Args[3]   //change to int!!!

	start := 0
	end := 4
	trains := 9

	fmt.Println(stationsGraph)
	fmt.Println(start)
	fmt.Println(end)

	paths := paths.PathFinder(stationsGraph, start, end)
	trains.MoveTrains(paths, idToStations, 9)


	//read arguments and check for errors


	//call pathFinder(), networkfile, start, end, (number of trains???)


	//call trainCender() send and track trains trough optimal routes.

	//end function

}
