package ohlc

import (
	"fmt"
	"io"
	"math"
	"strconv"
	"sync"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/jeremija/tplot/scale"
	"github.com/rivo/tview"
	"github.com/shopspring/decimal"
)

type OHLC struct {
	O, H, L, C decimal.Decimal
	Timestamp  time.Time
}

type OHLCPlot struct {
	O, H, L, C int
	ts         time.Time
}

func NewOHLCPlots(ohlcs []OHLC, scale *scale.Linear, h int) []OHLCPlot {
	var (
		min, max       decimal.Decimal
		minSet, maxSet bool
	)

	for _, ohlc := range ohlcs {
		if !minSet || ohlc.L.LessThan(min) {
			min = ohlc.L
			minSet = true
		}

		if !maxSet || ohlc.H.GreaterThan(max) {
			max = ohlc.H
			maxSet = true
		}
	}

	scale.SetRange(min, max)
	scale.SetSize(h)

	ret := make([]OHLCPlot, len(ohlcs))

	for i, ohlc := range ohlcs {
		ret[i] = OHLCPlot{
			O:  scale.Value(ohlc.O),
			H:  scale.Value(ohlc.H),
			L:  scale.Value(ohlc.L),
			C:  scale.Value(ohlc.C),
			ts: ohlc.Timestamp,
		}
	}

	return ret
}

type OHLCPanel struct {
	*tview.Box

	mu     sync.Mutex
	ohlcs  []OHLC
	offset int
	logger io.Writer
}

func New() *OHLCPanel {
	return &OHLCPanel{
		Box: tview.NewBox(),
	}
}

func (p *OHLCPanel) SetLogger(w io.Writer) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.logger = w
}

func (p *OHLCPanel) Logger() io.Writer {
	p.mu.Lock()
	defer p.mu.Unlock()

	if w := p.logger; w != nil {
		return w
	}

	return nopWriter{}
}

func (p *OHLCPanel) Offset() int {
	p.mu.Lock()
	defer p.mu.Unlock()

	return p.offset
}

func (p *OHLCPanel) AddOffset(offset int) {
	p.mu.Lock()
	defer p.mu.Unlock()

	offset = p.offset + offset

	p.setOffset(offset)
}

func (p *OHLCPanel) SetOffset(offset int) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.setOffset(offset)
}

func (p *OHLCPanel) setOffset(offset int) {
	if l := len(p.ohlcs); offset >= l {
		offset = l - 1
	}

	if offset < 0 {
		offset = 0
	}

	p.offset = offset
}

func (p *OHLCPanel) SetOHLC(ohlcs []OHLC) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.ohlcs = ohlcs
}

func (p *OHLCPanel) OHLCs() []OHLC {
	p.mu.Lock()
	defer p.mu.Unlock()

	return p.ohlcs
}

func (b *OHLCPanel) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return b.WrapInputHandler(
		func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
			switch event.Key() {
			case tcell.KeyEnd:
				b.SetOffset(0)
			case tcell.KeyHome:
				s := len(b.OHLCs())
				b.SetOffset(s - 1)
			case tcell.KeyPgUp:
				b.AddOffset(20)
			case tcell.KeyPgDn:
				b.AddOffset(-20)
			case tcell.KeyLeft:
				b.AddOffset(1)
			case tcell.KeyRight:
				b.AddOffset(-1)

			case tcell.KeyRune:
				switch event.Rune() {
				case 'h':
					b.AddOffset(1)
				case 'l':
					b.AddOffset(-1)
				}
			}
		},
	)
}

func (b *OHLCPanel) MouseHandler() func(action tview.MouseAction, event *tcell.EventMouse, setFocus func(p tview.Primitive)) (consumed bool, capture tview.Primitive) {
	return b.Box.WrapMouseHandler(func(action tview.MouseAction, ev *tcell.EventMouse, setFocus func(p tview.Primitive)) (bool, tview.Primitive) {
		switch action {
		case tview.MouseScrollUp:
			b.AddOffset(10)

			return true, b
		case tview.MouseScrollDown:
			b.AddOffset(-10)

			return true, b
		}

		return false, b
	})
}

