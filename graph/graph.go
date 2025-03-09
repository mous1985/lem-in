package graph

import "sort"

var (
	AntsFarm = make(map[string][]string)
	Paths    [][]string
)

// FindPaths takes the AntsFarm and the number of ants and returns all possible paths from the AntsFarm to the end,
// with the optimal distribution of ants on each path.
func FindPaths(data map[string][]string, ants int) ([][]string, []int) {
	AntsFarm = data

	var start string
	for room, options := range AntsFarm {
		if options[0] == "start" {
			start = room
			break
		}
	}

	findPossiblePaths([]string{start})

	var end string
	for room, options := range AntsFarm {
		if options[0] == "end" {
			end = room
			break
		}
	}

	return filterPaths(end, ants)
}

// findPossiblePaths takes the AntsFarm and gives us all possible paths to the end.
func findPossiblePaths(path []string) {
	currentRoom := path[len(path)-1]
	options := AntsFarm[currentRoom]

	// if we are at the end, add the path to Paths
	if len(options) > 0 && options[0] == "end" {
		Paths = append(Paths, path)
		return
	}

	// try to find the room that we haven't visited yet
	for _, option := range options {
		if option != "start" && !contains(path, option) {
			newPath := append([]string(nil), path...)
			newPath = append(newPath, option)
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
	minMoves := int(^uint(0) >> 1) // Initialize to max int value

	for _, shortPath := range Paths {
		if len(shortPath) == 2 {
			return [][]string{shortPath}, []int{ants}
		}

		paths := structurePaths(shortPath[1:len(shortPath)-1], [][]string{shortPath}, Paths)
		if allPathsEndAt(paths, endRoom) {
			newPaths, newDistribution, newMoves := calculateDistribution(paths, ants)
			if newMoves < minMoves {
				bestPaths, bestDistribution, minMoves = newPaths, newDistribution, newMoves
			}
		}
	}

	return bestPaths, bestDistribution
}

func sortPathsByLength(paths [][]string) {
	sort.Slice(paths, func(i, j int) bool {
		return len(paths[i]) < len(paths[j])
	})
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
	middleSet := make(map[string]struct{}, len(middle1))
	for _, room := range middle1 {
		middleSet[room] = struct{}{}
	}

	for _, long := range tempPaths {
		middle2 := long[1 : len(long)-1]
		contains := false

		for _, room := range middle2 {
			if _, found := middleSet[room]; found {
				contains = true
				break
			}
		}

		if !contains {
			paths = append(paths, long)
			for _, room := range middle2 {
				middleSet[room] = struct{}{}
			}
		}
	}
	return paths
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

	if antsArrived >= ants {
		if antsArrived > ants {
			return recount(paths, ants, antsArrived, distribution, moves)
		}
		return paths, distribution, moves
	}

	remainingAnts := ants - antsArrived

	if len(distribution) == 1 {
		distribution[0] += remainingAnts
		moves += remainingAnts
	} else {
		for i := 0; remainingAnts > 0; i = (i + 1) % len(distribution) {
			distribution[i]++
			remainingAnts--
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

	distribution := make([]int, len(paths))
	antsArrived := 0
	longestPath := len(paths[len(paths)-1])

	for i, path := range paths {
		antsOnPath := longestPath - len(path) + 1
		antsArrived += antsOnPath
		distribution[i] = antsOnPath
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
