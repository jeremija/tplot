package main

import (
	"math/rand"

	"github.com/gdamore/tcell/v2"
	"github.com/jeremija/tplot"
	"github.com/rivo/tview"
	"github.com/shopspring/decimal"
)

type BarWithAxis struct {
	*tview.Box

	*tplot.Axis
	*tplot.Bars
}

func NewBarsWithAxis() *BarWithAxis {
	return &BarWithAxis{
		Box:  tview.NewBox(),
		Bars: tplot.NewBars(),
		Axis: tplot.NewAxis(),
	}
}

func (b *BarWithAxis) Draw(screen tcell.Screen) {
	b.Box.DrawForSubclass(screen, b)

	scale := tplot.NewScaleLinear()

	b.Bars.SetScale(scale)
	b.Axis.SetScale(scale)

	data := b.Bars.Data()

	var rng tplot.Range

	for _, value := range data {
		rng = rng.Feed(value)
	}

	x, y, w, h := b.GetInnerRect()

	scale.SetSize(h)
	scale.SetRange(rng)

	axisW := b.Axis.CalcWidth()
	barsW := w - axisW

	// Hide axis when there's no space
	if barsW < 0 {
		axisW = 0
		barsW = w
	}

	var highlight tplot.DecimalValue

	if l := len(data); l > 0 {
		highlight.Decimal = data[l-1]
		highlight.Valid = true
	}

	b.Axis.SetHighlight(highlight)

	b.Axis.SetRect(x+barsW, y, axisW, h)
	b.Axis.Draw(screen)

	b.Bars.SetRect(x, y, barsW, h)
	b.Bars.Draw(screen)
}

func main() {
	barsWithAxis := NewBarsWithAxis()

	r := rand.New(rand.NewSource(100))

	size := 100
	data := make([]decimal.Decimal, size)

	for i := 0; i < size; i++ {
		data[i] = decimal.NewFromFloat(r.Float64())
	}

	barsWithAxis.SetSpacing(2)
	barsWithAxis.Bars.SetStyle(tcell.StyleDefault.Foreground(tcell.ColorBlue))
	barsWithAxis.Axis.SetStyle(tcell.StyleDefault.Foreground(tcell.ColorDarkBlue))
	barsWithAxis.SetData(data)

	if err := tview.NewApplication().SetRoot(barsWithAxis, true).Run(); err != nil {
		panic(err)
	}
}
