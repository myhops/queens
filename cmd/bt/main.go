package main

import (
	"fmt"
	"os"

	"github.com/myhops/queens/pkg/board1"
)

func main() {
	b := board1.NewBoard(10)

	q, err := b.Solve()
	if err != nil {
		os.Exit(1)
	}
	fmt.Printf("%+v\n", q)
}
