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





	
func PathFinder() {
	var paths [][]int

graph := [][]int{
    {1, 2, 4},   // 0: connected to 1,2,4
    {0, 7},      // 1: connected to 0,7
    {0, 3},      // 2: connected to 0,3
    {2, 7},      // 3: connected to 2,7
    {0, 5},      // 4: connected to 0,5
    {4, 6},      // 5: connected to 4,6
    {5, 7},      // 6: connected to 5,7
    {1, 3, 6},   // 7: connected to 1,3,6
}
start := 0
end := 7
	path := bfs(graph, start, end)
	paths = append(paths, path)
	for path != nil {
		graph = ignoreUsedStations(graph, path)
		path = bfs(graph, start, end)
		paths = append(paths, path)
	}
	fmt.Printf("Shortest path from node %d to %d: %d\n", start, end, path)
	fmt.Println("Paths from shortes to longest:", paths)
}

func ignoreUsedStations(graph [][]int, path []int) [][]int {
	for i := 1; i < len(path)-1; i++ {
		station := path[i]
		graph[station] = nil
	}
	return graph
}
