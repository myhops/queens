package main

import (
	"fmt"
	"log"
	"os"
	"runtime/pprof"
	"strings"

	"github.com/myhops/queens/pkg/board1"
)

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

func main() {
	f, err := os.Create("proffile")
	if err != nil {
		log.Fatal(err)
	}
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()

	r := strings.NewReader(board)
	a, i, err := board1.LoadAreas(r)
	if err != nil {
		panic(err)
	}

	g := board1.NewGame(i, i, a...)
	b, err := g.Solve()
	if err != nil {
		panic(err)
	}

	b.Print()

	fmt.Printf("solve called: %d\n", g.SolveCalled())
	fmt.Printf("boards used: %d\n", g.BoardPool.MaxEntries())

}
