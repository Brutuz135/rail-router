package paths

import (
)

	
func PathFinder(graph [][]int, start int, end int) [][]int {
	var paths [][]int


//just and old graph for testing, but this is how the connection graph looks like	
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

	// Find first shortest path from start to end.
	path := bfs(graph, start, end)

	// As long as a path is found:
	for path != nil {
		// Store this path.
		paths = append(paths, path)

		// Remove intermediate stations of this path from the graph
        	// (start and end are kept), so future paths cannot reuse them.
		graph = ignoreUsedStations(graph, path)

		// Find next shortest path in the modified graph.
		path = bfs(graph, start, end)
	}

	return paths
}

func ignoreUsedStations(graph [][]int, path []int) [][]int {
	// For a found path, remove all internal stations (indices 1..len-2)
    	// from the graph by setting their adjacency lists to nil.
    	// This prevents using these stations in later paths.
	for i := 1; i < len(path)-1; i++ {
		station := path[i]
		graph[station] = nil
	}
	return graph
}



func bfs(graph [][]int, start int, end int) []int {

	nNodes := len(graph)
	
	
    	// BFS data structures:
    	// pq: queue of nodes to visit (FIFO).
    	// visited: marks whether a node has been seen.
    	// parrent: stores predecessor of each node in the BFS tree to reconstruct path.
	var pq []int
	var visited []bool
	var parrent []int
	var current int
	pathFound := false
	var path []int

	// Initialize visited and parent slices for all nodes.
	for range nNodes {
		visited = append(visited, false)
		parrent = append(parrent, -1)
	}

	// Start BFS from the start node.	
	pq = append(pq, start)
	visited[start] = true
	parrent[start] = -1
	current = pq[0]

	// Standard BFS loop while queue is not empty.
	for len(pq) > 0 {

		// Take the node at the front of the queue.
		current = pq[0]

		// Explore all neighbours of the current node.
		for _, neighbour := range graph[current] {

			// If neighbour not visited yet, mark and enqueue it,
            		// and record current as its parent.
			if visited[neighbour] == false{
				visited[neighbour] = true
				parrent[neighbour] = current
				pq = append(pq, neighbour)
			}
		}

		// Pop the front element from the queue (move window by one).
		pq = append(pq[:0], pq[1:]...)//take the current out of the priorityQueue

		// If end node is reached, stop BFS.
		if current == end {
			pathFound = true
			break
		}
	
	}


	// If no path was found, return nil.
	if pathFound == false{
		return nil
	}
	
	// Reconstruct path from end back to start using the parent array.
	cur := end
	path = append(path, end)
	for cur != start {
		cur = parrent[cur]
		// Prepend the parent node to the path slice.
		path = append([]int{cur}, path...)
	}
	return path

}




