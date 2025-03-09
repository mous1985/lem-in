package functions

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
)

// NbrOfAntsInRoom finds the number of ants in the room
// and returns the number of ants and the data of the room
func NbrOfAntsInRoom(data []string) (int, []string) {
	for i, line := range data {
		if line[0] != '#' {
			nbrAnts, err := strconv.Atoi(line)
			if err != nil || nbrAnts < 1 {
				ErrorCheck(errors.New("ERROR: Invalid number of ants"))
			}
			return nbrAnts, data[i+1:]
		}
	}
	ErrorCheck(errors.New("ERROR: No valid number of ants found"))
	return 0, nil
}

// SeparData separates room data into two categories: room information and relationships.
func SeparData(roomData []string) ([]string, []string) {
	var dataEmplace, relationData []string

	for _, line := range roomData {
		if len(strings.Fields(line)) == 3 || strings.HasPrefix(line, "##") {
			dataEmplace = append(dataEmplace, line)
		} else if line[0] == '#' {
			continue
		} else if strings.Count(line, "-") == 1 {
			relationData = append(relationData, line)
		} else {
			ErrorCheck(errors.New("ERROR: Invalid data format"))
		}
	}

	return dataEmplace, relationData
}

// CoordoneMapRoom creates a map of rooms and their coordinates
func CoordoneMapRoom(lignes []string) map[string][]int {
	roomMap := make(map[string][]int)
	var startPoint, endPoint int
	var startPointer, endPointer *int

	for i, line := range lignes {
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
		ErrorCheck(errors.New("ERROR: Missing start or end room"))
		return nil
	}

	for i := startPoint; i <= endPoint; i++ {
		if strings.HasPrefix(lignes[i], "##") {
			continue
		}
		room := strings.Fields(lignes[i])
		if len(room) == 3 {
			coorX, err := strconv.Atoi(room[1])
			ErrorCheck(err)
			coorY, err := strconv.Atoi(room[2])
			ErrorCheck(err)
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
				ErrorCheck(errors.New(key1 + " and " + key2 + " have the same coordinates"))
			}
		}
	}

	return roomMap
}

// MapRoomConnections creates a map of connections between rooms
func MapRoomConnections(dataBrute []string, coordoneMap map[string][]int) map[string][]string {
	originalMap := make(map[string][]string)

	for key, v := range coordoneMap {
		switch v[0] {
		case 0:
			originalMap[key] = append(originalMap[key], "start")
		case 2:
			originalMap[key] = append(originalMap[key], "end")
		default:
			originalMap[key] = append(originalMap[key], "middle")
		}
	}

	for _, v := range dataBrute {
		connection := strings.Split(v, "-")
		if len(connection) == 2 {
			if _, ok := originalMap[connection[0]]; ok {
				originalMap[connection[0]] = append(originalMap[connection[0]], connection[1])
			}
			if _, ok := originalMap[connection[1]]; ok {
				originalMap[connection[1]] = append(originalMap[connection[1]], connection[0])
			}
		} else {
			ErrorCheck(errors.New("ERROR: Invalid connection format"))
		}
	}

	return originalMap
}

// SupprSsissu removes dead-end paths
func SupprSsissu(mapDepart map[string][]string) map[string][]string {
	for key, value := range mapDepart {
		if len(value) == 1 && value[0] != "start" && value[0] != "end" {
			delete(mapDepart, key)
		} else if len(value) == 2 && value[0] != "start" && value[0] != "end" {
			for key1, values := range mapDepart {
				for i, v := range values {
					if v == key {
						mapDepart[key1] = Remove(values, i)
					}
				}
			}
			delete(mapDepart, key)
			return SupprSsissu(mapDepart)
		}
	}

	return mapDepart
}

// Remove removes a value from a slice of strings
func Remove(slice []string, s int) []string {
	return append(slice[:s], slice[s+1:]...)
}

// ErrorCheck prints an error message if err is not nil
func ErrorCheck(e error) {
	if e != nil {
		fmt.Println("ERROR: Invalid data format")
		log.Fatal(e)
	}
}
