package tplot

import (
	"github.com/gdamore/tcell/v2"
)

type Ticks struct {
	*base
}

var DefaultTicksRunes = []rune{'⎽', '⎼', '—', '⎻', '⎺'}

func NewTicks() *Ticks {
	bars := NewBars()

	bars.SetRunes(DefaultTicksRunes)

	return &Ticks{
		base: newBase(DefaultTicksRunes),
	}
}

func (b *Ticks) Draw(screen tcell.Screen) {
	b.DrawForSubclass(screen, b)

	data := b.DataSlice()
	scale := b.scale
	spacing := b.spacing
	runes := b.runes
	style := b.style
	x, y, w, h := b.GetInnerRect()

	if h == 0 || w == 0 {
		return
	}

	if len(runes) == 0 {
		runes = []rune{'-'}
	}

	// We are using special block characters to display quarters so we need
	// to resize our scale after we've drawn the axis.
	numFractions := len(b.runes)

	rng := b.calcRange(data)
	scale.SetRange(rng)
	// If we're sharing the scale with other components that can't use the
	// fractions.
	scale = scale.Copy()
	scale.SetSize(h * numFractions)

	for i, dec := range data {
		v := scale.Value(dec)
		xx := x + i*spacing + (w - len(data)*spacing)

		fullSteps := v / numFractions
		rem := v % numFractions

		ch := runes[rem]

		if rem > 0 {
			yy := y + h - fullSteps - 1
			screen.SetContent(xx, yy, ch, nil, style)
		}
	}
}
