package tplot

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/shopspring/decimal"
)

type base struct {
	*tview.Box
	scale       Scale
	style       tcell.Style
	spacing     int
	rng         Range
	runes       []rune
	data        []decimal.Decimal
	sliceMethod SliceMethod
}

// SliceMethod describes the slicing method when the number of items in the
// slice is too large to be rendered on the screen.
type SliceMethod int

const (
	// Last means only the last items will be kept in the slice.
	Last SliceMethod = iota
	// First means only the first items will be kept in the slice.
	First
)

func newBase(runes []rune) *base {
	return &base{
		Box:     tview.NewBox(),
		scale:   NewScaleLinear(),
		spacing: 1,
		runes:   runes,
	}
}

func (b *base) SetSliceMethod(method SliceMethod) {
	b.sliceMethod = method
}

func (b *base) SliceMethod() SliceMethod {
	return b.sliceMethod
}

func (b *base) SetSpacing(spacing int) {
	if spacing <= 0 {
		spacing = 1
	}

	b.spacing = spacing
}

func (b *base) Spacing() int {
	return b.spacing
}

func (b *base) SetStyle(style tcell.Style) {
	b.style = style
}

func (b *base) Style() tcell.Style {
	return b.style
}

func (b *base) calcRange(values []decimal.Decimal) (rng Range) {
	for _, dec := range values {
		rng = rng.Feed(dec)
	}

	return rng
}

func (b *base) Data() []decimal.Decimal {
	return b.data
}

// DataSlice returns data, but only the items that
// fit on the screen.
func (b *base) DataSlice() []decimal.Decimal {
	_, _, w, _ := b.GetInnerRect()
	maxCount := w / b.spacing
	data := b.data

	if l := len(data); l > maxCount {
		if b.sliceMethod == Last {
			data = data[l-maxCount:]
		} else {
			data = data[:maxCount]
		}
	}

	return data
}

func (b *base) SetData(data []decimal.Decimal) {
	b.data = data
	b.rng = b.calcRange(data)
}

func (b *base) Values() []decimal.Decimal {
	return b.data
}

func (b *base) SetScale(scale Scale) {
	b.scale = scale
}

func (b *base) Scale() Scale {
	return b.scale
}

func (b *base) SetRunes(runes []rune) {
	b.runes = runes
}

func (b *base) Runes() []rune {
	return b.runes
}
