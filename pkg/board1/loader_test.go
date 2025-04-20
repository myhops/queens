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

	g := NewGame(i, i, a...)
	b, err := g.Solve()
	if err != nil {
		t.Fatal(err)
	}

	b.Print()
	t.Error()
}

func TestLoad2(t *testing.T) {
	const board = `000000222
100600022
104600002
104678022
104778922
114778922
334778922
333778222
322222222
`

	r := strings.NewReader(board)
	a, i, err := LoadAreas(r)
	if err != nil {
		t.Fatal(err)
	}
	// t.Logf("dim %d, areas: %#v", i, a)
	// t.Error()

	g := NewGame(i, i, a...)
	b, err := g.Solve()
	if err != nil {
		t.Fatal(err)
	}

	b.Print()

	t.Errorf("solve called: %d", g.SolveCalled())
	t.Errorf("boards used: %d", g.BoardPool.MaxEntries())
	t.Error()
}


func TestLoad3(t *testing.T) {
	const board = `0	0	0	0	0	0	2	2	2
1	0	0	6	0	0	0	2	2
1	0	4	6	0	0	0	0	2
1	0	4	6	7	8	0	2	2
1	0	4	7	7	8	9	2	2
1	1	4	7	7	8	9	2	2
3	3	4	7	7	8	9	2	2
3	3	3	7	7	8	9	2	2
3	2	2	2	2	2	2	2	2
`

	r := strings.NewReader(board)
	a, i, err := LoadAreas(r)
	if err != nil {
		t.Fatal(err)
	}
	// t.Logf("dim %d, areas: %#v", i, a)
	// t.Error()

	g := NewGame(i, i, a...)
	b, err := g.Solve()
	if err != nil {
		t.Fatal(err)
	}

	b.Print()

	t.Errorf("solve called: %d", g.SolveCalled())
	t.Errorf("boards used: %d", g.BoardPool.MaxEntries())
	t.Error()
}
