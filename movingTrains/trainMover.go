package trains

import (
	"os"
)

func MoveTrains(paths [][]int, IdToStations (map[int]string), nTrains int) {

	//send trains trough shortes possible paths, until path is full

	//write every send on a file

	//send them onwards, and repeat.

	//print out the file, go to the next line on the file

	//repeat until all trains have reached the destination.

	os.Create("train_truns.txt")

	
	// write all used paths to the file in turn order. 
	for i := 0; nTrains > i; {
		for j, path := range paths {
			if len(path) > len(paths[0]) * (nTrains - i) {//checks if it would more efficient to not use the path in question. Needs to be checked, how does it work on different situations.
				break 
			}
			i++
			for _, station := range path{
				train := "T" + string(rune(i+j))
				str := train + IdToStations[station] + "\n"
				os.WriteFile("train_truns.txt", []byte(str), 7777)
			}
		}
	}


}
