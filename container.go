package tplot

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// Container is a dynamic box that can dynamically render different primitives.
type Container struct {
	*tview.Box

	item         tview.Primitive
	inputHandler func(event *tcell.EventKey, setFocus func(p tview.Primitive))
	mouseHandler func(action tview.MouseAction, event *tcell.EventMouse, setFocus func(p tview.Primitive)) (consumed bool, capture tview.Primitive)
}

func NewContainer() *Container {
	box := tview.NewBox().SetBorder(true)

	return &Container{
		Box: box,
	}
}

func (c *Container) Draw(screen tcell.Screen) {
	c.Box.DrawForSubclass(screen, c)

	if c.item != nil {
		c.item.SetRect(c.Box.GetInnerRect())
		c.item.Draw(screen)
	}
}

func (c *Container) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return c.Box.WrapInputHandler(func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
		if handler := c.inputHandler; handler != nil {
			handler(event, setFocus)
		}
	})
}

func (c *Container) MouseHandler() func(
	action tview.MouseAction, event *tcell.EventMouse, setFocus func(p tview.Primitive),
) (consumed bool, capture tview.Primitive) {
	return c.Box.WrapMouseHandler(
		func(
			action tview.MouseAction,
			event *tcell.EventMouse,
			setFocus func(p tview.Primitive),
		) (consumed bool, capture tview.Primitive) {
			if !c.InRect(event.Position()) {
				return false, nil
			}

			if handler := c.mouseHandler; handler != nil {
				return handler(action, event, setFocus)
			}

			return false, nil
		})
}

func (c *Container) SetPrimitive(item tview.Primitive) {
	c.item = item
	c.inputHandler = item.InputHandler()
	c.mouseHandler = item.MouseHandler()
}

func (c *Container) HasFocus() bool {
	if c.item != nil {
		return c.item.HasFocus()
	}

	return c.Box.HasFocus()
}

func (c *Container) Focus(delegate func(p tview.Primitive)) {
	if c.item != nil {
		delegate(c.item)
		return
	}

	c.Box.Focus(delegate)
}
