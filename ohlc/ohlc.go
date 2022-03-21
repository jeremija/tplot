package ohlc

import (
	"fmt"
	"io"
	"strconv"
	"sync"

	"github.com/gdamore/tcell/v2"
	"github.com/jeremija/tplot/scale"
	"github.com/rivo/tview"
	"github.com/shopspring/decimal"
)

// OHLC is a Box component that can render OHLC data.
type OHLC struct {
	*tview.Box

	mu      sync.Mutex
	items   []Item
	offset  int
	spacing int
	logger  io.Writer
	runes   Runes

	volSize int
}

// New creates a new instance of the OHLC component.
func New() *OHLC {
	return &OHLC{
		Box:     tview.NewBox(),
		runes:   DefaultRunes,
		spacing: 1,
		volSize: 10,
	}
}

// SetLogger sets the logger for debugging.
func (o *OHLC) SetLogger(w io.Writer) {
	o.mu.Lock()
	defer o.mu.Unlock()

	o.logger = w
}

// Logger returns the current logger set. Used for debugging.
func (o *OHLC) Logger() io.Writer {
	o.mu.Lock()
	defer o.mu.Unlock()

	if w := o.logger; w != nil {
		return w
	}

	return nopWriter{}
}

// Offset returns the current offset.
func (o *OHLC) Offset() int {
	o.mu.Lock()
	defer o.mu.Unlock()

	return o.offset
}

// AddOffset adds delta to the current offset.
func (o *OHLC) AddOffset(delta int) {
	o.mu.Lock()
	defer o.mu.Unlock()

	delta = o.offset + delta

	o.setOffset(delta)
}

// SetOffset sets the scroll offset for OHLC data.
func (o *OHLC) SetOffset(offset int) {
	o.mu.Lock()
	defer o.mu.Unlock()

	o.setOffset(offset)
}

// SetRunes sets the runes used to plot the chart.
func (o *OHLC) SetRunes(runes Runes) {
	o.mu.Lock()
	defer o.mu.Unlock()

	o.runes = runes
}

// Runes returns the current set of runes used to plot the chart.
func (o *OHLC) Runes() Runes {
	o.mu.Lock()
	defer o.mu.Unlock()

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
func (o *OHLC) SetItems(items []Item) {
	o.mu.Lock()
	defer o.mu.Unlock()

	o.items = items
}

// Items returns the current OHLC data.
func (o *OHLC) Items() []Item {
	o.mu.Lock()
	defer o.mu.Unlock()

	return o.items
}

// SetSpacing sets the chart spacing.
func (o *OHLC) SetSpacing(spacing int) {
	o.mu.Lock()
	defer o.mu.Unlock()

	if spacing <= 0 {
		spacing = 1
	}

	o.spacing = spacing
}

// Spacing returns the current spacing.
func (o *OHLC) Spacing() int {
	o.mu.Lock()
	defer o.mu.Unlock()

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

// drawOHLCAxisY draws the Y axis.
func (o *OHLC) drawAxisY(screen tcell.Screen, log io.Writer, scale *scale.Linear, r rect, color tcell.Color, getValue func() (decimal.Decimal, bool)) int {
	var current int

	currentVal, lastValOK := getValue()
	if lastValOK {
		current = scale.Value(currentVal)
	}

	// numDecs contains the number of decimals spot to represent the axis.
	numDecs := scale.NumDecimals() + 2

	maxWidth := 0

	vals := make([]string, r.h)

	for i := 0; i < r.h; i++ {
		val := scale.Reverse(i)

		if lastValOK && i == current {
			val = currentVal
		}

		valFloat, _ := val.Float64()

		// No decimals, ok for large values, might need some readjusting.
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

		if lastValOK && i == current {
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

func (o *OHLC) getOHLCRect() rect {
	x, y, w, h := o.GetInnerRect()

	if o.volSize < h {
		h -= o.volSize
	}

	return rect{x: x, y: y, w: w, h: h}
}

func (o *OHLC) getVolRect() rect {
	x, y, w, h := o.GetInnerRect()
	ohlcRect := o.getOHLCRect()

	if o.volSize >= h {
		return rect{x: x, y: y + h, w: w, h: 0}
	}

	h = o.volSize
	y += ohlcRect.h

	return rect{x: x, y: y, w: w, h: h}
}

// Draw implements tview.Primitive.
func (o *OHLC) Draw(screen tcell.Screen) {
	ohlcRect := o.getOHLCRect()
	volRect := o.getVolRect()
	ohlcScale := scale.NewLinear()
	volScale := scale.NewLinear()
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

	scaled := newScaledItems(items, ohlcScale, volScale)
	logger := o.Logger()

	axisYWidth := 0

	var lastItem *Item

	if l := len(items); l > 0 {
		lastItem = &items[l-1]
	}

	if len(scaled) > 0 {
		w1 := o.drawAxisY(screen, logger, ohlcScale, ohlcRect, tcell.ColorDarkCyan, func() (decimal.Decimal, bool) {
			if lastItem == nil {
				return decimal.Zero, false
			}

			return lastItem.C, true
		})

		w2 := o.drawAxisY(screen, logger, volScale, volRect, tcell.ColorDarkBlue, func() (decimal.Decimal, bool) {
			if lastItem == nil {
				return decimal.Zero, false
			}

			return lastItem.V, true
		})

		if w2 > w1 {
			axisYWidth = w2
		} else {
			axisYWidth = w1
		}
	}

	width -= axisYWidth

	// We need to readjust the maxCount because we drew the axis.
	maxCount = width / spacing

	if l := len(scaled); l > maxCount {
		scaled = scaled[l-maxCount:]
		items = items[l-maxCount:]
	}

	if width < 0 {
		return
	}

	fmt.Fprintln(logger, "log start", len(scaled))

	runes := o.Runes()

	// We are using special block characters to display quarters so we need
	// to resize our scale after we've drawn the axis.
	numVolFractions := 4
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
			yy := ohlcRect.y + ohlcRect.h - j

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

			for j := 0; j < volFullSteps; j++ {
				yy := volRect.y + volRect.h - j - 1
				screen.SetContent(xx, yy, runes.Block4, nil, volStyle)
			}

			if volRemFrac > 0 {
				var ch rune

				switch volRemFrac {
				case 1:
					ch = runes.Block1
				case 2:
					ch = runes.Block2
				case 3:
					ch = runes.Block3
				default:
					ch = ' '
				}

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
