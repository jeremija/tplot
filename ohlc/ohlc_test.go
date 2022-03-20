package ohlc_test

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/jeremija/tplot/ohlc"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestOHLC(t *testing.T) {
	p := ohlc.New()
	scr := newScreen()

	ohlcs := []ohlc.OHLC{
		{decimal.New(10, 0), decimal.New(30, 0), decimal.New(5, 0), decimal.New(15, 0), time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)},
		{decimal.New(15, 0), decimal.New(15, 0), decimal.New(15, 0), decimal.New(15, 0), time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)},
		{decimal.New(15, 0), decimal.New(20, 0), decimal.New(5, 0), decimal.New(15, 0), time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)},
		{decimal.New(20, 0), decimal.New(20, 0), decimal.New(10, 0), decimal.New(10, 0), time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)},
		{decimal.New(15, 0), decimal.New(20, 0), decimal.New(5, 0), decimal.New(15, 0), time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)},
	}

	p.SetRect(0, 0, 20, 15)

	p.SetOHLC(ohlcs)

	p.Draw(scr)

	exp := `          ╷    30.00
          │    28.21
          │    26.43
          │    24.64
          │    22.86
          │    21.07
          │ ╷╻╷19.29
          │ │┃│17.50
          │ │┃│15.71
          ╽─┼┃┼13.93
          ┃ │┃│12.14
          ┃ │┃│10.36
          ╿ │╹│ 8.57
          │ │ │ 6.79
          ╵ ╵ ╵ 5.00`

	fmt.Println("===")
	fmt.Println(scr.Content())
	fmt.Println("===")

	assert.Equal(t, exp, scr.Content())
}

type screen struct {
	tcell.Screen

	content [][]rune
}

func newScreen() *screen {
	s := &screen{}

	s.Clear()

	return s
}

func (s *screen) Clear() {
	s.content = nil
}

// func (s *screen) Size() (int, int) {
// 	return 30, 40
// }

func (s *screen) SetContent(x int, y int, mainc rune, combc []rune, style tcell.Style) {
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

func (s *screen) Content() string {
	ret := make([]string, len(s.content))

	for i, row := range s.content {
		ret[i] = strings.TrimSuffix(string(row), " ")
	}

	return strings.Join(ret, "\n")
}
