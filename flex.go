package tplot

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Direction int

const (
	DirectionVertical   Direction = tview.FlexRow
	DirectionHorizontal Direction = tview.FlexColumn
)

type Focus int

const (
	FocusLeft Focus = iota
	FocusDown
	FocusUp
	FocusRight
)

// Flex is a wrapper around tview.Flex that allows easy keyboard navigation.
type Flex struct {
	*tview.Flex
	focused   tview.Primitive
	bindings  *FlexBindings
	direction Direction
}

type FlexBindings struct {
	FocusLeft  tcell.Key
	FocusDown  tcell.Key
	FocusUp    tcell.Key
	FocusRight tcell.Key
}

var DefaultFlexBindings = &FlexBindings{
	FocusLeft:  tcell.KeyCtrlH,
	FocusDown:  tcell.KeyCtrlJ,
	FocusUp:    tcell.KeyCtrlK,
	FocusRight: tcell.KeyCtrlL,
}

func NewFlex() *Flex {
	return &Flex{
		Flex:     tview.NewFlex(),
		bindings: DefaultFlexBindings,
	}
}

func (f *Flex) SetDirection(direction Direction) {
	f.direction = direction
	f.Flex.SetDirection(int(direction))
}

func (f *Flex) Direction() Direction {
	return f.direction
}

func (f *Flex) SetBindings(bindings *FlexBindings) {
	f.bindings = bindings
}

func (f *Flex) Bindings() *FlexBindings {
	return f.bindings
}

func (f *Flex) FocusedItem() (i int, count int, item tview.Primitive) {
	count = f.Flex.GetItemCount()

	for i = 0; i < count; i++ {
		item := f.Flex.GetItem(i)

		if item.HasFocus() {
			return i, count, item
		}
	}

	return -1, count, nil
}

func (f *Flex) wrapSetFocus(setFocus func(p tview.Primitive)) func(tview.Primitive) {
	return func(p tview.Primitive) {
		setFocus(p)

		_, _, item := f.FocusedItem()
		f.focused = item
	}
}

func (f *Flex) move(move Focus, setFocus func(p tview.Primitive)) (consumed bool) {
	setFocus = f.wrapSetFocus(setFocus)

	i, count, item := f.FocusedItem()
	if item == nil {
		return false
	}

	if block, ok := item.(*Tile); ok {
		if consumed := block.move(move, setFocus); consumed {
			return true
		}
	}

	dec := func() bool {
		i--
		if i < 0 {
			return false
		}

		setFocus(f.Flex.GetItem(i))
		return true
	}

	inc := func() bool {
		i++
		if i >= count {
			return false
		}

		setFocus(f.Flex.GetItem(i))
		return true
	}

	if f.direction == DirectionHorizontal {
		switch move {
		case FocusLeft:
			return dec()
		case FocusRight:
			return inc()
		default:
			return false
		}
	} else {
		switch move {
		case FocusUp:
			return dec()
		case FocusDown:
			return inc()
		default:
			return false
		}
	}
}

func (b *Flex) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return b.Flex.WrapInputHandler(
		func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
			switch event.Key() {
			case b.bindings.FocusLeft:
				b.move(FocusLeft, setFocus)
			case b.bindings.FocusDown:
				b.move(FocusDown, setFocus)
			case b.bindings.FocusUp:
				b.move(FocusUp, setFocus)
			case b.bindings.FocusRight:
				b.move(FocusRight, setFocus)
			default:
				handler := b.Flex.InputHandler()
				handler(event, setFocus)
			}
		},
	)
}

func (b *Flex) Focus(setFocus func(p tview.Primitive)) {
	count := b.Flex.GetItemCount()

	for i := 0; i < count; i++ {
		item := b.Flex.GetItem(i)
		if item == b.focused {
			setFocus(item)
			return
		}
	}

	b.Flex.Focus(setFocus)
}
