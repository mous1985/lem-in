package graph

var (
	AntsFarm = make(map[string][]string)
	Paths    [][]string
)

// FindPaths takes the AntsFarm and the number of ants and returns all possible paths from the AntsFarm to the end,
// with the optimal distribution of ants on each path.
func FindPaths(data map[string][]string, ants int) ([][]string, []int) {
	AntsFarm = data

	for start, options := range AntsFarm {
		if options[0] == "start" {
			findPossiblePaths([]string{start})
			break
		}
	}

	var paths [][]string
	var distribution []int

	for end, options := range AntsFarm {
		if options[0] == "end" {
			paths, distribution = filterPaths(end, ants)
		}
	}

	return paths, distribution
}

// findPossiblePaths takes the AntsFarm and gives us all possible paths to the end.
func findPossiblePaths(path []string) {
	options := AntsFarm[path[len(path)-1]]

	// if we are at the end, add the path to Paths
	if options[0] == "end" {
		Paths = append(Paths, path)
		return
	}

	// try to find the room that we haven't visited yet
	for _, option := range options[1:] {
		if !contains(path, option) {
			newPath := append(path, option)
			findPossiblePaths(newPath)
		}
	}
}

func contains(slice []string, item string) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}

/*
	filterPaths sorts the found paths to select the best path based on certain criteria

and returns this path along with the distribution of ants
*/
func filterPaths(endRoom string, ants int) ([][]string, []int) {
	if len(Paths) == 1 {
		return Paths, []int{ants}
	}

	sortPathsByLength(Paths)

	var bestPaths [][]string
	var bestDistribution []int
	var minMoves int

	for _, shortPath := range Paths {
		if len(shortPath) == 2 {
			return [][]string{shortPath}, []int{ants}
		}

		paths := structurePaths(shortPath[1:len(shortPath)-1], [][]string{shortPath}, Paths)
		if allPathsEndAt(paths, endRoom) {
			if minMoves == 0 {
				bestPaths, bestDistribution, minMoves = calculateDistribution(paths, ants)
			} else {
				newPaths, newDistribution, newMoves := calculateDistribution(paths, ants)
				if newMoves < minMoves {
					bestPaths, bestDistribution, minMoves = newPaths, newDistribution, newMoves
				}
			}
		}
	}

	return bestPaths, bestDistribution
}

func sortPathsByLength(paths [][]string) {
	for i := 1; i < len(paths); i++ {
		for j := i; j > 0 && len(paths[j]) < len(paths[j-1]); j-- {
			paths[j], paths[j-1] = paths[j-1], paths[j]
		}
	}
}

func allPathsEndAt(paths [][]string, endRoom string) bool {
	for _, path := range paths {
		if path[len(path)-1] != endRoom {
			return false
		}
	}
	return true
}

/*
	structurePaths compares each path to middle1. If a path contains a room that is also present in middle1,

the comparison is interrupted and the loop continues with the next path.
If no path contains a room that is also present in middle1,
the path is added to paths and the rooms from middle2 are added to middle1.
*/
func structurePaths(middle1 []string, paths [][]string, tempPaths [][]string) [][]string {
	for _, long := range tempPaths {
		middle2 := long[1 : len(long)-1]

		if !containsAny(middle1, middle2) {
			paths = append(paths, long)
			middle1 = append(middle1, middle2...)
		}
	}
	return paths
}

func containsAny(slice1, slice2 []string) bool {
	for _, item1 := range slice1 {
		for _, item2 := range slice2 {
			if item1 == item2 {
				return true
			}
		}
	}
	return false
}

// calculateDistribution calculates the optimal distribution of ants along different paths in the AntsFarm
func calculateDistribution(paths [][]string, ants int) ([][]string, []int, int) {
	antsArrived, distribution := moveAnts(paths)
	moves := len(paths[len(paths)-1]) - 1

	for _, path := range paths {
		if len(path)-1 > moves {
			moves = len(path) - 1
		}
	}

	if antsArrived > ants {
		return recount(paths, ants, antsArrived, distribution, moves)
	} else if antsArrived == ants {
		return paths, distribution, moves
	}

	ants -= antsArrived

	if len(distribution) == 1 {
		distribution[0] += ants
		moves += ants
	} else {
		for i := 0; ants > 0; i = (i + 1) % len(distribution) {
			distribution[i]++
			ants--
			if i == 0 {
				moves++
			}
		}
	}

	return paths, distribution, moves
}

/*
moveAnts takes a given path and returns the base amount of ants that finish
and how these ants are distributed on each path.
*/
func moveAnts(paths [][]string) (int, []int) {
	if len(paths) == 1 {
		return 1, []int{1}
	}

	var distribution []int
	var antsArrived int
	longestPath := len(paths[len(paths)-1])

	for _, path := range paths {
		if len(path) > longestPath {
			longestPath = len(path)
		}
	}

	for _, path := range paths {
		antsArrived += longestPath - len(path) + 1
		distribution = append(distribution, longestPath-len(path)+1)
	}

	return antsArrived, distribution
}

func recount(paths [][]string, ants int, antsArrived int, distribution []int, moves int) ([][]string, []int, int) {
	for i := len(distribution) - 1; i >= 0; i-- {
		distribution[i]--
		antsArrived--

		if distribution[i] == 0 {
			distribution = append(distribution[:i], distribution[i+1:]...)
			paths = append(paths[:i], paths[i+1:]...)
		}

		if antsArrived == ants {
			if i == 0 {
				moves--
			}
			return paths, distribution, moves
		}

		if i == 0 {
			moves--
			i = len(distribution)
		}
	}
	return paths, distribution, moves
}
