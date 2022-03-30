package main

import (
	"time"

	"github.com/jeremija/tplot"
	"github.com/rivo/tview"
	"github.com/shopspring/decimal"
)

func main() {
	size := 1001
	ohlcs := make([]tplot.OHLC, size)
	tickData := make([]decimal.Decimal, size)
	ts := time.Now().Truncate(time.Minute)

	for i := 0; i < size; i++ {
		v := decimal.NewFromInt(int64(i))
		ts := ts.Add(time.Minute)

		ohlcs[i] = tplot.OHLC{
			Timestamp: ts,
			O:         v.Add(decimal.NewFromInt(10)),
			H:         v.Add(decimal.NewFromInt(15)),
			L:         v.Sub(decimal.NewFromInt(15)),
			C:         v.Sub(decimal.NewFromInt(10)),
			V:         v.Mul(decimal.NewFromInt(500)),
		}

		tickData[i] = ohlcs[i].V
	}

	app := tview.NewApplication()

	factory := func() tview.Primitive {
		list := tview.NewList()

		c := tplot.NewContainer()

		list.AddItem("Bar", "Bar Chart", 'b', func() {
			bars := tplot.NewBars()
			bars.SetData(tickData)
			bars.SetSpacing(2)

			c.SetPrimitive(bars)
			app.SetFocus(bars)
		})
		list.AddItem("Tick", "Tick Chart", 't', func() {
			ticks := tplot.NewTicks()
			ticks.SetData(tickData)
			ticks.SetSpacing(2)

			c.SetPrimitive(ticks)
			app.SetFocus(ticks)
		})
		list.AddItem("OHLC", "OHLC Candles", 'o', func() {
			candles := tplot.NewOHLCCandles()
			candles.SetData(ohlcs)
			candles.SetSpacing(2)

			c.SetPrimitive(candles)
			app.SetFocus(candles)
		})

		c.SetPrimitive(list)

		return c
	}

	bl := tplot.NewTile()
	bl.SetDirection(tplot.DirectionVertical)
	bl.SetFactory(factory)

	if err := app.SetRoot(bl, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}