func (p *OHLCPanel) drawAxisY(screen tcell.Screen, scale *scale.Linear, log io.Writer) int {
	x, y, width, height := p.GetInnerRect()

	min, _ := scale.Range()
	step := scale.Step()

	// numDecs contains the number of decimals spot to represent the axis.
	numDecs := 2

	stepFloat, _ := step.Float64()

	if f := math.Log10(stepFloat); f < 0 {
		numDecs = int(math.Abs(f)) + 2
	}

	maxWidth := 0

	vals := make([]string, height)

	for i := 0; i < height; i++ {
		val := min.Add(decimal.New(int64(i), 0).Mul(step))

		valFloat, _ := val.Float64()

		// No decimals, ok for large values, might need some readjusting.
		valStr := strconv.FormatFloat(valFloat, 'f', numDecs, 64)

		if l := len(valStr); l > maxWidth {
			maxWidth = l
		}

		fmt.Fprintln(log, "value", i, string(valStr), step)

		vals[i] = valStr
	}

	if width < maxWidth {
		return 0 // hide axis when no room.
	}

	style := tcell.StyleDefault.Foreground(tcell.ColorDarkCyan)

	for i, valStr := range vals {
		yy := y + height - i - 1

		for i, val := range valStr {
			xx := x + width - len(valStr) + i

			fmt.Fprintln(log, "axis y", xx, yy, string(val))
			screen.SetContent(xx, yy, val, nil, style)
		}
	}

	return maxWidth
}

// Draw draws this primitive onto the screen. Implementers can call the
// screen's ShowCursor() function but should only do so when they have focus.
// (They will need to keep track of this themselves.)
func (p *OHLCPanel) Draw(screen tcell.Screen) {
	p.DrawForSubclass(screen, p)

	x, y, width, height := p.GetInnerRect()

	scale := scale.NewLinear()

	ohlcW := 1

	rawOHLCs := p.OHLCs()

	offset := p.Offset()
	if l := len(rawOHLCs); offset > l {
		rawOHLCs = nil
	} else {
		rawOHLCs = rawOHLCs[:l-offset]
	}

	if l := len(rawOHLCs); l > 0 {
		ohlc := rawOHLCs[l-1]

		title := fmt.Sprintf(" O=%s H=%s L=%s C=%s TS=%s", ohlc.O, ohlc.H, ohlc.L, ohlc.C, ohlc.Timestamp.Format("2006-01-02T15:04:05"))

		p.SetTitle(title)

		p.DrawForSubclass(screen, p)
	}

	maxCount := width / ohlcW

	if l := len(rawOHLCs); l > maxCount {
		rawOHLCs = rawOHLCs[l-maxCount:]
	}

	ohlcs := NewOHLCPlots(rawOHLCs, scale, height)
	logger := p.Logger()

	axisYWidth := 0

	if len(ohlcs) > 0 {
		axisYWidth = p.drawAxisY(screen, scale, logger)
	}

	width -= axisYWidth

	// We need to readjust the maxCount because we drew the axis.
	maxCount = width / ohlcW

	if l := len(rawOHLCs); l > maxCount {
		ohlcs = ohlcs[l-maxCount:]
	}

	if width < 0 {
		return
	}

	fmt.Fprintln(logger, "log start", len(ohlcs))

	for i, ohlc := range ohlcs {
		o, h, l, c := ohlc.O, ohlc.H, ohlc.L, ohlc.C

		style := tcell.StyleDefault

		style = style.Foreground(tcell.ColorRed)

		a, b := o, c
		if b >= a {
			style = style.Foreground(tcell.ColorGreen)
			a, b = b, a
		}

		xx := x + i*ohlcW + (width - len(ohlcs)*ohlcW)
		_ = y

		fmt.Fprintln(logger, "ohlc", i, o, h, l, c)

		// Sample values:
		//
		//     ╽   ╷ ╷     │
		//     ┃╻  │ ╵ ┃ ─ ┼ ┼ ┬ ┴
		//     ╿╹  │
		//     │
		//
		// See more at: https://en.wikipedia.org/wiki/Box-drawing_character

		for j := h; j >= l; j-- {
			style := style
			yy := y + height - j

			isHigh := j == h
			isLow := j == l
			isOpen := j == a
			isClose := j == b
			isThick := j < a && j > b

			var ch rune

			switch {
			case isHigh && isLow && isOpen && isClose:
				ch = '─'
			case isHigh && isOpen && isClose:
				ch = '┬'
			case isLow && isOpen && isClose:
				ch = '┴'
			case isHigh && isOpen:
				ch = '╻'
				style = style.Bold(true)
			case isLow && isClose:
				ch = '╹'
				style = style.Bold(true)
			case isHigh:
				ch = '╷'
			case isLow:
				ch = '╵'
			case isOpen && isClose:
				ch = '┼'
			case isOpen:
				ch = '╽'
			case isClose:
				ch = '╿'
			case isThick:
				style = style.Bold(true)
				ch = '┃'
			default:
				ch = '│'
			}

			fmt.Fprintln(logger, "coord", xx, yy, string(ch))

			screen.SetContent(xx, yy, ch, nil, style)
		}
	}
}

type nopWriter struct{}

func (n nopWriter) Write(b []byte) (int, error) {
	return len(b), nil
}
