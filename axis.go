package tplot

import (
	"strconv"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// Axis represents a chart Axis. The current implementation only supports a
// vertical axis.
type Axis struct {
	*tview.Box
	factory        DecimalFactory
	scale          Scale
	style          tcell.Style
	highlightStyle tcell.Style
	highlight      DecimalValue
}

// NewAxis creates a new instance of Axis.
func NewAxis(factory DecimalFactory) *Axis {
	return &Axis{
		factory:        factory,
		Box:            tview.NewBox(),
		style:          tcell.StyleDefault,
		scale:          NewScaleLinear(factory),
		highlightStyle: tcell.StyleDefault,
	}
}

// SetStyle sets a default axis style.
func (a *Axis) SetStyle(style tcell.Style) {
	a.style = style
}

// Style returns the default axis style.
func (a *Axis) Style() tcell.Style {
	return a.style
}

// SetHighlightStyle sets the style for highlighting an item on the axis.
func (a *Axis) SetHighlightStyle(highlightStyle tcell.Style) {
	a.highlightStyle = highlightStyle
}

// HighlightStyle returns the style for highlighting an item on the axis.
func (a *Axis) HighlightStyle() tcell.Style {
	return a.highlightStyle
}

// SetHighlight sets the value for highlighting.
func (a *Axis) SetHighlight(highlight DecimalValue) {
	a.highlight = highlight
}

// Highlight returns the currently highlighted item.
func (a *Axis) Highlight() DecimalValue {
	return a.highlight
}

// SetScale sets the axis acale.
func (a *Axis) SetScale(scale Scale) {
	a.scale = scale
}

// Scale returns the axis scale. May be nil.
func (a *Axis) Scale() Scale {
	return a.scale
}

// CalcWidth calculates the width of the axis.
func (a *Axis) CalcWidth() int {
	// numDecs contains the number of decimals spot to represent the axis.
	numDecs := a.scale.NumDecimals() + 2

	rng := a.scale.Range()

	size := len(rng.Max.Round().String()) + 1 + numDecs

	return size
}

// Draw implements tview.Primitive.
func (a *Axis) Draw(screen tcell.Screen) {
	a.Box.DrawForSubclass(screen, a)

	x, y, w, h := a.GetInnerRect()

	highlight := a.highlight
	scale := a.scale

	var highlightInt int

	if a.highlight.Valid {
		highlightInt = scale.Value(highlight.Decimal)
	}

	// numDecs contains the number of decimals spot to represent the axis.
	numDecs := scale.NumDecimals() + 2

	maxWidth := 0

	vals := make([]string, h)

	for i := 0; i < h; i++ {
		rev := scale.Reverse(i)

		if highlight.Valid && i == highlightInt {
			rev = highlight.Decimal
		}

		valFloat := rev.Float64()

		valStr := strconv.FormatFloat(valFloat, 'f', numDecs, 64)

		if l := len(valStr); l > maxWidth {
			maxWidth = l
		}

		vals[i] = valStr
	}

	// Hide axis when no room.
	if w < maxWidth {
		return
	}

	for i, valStr := range vals {
		yy := y + h - i - 1

		currentStyle := a.style

		if highlight.Valid && i == highlightInt {
			currentStyle = a.highlightStyle
		}

		for i, val := range valStr {
			xx := x + w - len(valStr) + i

			screen.SetContent(xx, yy, val, nil, currentStyle)
		}
	}
}
