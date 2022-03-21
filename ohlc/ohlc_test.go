package ohlc_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/jeremija/tplot/ohlc"
	"github.com/jeremija/tplot/test"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestOHLC(t *testing.T) {
	p := ohlc.New()
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

	items := []ohlc.Item{
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
        │      28.68
        │      27.37
        │      26.05
        │      24.74
        │      23.42
        │      22.11
        │      20.79
        │ ╷╻╷  19.47
        │ │┃│  18.16
        │ │┃│  16.84
        │ │┃│  15.00
        ╽─┼┃┼  14.21
        ┃ │┃│  12.89
        ┃ │┃│  11.58
        ┃ │┃│  10.26
        ╿ │╹│   8.95
        │ │ │   7.63
        │ │ │   6.32
        ╵ ╵ ╵   5.00
            █1000.00
            █ 911.11
            █ 822.22
         ▆  █ 733.33
         █  █ 644.44
         █  █ 555.56
        ▆█  █ 466.67
        ██  █ 377.78
        ██ ▃█ 288.89
        ██▃██ 200.00`

	fmt.Println("===")
	fmt.Println(scr.Content())
	fmt.Println("===")

	assert.Equal(t, exp, "\n"+scr.Content())
}
