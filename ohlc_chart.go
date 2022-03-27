package tplot

import (
	"fmt"
	"io"
	"math"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/shopspring/decimal"
)

// OHLCChart is a Box component that can render OHLCChart data.
type OHLCChart struct {
	*tview.Box

	ohlcCandles *OHLCCandles
	ohlcAxis    *Axis

	volBars *Bars
	volAxis *Axis

	items     []OHLCItem
	offset    int
	spacing   int
	logger    io.Writer
	ohlcRunes OHLCRunes
	volRunes  []rune

	volFrac float64
}

// NewOHLCChart creates a new instance of the OHLC component.
func NewOHLCChart() *OHLCChart {
	ohlc := &OHLCChart{
		Box: tview.NewBox(),

		ohlcCandles: NewOHLCCandles(),
		ohlcAxis:    NewAxis(),

		volBars: NewBars(),
		volAxis: NewAxis(),

		ohlcRunes: DefaultOHLCRunes,
		volRunes:  DefaultBarRunes,
		spacing:   1,
		volFrac:   0.2,
	}

	return ohlc
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

// AddOffset adds delta to the current offset.
func (o *OHLCChart) AddOffset(delta int) {
	delta = o.offset + delta

	o.setOffset(delta)
}

// SetOffset sets the scroll offset for OHLC data.
func (o *OHLCChart) SetOffset(offset int) {
	o.setOffset(offset)
}

// SetRunes sets the runes used to plot the chart.
func (o *OHLCChart) SetRunes(ohlcRunes OHLCRunes, volRunes []rune) {
	o.ohlcRunes = ohlcRunes
	o.volRunes = volRunes
}

// Runes returns the current set of runes used to plot the chart.
func (o *OHLCChart) Runes() (ohlcRunes OHLCRunes, volRunes []rune) {
	return o.ohlcRunes, o.volRunes
}

// setOffset sets the offset, but it ensures it's always less than the size of
// the items and is never negative.
func (o *OHLCChart) setOffset(offset int) {
	if l := len(o.items); offset >= l {
		offset = l - 1
	}

	if offset < 0 {
		offset = 0
	}

	o.offset = offset
}

// SetItems sets the OHLC data.
func (o *OHLCChart) SetItems(items OHLCItems) {
	o.items = items
}

// Items returns the current OHLC data.
func (o *OHLCChart) Items() OHLCItems {
	return o.items
}

// SetSpacing sets the chart spacing.
func (o *OHLCChart) SetSpacing(spacing int) {
	if spacing <= 0 {
		spacing = 1
	}

	o.spacing = spacing
}

// Spacing returns the current spacing.
func (o *OHLCChart) Spacing() int {
	return o.spacing
}

// MouseHandler implements tview.Primitive.
func (o *OHLCChart) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return o.WrapInputHandler(
		func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
			switch event.Key() {
			case tcell.KeyEnd:
				o.SetOffset(0)
			case tcell.KeyHome:
				s := len(o.Items())
				o.SetOffset(s - 1)
			case tcell.KeyPgUp:
				o.AddOffset(20)
			case tcell.KeyPgDn:
				o.AddOffset(-20)
			case tcell.KeyLeft:
				o.AddOffset(1)
			case tcell.KeyRight:
				o.AddOffset(-1)

			case tcell.KeyRune:
				switch event.Rune() {
				case 'h':
					o.AddOffset(1)
				case 'l':
					o.AddOffset(-1)
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

	v := int(math.Floor(float64(h) * o.volFrac))

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
	x, y, w, h := o.GetInnerRect()
	ohlcRect := o.ohlcRect()

	volSize := o.volSize()

	h = volSize
	y += ohlcRect.h

	return rect{x: x, y: y, w: w, h: h}
}

// Draw implements tview.Primitive.
func (o *OHLCChart) Draw(screen tcell.Screen) {
	ohlcRect := o.ohlcRect()
	volRect := o.volRect()
	ohlcScale := NewScaleLinear()
	volScale := NewScaleLinear()
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

	ohlcRange := items.FindOHLCRange()
	ohlcScale.SetRange(ohlcRange)

	volRange := items.FindVolumeRange()
	volScale.SetRange(volRange)

	o.ohlcAxis.SetScale(ohlcScale)
	o.ohlcAxis.SetStyle(tcell.StyleDefault.Foreground(tcell.ColorDarkCyan))

	o.volAxis.SetScale(volScale)
	o.volAxis.SetStyle(tcell.StyleDefault.Foreground(tcell.ColorDarkBlue))

	drawYAxis := true
	axisYWidth := 0

	if len(items) > 0 {
		w1 := o.ohlcAxis.CalcWidth()
		w2 := o.volAxis.CalcWidth()

		if w2 > w1 {
			axisYWidth = w2
		} else {
			axisYWidth = w1
		}

		drawYAxis = width-axisYWidth > 0

		if drawYAxis {
			width -= axisYWidth

			o.ohlcAxis.SetRect(ohlcRect.x+width, ohlcRect.y, axisYWidth, ohlcRect.h)
			o.volAxis.SetRect(volRect.x+width, volRect.y, axisYWidth, volRect.h)

			// We need to readjust the maxCount after taking account the axis width.
			maxCount = width / spacing

			// If we didn't have enough space, we need to make the slice smaller and
			// find min/max again.
			if l := len(items); l > maxCount {
				items = items[l-maxCount:]

				ohlcRange = items.FindOHLCRange()
				ohlcScale.SetRange(ohlcRange)

				volRange = items.FindVolumeRange()
				volScale.SetRange(volRange)
			}
		}
	}

	logger := o.Logger()

	var lastItem *OHLCItem

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

		o.volAxis.SetHighlight(lastV)
		o.volAxis.Draw(screen)
	}

	if width < 0 {
		return
	}

	fmt.Fprintln(logger, "log start", len(items))

	o.ohlcCandles.SetPositiveStyle(tcell.StyleDefault.Foreground(tcell.ColorGreen))
	o.ohlcCandles.SetNegativeStyle(tcell.StyleDefault.Foreground(tcell.ColorRed))
	o.ohlcCandles.SetRect(ohlcRect.x, ohlcRect.y, width, ohlcRect.h)
	o.ohlcCandles.SetRunes(o.ohlcRunes)
	o.ohlcCandles.SetScale(ohlcScale)
	o.ohlcCandles.SetOHLCItems(items)
	o.ohlcCandles.Draw(screen)

	volValues := make([]decimal.Decimal, len(items))

	for i, item := range items {
		volValues[i] = item.V
	}

	o.volBars.SetStyle(tcell.StyleDefault.Foreground(tcell.ColorDarkBlue))
	o.volBars.SetRect(volRect.x, volRect.y, width, volRect.h)
	o.volBars.SetRunes(o.volRunes)
	o.volBars.SetScale(volScale)
	o.volBars.SetValues(volValues)
	o.volBars.Draw(screen)
}

type rect struct {
	x, y, w, h int
}

type nopWriter struct{}

func (n nopWriter) Write(b []byte) (int, error) {
	return len(b), nil
}
