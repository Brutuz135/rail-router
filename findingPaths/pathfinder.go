package paths

import (
	//	"container/heap"
	"fmt"

	// "math"
)

//func FindPaths {

//Initialize the map

//find the shortest route using Djikstras algorithm and length set to 1 for every connection

//Take out the fastest path, and it stations form the map, and find second fastest map with no same stations with the fastest

//Repeat until no more paths can be found.

//Return the paths in order

//}

func bfs(graph [][]int, start int, end int) []int {

	nNodes := len(graph)
//	fmt.Printf("length of nodes: %d\n", nNodes)
		


	
	var pq []int
	var visited []bool
	var parrent []int
	var current int
	pathFound := false
	var path []int


	for range nNodes {
//		fmt.Println("loops")
		visited = append(visited, false)
		parrent = append(parrent, -1)
	}

	pq = append(pq, start)
//	fmt.Println("pq:", pq)
	
	visited[start] = true
	parrent[start] = -1
	current = pq[0]
//	fmt.Println("visited: ", visited)
//	fmt.Println("parrent: ", parrent)
//	fmt.Println("lenght pq:", len(pq))

	for len(pq) > 0 {

		current = pq[0]
		fmt.Println("loops len")
		fmt.Println("pq: ", pq)
		fmt.Println("current", current)
	for _, neighbour := range graph[current] {
			println("loops neighbour")
			fmt.Println("visited: ", visited)
			if visited[neighbour] == false{
				println("if")
				visited[neighbour] = true
				parrent[neighbour] = current
				pq = append(pq, neighbour)
			}
		}


		pq = append(pq[:0], pq[1:]...)//take the current out of the priorityQueue

		if current == end {
			fmt.Println("end found")
			pathFound = true
			break
		}
	
	}


	
	if pathFound == false{
		return nil
	}
	
	cur := end
	path = append(path, end)
	for cur != start {
		cur = parrent[cur]
		path = append([]int{cur}, path...)
	}
	return path

}





	
func PathFinder(graph [][]int, start int, end int) [][]int {
	var paths [][]int

/*graph := [][]int{
    {5, 9, 1},      // 0: jungle connected to grasslands, farms, green_belt
    {0, 2},         // 1: green_belt connected to jungle, village
    {1, 3},         // 2: village connected to green_belt, mountain
    {4, 2, 8, 9},   // 3: mountain connected to treetop, village, wetlands, farms
    {3, 13},        // 4: treetop connected to mountain, desert
    {0, 6},         // 5: grasslands connected to jungle, suburbs
    {5, 7},         // 6: suburbs connected to grasslands, clouds
    {6, 8},         // 7: clouds connected to suburbs, wetlands
    {7, 13, 3},     // 8: wetlands connected to clouds, desert, mountain
    {0, 10, 3},     // 9: farms connected to jungle, downtown, mountain
    {9, 11},        // 10: downtown connected to farms, metropolis
    {10, 12},       // 11: metropolis connected to downtown, industrial
    {11, 13},       // 12: industrial connected to metropolis, desert
    {12, 8, 4},     // 13: desert connected to industrial, wetlands, treetop
}

start := 0
end := 13
*/
	path := bfs(graph, start, end)
	paths = append(paths, path)
	for path != nil {
		graph = ignoreUsedStations(graph, path)
		path = bfs(graph, start, end)
		paths = append(paths, path)
	}
	fmt.Printf("Shortest path from node %d to %d: %d\n", start, end, path)
	fmt.Println("Paths from shortes to longest:", paths)

	return paths
}

func ignoreUsedStations(graph [][]int, path []int) [][]int {
	for i := 1; i < len(path)-1; i++ {
		station := path[i]
		graph[station] = nil
	}
	return graph
}
