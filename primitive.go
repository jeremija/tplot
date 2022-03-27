package tplot

import "github.com/gdamore/tcell/v2"

// Primitive is a tview.Primitive that can set and retrieve a
// Scale.
type Primitive interface {
	// SetRect is the tview.Primitive.SetRect method.
	SetRect(x, y, w, h int)
	// SetRect is the tview.Primitive.Draw method.
	Draw(screen tcell.Screen)
	// SetScale sets the scale to the Primitive.
	SetScale(Scale)
	// Scale returns the current scale used by the Primitive.
	Scale() Scale
}
