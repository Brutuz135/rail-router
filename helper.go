package main

import (
	"os"
	"regexp"
	"slices"
	"strings"

	//"golang.org/x/telemetry/counter"
)

func parseMap(mapFile string) ([][]int, (map[int]string), (map[string]int)) {
	var counter int
	IdToStations := make(map[int]string)
	StationsToId := make(map[string]int)

	data, _ := os.ReadFile(mapFile)
	mapInfo := string(data)
	

	idx := strings.Index(mapInfo, "connections:")
	stationsStr := mapInfo[:idx]
	connectionsStr := mapInfo[idx + 12:]

	re := regexp.MustCompile(`([A-Za-z0-9_]+)\s*,\s*(\d+)\s*,\s*(\d+)`)
	matchStationInfo := re.FindAllStringSubmatch(stationsStr, -1)
	
	for i, m := range matchStationInfo {
		station := m[1]
		counter += 1
		// x := m[2]
		// y := m[3]
		IdToStations[i] = station
		StationsToId[station] = i
	}
	stationGraph := make([][]int, counter)

	re1 := regexp.MustCompile(`([A-Za-z0-9_]+)\s*-\s*([A-Za-z0-9_]+)`)
	matchSationConnections := re1.FindAllStringSubmatch(connectionsStr, -1)

	for _, m := range matchSationConnections {
		station1 := StationsToId[m[1]]
		station2 := StationsToId[m[2]]
		

		if stationGraph[station1] == nil {
			stationGraph[station1] = append(stationGraph[station1], station2)

		} else if !slices.Contains(stationGraph[station1], station2) {
			stationGraph[station1] = append(stationGraph[station1], station2)
		}

		if stationGraph[station2] == nil {
			stationGraph[station2] = append(stationGraph[station2], station1)

		} else if !slices.Contains(stationGraph[station2], station1) {
			stationGraph[station2] = append(stationGraph[station2], station1)
		}
	}
	

	
	return stationGraph, IdToStations, StationsToId
}
