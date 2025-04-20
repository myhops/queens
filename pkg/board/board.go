package pkg

import (
	"errors"
	"fmt"
)

type fieldState int

var (
	ErrNotEmpty  = errors.New("field not empty")
	ErrBadFieldX = errors.New("invalid x value")
	ErrBadFieldY = errors.New("invalid x value")
)

const (
	IsEmpty fieldState = iota
	HasCrown
	IsBlocked
)



type Board struct {
	fields [][]fieldState // x, y
	areas  []Area
}

type Area []Position

func (a Area) has(x, y int) bool {
	for _, p := range a {
		if p.x == x && p.y == y {
			return true
		}
	}
	return false
}

func findAreaWith(areas []Area, x, y int) Area {
	for _, a := range areas {
		if a.has(x, y) {
			return a
		}
	}
	return nil
}

// cloneFields make a shallow copy of the board with a deep copy of the fields
func cloneFields(b *Board) *Board {
	res := &Board{}
	*res = *b

	// Clone the fields
	fields := make([][]fieldState, len(b.fields))
	l := len(fields[0])
	for i := range fields {
		fields[i] = make([]fieldState, l)
		copy(fields[i], b.fields[i])
	}
	res.fields = fields
	return res
}

func PutCrown(b *Board, x, y int) (*Board, error) {
	// Check if x and y are valid.
	if x < 0 || x > len(b.fields) {
		return nil, ErrBadFieldX
	}
	if y < 0 || x > len(b.fields[0]) {
		return nil, ErrBadFieldY
	}
	// Field must be empty
	if b.fields[x][y] != IsEmpty {
		return nil, ErrNotEmpty
	}

	// create a new board with cloned fields
	nb := cloneFields(b)
	nb.putCrown(x, y)

	return nb, nil
}

func (b *Board) putCrown(x, y int) {
	// Place the crown
	b.fields[x][y] = HasCrown

	// Block row cells that do not have a crown
	for r := 0; x < len(b.fields); x++ {
		if b.fields[r][y] != HasCrown {
			b.fields[r][y] = IsBlocked
		}
	}
	// Block column cells that do not have a crown
	for c := 0; y < len(b.fields[x]); x++ {
		if b.fields[x][c] != HasCrown {
			b.fields[x][c] = IsBlocked
		}
	}
	// Find the area that has x,y in it.
	a := findAreaWith(b.areas, x, y)
	if a == nil {
		panic(fmt.Sprintf("area not found for %d, %d", x, y))
	}
	// Block the area cells
	for _, p := range a {
		if b.fields[p.x][p.y] == IsEmpty {
			b.fields[p.x][p.y] = IsBlocked
		}
	}
}

func newFields(size int) [][]fieldState {
	res := make([][]fieldState, size)
	for i := 0; i < size; i++ {
		res[i] = make([]fieldState, size)
	}
	return res
}

func (p Position) onField(size int) bool {
	return 0 < p.x && p.x < size && 0 < p.y && p.y < size
}

func validateAreasForOverlap(size int, areas []Area) error {
	fields := newFields(size)

	for ia := range areas {
		for ip := range areas[ia] {
			pos := areas[ia][ip]
			if !pos.onField(size) {
				return fmt.Errorf("validateAreas: areas[%d][%d] is outside the field", ia, ip)
			}
			// Check if place not already blocked
			if fields[pos.x][pos.y] != IsEmpty {
				return fmt.Errorf("validateAreas: areas[%d][%d] overlaps", ia, ip)
			}
			fields[pos.x][pos.y] = IsBlocked
		}
	}
	return nil
}

func NewBoard(size int, crowns []Position, areas []Area) (*Board, error) {
	if err := validateAreasForOverlap(size, areas); err != nil {
		return nil, err
	}

	// Create a new board.
	b := &Board{
		fields: newFields(size),
		areas: areas,
	}

	return b, nil
}
