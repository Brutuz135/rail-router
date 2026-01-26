package main

import (
	"bytes"
	"os"
	"strings"
	"testing"

	paths "stations/findingPaths"
	helpers "stations/helpers"
	moves "stations/movingTrains"
)

// TestCase defines a single scenario to test
type TestCase struct {
	Name           string
	MapContent     string
	StartStation   string
	EndStation     string
	NumTrains      int
	MaxTurns       int    // Optional: if > 0, asserts turns <= MaxTurns
	MinRoutes      int    // Optional: if > 0, asserts found paths >= MinRoutes
	ExpectedOutput string // Optional: if non-empty, checks if output contains specific sequences
}

func TestTrainSystem(t *testing.T) {
	tests := []TestCase{
		{
			Name: "London 2 Trains (Multiple Routes)",
			MapContent: `stations:
waterloo,3,1
victoria,6,7
euston,11,23
st_pancras,5,15
connections:
waterloo-victoria
waterloo-euston
st_pancras-euston
victoria-st_pancras`,
			StartStation: "waterloo",
			EndStation:   "st_pancras",
			NumTrains:    2,
			MinRoutes:    2, // Expects > 1 route
		},
		{
			Name: "London 3 Trains (Multiple Routes)",
			MapContent: `stations:
waterloo,3,1
victoria,6,7
euston,11,23
st_pancras,5,15
connections:
waterloo-victoria
waterloo-euston
st_pancras-euston
victoria-st_pancras`,
			StartStation: "waterloo",
			EndStation:   "st_pancras",
			NumTrains:    3,
			MinRoutes:    2,
		},
		{
			Name: "London 100 Trains (Check Scalability)",
			MapContent: `stations:
waterloo,3,1
victoria,6,7
euston,11,23
st_pancras,5,15
connections:
waterloo-victoria
waterloo-euston
st_pancras-euston
victoria-st_pancras`,
			StartStation: "waterloo",
			EndStation:   "st_pancras",
			NumTrains:    100,
			MinRoutes:    2,
		},
		{
			Name: "Bond Square to Space Port (4 Trains, Max 6 Turns)",
			MapContent: `stations:
bond_square,20,6
apple_avenue,7,7
orange_junction,6,1
space_port,1,11
connections:
bond_square-apple_avenue
apple_avenue-orange_junction
orange_junction-space_port`,
			StartStation: "bond_square",
			EndStation:   "space_port",
			NumTrains:    4,
			MaxTurns:     6,
		},
		{
			Name: "Jungle to Desert (10 Trains, Max 8 Turns)",
			MapContent: `stations:
jungle,5,16
green_belt,6,1
village,5,7
mountain,9,16
treetop,0,4
grasslands,15,13
suburbs,4,9
clouds,0,0
wetlands,2,12
farms,11,10
downtown,4,4
metropolis,3,20
industrial,1,18
desert,9,0
connections:
jungle-grasslands
mountain-treetop
clouds-wetlands
downtown-metropolis
green_belt-village
suburbs-clouds
industrial-desert
jungle-farms
village-mountain
wetlands-desert
grasslands-suburbs
jungle-green_belt
farms-downtown
treetop-desert
metropolis-industrial
mountain-wetlands
farms-mountain`,
			StartStation: "jungle",
			EndStation:   "desert",
			NumTrains:    10,
			MaxTurns:     8,
		},
		{
			Name: "Beginning to Terminus (20 Trains, Max 11 Turns)",
			MapContent: `stations:
beginning,0,0
near,1,0
far,1,3
terminus,0,3
connections:
beginning-near
beginning-terminus
near-far
terminus-far`,
			StartStation: "beginning",
			EndStation:   "terminus",
			NumTrains:    20,
			MaxTurns:     11,
		},
		{
			Name: "Two to Four (4 Trains, Max 6 Turns)",
			MapContent: `stations:
one,1,1
two,2,2
three,3,3
four,4,4
five,5,5
six,6,6
connections:
two-three
five-one
three-one
two-five
one-four
six-two
one-six`,
			StartStation: "two",
			EndStation:   "four",
			NumTrains:    4,
			MaxTurns:     6,
		},
		{
			Name: "Small to Large (9 Trains, Max 8 Turns)",
			MapContent: `stations:
small,4,0
large,4,6
00,0,0
01,0,1
02,0,2
03,0,3
04,0,4
05,0,5
10,1,0
11,1,1
12,1,2
13,1,3
14,1,4
15,1,5
20,2,0
21,2,1
22,2,2
23,2,3
24,2,4
25,2,5
30,3,0
31,3,1
32,3,2
33,3,3
34,3,4
35,3,5
36,3,6
connections:
24-25
24-23
23-12
small-32
32-33
33-34
34-35
35-36
36-22
small-10
10-11
10-20
11-12
11-14
12-large
12-03
small-13
13-14
14-15
small-00
00-01
01-02
02-03
03-04
20-21
20-25
21-15
21-22
21-30
22-large
25-30
30-31
31-large
04-05
05-large`,
			StartStation: "small",
			EndStation:   "large",
			NumTrains:    9,
			MaxTurns:     8,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			//  Create Temp File
			tmpFile, err := os.CreateTemp("", "test_map_*.txt")
			if err != nil {
				t.Fatalf("Failed to create temp file: %v", err)
			}
			defer os.Remove(tmpFile.Name())

			if _, err := tmpFile.WriteString(tc.MapContent); err != nil {
				t.Fatalf("Failed to write map content: %v", err)
			}
			tmpFile.Close()

			//  Parse Map
			graph, idToStations, stationsToId, err := helpers.ParseMap(tmpFile.Name())
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
			allPaths := paths.PathFinder(graph, startID, endID)

			// ASSERT: Minimum Routes
			if tc.MinRoutes > 0 {
				if len(allPaths) < tc.MinRoutes {
					t.Errorf("Expected at least %d routes, got %d", tc.MinRoutes, len(allPaths))
				}
			}

			// If we found NO paths, we can't test movement
			if len(allPaths) == 0 {
				if tc.MinRoutes > 0 {
					t.FailNow()
				}
				return // If we didn't expect paths, we are done
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
			lines := strings.Split(output, "\n")

			// Filter empty lines just in case
			var turnLines []string
			for _, line := range lines {
				if strings.TrimSpace(line) != "" {
					turnLines = append(turnLines, line)
				}
			}
			turnsTaken := len(turnLines)

			// ASSERT: Max Turns
			if tc.MaxTurns > 0 {
				if turnsTaken > tc.MaxTurns {
					t.Errorf("Too slow! Expected max %d turns, took %d turns.\nOutput sample:\n%s", tc.MaxTurns, turnsTaken, output)
				} else {
					t.Logf("Success: Completed in %d turns (Limit: %d)", turnsTaken, tc.MaxTurns)
				}
			}
		})
	}
}
