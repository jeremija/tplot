package scale

import (
	"math"

	"github.com/shopspring/decimal"
)

// Linear represents a linear scale.
type Linear struct {
	min  decimal.Decimal
	max  decimal.Decimal
	size int
}

// NewLinear constructs a new linear scale.
func NewLinear() *Linear {
	return &Linear{}
}

// Size returns the scale size.
func (a *Linear) Size() int {
	return a.size
}

// SetRange sets the scale range.
func (a *Linear) SetRange(min, max decimal.Decimal) {
	a.min = min
	a.max = max
}

// SetSize sets the size.
func (a *Linear) SetSize(size int) {
	a.size = size
}

func (a *Linear) Reverse(i int) decimal.Decimal {
	return a.min.Add(decimal.New(int64(i), 0).Mul(a.step()))
}

func (a *Linear) NumDecimals() int {
	step := a.step()
	stepFloat, _ := step.Float64()

	numDecs := 0

	if f := math.Log10(stepFloat); f < 0 {
		numDecs = int(math.Abs(f))
	}

	return numDecs
}

func (a *Linear) step() decimal.Decimal {
	step := decimal.Zero

	if s := a.size - 1; s > 0 {
		step = a.max.Sub(a.min).Div(decimal.New(int64(s), 0))
	}

	return step
}

// Range returns the scale range.
func (a *Linear) Range() (min, max decimal.Decimal) {
	return a.min, a.max
}

func (a *Linear) scale() decimal.Decimal {
	if a.min.Equal(a.max) {
		return decimal.Zero
	}

	return decimal.New(int64(a.size-1), 0).Div(a.max.Sub(a.min))
}

// Value returns a scaled value from decimal.
func (a *Linear) Value(v decimal.Decimal) int {
	scale := a.scale()

	ret := v.Sub(a.min).Mul(scale).IntPart() + 1

	return int(ret)
}
