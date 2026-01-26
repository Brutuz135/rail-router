package helpers

import (
	"fmt"
	"os"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

// ParseMap reads a map file and returns:
// 1. stationGraph: adjacency list of the graph
// 2. idToStations: mapping from int ID → station name
// 3. stationsToId: mapping from station name → int ID
func ParseMap(mapFile string) ([][]int, map[int]string, map[string]int, error) {
	data, err := os.ReadFile(mapFile)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to read map file: %w", err)
	}

	content := string(data)

	// Split stations and connections
	idx := strings.Index(content, "connections:")
	if idx == -1 {
		return nil, nil, nil, fmt.Errorf("invalid map format: missing connections section")
	}

	stationsStr := content[:idx]
	connectionsStr := content[idx+len("connections:"):]

	idToStations := make(map[int]string)
	stationsToId := make(map[string]int)
	coordMap := make(map[[2]int]string) // detect duplicate coordinates

	// Parse stations info: name,x,y STRICTLY line-by-line:
	reStationLine := regexp.MustCompile(`^([A-Za-z0-9_]+)\s*,\s*(\d+)\s*,\s*(\d+)\s*(?:#.*)?$`)

	for _, line := range strings.Split(stationsStr, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") || line == "stations:" {
			continue
		}

		// Detect invalid station names
		m := reStationLine.FindStringSubmatch(line)
		if m == nil {
			return nil, nil, nil, fmt.Errorf("invalid station name")
		}

		station := m[1]

		// Detect duplicate station names
		if _, exists := stationsToId[station]; exists {
			return nil, nil, nil, fmt.Errorf("duplicate station name %q", station)
		}

		// Parse and validate coordinates (positive integers only)
		x, err1 := strconv.Atoi(m[2])
		y, err2 := strconv.Atoi(m[3])
		if err1 != nil || err2 != nil || x < 0 || y < 0 {
			return nil, nil, nil, fmt.Errorf("invalid coordinates for station %q", station)
		}

		// Detect  duplicate coordinates
		coord := [2]int{x, y}
		if other, exists := coordMap[coord]; exists {
			return nil, nil, nil, fmt.Errorf("stations %q and %q share the same coordinates", station, other)
		}
		coordMap[coord] = station

		// Assign IDs in encounter order
		id := len(idToStations)
		idToStations[id] = station
		stationsToId[station] = id
	}

	// Create adjacency list
	stationGraph := make([][]int, len(idToStations))

	// Parse connections: A-B line-by-line: skip empty lines;
	reConnLine := regexp.MustCompile(`^([A-Za-z0-9_]+)\s*-\s*([A-Za-z0-9_]+)\s*(?:#.*)?$`)

	seenEdges := make(map[[2]int]bool) // detect duplicates (incl reverse)

	for _, line := range strings.Split(connectionsStr, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		m := reConnLine.FindStringSubmatch(line)
		if m == nil {
			return nil, nil, nil, fmt.Errorf("invalid connection")
		}

		fromName := m[1]
		toName := m[2]

		a, okA := stationsToId[fromName]
		b, okB := stationsToId[toName]

		// Connection references non-existent station
		if !okA {
			return nil, nil, nil, fmt.Errorf("connection references unknown station %q", fromName)
		}
		if !okB {
			return nil, nil, nil, fmt.Errorf("connection references unknown station %q", toName)
		}

		// Duplicate connections (including reverse)
		key := [2]int{a, b}
		if a > b {
			key = [2]int{b, a} // normalize
		}
		if seenEdges[key] {
			return nil, nil, nil, fmt.Errorf("duplicate connection between %q and %q", fromName, toName)
		}
		seenEdges[key] = true

		// Add to adjacency list (undirected)
		if !slices.Contains(stationGraph[a], b) {
			stationGraph[a] = append(stationGraph[a], b)
		}
		if !slices.Contains(stationGraph[b], a) {
			stationGraph[b] = append(stationGraph[b], a)
		}
	}

	return stationGraph, idToStations, stationsToId, nil
}
