package paths

import (
	"fmt"

)

	
func PathFinder(graph [][]int, start int, end int) [][]int {
	var paths [][]int

	// Find first shortest path from start to end.
	path := bfs(graph, start, end)

	// As long as a path is found:
	for path != nil {
		// Store this path.
		paths = append(paths, path)

		// Remove intermediate stations of this path from the graph
		// Remove also possible connection between start and end
        	// (start and end are kept), so future paths cannot reuse them.
		graph = ignoreUsedStations(graph, path)
		fmt.Println(graph)


		// Find next shortest path in the modified graph.
		path = bfs(graph, start, end)
	}


	return paths
}

func ignoreUsedStations(graph [][]int, path []int) [][]int {
	PathLen := len(path)
	//if path is only start and end, remove the connection between them
	if PathLen == 2 {
		start := path[0]
		end := path[1]
		for i, connection := range graph[start] {
			if connection == end {
				row := graph[start]
				row = append(row[:i], row[i+1:]...)
				graph[start] = row
			}
		}	

		for i, connection := range graph[end] {
			if connection == start {
				row := graph[end]
				row = append(row[:i], row[i+1:]...)
				graph[end] = row
			}
		}
		fmt.Println(path)
		fmt.Println(PathLen)
	}

	// For a found path, remove all internal stations (indices 1..len-2)
    	// from the graph by setting their adjacency lists to nil.
    	// This prevents using these stations in later paths.
	for i := 1; i < PathLen-1; i++ {
		station := path[i]
		graph[station] = nil
	}
	return graph
}



func bfs(graph [][]int, start int, end int) []int {

	nNodes := len(graph)
	
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




