package board1

import (
	"fmt"
	"iter"
	"slices"
)

type State int

const (
	Empty State = iota
	Blocked
	Queen
)

func (s State) String() string {
	switch s {
	case Empty:
		return " "
	case Blocked:
		return "X"
	case Queen:
		return "Q"
	default:
		return "?"
	}
}

type Board struct {
	Fields [][]State // row, col (row 0 is the top)

	Areas []Area
}

type Position []int // row, col

type Area []Position

func NewBoard(size int, areas ...Area) *Board {
	b := &Board{
		Fields: make([][]State, size),
	}

	for i := range b.Fields {
		b.Fields[i] = make([]State, size)
	}

	// Create our own copy
	b.Areas = slices.Clone(areas)

	return b
}

func (b *Board) PlaceQueen(p Position) error {
	row := p[0]
	col := p[1]
	if b.Fields[row][col] != Empty {
		return fmt.Errorf("position (%d, %d) is occupied with %s", row, col, b.Fields[row][col].String())
	}
	b.Fields[row][col] = Queen
	b.blockAround(row, col)
	b.blockHorizontal(row)
	b.blockVertical(col)

	for _, area := range b.Areas {
		b.blockArea(area)
	}
	return nil
}

func (b *Board) blockArea(area Area) {
	for _, pos := range area {
		x, y := pos[0], pos[1]
		if b.Fields[x][y] == Empty {
			b.Fields[x][y] = Blocked
		}
	}
}

func (b *Board) blockAround(row, col int) {
	for i := max(row-1, 0); i < min(row+2, len(b.Fields)); i++ {
		for j := max(col-1, 0); j < min(col+2, len(b.Fields)); j++ {
			if b.Fields[i][j] == Empty {
				b.Fields[i][j] = Blocked
			}
		}
	}
}

func (b *Board) blockHorizontal(x int) {
	for i := range len(b.Fields) {
		if b.Fields[x][i] == Empty {
			b.Fields[x][i] = Blocked
		}
	}
}
func (b *Board) blockVertical(y int) {
	for i := 0; i < len(b.Fields); i++ {
		if b.Fields[i][y] == Empty {
			b.Fields[i][y] = Blocked
		}
	}
}

func (b *Board) ClearFields() {
	for i := range b.Fields {
		for j := range b.Fields[i] {
			b.Fields[i][j] = Empty
		}
	}
}

func (b *Board) FindEmpty(from Position) (Position, error) {
	if from[0] > len(b.Fields) || from[1] > len(b.Fields[0]) {
		return Position{0, 0}, fmt.Errorf("out of range")
	}
	for i := from[0]; i < len(b.Fields); i++ {
		for j := from[1]; j < len(b.Fields); j++ {
			if b.Fields[i][j] == Empty {
				return Position{i, j}, nil
			}
		}
	}
	return nil, fmt.Errorf("no empty position found")
}

func (b *Board) Print() {
	for i := range b.Fields {
		for j := range b.Fields[i] {
			fmt.Print(b.Fields[i][j].String())
		}
		fmt.Println()
	}
}

// Solve the board by placing n queens
func (b *Board) Solve() ([]Position, error) {
	var res []Position
	return b.solve(res)
}

func (b *Board) Remaining(p Position) iter.Seq[Position] {
	return func(yield func(p Position) bool) {
		l := len(b.Fields) * len(b.Fields[0])
		for i := range l {
			x := i / len(b.Fields)
			y := i % len(b.Fields[0])

			if !yield(Position{x, y}) {
				return
			}
		}
	}
}

func (b *Board) Next(p Position) (Position, error) {
	pos := p[0]*len(b.Fields) + p[1] + 1
	if pos >= len(b.Fields)*len(b.Fields[0]) {
		return Position{}, fmt.Errorf("no more positions")
	}
	return Position{pos / len(b.Fields), pos % len(b.Fields[0])}, nil
}

func (b *Board) solve(queens []Position) ([]Position, error) {
	// Test if done
	if len(queens) == len(b.Fields) {
		return queens, nil
	}

	last := Position{0, 0}
	for {
		// Clear the fields
		b.ClearFields()

		// Place the queens
		b.PlaceQueens(queens)
		
		try, err := b.FindEmpty(last)
		// No empty place found
		if err != nil {
			return nil, err
		}

		// found an empty slot
		// add the queen and try to solve the rest of the puzzle
		res, err := b.solve(append(queens, try))
		if err == nil {
			return res, nil
		}

		// if an error occured, advance last and continue
		l, err := b.Next(last)
		if err != nil {
			return nil, err
		}
		last = l
	}
}

func (b *Board) PlaceQueens(queens []Position) {
	for _, queen := range queens {
		b.PlaceQueen(queen)
	}
}
