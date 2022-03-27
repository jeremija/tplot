package tplot

import (
	"fmt"
	"io"
	"math"
	"strconv"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/shopspring/decimal"
)

// OHLC is a Box component that can render OHLC data.
type OHLC struct {
	*tview.Box

	items   []OHLCItem
	offset  int
	spacing int
	logger  io.Writer
	runes   OLHCRunes

	volFrac float64
}

// NewOHLC creates a new instance of the OHLC component.
func NewOHLC() *OHLC {
	return &OHLC{
		Box:     tview.NewBox(),
		runes:   DefaultOHLCRunes,
		spacing: 1,
		volFrac: 0.2,
	}
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

type value struct {
	decimal decimal.Decimal
	valid   bool
}

func (o *OHLC) getAxisWidth(scale Scale) int {
	// numDecs contains the number of decimals spot to represent the axis.
	numDecs := scale.NumDecimals() + 2

	_, max := scale.Range()

	size := len(max.Round(0).String()) + 1 + numDecs

	return size
}

// drawOHLCAxisY draws the Y axis.
func (o *OHLC) drawAxisY(screen tcell.Screen, log io.Writer, scale Scale, r rect, color tcell.Color, val value) int {
	var current int

	if val.valid {
		current = scale.Value(val.decimal)
	}

	// numDecs contains the number of decimals spot to represent the axis.
	numDecs := scale.NumDecimals() + 2

	maxWidth := 0

	vals := make([]string, r.h)

	for i := 0; i < r.h; i++ {
		rev := scale.Reverse(i)

		if val.valid && i == current {
			rev = val.decimal
		}

		valFloat, _ := rev.Float64()

		valStr := strconv.FormatFloat(valFloat, 'f', numDecs, 64)

		if l := len(valStr); l > maxWidth {
			maxWidth = l
		}

		fmt.Fprintln(log, "value", i, string(valStr))

		vals[i] = valStr
	}

	if r.w < maxWidth {
		return 0 // hide axis when no room.
	}

	style := tcell.StyleDefault.Foreground(color)

	for i, valStr := range vals {
		yy := r.y + r.h - i - 1

		currentStyle := style

		if val.valid && i == current {
			currentStyle = tcell.StyleDefault
		}

		for i, val := range valStr {
			xx := r.x + r.w - len(valStr) + i

			fmt.Fprintln(log, "axis y", xx, yy, string(val))
			screen.SetContent(xx, yy, val, nil, currentStyle)
		}
	}

	return maxWidth
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

	drawYAxis := true

	if len(items) > 0 {
		axisYWidth := 0

		w1 := o.getAxisWidth(ohlcScale)
		w2 := o.getAxisWidth(volScale)

		if w2 > w1 {
			axisYWidth = w2
		} else {
			axisYWidth = w1
		}

		drawYAxis = width-axisYWidth > 0

		if drawYAxis {
			width -= axisYWidth

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
		lastC := value{}
		lastV := value{}

		if lastItem != nil {
			lastC.decimal = lastItem.C
			lastC.valid = true
			lastV.decimal = lastItem.V
			lastV.valid = true
		}

		o.drawAxisY(screen, logger, ohlcScale, ohlcRect, tcell.ColorDarkCyan, lastC)
		o.drawAxisY(screen, logger, volScale, volRect, tcell.ColorDarkBlue, lastV)
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
