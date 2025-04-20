package main

import (
	"flag"
	"fmt"
	"os"
	"runtime/pprof"

	"github.com/myhops/queens/pkg/board1"
)

type Options struct {
	gameFile   string
	memProfile string
	cpuProfile string
}

func getOptions() *Options {
	o := &Options{}
	flag.StringVar(&o.gameFile, "game", "", "game file")
	flag.StringVar(&o.memProfile, "memprofile", "", "write memory profile to `file`")
	flag.StringVar(&o.cpuProfile, "cpuprofile", "", "write cpu profile to `file`")
	flag.Parse()
	return o
}

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

func run(args []string) error {
	o := getOptions()

	// Start profiling
	if o.cpuProfile != "" {
		f, err := os.Create(o.cpuProfile)
		if err != nil {
			return err
		}
		defer f.Close()
		if err := pprof.StartCPUProfile(f); err != nil {
			return err
		}
		defer pprof.StopCPUProfile()
	}

	// Load the board
	r, err := os.Open(o.gameFile)
	if err != nil {
		return err
	}
	defer r.Close()

	a, i, err := board1.LoadAreas(r)
	if err != nil {
		return err
	}
	g := board1.NewGame(i, i, a...)
	// solve
	s := &board1.AreaSolver{}
	b, err := g.Solve(s)
	if err != nil {
		return err
	}

	b.Print()
	fmt.Printf("solve called: %d\n", g.SolveCalled())
	fmt.Printf("boards used: %d\n", g.BoardPool.MaxEntries())

	if o.memProfile != "" {
		f, err := os.Create(o.memProfile)
		if err != nil {
			return err
		}
		defer f.Close()

		if err := pprof.WriteHeapProfile(f); err != nil {
			return err
		}
	}
	return nil
}

func main() {
	run(os.Args)
}
