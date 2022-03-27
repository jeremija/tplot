package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/jeremija/tplot"
	"github.com/rivo/tview"
	"github.com/shopspring/decimal"
)

func main() {
	bars := tplot.NewBars()
	axis := tplot.NewAxis()

	bars.SetSliceMethod(tplot.Last)

	bars.SetSpacing(2)
	bars.SetStyle(tcell.StyleDefault.Foreground(tcell.ColorDarkRed))
	axis.SetStyle(tcell.StyleDefault.Foreground(tcell.ColorDarkRed))

	barsWithAxis := tplot.NewAxisBox(axis, bars)
	barsWithAxis.SetPosition(tplot.Right)

	size := 101
	data := make([]decimal.Decimal, size)

	for i := 0; i < size; i++ {
		data[i] = decimal.NewFromInt(int64(i))
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
