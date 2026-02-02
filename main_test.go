package main

import (
	"bytes"
	"os"
	"strings"
	"testing"
	"regexp"

	paths "stations/findingPaths"
	helpers "stations/helpers"
	moves "stations/movingTrains"
)

// TestCase defines a single scenario to test
type TestCase struct {
	Name           string
	MapFile        string
	StartStation   string
	EndStation     string
	NumTrains      int
	MaxTurns       int    // Optional: if > 0, asserts turns <= MaxTurns
	MinRoutes      int    // Optional: if > 0, asserts found paths >= MinRoutes
	MaxRoutes      int    
	ExpectedOutput string // Optional: if non-empty, checks if output contains specific sequences
}

func TestTrainSystem(t *testing.T) {
	tests := []TestCase{
		{
			Name: "London 2 Trains (Multiple Routes)",
			MapFile: "maps/london.map",
			StartStation: "waterloo",
			EndStation:   "st_pancras",
			NumTrains:    2,
			MinRoutes:    2, // Expects > 1 route
		},
		{
			Name: "London 3 Trains (Multiple Routes)",
			MapFile: "maps/london.map",	
			StartStation: "waterloo",
			EndStation:   "st_pancras",
			NumTrains:    3,
			MinRoutes:    2,
		},
		{
			Name: "London 4 Trains (Multiple Routes)",
			MapFile: "maps/london.map",	
			StartStation: "waterloo",
			EndStation:   "st_pancras",
			NumTrains:    4,
			MinRoutes:    2,
		},
		{
			Name: "London 100 Trains (Check Scalability)",
			MapFile: "maps/london.map",
			StartStation: "waterloo",
			EndStation:   "st_pancras",
			NumTrains:    100,
			MinRoutes:    2,
		},
		{
			Name: "London 1 Train (Only 1 Route)",
			MapFile: "maps/london.map",	
			StartStation: "waterloo",
			EndStation:   "st_pancras",
			NumTrains:    1,
			MinRoutes:    1,
			MaxRoutes:    1,
		},
		{
			Name: "Bond Square to Space Port (4 Trains, Max 6 Turns)",
			MapFile: "maps/bond.map",	
			StartStation: "bond_square",
			EndStation:   "space_port",
			NumTrains:    4,
			MaxTurns:     6,
		},
		{
			Name: "Jungle to Desert (10 Trains, Max 8 Turns)",
			MapFile: "maps/jungle.map",
			StartStation: "jungle",
			EndStation:   "desert",
			NumTrains:    10,
			MaxTurns:     8,
		},
		{
			Name: "Beginning to Terminus (20 Trains, Max 11 Turns)",
			MapFile: "maps/terminus.map",	
			StartStation: "beginning",
			EndStation:   "terminus",
			NumTrains:    20,
			MaxTurns:     11,
		},
		{
			Name: "Two to Four (4 Trains, Max 6 Turns)",
			MapFile: "maps/numbers.map",	
			StartStation: "two",
			EndStation:   "four",
			NumTrains:    4,
			MaxTurns:     6,
		},
		{
			Name: "Beethoven to Part (9 Trains, Max 6 Turns)",
			MapFile: "maps/clasical.map",	
			StartStation: "beethoven",
			EndStation:   "part",
			NumTrains:    9,
			MaxTurns:     6,
		},
		{
			Name: "Small to Large (9 Trains, Max 8 Turns)",
			MapFile: "maps/small.map",	
			StartStation: "small",
			EndStation:   "large",
			NumTrains:    9,
			MaxTurns:     8,
		},




	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {


			//  Parse Map
			graph, idToStations, stationsToId, err := helpers.ParseMap(tc.MapFile)
			if err != nil {
				t.Fatalf("ParseMap failed: %v", err)
			}

			//  Resolve IDs
			startID, okStart := stationsToId[tc.StartStation]
			endID, okEnd := stationsToId[tc.EndStation]
			if !okStart || !okEnd {
				t.Fatalf("Start or End station not found in map")
			}
			
			//  Find Paths
			allPaths := paths.PathFinder(graph, startID, endID, tc.NumTrains)

			// ASSERT: Minimum Routes
			if tc.MinRoutes > 0 {
				if len(allPaths) < tc.MinRoutes {
					t.Errorf("Expected at least %d routes, got %d", tc.MinRoutes, len(allPaths))
				}
			}

			//  Capture Output for Turns Calculation
			oldStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			// Execute Movement
			moves.MoveTrains(allPaths, idToStations, tc.NumTrains)

			w.Close()
			os.Stdout = oldStdout

			var buf bytes.Buffer
			buf.ReadFrom(r)
			output := strings.TrimSpace(buf.String())
			if testing.Verbose() {
				t.Log("\n" + output)
			}

			validateMovements(
				t,
				output,
				tc.StartStation,
				tc.EndStation,
				tc.NumTrains,
			)

			lines := strings.Split(output, "\n")

			// Filter empty lines just in case
			var turnLines []string
			for _, line := range lines {
				if strings.TrimSpace(line) != "" {
					turnLines = append(turnLines, line)
				}
			}
			turnsTaken := len(turnLines)


			//Format Check
			matched := regexp.MustCompile(`^T\d+-[a-z0-9_]+$`).MatchString
		
			for _, line := range turnLines {
    				for _, part := range strings.Fields(line) {
        				if !matched(part) {
           					t.Error("Output in wrong format")
            					break
        				}
    				}
			}

			//ASSERT: Maximum routes
			if tc.MaxRoutes > 0 && tc.MaxRoutes <= len(allPaths) {
				if turnsTaken > len(allPaths[tc.MaxRoutes-1]) - 1 {
					t.Errorf("Expected no more than %d routes, got %d", tc.MaxRoutes, len(allPaths))
				}
			}



			// ASSERT: Max Turns
			if tc.MaxTurns > 0 {
				if turnsTaken > tc.MaxTurns {
					t.Errorf("Too slow! Expected max %d turns, took %d turns.\nOutput sample:\n%s", tc.MaxTurns, turnsTaken, output)
				} else {
					t.Logf("Success: Completed in %d turns (Limit: %d)", turnsTaken, tc.MaxTurns)
				}
			}

			if tc.NumTrains == 1 && len(allPaths) != 1 {
				t.Fatalf("Expected exactly 1 route, got %d", len(allPaths))
			}
		})
	}
}
