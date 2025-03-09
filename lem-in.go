package main

import (
	"fmt"
	"os"
	"strings"

	"lem.in/functions"
	"lem.in/graph"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("ERROR: Incorrect number of arguments!")
		return
	}

	fileName := "./examples/" + os.Args[1]
	input, err := os.ReadFile(fileName)
	if err != nil {
		fmt.Println("ERROR: Unable to read file:", err)
		return
	}

	lines := strings.Split(strings.TrimSpace(string(input)), "\n")
	var inputData []string
	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)
		if trimmedLine != "" {
			inputData = append(inputData, trimmedLine)
		}
	}

	numAnts, roomData := functions.GetNumberOfAnts(inputData)
	roomCoordinates, roomConnections := functions.SeparateRoomData(roomData)
	coordinateMap := functions.CreateRoomCoordinateMap(roomCoordinates)
	initialMap := functions.CreateRoomConnectionsMap(roomConnections, coordinateMap)
	finalMap := functions.RemoveDeadEnds(initialMap)

	path, distribution := graph.FindPaths(finalMap, numAnts)

	fmt.Println(string(input) + "\n")

	ants := make([][]int, len(distribution))
	antNumber := 1
	for i := range distribution {
		ants[i] = make([]int, distribution[i])
		for j := range ants[i] {
			ants[i][j] = antNumber
			antNumber++
		}
	}

	maxAnts := 0
	for _, antGroup := range ants {
		if len(antGroup) > maxAnts {
			maxAnts = len(antGroup)
		}
	}

	remainingAnts := numAnts
	for i := 1; remainingAnts > 0; i++ {
		connected := false
		for k := 0; k < maxAnts; k++ {
			for j := range ants {
				if len(ants[j]) > k && len(path[j]) > i-k && i-k > 0 {
					fmt.Printf("L%d-%s ", ants[j][k], path[j][i-k])
					if i-k == len(path[j])-1 {
						remainingAnts--
					}
					if len(path[j]) == 2 {
						connected = true
					}
				}
			}
		}
		if !connected {
			fmt.Println()
		}
	}
}
