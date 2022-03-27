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
	V          decimal.Decimal
}

// scaledItem is a lossy representation of Item that can be rendered in the
// terminal grid.
type scaledItem struct {
	ts         time.Time
	O, H, L, C int
	V          int
}

// newScaledItems creates scaledItems from Items and scale. The scale's
// range will be set by this function.
func newScaledItems(ohlcs []Item, scale *scale.Linear, volScale *scale.Linear) []scaledItem {
	ret := make([]scaledItem, len(ohlcs))

	for i, ohlc := range ohlcs {
		ret[i] = scaledItem{
			O:  scale.Value(ohlc.O),
			H:  scale.Value(ohlc.H),
			L:  scale.Value(ohlc.L),
			C:  scale.Value(ohlc.C),
			V:  volScale.Value(ohlc.V),
			ts: ohlc.Timestamp,
		}
	}

	return ret
}

type decRange struct {
	min, max decimal.Decimal
}

func findRanges(items []Item) (ohlcRange, volRange decRange) {
	var (
		min, max       decimal.Decimal
		minSet, maxSet bool

		minVol, maxVol       decimal.Decimal
		minVolSet, maxVolSet bool
	)

	for _, ohlc := range items {
		if !minSet || ohlc.L.LessThan(min) {
			min = ohlc.L
			minSet = true
		}

		if !maxSet || ohlc.H.GreaterThan(max) {
			max = ohlc.H
			maxSet = true
		}

		if !minVolSet || ohlc.V.LessThan(minVol) {
			minVol = ohlc.V
			minVolSet = true
		}

		if !maxVolSet || ohlc.V.GreaterThan(maxVol) {
			maxVol = ohlc.V
			maxVolSet = true
		}
	}

	return decRange{min, max}, decRange{minVol, maxVol}
}
