package board1

import "slices"

type Area []Position

func (a Area) Contains(row, col int) bool {
	for _, pos := range a {
		if pos[0] == row && pos[1] == col {
			return true
		}
	}
	return false
}

func SortAreasReverse(areas []Area) {
	slices.SortFunc(areas, func(a, b Area) int {
		return len(b) - len(a)
	})
}

