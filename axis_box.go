package tplot

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type AxisBox struct {
	*tview.Box

	axis     *Axis
	content  Primitive
	scale    Scale
	position Position
}

func NewAxisBox(axis *Axis, content Primitive) *AxisBox {
	return &AxisBox{
		Box:      tview.NewBox(),
		axis:     axis,
		content:  content,
		position: Left,
	}
}

func (a *AxisBox) SetPosition(position Position) {
	a.position = position
}

func (a *AxisBox) Position() Position {
	return a.position
}

func (a *AxisBox) Draw(screen tcell.Screen) {
	a.Box.DrawForSubclass(screen, a)

	scale := a.content.Scale()

	a.axis.SetScale(scale)

	x, y, w, h := a.Box.GetInnerRect()

	scale.SetSize(h)

	axisW := a.axis.CalcWidth()
	barsW := w - axisW

	// Hide axis when there's no space
	if barsW < 0 {
		axisW = 0
		barsW = w
	}

	// We need to draw the content first because the scale.Range
	// might change.
	if a.position == Right {
		a.content.SetRect(x, y, barsW, h)
		a.axis.SetRect(x+barsW, y, axisW, h)
	} else {
		a.content.SetRect(x+axisW, y, barsW, h)
		a.axis.SetRect(x, y, axisW, h)
	}

	a.content.Draw(screen)
	a.axis.Draw(screen)
}
