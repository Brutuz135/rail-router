package main

import (
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"testing"
)

// ────────────────────────────────────────────────────────────────
//  Helpers
// ────────────────────────────────────────────────────────────────

// runCLI executes `go run .` with the given args and returns stderr as string.
func runCLI(args ...string) (string, error) {
	cmd := exec.Command("go", append([]string{"run", "."}, args...)...)
	out, err := cmd.CombinedOutput()
	return string(out), err
}

// assertContains fails the test if substr is not found in s.
func assertContains(t *testing.T, s, substr string) {
	t.Helper()
	if !strings.Contains(s, substr) {
		t.Fatalf("expected output to contain %q, got: %q", substr, s)
	}
}

// assertError ensures the CLI exited with an error.
func assertError(t *testing.T, err error, stderr string) {
	t.Helper()
	if err == nil {
		t.Fatalf("expected CLI to exit with error, got success. stderr: %q", stderr)
	}
	if !strings.Contains(stderr, "Error") {
		t.Fatalf("expected stderr to contain 'Error', got: %q", stderr)
	}
}

// ────────────────────────────────────────────────────────────────
//  Algorithm tests (FindTrainTurns)
// ────────────────────────────────────────────────────────────────

func TestLondon_MultipleRoutes_2Trains(t *testing.T) {
	// Test multiple trains on a map with multiple routes
	turns, err := FindTrainTurns("maps/london.map", "waterloo", "st_pancras", 2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(turns) == 0 || len(turns[0]) < 2 {
		t.Fatalf("expected 2 moves in first turn, got %d", len(turns[0]))
	}

	firstTargets := []string{turns[0][0].Target, turns[0][1].Target}
	if firstTargets[0] == firstTargets[1] {
		t.Errorf("expected trains to take different first steps, got %v", firstTargets)
	}
}

func TestLondon_MultipleRoutes_3Trains(t *testing.T) {
	// Test multiple trains on a map with several routes
	turns, err := FindTrainTurns("maps/london.map", "waterloo", "st_pancras", 3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(turns) == 0 {
		t.Fatalf("expected at least one turn")
	}

	unique := map[string]bool{}
	for _, move := range turns[0] {
		unique[move.Target] = true
	}

	if len(unique) < 2 {
		t.Errorf("expected at least two different first steps, got %v", unique)
	}
}

func TestLondon_MultipleRoutes_4Trains(t *testing.T) {
	// Test with 4 trains to ensure multiple routes are used
	turns, err := FindTrainTurns("maps/london.map", "waterloo", "st_pancras", 4)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(turns) == 0 || len(turns[0]) < 2 {
		t.Fatalf("expected at least 2 moves in first turn")
	}

	unique := map[string]bool{}
	for _, move := range turns[0] {
		unique[move.Target] = true
	}

	if len(unique) < 2 {
		t.Errorf("expected at least two unique first steps, got %v", unique)
	}
}

func TestLondon_MultipleRoutes_100Trains(t *testing.T) {
	// Test many trains to ensure pathfinding scales
	turns, err := FindTrainTurns("maps/london.map", "waterloo", "st_pancras", 100)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(turns) == 0 {
		t.Fatalf("expected at least one turn")
	}

	unique := map[string]bool{}
	for _, move := range turns[0] {
		unique[move.Target] = true
	}

	if len(unique) < 2 {
		t.Errorf("expected at least two unique first steps, got %v", unique)
	}
}

func TestLondon_SingleRoute_1Train(t *testing.T) {
	// Test a single train on a map with only one route
	turns, err := FindTrainTurns("maps/london.map", "waterloo", "st_pancras", 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(turns) == 0 || len(turns[0]) != 1 {
		t.Fatalf("expected exactly 1 move in first turn, got %d", len(turns[0]))
	}

	for i, turn := range turns {
		if len(turn) > 1 {
			t.Errorf("turn %d: expected only one train to move, got %d", i+1, len(turn))
		}
	}
}

// Verify move formatting is correct: Tn-station
func TestLondon_TrainMovementFormat(t *testing.T) {
	turns, err := FindTrainTurns("maps/london.map", "waterloo", "st_pancras", 3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	format := regexp.MustCompile(`^T\d+-[a-z0-9_]+$`)

	// Check each move in each turn
	for ti, turn := range turns {
		for mi, move := range turn {
			s := "T" + strconv.Itoa(move.Train) + "-" + move.Target
			if !format.MatchString(s) {
				t.Errorf("turn %d move %d: got %q, want Tn-station", ti+1, mi+1, s)
			}
		}
	}
}

// ────────────────────────────────────────────────────────────────
//
//	CLI error‑handling tests
//
// ────────────────────────────────────────────────────────────────
func TestCLI_TooFewArgs(t *testing.T) {
	//providing fewer arguments
	stderr, err := runCLI("maps/london.map", "waterloo", "st_pancras")
	assertError(t, err, stderr)
}

func TestCLI_TooManyArgs(t *testing.T) {
	//providing extra arguments
	stderr, err := runCLI("maps/london.map", "waterloo", "st_pancras", "2", "extra")
	assertError(t, err, stderr)
}

func TestCLI_NonexistentStart(t *testing.T) {
	//using a nonexistent start station
	stderr, _ := runCLI("maps/london.map", "notarealstation", "st_pancras", "2")
	assertContains(t, stderr, "does not exist")
}

func TestCLI_NonexistentEnd(t *testing.T) {
	//using a nonexistent end station
	stderr, _ := runCLI("maps/london.map", "waterloo", "notarealstation", "2")
	assertContains(t, stderr, "does not exist")
}

func TestCLI_SameStartEnd(t *testing.T) {
	//using same station for start and end
	stderr, _ := runCLI("maps/london.map", "waterloo", "waterloo", "2")
	assertContains(t, stderr, "start and end stations are the same")
}

func TestCLI_NoPath(t *testing.T) {
	//using a map file with no path between start and end
	stderr, _ := runCLI("maps/nonexistentpath.map", "waterloo", "st_pancras", "2")
	assertContains(t, stderr, "no path between the start and end stations")
}

func TestCLI_InvalidNumTrains(t *testing.T) {
	//using an invalid number of trains
	stderr, _ := runCLI("maps/london.map", "waterloo", "st_pancras", "-1")
	assertContains(t, stderr, "number of trains is not a valid positive integer")
}

func TestCLI_InvalidCoordinates(t *testing.T) {
	//using a map file with invalid coordinates
	stderr, _ := runCLI("maps/invalid_coordinates.map", "waterloo", "st_pancras", "2")
	assertContains(t, stderr, "invalid coordinates for station")
}

func TestCLI_DuplicateCoordinates(t *testing.T) {
	//using a map file with duplicate coordinates
	stderr, _ := runCLI("maps/duplicate_coordinates.map", "waterloo", "st_pancras", "2")
	assertContains(t, stderr, "share the same coordinates")
}

func TestCLI_ConnectionToUnknownStation(t *testing.T) {
	//using a map file with a connection to a nonexistent station
	stderr, _ := runCLI("maps/nonexistent_station.map", "waterloo", "st_pancras", "2")
	assertContains(t, stderr, "connection references unknown station")
}

func TestCLI_DuplicateStationNames(t *testing.T) {
	//using a map file with duplicate station names
	stderr, _ := runCLI("maps/duplicate_stations.map", "waterloo", "st_pancras", "2")
	assertContains(t, stderr, "duplicate station name")
}

func TestCLI_InvalidStationNames(t *testing.T) {
	//using a map file with invalid station names
	stderr, _ := runCLI("maps/invalid_station_name.map", "waterloo", "st_pancras", "2")
	assertContains(t, stderr, "invalid station name")
}

func TestCLI_MissingStationsSection(t *testing.T) {
	//using a map file missing stations section
	stderr, _ := runCLI("maps/station_missing_section.map", "waterloo", "st_pancras", "2")
	assertContains(t, stderr, "missing stations section")
}

func TestCLI_MissingConnectionsSection(t *testing.T) {
	// Using a map file missing connections section
	stderr, _ := runCLI("maps/connections_missing_section.map", "waterloo", "st_pancras", "2")
	assertContains(t, stderr, "missing connections section")
}

func TestCLI_TooManyStations(t *testing.T) {
	// Using a large map to trigger the station limit error
	stderr, _ := runCLI("maps/10k+_parallel.map", "waterloo", "st_pancras", "4")
	assertContains(t, stderr, "map contains more than 10000 stations")
}
