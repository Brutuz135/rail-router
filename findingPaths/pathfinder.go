package paths

import "math"

//Finds vertex-disjoint paths that minimize total turns 
// by using greedy BFS with edge-blocking retry.
func PathFinder(graph [][]int, start, end, trains int) [][]int {
	// Find initial set of paths using greedy approach
	basePaths := greedyDisjoint(copyGraph(graph), start, end)
	if len(basePaths) == 0 {
		return nil
	}

	bestPath := basePaths
	bestTurns := estimateTurns(bestPath, trains)

	// Try alternatives by blocking each edge of the first path
	firstPath := basePaths[0] // the first greedy path found
	for i := 0; i < len(firstPath)-1; i++ {
		fromStation := firstPath[i]
		toStation := firstPath[i+1]

		// Create modified graph without this edge
		graphCopy := copyGraph(graph)
		removeUndirectedEdge(graphCopy, fromStation, toStation)

		// Find paths without using the blocked edge
		candidatePaths := greedyDisjoint(graphCopy, start, end)
		if len(candidatePaths) == 0 {
			continue
		}

		// Keep this path set if it completes faster
		turns := estimateTurns(candidatePaths, trains)
		if turns < bestTurns {
			bestTurns = turns
			bestPath = candidatePaths
		}
	}

	return bestPath
}

func greedyDisjoint(graph [][]int, start, end int) [][]int {
	var paths [][]int

	// Keep finding shortest paths until no more exist
	for {
		shortestPath := bfs(graph, start, end)
		if shortestPath == nil {
			break
		}

		// Save the found path
		paths = append(paths, shortestPath)

		// Block intermediate stations (not start/end) 
		for i := 1; i < len(shortestPath)-1; i++ {
			graph[shortestPath[i]] = nil
		}

		// If the path is only a direct edge (start→end), remove the edge
		if len(shortestPath) == 2 {
			removeUndirectedEdge(graph, shortestPath[0], shortestPath[1])
		}
	}
	return paths
}

// Predicts how many turns are needed to move all trains using the given paths.
func estimateTurns(paths [][]int, trains int) int {
	if len(paths) == 0 {
		return math.MaxInt
	}

	pathLengths := make([]int, len(paths))
	for i, path := range paths {
		pathLengths[i] = len(path) - 1 // edges
	}

	// Track how many trains assigned to each path
	assigned := make([]int, len(paths))
	maxCompletionTime := 0

	// Assign each train to the path with minimum cost
	for trainNumber := 0; trainNumber < trains; trainNumber++ {
		bestP := 0
		bestCost := pathLengths[0] + assigned[0]

		// Find path with lowest cost (length + current load)
		for path := 1; path < len(paths); path++ {
			cost := pathLengths[path] + assigned[path]
			if cost < bestCost {
				bestCost = cost
				bestP = path
			}
		}

		// Assign train to best path
		assigned[bestP]++

		// Update maximum completion time
		completionTime := pathLengths[bestP] + assigned[bestP] - 1 //Formula: (pathLength + trainsOnPath - 1)
		if completionTime > maxCompletionTime {
			maxCompletionTime = completionTime
		}
	}

	return maxCompletionTime
}

// Creates a deep copy of the adjacency list graph.
func copyGraph(graph [][]int) [][]int {
	graphCopy := make([][]int, len(graph))
	for i := range graph {
		if graph[i] == nil {
			continue
		}
		graphCopy[i] = make([]int, len(graph[i]))
		copy(graphCopy[i], graph[i])
	}
	return graphCopy
}

// removeUndirectedEdge removes the connection between two stations
// in both directions (since the graph is undirected).
func removeUndirectedEdge(graph [][]int, stationA, stationB int) {
	removeEdge(graph, stationA, stationB)
	removeEdge(graph, stationB, stationA)
}

// removeEdge removes a directed edge from 'from' to 'to' in the graph.
func removeEdge(graph [][]int, from, to int) {
	adjacencyList := graph[from]
	if adjacencyList == nil {
		return
	}

	// Filter out the 'to' station from adjacency list
	filetered := adjacencyList[:0]
	for _, neighbor := range adjacencyList {
		if neighbor != to {
			filetered = append(filetered, neighbor)
		}
	}

	// Store the filtered list (reallocate to avoid aliasing)
	graph[from] = append([]int(nil), filetered...)
}

// Standard BFS shortest path (unweighted)
func bfs(graph [][]int, start, end int) []int {
	numNodes := len(graph)
	visited := make([]bool, numNodes)
	parent := make([]int, numNodes)

	for i := range parent {
		parent[i] = -1
	}

	// Queue for BFS (stores nodes to visit)
	queue := make([]int, 0, numNodes)
	queue = append(queue, start)
	visited[start] = true

	pathFound := false
	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:] // dequeue

		// Found the destination
		if current == end {
			pathFound = true
			break
		}

		// Explore all unvisited neighbors
		for _, nb := range graph[current] {
			if !visited[nb] {
				visited[nb] = true
				parent[nb] = current
				queue = append(queue, nb)
			}
		}
	}

	if !pathFound {
		return nil
	}

	// Reconstruct path from end to start using parent pointers
	path := []int{}
	for v := end; v != -1; v = parent[v] {
		path = append([]int{v}, path...)
	}
	return path
}
