package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/jeremija/tplot"
	"github.com/rivo/tview"
)

func main() {
	var factory tplot.FloatFactory

	bars := tplot.NewBars(factory)
	axis := tplot.NewAxis(factory)

	bars.SetSliceMethod(tplot.Last)

	bars.SetSpacing(2)
	bars.SetStyle(tcell.StyleDefault.Foreground(tcell.ColorOrange))
	axis.SetStyle(tcell.StyleDefault.Foreground(tcell.ColorRed))

	barsWithAxis := tplot.NewAxisBox(axis, bars)
	barsWithAxis.SetPosition(tplot.Right)

	size := 101
	data := make([]tplot.Decimal, size)

	for i := 0; i < size; i++ {
		data[i] = tplot.Float(int64(i))
	}

	axis.SetHighlight(tplot.DecimalValue{
		Decimal: data[size-1],
		Valid:   true,
	})

	bars.SetData(data)

	if err := tview.NewApplication().SetRoot(barsWithAxis, true).Run(); err != nil {
		panic(err)
	}
}
