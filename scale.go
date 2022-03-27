package tplot

import (
	"github.com/shopspring/decimal"
)

// Scale describes a component that can scale decimal values to terminal grid.
type Scale interface {
	// Size returns the current size.
	Size() int
	// SetSize sets the size.
	SetSize(int)
	// SetRange sets the range of the values.
	SetRange(min, max decimal.Decimal)
	// Range returns the current range.
	Range() (min, max decimal.Decimal)
	// NumDecimals returns the number of decimals
	// needed to draw the axis.
	NumDecimals() int
	// Value returns the place on scale.
	Value(decimal.Decimal) int
	// Reverse returns the min value for the place on
	// scale.
	Reverse(int) decimal.Decimal
}
