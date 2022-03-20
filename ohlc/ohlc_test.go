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
		{date("2020-01-01"), d(10), d(30), d(5), d(15)},
		{date("2020-01-02"), d(15), d(15), d(15), d(15)},
		{date("2020-01-03"), d(15), d(20), d(5), d(15)},
		{date("2020-01-04"), d(20), d(20), d(10), d(10)},
		{date("2020-01-05"), d(15), d(20), d(5), d(15)},
	}

	p.SetRect(0, 0, 20, 15)

	p.SetItems(items)

	p.Draw(scr)

	exp := `          ╷    30.00
          │    28.21
          │    26.43
          │    24.64
          │    22.86
          │    21.07
          │ ╷╻╷19.29
          │ │┃│17.50
          │ │┃│15.00
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
