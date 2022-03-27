package tplot

import (
	"time"

	"github.com/shopspring/decimal"
)

// OHLC represents an OHLC item with a timestamp.
type OHLC struct {
	Timestamp  time.Time
	O, H, L, C decimal.Decimal
	V          decimal.Decimal
}
