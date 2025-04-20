package board1

import (
	"strings"
	"testing"
)

func TestLoad(t *testing.T) {
	const board = `0000000000
1222223330
1244553330
1244553330
1246653330
1276669920
1277889920
1277888920
1222222220
1111111111
`

	r := strings.NewReader(board)
	a, i, err := LoadAreas(r)
	if err != nil {
		t.Fatal(err)
	}
	// t.Logf("dim %d, areas: %#v", i, a)
	// t.Error()

	b := NewBoard(i, a...)
	_, err = b.Solve()
	if err != nil {
		t.Fatal(err)
	}

	b.Print()
	t.Error()
}
