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
		fmt.Println("ERROR:nombre d'argument incorrecte! ")
	}

	// ouvre le fichier spécifié dans le premier argument et lit son contenu
	fileName := "./examples/" + os.Args[1]
	input, err := os.ReadFile(fileName)
	functions.ErrorCheck(err)
	datas := strings.Split(string(input), "\n")
	var data []string
	for _, line := range datas {
		LineV := strings.TrimSpace(line)
		if len(LineV) > 0 {
			data = append(data, LineV)
		}
	}

	// récupere nombre de Ant
	nbrAnts, roomData := functions.NbrOfAntsInRoom(data)

	// séparé les data a partir des coordonées des salles et les relation entre elles
	dataEmplace, dataRelation := functions.SeparData(roomData)
	coordoneMap := functions.CoordoneMapRoom(dataEmplace)
	originalMap := functions.MapRoomConnections(dataRelation, coordoneMap)
	finaleMap := functions.SupprSsissu(originalMap)

	// calculer le chemin optimal pour les Ants
	path, distribution := graph.Parcour(finaleMap, nbrAnts)

	fmt.Println(string(input) + "\n")

	var Ants [][]int
	AntsNr := 1
	// Création d'une slice dans la variable Ants[][] qui contient des entiers représentAnts le nombre de Ants
	for i := 0; i < len(distribution); i++ {
		Ants = append(Ants, []int{})

		for j := 0; j < distribution[i]; j++ {
			Ants[i] = append(Ants[i], AntsNr)
			AntsNr++
		}
	}
	maxAnts := 0
	// Affectation de la plus grande valeur de la variable Ants à la variable maxAnts
	for i := 0; i < len(Ants); i++ {
		if len(Ants[i]) > maxAnts {
			maxAnts = len(Ants[i])
		}
	}

	connecter := false
	restante := nbrAnts
	// Création d'une boucle infinie qui imprime le output  et se termine lorsque  la variable restAntse est égal à 0
	for i := 1; i > 0; i++ {
		LaAnts := 0
		for k := 0; k < maxAnts; k++ {
			for j := 0; j < len(Ants); j++ {
				if (len(Ants[j]) > LaAnts) && (len(path[j]) > i-LaAnts) && (i-LaAnts > 0) {
					fmt.Print("L")
					fmt.Print(Ants[j][LaAnts])
					fmt.Print("-")
					fmt.Print(path[j][i-LaAnts])
					fmt.Print(" ")
					if i-LaAnts == len(path[j])-1 {
						restante--
					}
					if len(path[j]) == 2 {
						connecter = true
					} else {
						connecter = false
					}
				}
			}
			LaAnts++
		}
		if !connecter {
			fmt.Println()
		}

		if restante == 0 {
			if connecter {
				fmt.Println()
			}

			break
		}

	}
}
