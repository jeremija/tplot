package tplot

import (
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type OHLCCandles struct {
	*tview.Box
	scale Scale

	factory DecimalFactory

	data    []OHLC
	rng     Range
	spacing int

	positiveStyle tcell.Style
	negativeStyle tcell.Style

	runes OHLCRunes
}

func NewOHLCCandles(factory DecimalFactory) *OHLCCandles {
	style := tcell.StyleDefault

	return &OHLCCandles{
		Box:           tview.NewBox(),
		factory:       factory,
		scale:         NewScaleLinear(factory),
		spacing:       1,
		negativeStyle: style.Foreground(tcell.ColorRed),
		positiveStyle: style.Foreground(tcell.ColorGreen),
		runes:         DefaultOHLCRunes,
		rng:           NewRange(factory),
	}
}

func (o *OHLCCandles) SetSpacing(spacing int) {
	if spacing <= 0 {
		spacing = 1
	}

	o.spacing = spacing
}

func (o *OHLCCandles) Spacing() int {
	return o.spacing
}

func (o *OHLCCandles) SetPositiveStyle(positiveStyle tcell.Style) {
	o.positiveStyle = positiveStyle
}

func (o *OHLCCandles) PositiveStyle() tcell.Style {
	return o.positiveStyle
}

func (o *OHLCCandles) SetNegativeStyle(negativeStyle tcell.Style) {
	o.negativeStyle = negativeStyle
}

func (o *OHLCCandles) NegativeStyle() tcell.Style {
	return o.negativeStyle
}

func (o *OHLCCandles) SetScale(scale Scale) {
	o.scale = scale
}

func (o *OHLCCandles) Scale() Scale {
	return o.scale
}

// SetRunes sets the runes used to plot the chart.
func (o *OHLCCandles) SetRunes(runes OHLCRunes) {
	o.runes = runes
}

// Runes returns the current set of runes used to plot the chart.
func (o *OHLCCandles) Runes() OHLCRunes {
	return o.runes
}

func (o *OHLCCandles) calcRange(items []OHLC) Range {
	rng := NewRange(o.factory)

	for _, item := range items {
		rng = rng.Feed(item.L)
		rng = rng.Feed(item.H)
	}

	return rng
}

func (o *OHLCCandles) SetData(data []OHLC) {
	o.data = data
	o.rng = o.calcRange(data)
}

func (o *OHLCCandles) Draw(screen tcell.Screen) {
	o.Box.DrawForSubclass(screen, o)

	x, y, w, h := o.GetInnerRect()
	scale := o.scale
	runes := o.runes
	data := o.data
	spacing := o.spacing
	maxCount := w / spacing

	scale.SetSize(h)
	scale.SetRange(o.rng)

	if h == 0 || w == 0 {
		return
	}

	type scaledOHLC struct {
		ts         time.Time
		O, H, L, C int
	}

	if l := len(data); l > maxCount {
		data = data[l-maxCount:]

		scale.SetRange(o.calcRange(data))
	}

	scaled := make([]scaledOHLC, len(data))

	for i, item := range data {
		scaled[i] = scaledOHLC{
			O:  scale.Value(item.O),
			H:  scale.Value(item.H),
			L:  scale.Value(item.L),
			C:  scale.Value(item.C),
			ts: item.Timestamp,
		}
	}

	for i, ohlc := range scaled {
		open, high, low, cl := ohlc.O, ohlc.H, ohlc.L, ohlc.C

		style := tcell.StyleDefault

		style = style.Foreground(tcell.ColorRed)

		a, b := open, cl
		if b >= a {
			style = style.Foreground(tcell.ColorGreen)
			a, b = b, a
		}

		xx := x + i*spacing + (w - len(scaled)*spacing)

		for j := high; j >= low; j-- {
			style := style
			yy := y + h - j - 1

			isHigh := j == high
			isLow := j == low
			isOpen := j == a
			isClose := j == b
			isThick := j < a && j > b

			var ch rune

			switch {
			case isHigh && isLow && isOpen && isClose:
				ch = runes.Same
			case isHigh && isOpen && isClose:
				ch = runes.HighOpenClose
			case isLow && isOpen && isClose:
				ch = runes.LowOpenClose
			case isHigh && isOpen:
				style = style.Bold(true)
				ch = runes.HighOpen
			case isLow && isClose:
				style = style.Bold(true)
				ch = runes.LowClose
			case isHigh:
				ch = runes.High
			case isLow:
				ch = runes.Low
			case isOpen && isClose:
				ch = runes.OpenClose
			case isOpen:
				ch = runes.Open
			case isClose:
				ch = runes.Close
			case isThick:
				style = style.Bold(true)
				ch = runes.Thick
			default:
				ch = runes.Thin
			}

			screen.SetContent(xx, yy, ch, nil, style)
		}
	}
}
