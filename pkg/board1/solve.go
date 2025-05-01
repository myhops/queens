package board1

type AreaSolver struct {}

func (s *AreaSolver) Solve(g *Game) (*Board, error) {	
	g.solveCalled = 0
	b := g.BoardPool.Get()
	defer g.BoardPool.Put(b)

	// Sort the areas
	SortAreasReverse(g.Areas)

	return s.solveBoard(g, b, g.Cols)
}

func (s *AreaSolver) solveBoard(g *Game, b *Board, n int) (*Board, error) {
	if n == 0 {
		return b, nil
	}
	g.solveCalled++

	// Get the last area of areas
	a := g.Areas[n-1]

	for _, p := range a {
		row := p[0]
		col := p[1]
		if b.Get(row, col) != Empty {
			continue
		}

		// Create copy of board
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

type SimpleSolver struct {}

func (s *SimpleSolver) Solve(g *Game) (*Board, error) {
	g.solveCalled = 0
	b := g.BoardPool.Get()
	defer g.BoardPool.Put(b)

	return s.solveBoard(g, b, g.Cols)
}	

func (s *SimpleSolver) solveBoard(g *Game, b *Board, n int) (*Board, error) {
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
