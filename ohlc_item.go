package tplot

import (
	"time"

	"github.com/shopspring/decimal"
)

// OHLCItem represents an OHLC item with a timestamp.
type OHLCItem struct {
	Timestamp  time.Time
	O, H, L, C decimal.Decimal
	V          decimal.Decimal
}

type OHLCItems []OHLCItem

func (o OHLCItems) FindValueRange() (min, max decimal.Decimal) {
	var minSet, maxSet bool

	for _, ohlc := range o {
		if !minSet || ohlc.L.LessThan(min) {
			min = ohlc.L
			minSet = true
		}

		if !maxSet || ohlc.H.GreaterThan(max) {
			max = ohlc.H
			maxSet = true
		}
	}

	return min, max
}

func (o OHLCItems) FindVolumeRange() (min, max decimal.Decimal) {
	var minSet, maxSet bool

	for _, ohlc := range o {
		if !minSet || ohlc.V.LessThan(min) {
			min = ohlc.V
			minSet = true
		}

		if !maxSet || ohlc.V.GreaterThan(max) {
			max = ohlc.V
			maxSet = true
		}
	}

	return min, max
}

// scaledOHLCItem is a lossy representation of Item that can be rendered in the
// terminal grid.
type scaledOHLCItem struct {
	ts         time.Time
	O, H, L, C int
	V          int
}

// newScaledOHLCItems creates scaledItems from Items and scale. The scale's
// range will be set by this function.
func newScaledOHLCItems(ohlcs []OHLCItem, scale Scale, volScale Scale) []scaledOHLCItem {
	ret := make([]scaledOHLCItem, len(ohlcs))

	for i, ohlc := range ohlcs {
		ret[i] = scaledOHLCItem{
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
