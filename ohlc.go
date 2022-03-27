package tplot

import (
	"fmt"
	"io"
	"math"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// OHLC is a Box component that can render OHLC data.
type OHLC struct {
	*tview.Box

	ohlcAxis *Axis
	volAxis  *Axis

	items   []OHLCItem
	offset  int
	spacing int
	logger  io.Writer
	runes   OLHCRunes

	volFrac float64
}

// NewOHLC creates a new instance of the OHLC component.
func NewOHLC() *OHLC {
	ohlc := &OHLC{
		Box: tview.NewBox(),

		ohlcAxis: NewAxis(),
		volAxis:  NewAxis(),

		runes:   DefaultOHLCRunes,
		spacing: 1,
		volFrac: 0.2,
	}

	return ohlc
}

// SetLogger sets the logger for debugging.
func (o *OHLC) SetLogger(w io.Writer) {
	o.logger = w
}

// Logger returns the current logger set. Used for debugging.
func (o *OHLC) Logger() io.Writer {
	if w := o.logger; w != nil {
		return w
	}

	return nopWriter{}
}

// Offset returns the current offset.
func (o *OHLC) Offset() int {
	return o.offset
}

// AddOffset adds delta to the current offset.
func (o *OHLC) AddOffset(delta int) {
	delta = o.offset + delta

	o.setOffset(delta)
}

// SetOffset sets the scroll offset for OHLC data.
func (o *OHLC) SetOffset(offset int) {
	o.setOffset(offset)
}

// SetRunes sets the runes used to plot the chart.
func (o *OHLC) SetRunes(runes OLHCRunes) {
	o.runes = runes
}

// Runes returns the current set of runes used to plot the chart.
func (o *OHLC) Runes() OLHCRunes {
	return o.runes
}

// setOffset sets the offset, but it ensures it's always less than the size of
// the items and is never negative.
func (o *OHLC) setOffset(offset int) {
	if l := len(o.items); offset >= l {
		offset = l - 1
	}

	if offset < 0 {
		offset = 0
	}

	o.offset = offset
}

// SetItems sets the OHLC data.
func (o *OHLC) SetItems(items OHLCItems) {
	o.items = items
}

// Items returns the current OHLC data.
func (o *OHLC) Items() OHLCItems {
	return o.items
}

// SetSpacing sets the chart spacing.
func (o *OHLC) SetSpacing(spacing int) {
	if spacing <= 0 {
		spacing = 1
	}

	o.spacing = spacing
}

// Spacing returns the current spacing.
func (o *OHLC) Spacing() int {
	return o.spacing
}

// MouseHandler implements tview.Primitive.
func (o *OHLC) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
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
func (o *OHLC) MouseHandler() func(action tview.MouseAction, event *tcell.EventMouse, setFocus func(p tview.Primitive)) (consumed bool, capture tview.Primitive) {
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

func (o *OHLC) volSize() int {
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

func (o *OHLC) ohlcRect() rect {
	x, y, w, h := o.GetInnerRect()

	h -= o.volSize()

	return rect{x: x, y: y, w: w, h: h}
}

func (o *OHLC) volRect() rect {
	x, y, w, h := o.GetInnerRect()
	ohlcRect := o.ohlcRect()

	volSize := o.volSize()

	h = volSize
	y += ohlcRect.h

	return rect{x: x, y: y, w: w, h: h}
}

// Draw implements tview.Primitive.
func (o *OHLC) Draw(screen tcell.Screen) {
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

	ohlcMin, ohlcMax := items.FindValueRange()
	volMin, volMax := items.FindVolumeRange()

	ohlcScale.SetRange(ohlcMin, ohlcMax)
	volScale.SetRange(volMin, volMax)

	o.ohlcAxis.SetScale(ohlcScale)
	o.ohlcAxis.SetStyle(tcell.StyleDefault.Foreground(tcell.ColorDarkCyan))

	o.volAxis.SetScale(volScale)
	o.volAxis.SetStyle(tcell.StyleDefault.Foreground(tcell.ColorDarkBlue))

	drawYAxis := true

	if len(items) > 0 {
		axisYWidth := 0

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

				ohlcMin, ohlcMax = items.FindValueRange()
				volMin, volMax = items.FindVolumeRange()

				ohlcScale.SetRange(ohlcMin, ohlcMax)
				volScale.SetRange(volMin, volMax)
			}
		}
	}

	scaled := newScaledOHLCItems(items, ohlcScale, volScale)
	logger := o.Logger()

	var lastItem *OHLCItem

	if l := len(items); l > 0 {
		lastItem = &items[l-1]
	}

	if len(scaled) > 0 && drawYAxis {
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

		// o.drawAxisY(screen, logger, ohlcScale, ohlcRect, tcell.ColorDarkCyan, lastC)
		// o.drawAxisY(screen, logger, volScale, volRect, tcell.ColorDarkBlue, lastV)
	}

	if width < 0 {
		return
	}

	fmt.Fprintln(logger, "log start", len(scaled))

	runes := o.Runes()

	// We are using special block characters to display quarters so we need
	// to resize our scale after we've drawn the axis.
	numVolFractions := len(runes.Blocks)
	if numVolFractions == 0 {
		numVolFractions = 1
	}

	volScale.SetSize(volRect.h * numVolFractions)

	for i, ohlc := range scaled {
		o, h, l, c := ohlc.O, ohlc.H, ohlc.L, ohlc.C
		v := volScale.Value(items[i].V)

		style := tcell.StyleDefault

		style = style.Foreground(tcell.ColorRed)

		a, b := o, c
		if b >= a {
			style = style.Foreground(tcell.ColorGreen)
			a, b = b, a
		}

		xx := ohlcRect.x + i*spacing + (width - len(scaled)*spacing)

		fmt.Fprintln(logger, "ohlc", i, o, h, l, c)

		for j := h; j >= l; j-- {
			style := style
			yy := ohlcRect.y + ohlcRect.h - j - 1

			isHigh := j == h
			isLow := j == l
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

			fmt.Fprintln(logger, "coord", xx, yy, string(ch))

			screen.SetContent(xx, yy, ch, nil, style)
		}

		if volRect.h > 0 {
			volFullSteps := v / numVolFractions
			volRemFrac := v % numVolFractions

			volStyle := style.Foreground(tcell.ColorDarkBlue)

			fullBlock := '█'

			if ll := len(runes.Blocks); ll > 0 {
				fullBlock = runes.Blocks[ll-1]
			}

			for j := 1; j <= volFullSteps; j++ {
				yy := volRect.y + volRect.h - j
				screen.SetContent(xx, yy, fullBlock, nil, volStyle)
			}

			if volRemFrac > 0 {
				ch := runes.Blocks[volRemFrac-1]

				yy := volRect.y + volRect.h - volFullSteps - 1
				screen.SetContent(xx, yy, ch, nil, volStyle)
			}
		}
	}
}

type rect struct {
	x, y, w, h int
}

type nopWriter struct{}

func (n nopWriter) Write(b []byte) (int, error) {
	return len(b), nil
}
