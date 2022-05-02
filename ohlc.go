package tplot

import (
	"time"
)

// OHLC represents an OHLC item with a timestamp.
type OHLC struct {
	Timestamp  time.Time
	O, H, L, C Decimal
	V          Decimal
}
