package main

import (
	"math/rand"

	"github.com/gdamore/tcell/v2"
	"github.com/jeremija/tplot"
	"github.com/rivo/tview"
)

func main() {
	var factory tplot.FloatFactory

	bars := tplot.NewBars(factory)

	r := rand.New(rand.NewSource(100))

	size := 100
	data := make([]tplot.Decimal, size)

	for i := 0; i < size; i++ {
		data[i] = tplot.Float(r.Float64())
	}

	bars.SetSpacing(2)
	bars.SetStyle(tcell.StyleDefault.Foreground(tcell.ColorBlue))
	bars.SetData(data)

	if err := tview.NewApplication().SetRoot(bars, true).Run(); err != nil {
		panic(err)
	}
}
