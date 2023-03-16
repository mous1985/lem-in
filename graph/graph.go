package graph

var Fourmiliere = make(map[string][]string)
var Parcours [][]string

// Parcour prend la Fourmiliere et le nombre de Fourmis et retourne tous les chemins possibles de la Fourmiliere jusqu'à la fin,
//avec la distribution optimale de Fourmis sur chaque chemin.
func Parcour(data map[string][]string, Fourmis int) ([][]string, []int) {
	Fourmiliere = data

	for start, options := range Fourmiliere {
		if options[0] == "start" {
			cheminPossible([]string{start})
			break
		}
	}

	var Chemin [][]string
	var distribution []int

	for end, options := range Fourmiliere {
		if options[0] == "end" {
			Chemin, distribution = filter(end, Fourmis)
		}
	}

	return Chemin, distribution
}

// cheminPossible prend la Fourmiliere et nous donne tous les chemins possibles jusqu'à la fin.
func cheminPossible(Chemin []string) {
	options := Fourmiliere[Chemin[len(Chemin)-1]]

	// si on est a la fin ajoute le chemain a Parcours
	if options[0] == "end" {
		Parcours = append(Parcours, Chemin)
	} else {
		// essaye de trouver la room que nous n'avons pas encore parcourue
	boucleExt:
		for i := 1; i < len(options); i++ {

			for _, oldRoom := range Chemin {
				if oldRoom == options[i] {
					continue boucleExt
				}
			}
			newParcour := append(Chemin, options[i])
			test := make([]string, len(newParcour))
			for i := 0; i < len(newParcour); i++ {
				test[i] = newParcour[i]
			}
			cheminPossible(test)

		}
	}
}

/* filter trie les chemins trouvés pour sélectionner le bon chemin en fonction de certains critères
et retourner ce chemin ainsi que la distribution des fourmis */

func filter(endroom string, Fourmis int) ([][]string, []int) {
	var leChemin [][]string
	var laBonneDistribution []int
	var move int

	if len(Parcours) == 1 {
		return Parcours, []int{Fourmis}
	}

	for i := 1; i < len(Parcours); i++ {
		if i < 1 {
			continue
		} else if len(Parcours[i]) < len(Parcours[i-1]) {
			Parcours[i], Parcours[i-1] = Parcours[i-1], Parcours[i]
			i -= 2
		}
	}

	var tempParcour [][]string
	tempParcour = append(tempParcour, Parcours...)

	// l'étiquette `boucleExtest` utilisée pour marquer la boucle externe,
	//qui parcourt les chemins disponibles pour trouver le chemin le plus court.
boucleExt:
	for _, cheminCourt := range Parcours {

		if len(cheminCourt) == 2 {
			return [][]string{cheminCourt}, []int{Fourmis}
		}

		Chemin := strcutureParcours(cheminCourt[1:len(cheminCourt)-1], [][]string{cheminCourt}, tempParcour)
		for _, arr := range Chemin {
			if arr[len(arr)-1] != endroom {
				continue boucleExt
			}
		}
		if move < 1 {
			leChemin, laBonneDistribution, move = formule(endroom, Chemin, Fourmis)
		} else {
			newChemin, newDistribution, newMoves := formule(endroom, Chemin, Fourmis)
			if newMoves < move {
				laBonneDistribution, move, leChemin = newDistribution, newMoves, newChemin
			}
		}
	}

	return leChemin, laBonneDistribution
}

/* strcutureParcours compare chaque parcours à milieu1. Si un parcours contient une pièce qui est également présente dans milieu1,
la comparaison est interrompue et la boucle continue avec le parcours suivant.
Si aucun parcours ne contient de pièce qui est également présente dans milieu1,
le parcours est ajouté à Chemin et les pièces de milieu2 sont ajoutées à milieu1.*/

func strcutureParcours(milieu1 []string, Chemin [][]string, tempParcour [][]string) [][]string {
	var stoppeur bool
	for _, long := range tempParcour {
		milieu2 := long[1 : len(long)-1]

		for i, room1 := range milieu1 {
			for _, room2 := range milieu2 {
				if room2 == room1 {
					stoppeur = true
					break
				}
			}
			if stoppeur {
				stoppeur = false
				break
			}

			if i == len(milieu1)-1 {
				Chemin = append(Chemin, long)
				mid := make([]string, len(milieu1)+len(milieu2))
				i := 0
				for i < len(milieu1) {
					mid[i] = milieu1[i]
					i++
				}

				for j := 0; j < len(milieu2); j++ {
					mid[i] = milieu2[j]
					i++
				}

				milieu1 = mid
			}
		}
	}
	return Chemin
}

//la fonction formule calcule la distribution optimale de Fourmis le long de différents Chemins dans  la fourmiliere
func formule(endroom string, option [][]string, Fourmis int) ([][]string, []int, int) {

	fourmiArrive, distribution := moveFourmis(option)
	moves := len(option[len(option)-1]) - 1

	for _, arr := range option {
		if len(arr)-1 > moves {
			moves = len(arr) - 1
		}
	}

	if fourmiArrive > Fourmis {
		return decompte(option, Fourmis, fourmiArrive, distribution, moves)
	} else if fourmiArrive == Fourmis {
		return option, distribution, moves
	}

	Fourmis = Fourmis - fourmiArrive
	base := make([]int, len(distribution))
	copy(base, distribution)

	if len(distribution) == 1 {
		distribution[0] += Fourmis
		moves += Fourmis
	} else {
		for i := 0; i < len(distribution); i++ {
			if i == 0 {
				moves++
			}
			distribution[i]++
			Fourmis--
			if Fourmis == 0 {
				break
			}

			if Fourmis > 0 && i == len(distribution)-1 {
				i = -1
			} else if Fourmis == 0 {
				break
			}
		}
	}

	return option, distribution, moves
}

/*moveFourmis prend un Chemin donné et renvoie la quantité de base de Fourmis qui se terminent
et comment ces fourmis sont réparties sur chaque chemin.*/
func moveFourmis(Chemin [][]string) (int, []int) {
	if len(Chemin) == 1 {
		return 1, []int{1}
	}

	var distribution []int
	var fourmiArriver int
	longParcour := len(Chemin[len(Chemin)-1])

	for _, arr := range Chemin {
		if len(arr) > longParcour {
			longParcour = len(arr)
		}
	}

	for i := range Chemin {
		fourmiArriver += longParcour - len(Chemin[i]) + 1
		distribution = append(distribution, longParcour-len(Chemin[i])+1)
	}

	return fourmiArriver, distribution
}

func decompte(option [][]string, Fourmis int, fourmiArrive int, distribution []int, moves int) ([][]string, []int, int) {
	for i := len(distribution) - 1; i > -1; i-- {
		distribution[i] -= 1
		fourmiArrive--

		if distribution[i] == 0 {
			tempDis, tempChemin := distribution[:i], option[:i]
			tempDis, tempChemin = append(tempDis, distribution[i+1:]...), append(tempChemin, option[i+1:]...)
			distribution, option = tempDis, tempChemin
			i--
		}
		if fourmiArrive == Fourmis {
			if i == 0 {
				moves--
			}
			return option, distribution, moves
		}
		if i == 0 {
			moves--
			i = len(distribution) - 1
		}
	}
	return option, distribution, moves

}
