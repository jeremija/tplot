package main

import (
	"math/rand"

	"github.com/gdamore/tcell/v2"
	"github.com/jeremija/tplot"
	"github.com/rivo/tview"
	"github.com/shopspring/decimal"
)

func main() {
	ticks := tplot.NewTicks()

	r := rand.New(rand.NewSource(100))

	size := 100
	data := make([]decimal.Decimal, size)

	for i := 0; i < size; i++ {
		data[i] = decimal.NewFromFloat(r.Float64())
	}

	ticks.SetSpacing(2)
	ticks.SetStyle(tcell.StyleDefault.Foreground(tcell.ColorBlue))
	ticks.SetData(data)

	if err := tview.NewApplication().SetRoot(ticks, true).Run(); err != nil {
		panic(err)
	}
}
