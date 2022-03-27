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

func (o OHLCItems) FindOHLCRange() Range {
	var rng Range

	for _, ohlc := range o {
		rng = rng.Feed(ohlc.L)
		rng = rng.Feed(ohlc.H)
	}

	return rng
}

func (o OHLCItems) FindVolumeRange() Range {
	var rng Range

	for _, ohlc := range o {
		rng = rng.Feed(ohlc.V)
	}

	return rng
}

// scaledOHLCItem is a lossy representation of Item that can be rendered in the
// terminal grid.
type scaledOHLCItem struct {
	ts         time.Time
	O, H, L, C int
}

// newScaledOHLCItems creates scaledItems from Items and scale. The scale's
// range will be set by this function.
func newScaledOHLCItems(ohlcs []OHLCItem, scale Scale) []scaledOHLCItem {
	ret := make([]scaledOHLCItem, len(ohlcs))

	for i, ohlc := range ohlcs {
		ret[i] = scaledOHLCItem{
			O:  scale.Value(ohlc.O),
			H:  scale.Value(ohlc.H),
			L:  scale.Value(ohlc.L),
			C:  scale.Value(ohlc.C),
			ts: ohlc.Timestamp,
		}
	}

	return ret
}
