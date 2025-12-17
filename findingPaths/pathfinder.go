package paths

import (
	"container/heap"
	"fmt"
	"math"
)

type priorityQueue []*node
//func FindPaths {

//Initialize the map

//find the shortest route using Djikstras algorithm and length set to 1 for every connection

//Take out the fastest path, and it stations form the map, and find second fastest map with no same stations with the fastest

//Repeat until no more paths can be found.

//Return the paths in order

//}

func bfs(graph [][]int, start int, end int) []int {





	
func main() {

	graph := [][]int{
		{1, 2},    // waterloo
		{0, 3},    // victoria
		{0, 3},    // euston
		{1, 2},    // st_pancras
	}

	fmt.Println("The given nodes are:", graph)
	start := 0
	end := 4
	path := bfs(graph, start, end)
	fmt.Printf("Shortest path from node %d to %d: %d\n", start, end, path)
}
