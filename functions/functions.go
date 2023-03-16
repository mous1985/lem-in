package functions

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
)

//trouvé le nombre de fourmis
func NombreDefourmis(data []string) (int, []string) { //je parse les lignes stocker dans le tableau data
	var nbrFourmis int
	var roomData []string
	var err error

	for i := 0; i < len(data); i++ {
		if data[i][0] != '#' {
			nbrFourmis, err = strconv.Atoi(data[i])
			if nbrFourmis < 1 {
				err = errors.New("EROOR:pas assez de fourmis")
			}
			ErrorCheck(err)
			roomData = data[i+1:]
			break
		}
	}

	return nbrFourmis, roomData //je renvois le nbr de fourmis,le restant de mon []data
}

// split les lignes restante en room et en relation
func SeparData(roomData []string) ([]string, []string) {

	var dataEmplace []string

	//sparé les coordonées de la room
	for i := 0; i < len(roomData); i++ {
		if len(strings.Split(roomData[i], " ")) == 3 {
			dataEmplace = append(dataEmplace, roomData[i])
		} else if roomData[i][0] == '#' {
			if roomData[i][1] == '#' {
				dataEmplace = append(dataEmplace, roomData[i])
			}
		} else {
			roomData = roomData[i:]
			break
		}
	}

	//les relations entre les rooms
	var realationData []string

	for _, v := range roomData {
		if len(strings.Split(v, "-")) == 2 {
			realationData = append(realationData, v)
		} else if v[0] == '#' && v[1] != '#' {

		} else {
			var err error = nil
			err = errors.New("ERROR:invalid data format")
			ErrorCheck(err)
		}
	}

	return dataEmplace, realationData
}

//la Map des room et leurs coordonnées
func CoordoneMapRoom(lignes []string) map[string][]int {
	roomMap := make(map[string][]int)
	var err error

	//check pour la room ##start et ##end emplacement
	var startPoint int
	var endPoint int
	var startPointer *int
	var endPointer *int

	for i, v := range lignes {
		if v == "##start" {
			startPoint = i + 1
			startPointer = &startPoint
		} else if v == "##end" {
			endPoint = i + 1
			endPointer = &endPoint
		}
	}
	//si l'une des lignes ##start ou ##end n'existe pas, affiche error
	if startPointer == nil || endPointer == nil {
		err = errors.New("ERROR:il manque la room start ou la room end")
		ErrorCheck(err)
	}

	//met chaque emplacement et ses coordonnées dans la roomMap

	for i := 0; i < len(lignes); i++ {
		room := strings.Split(lignes[i], " ")
		if len(room) == 3 {
			coorX, err := strconv.Atoi(room[1])
			ErrorCheck(err)
			coorY, err := strconv.Atoi(room[2])
			ErrorCheck(err)
			if i == startPoint {
				roomMap[room[0]] = append(roomMap[room[0]], 0, coorX, coorY)
			} else if i == endPoint {
				roomMap[room[0]] = append(roomMap[room[0]], 2, coorX, coorY)
			} else if lignes[i][:2] != "##" {
				roomMap[room[0]] = append(roomMap[room[0]], 1, coorX, coorY)
			}
		}
	}

	//vérification des salles avec les meme coordonnées
	for key1, v1 := range roomMap {
		for key2, v2 := range roomMap {
			if key1 == key2 {
				break
			} else if v1[1] == v2[1] && v1[2] == v2[2] {
				err := errors.New(key1 + " et " + key2 + " ont les meme coordonnées ")
				ErrorCheck(err)
			}
		}
	}

	return roomMap
}

//la map des connexions entre les pièces
func MapRoomConnections(dataBrute []string, coordoneMap map[string][]int) map[string][]string {
	var err error
	originalMap := make(map[string][]string)

	for key, v := range coordoneMap {
		switch v[0] {
		case 0:
			originalMap[key] = append(originalMap[key], "start")
		case 2:
			originalMap[key] = append(originalMap[key], "end")
		default:
			originalMap[key] = append(originalMap[key], "milieu")
		}
	}

	//ajouter les relations a la map
	for _, v := range dataBrute {
		connection := strings.Split(v, "-")
		if originalMap[connection[0]] != nil && originalMap[connection[1]] != nil {
			originalMap[connection[0]] = append(originalMap[connection[0]], connection[1])
			originalMap[connection[1]] = append(originalMap[connection[1]], connection[0])
		} else {
			err = errors.New("ERROR:manque emplacement")
			ErrorCheck(err)
		}
	}
	return originalMap
}

//supprimé les chemains sans issus
func SupprSsissu(mapDepart map[string][]string) map[string][]string {
	var err error

	for key, value := range mapDepart {
		if len(value) == 1 {
			if value[0] == "start" || value[0] == "end" {
				err = errors.New("ERROR:il manque un chemain ou la salle start")
				ErrorCheck(err)
			} else {
				delete(mapDepart, key)
			}
		} else if len(value) == 2 && value[0] != "start" && value[0] != "end" {
			for key1, values := range mapDepart {
				for i, v := range values {
					if v == key {
						values = Remove(values, i)
						mapDepart[key1] = values
					}
				}
			}
			delete(mapDepart, key)
			SupprSsissu(mapDepart)
		}
	}

	return mapDepart
}

//remove supprime une valeur nommée de la string
func Remove(slice []string, s int) []string {
	return append(slice[:s], slice[s+1:]...)
}

//Return  error si err!=nil
func ErrorCheck(e error) {
	if e != nil {
		fmt.Println("ERROR:invalid data format")
		log.Fatal(e)
	}
}
