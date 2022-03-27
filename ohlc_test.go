package tplot_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/jeremija/tplot"
	"github.com/jeremija/tplot/test"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestOHLC(t *testing.T) {
	p := tplot.NewOHLC()
	scr := test.NewScreen()

	d := func(val int64) decimal.Decimal {
		return decimal.New(val, 0)
	}

	date := func(str string) time.Time {
		ts, err := time.Parse("2006-01-02", str)
		if err != nil {
			panic(err)
		}

		return ts
	}

	items := []tplot.OHLCItem{
		{date("2020-01-01"), d(10), d(30), d(5), d(15), d(500)},
		{date("2020-01-02"), d(15), d(15), d(15), d(15), d(750)},
		{date("2020-01-03"), d(15), d(20), d(5), d(15), d(200)},
		{date("2020-01-04"), d(20), d(20), d(10), d(10), d(300)},
		{date("2020-01-05"), d(15), d(20), d(5), d(15), d(1000)},
	}

	p.SetRect(0, 0, 20, 30)

	p.SetItems(items)

	p.Draw(scr)

	exp := `
        ╷      30.00
        │      28.91
        │      27.83
        │      26.74
        │      25.65
        │      24.57
        │      23.48
        │      22.39
        │      21.30
        │      20.22
        │ ╷╻╷  19.13
        │ │┃│  18.04
        │ │┃│  16.96
        │ │┃│  15.87
        ╽─┼┃┼  15.00
        ┃ │┃│  13.70
        ┃ │┃│  12.61
        ┃ │┃│  11.52
        ┃ │┃│  10.43
        ╿ │╹│   9.35
        │ │ │   8.26
        │ │ │   7.17
        │ │ │   6.09
        ╵ ╵ ╵   5.00
            ▆1000.00
            █ 840.00
         ▆  █ 680.00
         █  █ 520.00
        ██  █ 360.00
        ██ ▄█ 200.00`

	fmt.Println("===")
	fmt.Println(scr.Content())
	fmt.Println("===")

	assert.Equal(t, exp, "\n"+scr.Content())
}
