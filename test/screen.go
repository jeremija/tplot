package test

import (
	"strings"

	"github.com/gdamore/tcell/v2"
)

// Screen is a mock implementation of tcell.Screen designed for testing.
type Screen struct {
	tcell.Screen

	content [][]rune
}

// NewScreen creates a new Screen mock.
func NewScreen() *Screen {
	return &Screen{}
}

// Clear implements tcell.Screen.
func (s *Screen) Clear() {
	s.content = nil
}

// SetContent implements tcell.Screen.
func (s *Screen) SetContent(x int, y int, mainc rune, combc []rune, style tcell.Style) {
	s.setRune(x, y, mainc)

	for i, c := range combc {
		s.setRune(x+i+1, y, c)
	}
}

func (s *Screen) setRune(x, y int, mainc rune) {
	yy := len(s.content)

	if yy <= y {
		v := make([][]rune, y+1)
		copy(v, s.content)
		s.content = v
	}

	row := s.content[y]
	xx := len(row)

	if xx <= x {
		v := make([]rune, x+1)
		copy(v, row)
		s.content[y] = v
	}

	s.content[y][x] = mainc
}

// Content returns the current content as string. All trailing spaces will be
// trimmed.
func (s *Screen) Content() string {
	ret := make([]string, len(s.content))

	for i, row := range s.content {
		ret[i] = strings.TrimRight(string(row), " ")
	}

	return strings.Join(ret, "\n")
}
