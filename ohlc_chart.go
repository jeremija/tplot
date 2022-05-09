package tplot

import (
	"fmt"
	"io"
	"math"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// OHLCChart is a Box component that can render OHLCChart data.
type OHLCChart struct {
	*tview.Box

	factory DecimalFactory

	ohlcCandles *OHLCCandles
	ohlcAxis    *Axis

	volumeBars *Bars
	volumeAxis *Axis

	items  []OHLC
	offset int
	logger io.Writer

	// volumeHeightFraction is a number between 0 and 1 that determines how much
	// of the layout should be taken by the OHLC chart.
	volumeHeightFraction float64
}

// NewOHLCChart creates a new instance of the OHLC component.
func NewOHLCChart(factory DecimalFactory) *OHLCChart {
	ohlc := &OHLCChart{
		Box: tview.NewBox(),

		factory: factory,

		ohlcCandles: NewOHLCCandles(factory),
		ohlcAxis:    NewAxis(factory),

		volumeBars: NewBars(factory),
		volumeAxis: NewAxis(factory),

		volumeHeightFraction: 0.2,
	}

	ohlc.SetVolumeBarsStyle(tcell.StyleDefault.Foreground(tcell.ColorDarkBlue))

	return ohlc
}

func (o *OHLCChart) SetVolumeHeight(fraction float64) {
	if fraction < 0 {
		fraction = 0
	}

	if fraction > 1 {
		fraction = 1
	}

	o.volumeHeightFraction = fraction
}

func (o *OHLCChart) SetVolumeBarsRunes(runes []rune) {
	o.volumeBars.SetRunes(runes)
}

func (o *OHLCChart) VolumeBarsRunes() []rune {
	return o.volumeBars.Runes()
}

func (o *OHLCChart) SetOHLCCandlesRunes(runes OHLCRunes) {
	o.ohlcCandles.SetRunes(runes)
}

func (o *OHLCChart) OHLCCandlesRunes() OHLCRunes {
	return o.ohlcCandles.Runes()
}

func (o *OHLCChart) SetPositiveStyle(positiveStyle tcell.Style) {
	o.ohlcCandles.SetPositiveStyle(positiveStyle)
}

func (o *OHLCChart) PositiveStyle() tcell.Style {
	return o.ohlcCandles.PositiveStyle()
}

func (o *OHLCChart) SetNegativeStyle(negativeStyle tcell.Style) {
	o.ohlcCandles.SetNegativeStyle(negativeStyle)
}

func (o *OHLCChart) NegativeStyle() tcell.Style {
	return o.ohlcCandles.NegativeStyle()
}

func (o *OHLCChart) SetVolumeBarsStyle(style tcell.Style) {
	o.volumeBars.SetStyle(style)
}

func (o *OHLCChart) VolumeBarsStyle() tcell.Style {
	return o.volumeBars.Style()
}

func (o *OHLCChart) SetVolumeAxisStyle(style tcell.Style) {
	o.volumeAxis.SetStyle(style)
}

func (o *OHLCChart) VolumeAxisStyle() tcell.Style {
	return o.volumeAxis.Style()
}

func (o *OHLCChart) SetOHLCAxisStyle(style tcell.Style) {
	o.ohlcAxis.SetStyle(style)
}

func (o *OHLCChart) OHLCAxisStyle() tcell.Style {
	return o.ohlcAxis.Style()
}

// SetLogger sets the logger for debugging.
func (o *OHLCChart) SetLogger(w io.Writer) {
	o.logger = w
}

// Logger returns the current logger set. Used for debugging.
func (o *OHLCChart) Logger() io.Writer {
	if w := o.logger; w != nil {
		return w
	}

	return nopWriter{}
}

// Offset returns the current offset.
func (o *OHLCChart) Offset() int {
	return o.offset
}

// SetOffset sets the scroll offset for OHLC data. It ensures it's always less
// than the size of the items and is never negative.
func (o *OHLCChart) SetOffset(offset int) {
	if l := len(o.items); offset >= l {
		offset = l - 1
	}

	if offset < 0 {
		offset = 0
	}

	o.offset = offset
}

func (o *OHLCChart) AddOffset(delta int) {
	o.SetOffset(o.offset + delta)
}

// SetItems sets the OHLC data.
func (o *OHLCChart) SetItems(items []OHLC) {
	o.items = items
}

// Items returns the current OHLC data.
func (o *OHLCChart) Items() []OHLC {
	return o.items
}

// SetSpacing sets the chart spacing.
func (o *OHLCChart) SetSpacing(spacing int) {
	o.ohlcCandles.SetSpacing(spacing)
	o.volumeBars.SetSpacing(spacing)
}

func (o *OHLCChart) AddSpacing(delta int) {
	o.ohlcCandles.SetSpacing(o.ohlcCandles.Spacing() + delta)
	o.volumeBars.SetSpacing(o.volumeBars.Spacing() + delta)
}

// Spacing returns the current spacing.
func (o *OHLCChart) Spacing() int {
	return o.ohlcCandles.Spacing()
}

// MouseHandler implements tview.Primitive.
func (o *OHLCChart) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	moveHome := func() {
		s := len(o.Items())
		o.SetOffset(s - 1)
	}

	moveEnd := func() {
		o.SetOffset(0)
	}

	moveLeft := func() {
		o.AddOffset(1)
	}

	moveRight := func() {
		o.AddOffset(-1)
	}

	moveLeftLong := func() {
		o.SetOffset(o.offset + 20)
	}

	moveRightLong := func() {
		o.SetOffset(o.offset - 20)
	}

	return o.WrapInputHandler(
		func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
			switch event.Key() {
			case tcell.KeyEnd:
				moveEnd()
			case tcell.KeyHome:
				moveHome()
			case tcell.KeyPgUp:
				moveLeftLong()
			case tcell.KeyPgDn:
				moveRightLong()
			case tcell.KeyLeft:
				moveLeft()
			case tcell.KeyRight:
				moveRight()

			case tcell.KeyRune:
				if event.Modifiers()&tcell.ModAlt > 0 {
					switch event.Rune() {
					case 'b':
						moveLeftLong()
					case 'f':
						moveRightLong()

						return
					}
				}

				switch event.Rune() {
				case '0':
					o.SetSpacing(1)
				case '=':
					o.AddSpacing(1)
				case '-':
					o.AddSpacing(-1)
				case 'b':
					moveLeftLong()
				case 'w', 'e':
					moveRightLong()
				case 'g', '^':
					moveHome()
				case 'G', '$':
					moveEnd()
				case 'h':
					moveLeft()
				case 'l':
					moveRight()
				}
			}
		},
	)
}

