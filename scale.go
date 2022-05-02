package tplot

// Scale describes a component that can scale decimal values to terminal grid.
type Scale interface {
	Copy() Scale
	// Size returns the current size.
	Size() int
	// SetSize sets the size.
	SetSize(int)
	// SetRange sets the range of the values.
	SetRange(Range)
	// Range returns the current range.
	Range() Range
	// NumDecimals returns the number of decimals
	// needed to draw the axis.
	NumDecimals() int
	// Value returns the place on scale.
	Value(Decimal) int
	// Reverse returns the min value for the place on
	// scale.
	Reverse(int) Decimal
}
