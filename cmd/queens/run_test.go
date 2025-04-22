package main

import "testing"

func TestArea(t *testing.T) {
	args := []string{"bt", "-game", "../../2025-4-22.txt", "-solver", "simple"}

	err := run(args)
	if err != nil {
		t.Error(err)
	}
}
