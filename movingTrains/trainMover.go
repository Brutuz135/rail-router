package moves

import (
	"fmt"
	"os"
	"strconv"
)

func MoveTrains(paths [][]int, IdToStations map[int]string, nTrains int) {

	//send trains trough shortes possible paths, until path is full
	//writes paths to the file acording of how many train there are. Incomplete. Needed: write the path column by column and and start them from the specific row when the train is sent.
	f, _ := os.Create("train_truns.txt")
	defer f.Close()

	// loop is going until there is no more trains to send
	for i := 0; nTrains > i; {

		//range paths and work with them one by one
		for _, path := range paths {

			//checks if it would more efficient to not use the path in question. Needs to be checked, how does it work on different situations.
			if len(path) > len(paths[0])*(nTrains-i) {
				break
			}

			i++

			//gives the train a number
			train := "T" + strconv.Itoa(i+1)

			// ranges trough a current path and writes it to the file with the train number
			for _, station := range path {
				str := train + " " + IdToStations[station] + "\n"
				fmt.Fprint(f, str)
			}
		}

		// if no more trains stop sending
		if i >= nTrains {
			break
		}
	}
}
