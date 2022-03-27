package tplot

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/shopspring/decimal"
)

type Bars struct {
	*tview.Box
	scale   Scale
	style   tcell.Style
	spacing int
	rng     Range
	runes   []rune
	data    []decimal.Decimal
}

var DefaultBarRunes = []rune{'▃', '▄', '▆', '█'}

func NewBars() *Bars {
	return &Bars{
		Box:     tview.NewBox(),
		scale:   NewScaleLinear(),
		style:   tcell.StyleDefault,
		spacing: 1,
		runes:   DefaultBarRunes,
	}
}

func (b *Bars) SetSpacing(spacing int) {
	if spacing <= 0 {
		spacing = 1
	}

	b.spacing = spacing
}

func (b *Bars) Spacing() int {
	return b.spacing
}

func (b *Bars) SetStyle(style tcell.Style) {
	b.style = style
}

func (b *Bars) Style() tcell.Style {
	return b.style
}

func (b *Bars) calcRange(values []decimal.Decimal) (rng Range) {
	for _, dec := range values {
		rng = rng.Feed(dec)
	}

	return rng
}

func (b *Bars) SetData(data []decimal.Decimal) {
	b.data = data
	b.rng = b.calcRange(data)
}

func (b *Bars) Values() []decimal.Decimal {
	return b.data
}

func (b *Bars) SetScale(scale Scale) {
	b.scale = scale
}

func (b *Bars) Scale() Scale {
	return b.scale
}

func (b *Bars) SetRunes(runes []rune) {
	b.runes = runes
}

func (b *Bars) Runes() []rune {
	return b.runes
}

func (b *Bars) Draw(screen tcell.Screen) {
	b.DrawForSubclass(screen, b)

	data := b.data
	scale := b.scale
	spacing := b.spacing
	runes := b.runes
	style := b.style
	x, y, w, h := b.GetInnerRect()

	if h == 0 || w == 0 {
		return
	}

	scale = scale.Copy()

	maxCount := w / spacing

	if l := len(data); l > maxCount {
		data = data[l-maxCount:]
	}

	if len(runes) == 0 {
		runes = []rune{'█'}
	}

	// We are using special block characters to display quarters so we need
	// to resize our scale after we've drawn the axis.
	numVolFractions := len(b.runes)

	scale.SetRange(b.rng)
	scale.SetSize(h * numVolFractions)

	for i, dec := range data {
		v := scale.Value(dec)
		xx := x + i*spacing + (w - len(data)*spacing)

		volFullSteps := v / numVolFractions
		volRemFrac := v % numVolFractions

		fullBlock := runes[len(runes)-1]

		for j := 1; j <= volFullSteps; j++ {
			yy := y + h - j
			screen.SetContent(xx, yy, fullBlock, nil, style)
		}

		if volRemFrac > 0 {
			ch := runes[volRemFrac-1]

			yy := y + h - volFullSteps - 1
			screen.SetContent(xx, yy, ch, nil, style)
		}
	}
}
