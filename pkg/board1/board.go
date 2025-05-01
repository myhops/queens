package board1

import (
	"errors"
	"fmt"
	"sync"
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

type BoardPool struct {
	pool sync.Pool
	rows int
	cols int

	maxEntries int
}

func NewBoardPool(rows, cols int) *BoardPool {
	bp := &BoardPool{
		rows: rows,
		cols: cols,
		pool: sync.Pool{
			New: func() any {
				return &Board{
					Fields: make([]State, rows*cols),
					Rows:   rows,
					Cols:   cols,
				}
			},
		},
	}
	bp.pool.New = bp.New
	return bp
}

func (p *BoardPool) New() any {
	b := &Board{
		Fields: make([]State, p.rows*p.cols),
		Rows:   p.rows,
		Cols:   p.cols,
	}
	p.maxEntries++
	return b
}

func (p *BoardPool) Get() *Board {
	bp := p.pool.Get().(*Board)
	clear(bp.Fields)
	return bp
}

func (p *BoardPool) Put(b *Board) {
	p.pool.Put(b)
}

func (p *BoardPool) MaxEntries() int {
	return p.maxEntries
}

type Game struct {
	Areas []Area
	Rows  int
	Cols  int

	BoardPool *BoardPool

	solveCalled int
	queenPlaced int64
}

func NewGame(rows, cols int, areas ...Area) *Game {
	return &Game{
		Rows:      rows,
		Cols:      cols,
		BoardPool: NewBoardPool(rows, cols),
		Areas:     areas,
	}
}

type Board struct {
	Fields []State
	Rows   int
	Cols   int
}

// Put places a state on the board.
// Panics if the index is out of range
func (b *Board) Put(row, col int, s State) {
	if row < 0 || row >= b.Rows || col < 0 || col >= b.Cols {
		panic("index out of range")
	}
	b.Fields[row*b.Cols+col] = s
}

func (b *Board) Get(row, col int) State {
	if row < 0 || row >= b.Rows || col < 0 || col >= b.Cols {
		panic("index out of range")
	}
	return b.Fields[row*b.Cols+col]
}

type Position []int // row, col

func (g *Game) PlaceQueen(b *Board, row, col int) error {
	if b.Get(row, col) != Empty {
		return fmt.Errorf("position (%d, %d) is occupied with %s", row, col, b.Get(row, col).String())
	}
	b.Put(row, col, Queen)
	b.blockAround(row, col)
	b.blockRow(row)
	b.blockColumn(col)

	if a := g.inArea(row, col); a != nil {
		b.blockArea(a)
	}
	g.queenPlaced++
	return nil
}

func (g *Game) QueenPlaced() int64 {
	return g.queenPlaced
}

func (g *Game) inArea(row, col int) Area {
	for _, area := range g.Areas {
		if area.Contains(row, col) {
			return area
		}
	}
	return nil
}

func (b *Board) blockArea(area Area) {
	for _, pos := range area {
		x, y := pos[0], pos[1]
		if b.Get(x, y) == Empty {
			b.Put(x, y, Blocked)
		}
	}
}

func (b *Board) blockAround(row, col int) {
	for i := max(row-1, 0); i < min(row+2, b.Rows); i++ {
		for j := max(col-1, 0); j < min(col+2, b.Cols); j++ {
			if b.Get(i, j) == Empty {
				b.Put(i, j, Blocked)
			}
		}
	}
}

// blockRow blockx all fields in row row
func (b *Board) blockRow(row int) {
	for i := range b.Cols {
		if b.Get(row, i) == Empty {
			b.Put(row, i, Blocked)
		}
	}
}

func (b *Board) blockColumn(col int) {
	for i := 0; i < b.Rows; i++ {
		if b.Get(i, col) == Empty {
			b.Put(i, col, Blocked)
		}
	}
}

func (b *Board) ClearFields() {
	clear(b.Fields)
}

var ErrNoEmptyFound = errors.New("no empty position found")

func (b *Board) FindEmpty(row, col int) (int, int, error) {
	if row >= b.Rows || col >= b.Cols {
		return 0, 0, fmt.Errorf("out of range")
	}
	for i := row; i < b.Rows; i++ {
		for j := col; j < b.Cols; j++ {
			if b.Get(i, j) == Empty {
				return i, j, nil
			}
		}
	}
	return 0, 0, ErrNoEmptyFound
}

func (b *Board) Print() {
	for i := range b.Rows {
		for j := range b.Cols {
			fmt.Print(b.Get(i, j).String())
		}
		fmt.Println()
	}
}

func (b *Board) Next(row, col int) (int, int, error) {
	maxPos := b.Rows * b.Cols
	pos := row*b.Cols + col + 1

	if pos >= maxPos {
		return 0, 0, fmt.Errorf("no more positions")
	}

	row = pos / b.Cols
	col = pos % b.Cols
	return row, col, nil
}

func (b *Board) FindNextEmpty(row, col int) (int, int, error) {
	row, col, err := b.Next(row, col)
	if err != nil {
		return 0, 0, err
	}
	return b.FindEmpty(row, col)
}

func (g *Game) PlaceQueens(b *Board, queens []Position) {
	for _, queen := range queens {
		g.PlaceQueen(b, queen[0], queen[1])
	}
}

type Solver interface {
	Solve(g *Game) (*Board, error)
}

func (g *Game) Solve(solver Solver) (*Board, error) {
	g.solveCalled = 0
	b := g.BoardPool.Get()
	defer g.BoardPool.Put(b)

	return solver.Solve(g)
}

func (b *Board) CopyFrom(bc *Board) {
	copy(b.Fields, bc.Fields)
	b.Rows = bc.Rows
	b.Cols = bc.Cols
}

func (g *Game) SolveCalled() int {
	return g.solveCalled
}

var ErrNoSolution = errors.New("no solution found")

func (g *Game) solveBoard(b *Board, n int) (*Board, error) {
	if n == 0 {
		return b, nil
	}
	g.solveCalled++

	for row, col, err := b.FindEmpty(0, 0); err == nil; row, col, err = b.FindNextEmpty(row, col) {
		// Create a new board
		nb := g.BoardPool.Get()
		nb.CopyFrom(b)
		// Place the queen on the empty place, return on error
		if err := g.PlaceQueen(nb, row, col); err != nil {
			return nil, err
		}
		// Try to solve this board
		res, err := g.solveBoard(nb, n-1)
		if err == nil {
			return res, nil
		}
		// No solution found, return the board to the pool
		g.BoardPool.pool.Put(nb)
	}
	return nil, ErrNoSolution
}
