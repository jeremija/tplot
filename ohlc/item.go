package ohlc

import (
	"time"

	"github.com/jeremija/tplot/scale"
	"github.com/shopspring/decimal"
)

// Item represents an OHLC item with a timestamp.
type Item struct {
	Timestamp  time.Time
	O, H, L, C decimal.Decimal
}

// scaledItem is a lossy representation of Item that can be rendered in the
// terminal grid.
type scaledItem struct {
	ts         time.Time
	O, H, L, C int
}

// newScaledItems creates scaledItems from Items and scale. The scale's
// range will be set by this function.
func newScaledItems(ohlcs []Item, scale *scale.Linear) []scaledItem {
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

	ret := make([]scaledItem, len(ohlcs))

	for i, ohlc := range ohlcs {
		ret[i] = scaledItem{
			O:  scale.Value(ohlc.O),
			H:  scale.Value(ohlc.H),
			L:  scale.Value(ohlc.L),
			C:  scale.Value(ohlc.C),
			ts: ohlc.Timestamp,
		}
	}

	return ret
}
