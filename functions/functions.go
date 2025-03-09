package functions

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
)

// GetNumberOfAnts finds the number of ants in the room
// and returns the number of ants and the data of the room
func GetNumberOfAnts(data []string) (int, []string) {
	for i, line := range data {
		if line[0] != '#' {
			nbrAnts, err := strconv.Atoi(line)
			if err != nil || nbrAnts < 1 {
				CheckError(errors.New("ERROR: Invalid number of ants"))
			}
			return nbrAnts, data[i+1:]
		}
	}
	CheckError(errors.New("ERROR: No valid number of ants found"))
	return 0, nil
}

// SeparateRoomData separates room data into two categories: room information and relationships.
func SeparateRoomData(roomData []string) ([]string, []string) {
	var roomInfo, relationData []string

	for _, line := range roomData {
		if len(strings.Fields(line)) == 3 || strings.HasPrefix(line, "##") {
			roomInfo = append(roomInfo, line)
		} else if line[0] == '#' {
			continue
		} else if strings.Count(line, "-") == 1 {
			relationData = append(relationData, line)
		} else {
			CheckError(errors.New("ERROR: Invalid data format"))
		}
	}

	return roomInfo, relationData
}

// CreateRoomCoordinateMap creates a map of rooms and their coordinates
func CreateRoomCoordinateMap(lines []string) map[string][]int {
	roomMap := make(map[string][]int)
	var startPoint, endPoint int
	var startPointer, endPointer *int

	for i, line := range lines {
		switch line {
		case "##start":
			startPoint = i + 1
			startPointer = &startPoint
		case "##end":
			endPoint = i + 1
			endPointer = &endPoint
		}
		if startPointer != nil && endPointer != nil {
			break
		}
	}

	if startPointer == nil || endPointer == nil {
		CheckError(errors.New("ERROR: Missing start or end room"))
		return nil
	}

	for i := startPoint; i <= endPoint; i++ {
		if strings.HasPrefix(lines[i], "##") {
			continue
		}
		room := strings.Fields(lines[i])
		if len(room) == 3 {
			coorX, err := strconv.Atoi(room[1])
			CheckError(err)
			coorY, err := strconv.Atoi(room[2])
			CheckError(err)
			position := 1
			if i == startPoint {
				position = 0
			} else if i == endPoint {
				position = 2
			}
			roomMap[room[0]] = []int{position, coorX, coorY}
		}
	}

	for key1, v1 := range roomMap {
		for key2, v2 := range roomMap {
			if key1 != key2 && v1[1] == v2[1] && v1[2] == v2[2] {
				CheckError(errors.New(key1 + " and " + key2 + " have the same coordinates"))
			}
		}
	}

	return roomMap
}

// CreateRoomConnectionsMap creates a map of connections between rooms
func CreateRoomConnectionsMap(rawData []string, coordinateMap map[string][]int) map[string][]string {
	connectionsMap := make(map[string][]string)

	for key, v := range coordinateMap {
		switch v[0] {
		case 0:
			connectionsMap[key] = append(connectionsMap[key], "start")
		case 2:
			connectionsMap[key] = append(connectionsMap[key], "end")
		default:
			connectionsMap[key] = append(connectionsMap[key], "middle")
		}
	}

	for _, v := range rawData {
		connection := strings.Split(v, "-")
		if len(connection) == 2 {
			if _, ok := connectionsMap[connection[0]]; ok {
				connectionsMap[connection[0]] = append(connectionsMap[connection[0]], connection[1])
			}
			if _, ok := connectionsMap[connection[1]]; ok {
				connectionsMap[connection[1]] = append(connectionsMap[connection[1]], connection[0])
			}
		} else {
			CheckError(errors.New("ERROR: Invalid connection format"))
		}
	}

	return connectionsMap
}

// RemoveDeadEnds removes dead-end paths
func RemoveDeadEnds(initialMap map[string][]string) map[string][]string {
	for key, value := range initialMap {
		if len(value) == 1 && value[0] != "start" && value[0] != "end" {
			delete(initialMap, key)
		} else if len(value) == 2 && value[0] != "start" && value[0] != "end" {
			for key1, values := range initialMap {
				for i, v := range values {
					if v == key {
						initialMap[key1] = RemoveElement(values, i)
					}
				}
			}
			delete(initialMap, key)
			return RemoveDeadEnds(initialMap)
		}
	}

	return initialMap
}

// RemoveElement removes a value from a slice of strings
func RemoveElement(slice []string, s int) []string {
	return append(slice[:s], slice[s+1:]...)
}

// CheckError prints an error message if err is not nil
func CheckError(e error) {
	if e != nil {
		fmt.Println("ERROR: Invalid data format")
		log.Fatal(e)
	}
}
