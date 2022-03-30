package tplot

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type BlockBindings struct {
	*FlexBindings
	SplitH tcell.Key
	SplitV tcell.Key
	Delete tcell.Key
	Reset  tcell.Key
}

var DefaultBlockBindings = &BlockBindings{
	SplitH:       tcell.KeyCtrlS,
	SplitV:       tcell.KeyCtrlV,
	Delete:       tcell.KeyCtrlX,
	FlexBindings: DefaultFlexBindings,
	Reset:        tcell.KeyCtrlO,
}

// Tile is a wrapper around Flex that allows the user to dynamically create
// and remove splits, and navigate thorugh them, similar to a tiling window
// manager. See DefaultBlockBindings for default keybindings.
//
// After calling NewTile, the user should call SetFactory.
type Tile struct {
	*Flex
	root    bool            // root is false for all sub-tiles created via key-bindings.
	focused tview.Primitive // focused primitive within the Tile.

	factory  func() tview.Primitive // factory for creating new primitives.
	bindings *BlockBindings         // bindings is key-bindings configuration.
}

func NewTile() *Tile {
	flex := NewFlex()

	return &Tile{
		root:     true,
		Flex:     flex,
		bindings: DefaultBlockBindings,
	}
}

func (t *Tile) SetBindings(bindings *BlockBindings) {
	t.bindings = bindings
	t.Flex.SetBindings(bindings.FlexBindings)
}

func (t *Tile) SetFactory(factory func() tview.Primitive) {
	t.factory = factory

	if t.Flex.GetItemCount() == 0 && factory != nil {
		p = factory()
		t.Flex.AddItem(p, 0, 1, true)
	}
}

func (t *Tile) Factory() func() tview.Primitive {
	return t.factory
}

func (t *Tile) remove(setFocus func(p tview.Primitive)) (hasMore bool) {
	i, count, item := t.Flex.FocusedItem()

	if block, ok := item.(*Tile); ok {
		if hasMore := block.remove(setFocus); hasMore {
			return true
		}
	}

	t.Flex.RemoveItem(item)

	i--
	count--

	if i < 0 {
		i = 0
	}

	if count > i {
		it := t.Flex.GetItem(i)
		setFocus(it)

		return true
	}

	if count == 0 && t.root {
		it := t.factory()
		t.Flex.AddItem(it, 0, 1, true)
		setFocus(it)
	}

	return false
}

func (t *Tile) split(direction Direction, setFocus func(p tview.Primitive)) {
	_, _, item := t.Flex.FocusedItem()

	if block, ok := item.(*Tile); ok {
		block.split(direction, setFocus)
		return
	}

	it := t.Flex.GetItem(0)
	if _, ok := it.(*Tile); !ok {
		t.Flex.RemoveItem(it)

		t.direction = direction
		t.Flex.SetDirection(direction)

		f1 := NewTile()
		f1.SetDirection(direction)
		f1.AddItem(it, 0, 1, true)
		f1.SetFactory(t.factory)

		f1.root = false
		t.Flex.AddItem(f1, 0, 1, false)
	}

	f2 := NewTile()
	f2.SetDirection(direction)
	f2.SetFactory(t.factory)

	f2.root = false
	t.Flex.AddItem(f2, 0, 1, true)

	setFocus(f2)
}

func (t *Tile) reset(setFocus func(p tview.Primitive)) {
	i, _, item := t.Flex.FocusedItem()

	if block, ok := item.(*Tile); ok {
		block.reset(setFocus)

		return
	}

	if i < 0 {
		return
	}

	p := t.Flex.GetItem(i)
	t.Flex.RemoveItem(p)

	if t.factory != nil {
		p := t.factory()
		t.Flex.AddItem(p, 0, 1, true)
		setFocus(p)
	}
}

func (t *Tile) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return t.Flex.WrapInputHandler(
		func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
			switch event.Key() {
			case t.bindings.SplitH:
				t.split(DirectionHorizontal, setFocus)
			case t.bindings.SplitV:
				t.split(DirectionVertical, setFocus)
			case t.bindings.Delete:
				t.remove(setFocus)
			case t.bindings.Reset:
				t.reset(setFocus)
			default:
				handler := t.Flex.InputHandler()
				handler(event, setFocus)
			}
		},
	)
}

var _ tview.Primitive = &Tile{}
