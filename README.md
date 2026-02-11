# Stations

A command-line Go program that simulates trains moving through a network of stations. Given a network map, a start station, an end station, and a number of trains, the program calculates valid routes and prints turn-by-turn train movements while respecting path constraints.

---

## Features

- Parses a custom **network map file** describing stations and connections
- Validates input arguments and map format with clear error messages
- Finds valid paths between two stations
- Simulates multiple trains moving concurrently
- Prints movements turn-by-turn in the required format
- Handles edge cases (invalid stations, no paths, invalid train count)

---

## Installation

```bash
# Clone the repository and build

git clone <https://gitea.kood.tech/marynjoroge/pathfinder >
cd pathfinder
go build -o pathfinder .
```

Requirements:

- Go 1.20+

---

## Usage

```bash
go run . [path to network map] [start station] [end station] [number of trains]
```

### Example

```bash
go run . maps/london.map waterloo st_pancras 2
```

Output:

```
T1-victoria T2-euston
T1-st_pancras T2-st_pancras
```

Each line represents **one turn**. Each entry follows the format:

```
T<train_number>-<station_name>
```

---

## Network Map Format

The map file must follow this structure:

```
stations:
waterloo, 3 , 1
victoria,6,7
euston,11,23
st_pancras,5,15

connections:
waterloo -victoria
waterloo- euston
st_pancras-euston
victoria-st_pancras
```

### Rules

- Station names must be **lowercase**
- Sections `stations:` and `connections:` are mandatory
- Connections are **bidirectional**
- Duplicate stations or invalid connections will result in errors

---

## Error Handling

The program fails fast and reports **one clear error at a time**.

Examples:

```bash
# Missing arguments
Error: expected 4 arguments, got 3
Usage: go run . [path to file containing network map] [start station] [end station] [number of trains]

# Invalid train count
Error: number of trains is not a valid positive integer

# Missing stations section in map
Error parsing map: invalid map format: missing stations section

# Invalid station name
Error: start station "Waterloo" does not exist
```
---

## Project Structure

```
stations/
├── main.go                 # Program entry point
├── go.mod
├── README.md
├── helpers/
│   └── parse_map.go        # Map parsing & validation
├── findingPaths/
│   └── path_finder.go      # Path finding logic
├── movingTrains/
│   └── move_trains.go      # Train movement simulation
├── maps/
│   └── network.map        # Example network map
├── main_test.go            # High-level integration tests
```

---

## Testing

Run all tests with:

```bash
go test -v
```

To print debug output inside your code during tests, use:
```
t.Log("stdErr", stderr)

```
## Example Output
A successful test run looks like this:
```
=== RUN   TestLondon_MultipleRoutesFor2Trains
--- PASS: TestLondonNetwork_MultipleRoutesFor2Trains (0.01s)
=== RUN   TestLondonNetwork_MultipleRoutesFor3Trains
--- PASS: TestLondonNetwork_MultipleRoutesFor3Trains (0.01s)
=== RUN   TestLondon_SingleRouteFor1Train
--- PASS: TestLondon_SingleRouteFor1Train (0.00s)
=== RUN   TestCLI_InvalidCoordinates
--- PASS: TestCLI_InvalidCoordinates (0.00s)
=== RUN   TestCLI_MissingStationsSection
--- PASS: TestCLI_MissingStationsSection (0.00s)
PASS
ok      stations  0.123s

```

Tests cover:

- Pathfinding & Movement
- Map Parsing
- CLI Behavior

If any test fails, Go prints a clear message showing which test failed and why.

---

## Troubleshooting Tests
```
If a test fails, you can re-run only that specific test:

go test -v -run TestCLI_TooManyStations
```

## Bonus Features
### 1. Advanced Error Handling
The program provides clear, human‑readable error messages that help users understand exactly what went wrong.

Example:
```
Error: start station "water" does not exist

```
### 2. Precise, Context‑Aware Error Reporting
Errors identify the specific problematic entity or line, making debugging significantly easier.

Examples:

```
Error: parsing map: stations "euston" and "waterloo" share the same coordinates
```

### 3. Comprehensive Test Suite
A full suite of automated tests validates:

- Map parsing

- Error handling

- Path‑finding correctness

- Train movement rules

- Turn limits

- CLI behavior

## Acknowledgments

Built as part of the **Kood/Sisu**, focusing on algorithms, parsing, and simulation.
