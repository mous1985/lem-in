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

	//ouvre le fichier spécifié dans le premier argument et lit son contenu
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

	//récupere nombre de fourmis
	nbrFourmis, roomData := functions.NombreDefourmis(data)

	//séparé les data a partir des coordonées des salles et les relation entre elles
	dataEmplace, dataRelation := functions.SeparData(roomData)
	coordoneMap := functions.CoordoneMapRoom(dataEmplace)
	originalMap := functions.MapRoomConnections(dataRelation, coordoneMap)
	finaleMap := functions.SupprSsissu(originalMap)

	//calculer le chemin optimal pour les fourmis
	path, distribution := graph.Parcour(finaleMap, nbrFourmis)

	fmt.Println(string(input) + "\n")

	var fourmis [][]int
	antNr := 1
	// Création d'une slice dans la variable fourmis[][] qui contient des entiers représentant le nombre de fourmis
	for i := 0; i < len(distribution); i++ {
		fourmis = append(fourmis, []int{})

		for j := 0; j < distribution[i]; j++ {
			fourmis[i] = append(fourmis[i], antNr)
			antNr++
		}
	}
	maxAnt := 0
	//Affectation de la plus grande valeur de la variable fourmis à la variable maxAnt
	for i := 0; i < len(fourmis); i++ {
		if len(fourmis[i]) > maxAnt {
			maxAnt = len(fourmis[i])
		}
	}

	connecter := false
	restante := nbrFourmis
	//Création d'une boucle infinie qui imprime le output  et se termine lorsque  la variable restante est égal à 0
	for i := 1; i > 0; i++ {
		Lafourmis := 0
		for k := 0; k < maxAnt; k++ {
			for j := 0; j < len(fourmis); j++ {
				if (len(fourmis[j]) > Lafourmis) && (len(path[j]) > i-Lafourmis) && (i-Lafourmis > 0) {
					fmt.Print("L")
					fmt.Print(fourmis[j][Lafourmis])
					fmt.Print("-")
					fmt.Print(path[j][i-Lafourmis])
					fmt.Print(" ")
					if i-Lafourmis == len(path[j])-1 {
						restante--
					}
					if len(path[j]) == 2 {
						connecter = true
					} else {
						connecter = false
					}
				}

			}
			Lafourmis++
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
