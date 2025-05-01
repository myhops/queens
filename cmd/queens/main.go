package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
	"runtime/pprof"
	"time"

	"github.com/myhops/queens/pkg/board1"
)

type Options struct {
	gameFile   string
	memProfile string
	cpuProfile string
	solver     string
	sheet      string
}

func getOptions(args []string) *Options {
	o := &Options{}
	fs := flag.NewFlagSet(args[0], flag.ExitOnError)
	fs.StringVar(&o.gameFile, "game", "", "game file")
	fs.StringVar(&o.memProfile, "memprofile", "", "write memory profile to `file`")
	fs.StringVar(&o.cpuProfile, "cpuprofile", "", "write cpu profile to `file`")
	fs.StringVar(&o.solver, "solver", "area", "solver to use")
	fs.StringVar(&o.sheet, "sheet", "queens", "Google sheet to use")
	fs.Parse(args[1:])
	return o
}

func getSolver(s string) board1.Solver {
	switch s {
	case "simple":
		return &board1.SimpleSolver{}
	case "area":
		return &board1.AreaSolver{}
	default:
		return &board1.AreaSolver{}
	}
}

func loadAreas(gameFile string) ([]board1.Area, int, error) {
	// Load the board
	r, err := os.Open(gameFile)
	if err != nil {
		return nil, 0, err
	}
	defer r.Close()

	a, i, err := board1.LoadAreas(r)
	if err != nil {
		return nil, 0, err
	}
	return a, i, nil
}

func run(args []string) error {
	logger := slog.Default().With(
		"method", "run",
	)
	defer func(start time.Time) {
		fmt.Printf("queens took %v\n", time.Since(start))
	}(time.Now())

	o := getOptions(args)

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

	a, i, err := loadAreas(o.gameFile)
	if err != nil {
		return err
	}
	logger.Debug("loaded areas", "count", len(a), "board_size", i)
	if len(a) != i {
		return fmt.Errorf("areas and board size do not match, areas: %d, board size: %d", len(a), i)
	}

	g := board1.NewGame(i, i, a...)
	// solve

	s := getSolver(o.solver)
	// run in func to ease timing
	b, err := func() (*board1.Board, error) {
		defer func(start time.Time) {
			fmt.Printf("solve took %v\n", time.Since(start))
		}(time.Now())

		b, err := g.Solve(s)
		if err != nil {
			return nil, err
		}
		return b, nil
	}()

	if err != nil {
		return err
	}

	b.Print()
	fmt.Printf("number of times queen placed: %d\n", g.QueenPlaced())
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
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug}))
	logger = logger.With(
		"application", "queens")
	slog.SetDefault(logger)

	run(os.Args)
}