// MouseHandler implements tview.Primitive.
func (o *OHLCChart) MouseHandler() func(action tview.MouseAction, event *tcell.EventMouse, setFocus func(p tview.Primitive)) (consumed bool, capture tview.Primitive) {
	return o.Box.WrapMouseHandler(func(action tview.MouseAction, ev *tcell.EventMouse, setFocus func(p tview.Primitive)) (bool, tview.Primitive) {
		switch action {
		case tview.MouseScrollUp:
			o.AddOffset(10)

			return true, o
		case tview.MouseScrollDown:
			o.AddOffset(-10)

			return true, o
		}

		return false, o
	})
}

func (o *OHLCChart) volSize() int {
	_, _, _, h := o.GetInnerRect()

	v := int(math.Floor(float64(h) * o.volumeHeightFraction))

	if v < 0 {
		return 0
	}

	if v > h {
		return h
	}

	return v
}

func (o *OHLCChart) ohlcRect() rect {
	x, y, w, h := o.GetInnerRect()

	h -= o.volSize()

	return rect{x: x, y: y, w: w, h: h}
}

func (o *OHLCChart) volRect() rect {
	x, y, w, _ := o.GetInnerRect()
	ohlcRect := o.ohlcRect()

	volSize := o.volSize()

	h := volSize
	y += ohlcRect.h

	return rect{x: x, y: y, w: w, h: h}
}

func (o *OHLCChart) ohlcRange(items []OHLC) Range {
	rng := NewRange(o.factory)

	for _, ohlc := range items {
		rng = rng.Feed(ohlc.L)
		rng = rng.Feed(ohlc.H)
	}

	return rng
}

func (o *OHLCChart) volumeRange(items []OHLC) Range {
	rng := NewRange(o.factory)

	for _, ohlc := range items {
		rng = rng.Feed(ohlc.V)
	}

	return rng
}

