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
go run . maps/network.map waterloo st_pancras 4
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
├── paths/
│   └── path_finder.go      # Path finding logic
├── moves/
│   └── move_trains.go      # Train movement simulation
├── maps/
│   └── network.map        # Example network map
├── main_test.go            # High-level integration tests
├── train_system_test.go   # Train movement validation tests
```

---

## Testing

Run all tests with:

```bash
go test -v
```

Tests cover:

- Valid and invalid map files
- Different train counts
- Single and multiple routes
- Scalability (large number of trains)
- Correct movement formatting and turn limits

---

## Acknowledgments

Built as part of the **Kood/Sisu**, focusing on algorithms, parsing, and simulation.
