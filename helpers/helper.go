package helpers

import (
	"fmt"
	"os"
	"regexp"
	"slices"
	"strings"
)

// ParseMap reads a map file and returns:
// 1. stationGraph: adjacency list of the graph ([][]int)
// 2. idToStations: mapping from int ID → station name
// 3. stationsToId: mapping from station name → int ID
func ParseMap(mapFile string) ([][]int, map[int]string, map[string]int, error) {
	// Read file
	data, err := os.ReadFile(mapFile)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to read map file: %w", err)
	}

	content := string(data)

	// Split stations and connections sections
	idx := strings.Index(content, "connections:")
	if idx == -1 {
		return nil, nil, nil, fmt.Errorf("invalid map format: missing connections section")
	}

	stationsStr := content[:idx]
	connectionsStr := content[idx+len("connections:"):]

	// Regex to parse station info: name,x,y
	reStation := regexp.MustCompile(`([A-Za-z0-9_]+)\s*,\s*(\d+)\s*,\s*(\d+)`)
	matches := reStation.FindAllStringSubmatch(stationsStr, -1)

	idToStations := make(map[int]string)
	stationsToId := make(map[string]int)

	for i, match := range matches {
		station := match[1]
		idToStations := idToStations
		idToStations[i] = station
		stationsToId[station] = i
	}

	// Create adjacency list
	stationGraph := make([][]int, len(matches))

	// Regex to parse connections: A-B
	reConn := regexp.MustCompile(`([A-Za-z0-9_]+)\s*-\s*([A-Za-z0-9_]+)`)
	connMatches := reConn.FindAllStringSubmatch(connectionsStr, -1)

	for _, match := range connMatches {
		s1 := stationsToId[match[1]]
		s2 := stationsToId[match[2]]

		// Append connection if not already in slice
		if !slices.Contains(stationGraph[s1], s2) {
			stationGraph[s1] = append(stationGraph[s1], s2)
		}
		if !slices.Contains(stationGraph[s2], s1) {
			stationGraph[s2] = append(stationGraph[s2], s1)
		}
	}

	return stationGraph, idToStations, stationsToId, nil
}