// Draw implements tview.Primitive.
func (o *OHLCChart) Draw(screen tcell.Screen) {
	ohlcRect := o.ohlcRect()
	volRect := o.volRect()
	ohlcScale := NewScaleLinear(o.factory)
	volScale := NewScaleLinear(o.factory)
	spacing := o.Spacing()
	items := o.Items()
	offset := o.Offset()

	if l := len(items); offset > l {
		items = nil
	} else {
		items = items[:l-offset]
	}

	if l := len(items); l > 0 {
		ohlc := items[l-1]

		title := fmt.Sprintf(" O=%s H=%s L=%s C=%s V=%s TS=%s ", ohlc.O, ohlc.H, ohlc.L, ohlc.C, ohlc.V, ohlc.Timestamp.Format("2006-01-02T15:04:05"))

		o.SetTitle(title)
	}

	o.DrawForSubclass(screen, o)

	width := ohlcRect.w

	maxCount := width / spacing

	if l := len(items); l > maxCount {
		items = items[l-maxCount:]
	}

	ohlcScale.SetSize(ohlcRect.h)
	volScale.SetSize(volRect.h)

	ohlcRange := o.ohlcRange(items)
	ohlcScale.SetRange(ohlcRange)

	volRange := o.volumeRange(items)
	volScale.SetRange(volRange)

	o.ohlcAxis.SetScale(ohlcScale)
	o.ohlcAxis.SetStyle(tcell.StyleDefault.Foreground(tcell.ColorDarkCyan))

	o.volumeAxis.SetScale(volScale)
	o.volumeAxis.SetStyle(tcell.StyleDefault.Foreground(tcell.ColorDarkBlue))

	drawYAxis := true
	axisYWidth := 0

	if len(items) > 0 {
		w1 := o.ohlcAxis.CalcWidth()
		w2 := o.volumeAxis.CalcWidth()

		if w2 > w1 {
			axisYWidth = w2
		} else {
			axisYWidth = w1
		}

		widthWithoutYAxis := width - axisYWidth

		drawYAxis = widthWithoutYAxis > 0

		if drawYAxis {
			width = widthWithoutYAxis

			o.ohlcAxis.SetRect(ohlcRect.x+width, ohlcRect.y, axisYWidth, ohlcRect.h)
			o.volumeAxis.SetRect(volRect.x+width, volRect.y, axisYWidth, volRect.h)

			// We need to readjust the maxCount after taking account the axis width.
			maxCount = width / spacing

			// If we didn't have enough space, we need to make the slice smaller and
			// find min/max again.
			if l := len(items); l > maxCount {
				items = items[l-maxCount:]

				ohlcRange = o.ohlcRange(items)
				ohlcScale.SetRange(ohlcRange)

				volRange = o.volumeRange(items)
				volScale.SetRange(volRange)
			}
		}
	}

	var lastItem *OHLC

	if l := len(items); l > 0 {
		lastItem = &items[l-1]
	}

	if len(items) > 0 && drawYAxis {
		lastC := DecimalValue{}
		lastV := DecimalValue{}

		if lastItem != nil {
			lastC.Decimal = lastItem.C
			lastC.Valid = true
			lastV.Decimal = lastItem.V
			lastV.Valid = true
		}

		o.ohlcAxis.SetHighlight(lastC)
		o.ohlcAxis.Draw(screen)

		o.volumeAxis.SetHighlight(lastV)
		o.volumeAxis.Draw(screen)
	}

	if width < 0 {
		return
	}

	o.ohlcCandles.SetRect(ohlcRect.x, ohlcRect.y, width, ohlcRect.h)
	o.ohlcCandles.SetScale(ohlcScale)
	o.ohlcCandles.SetData(items)
	o.ohlcCandles.Draw(screen)

	volValues := make([]Decimal, len(items))

	for i, item := range items {
		volValues[i] = item.V
	}

	o.volumeBars.SetRect(volRect.x, volRect.y, width, volRect.h)
	o.volumeBars.SetScale(volScale)
	o.volumeBars.SetData(volValues)
	o.volumeBars.Draw(screen)
}

type rect struct {
	x, y, w, h int
}

type nopWriter struct{}

func (n nopWriter) Write(b []byte) (int, error) {
	return len(b), nil
}
