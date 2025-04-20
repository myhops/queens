package board1

import "testing"

func TestOne(t *testing.T) {
	// area := Area{
	// 	{0,0},
	// 	{0,1},
	// 	{0,2},
	// }

	board := NewBoard(5)

	q, err := board.Solve()
	if err != nil {
		t.Error(err)
	}

	board.ClearFields()
	board.PlaceQueens(q)
	board.Print()
	t.Error()
}
