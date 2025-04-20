package board1

import (
	"bufio"
	"io"
	"strings"
)

func LoadAreas(r io.Reader) ([]Area, int, error) {
	m := map[rune]Area{}

	s := bufio.NewScanner(r)
	var x int
	for s.Scan() {
		line := s.Text()
		// skip empty line
		if 	line == "" {
			continue
		}
		// remove all spaces
		line = strings.ReplaceAll(line, " ", "")
		line = strings.ReplaceAll(line, "\t", "")
		var y int
		for _, c := range line {
			m[c] = append(m[c], Position{x, y})
			y++
		}
		x++
	}
	res := make([]Area, 0, len(m))
	for _, a := range m {
		res = append(res, a)
	}
	return res, x, nil
}


