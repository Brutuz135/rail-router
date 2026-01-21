package main

import (
	"bytes"
	"os"
	"strings"
	"testing"

	helpers "stations/helpers"
	paths "stations/findingPaths"
	moves "stations/movingTrains"
)

func TestLondonNetworkMultiplePaths(t *testing.T) {
	// 1. Prepare test map file
	mapContent := `
stations:
waterloo,3,1
victoria,6,7
euston,11,23
st_pancras,5,15

connections:
waterloo-victoria
waterloo-euston
st_pancras-euston
victoria-st_pancras
`
	tmpFile, err := os.CreateTemp("", "network*.map")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.WriteString(mapContent); err != nil {
		t.Fatal(err)
	}
	tmpFile.Close()

	// 2. Parse map
	graph, idToStations, stationsToId, err := helpers.ParseMap(tmpFile.Name())
	if err != nil {
		t.Fatalf("ParseMap failed: %v", err)
	}

	start := stationsToId["waterloo"]
	end := stationsToId["st_pancras"]

	// 3. Find paths
	allPaths := paths.PathFinder(graph, start, end)
	if len(allPaths) < 2 {
		t.Fatalf("expected at least 2 paths, got %d", len(allPaths))
	}

	// 4. Capture stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// 5. Move trains
	moves.MoveTrains(allPaths, idToStations, 4)

	// Restore stdout
	w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := strings.TrimSpace(buf.String())

	lines := strings.Split(output, "\n")

	// 6. Assertions
	if len(lines) != 3 {
		t.Fatalf("expected 3 turns, got %d\nOutput:\n%s", len(lines), output)
	}
}
