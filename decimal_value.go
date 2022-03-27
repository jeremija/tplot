package tplot

import "github.com/shopspring/decimal"

// DecimalValue represents a decimal value that can be undefined, in which case
// Valid is set to false.
type DecimalValue struct {
	// Decimal value.
	Decimal decimal.Decimal
	// Valid is ture when Decimal is defined.
	Valid bool
}
