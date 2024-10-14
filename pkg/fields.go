package pkg

type Position struct{ x, y int }

type Fields struct {
	size int
}

func (f *Fields) Row(y int) []Position {
	positions := make([]Position, f.size)

	for x := 0; x < f.size; x++ {
		positions[x].x = x
		positions[x].y = y
	}

	return positions
}

func (f *Fields) Column(x int) []Position {
	positions := make([]Position, f.size)

	for y := 0; y < f.size; y++ {
		positions[y].x = x
		positions[y].y = y
	}

	return positions
}

func (f *Fields) Around(x, y int) []Position {
	positions := make([]Position, 0, 8)

	for i := max(x-1, 0); i < min(x+2, f.size); i++ {
		for j := max(y-1, 0); j < min(y+2, f.size); j++ {
			if i != x || j != y {
				positions = append(positions, Position{x: i, y: j})
			}
		}
	}
	return positions
}

