package tplot

import (
	"math"

	"github.com/shopspring/decimal"
)

// Linear represents a linear scale.
type ScaleLinear struct {
	rng  Range
	size int
}

var _ Scale = &ScaleLinear{}

// NewLinear constructs a new linear scale.
func NewScaleLinear() *ScaleLinear {
	return &ScaleLinear{}
}

func (a *ScaleLinear) Copy() Scale {
	b := *a
	return &b
}

// Size returns the scale size.
func (a *ScaleLinear) Size() int {
	return a.size
}

// SetRange sets the scale range.
func (a *ScaleLinear) SetRange(rng Range) {
	a.rng = rng
}

func (a *ScaleLinear) Range() Range {
	return a.rng
}

// SetSize sets the size.
func (a *ScaleLinear) SetSize(size int) {
	a.size = size
}

func (a *ScaleLinear) Reverse(i int) decimal.Decimal {
	return a.rng.Min.Add(decimal.New(int64(i), 0).Mul(a.step()))
}

func (a *ScaleLinear) NumDecimals() int {
	step := a.step()
	if step.IsZero() {
		return 0
	}

	stepFloat, _ := step.Float64()

	numDecs := 0

	if f := math.Log10(stepFloat); f < 0 {
		numDecs = int(math.Abs(f))
	}

	return numDecs
}

func (a *ScaleLinear) step() decimal.Decimal {
	step := decimal.Zero

	if s := a.size - 1; s > 0 {
		step = a.rng.Max.Sub(a.rng.Min).Div(decimal.NewFromInt(int64(s)))
	}

	return step
}

func (a *ScaleLinear) scale() decimal.Decimal {
	if a.rng.Min.Equal(a.rng.Max) {
		return decimal.Zero
	}

	return decimal.New(int64(a.size-1), 0).Div(a.rng.Max.Sub(a.rng.Min))
}

// Value returns a scaled value from decimal.
func (a *ScaleLinear) Value(v decimal.Decimal) int {
	scale := a.scale()

	ret := v.Sub(a.rng.Min).Mul(scale).IntPart()

	return int(ret)
}
