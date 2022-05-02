package tplot

import (
	"math"
	"strconv"
)

// DecimalValue represents a decimal value that can be undefined, in which case
// Valid is set to false.
type DecimalValue struct {
	// Decimal value.
	Decimal Decimal
	// Valid is ture when Decimal is defined.
	Valid bool
}

type Decimal interface {
	IsZero() bool
	Add(Decimal) Decimal
	Sub(Decimal) Decimal
	Mul(Decimal) Decimal
	Div(Decimal) Decimal
	Equal(Decimal) bool
	GreaterThan(Decimal) bool
	LessThan(Decimal) bool
	Float64() float64
	String() string
	Round() Decimal
	IntPart() int64
}

type DecimalFactory interface {
	Zero() Decimal
	NewFromInt64(int64) Decimal
}

type FloatFactory struct{}

func (f FloatFactory) Zero() Decimal {
	return Float(0)
}

func (f FloatFactory) NewFromInt64(i int64) Decimal {
	return Float(i)
}

// Float is an implemnetation of Decimal that uses float64.
type Float float64

var _ Decimal = Float(0)

func (f Float) IsZero() bool {
	return f == 0
}

func (f Float) Add(other Decimal) Decimal {
	return Float(f + other.(Float))
}

func (f Float) Sub(other Decimal) Decimal {
	return Float(f - other.(Float))
}

func (f Float) Mul(other Decimal) Decimal {
	return Float(f * other.(Float))
}

func (f Float) Div(other Decimal) Decimal {
	return Float(f / other.(Float))
}

func (f Float) Equal(other Decimal) bool {
	return f == other.(Float)
}

func (f Float) GreaterThan(other Decimal) bool {
	return f > other.(Float)
}

func (f Float) LessThan(other Decimal) bool {
	return f < other.(Float)
}

func (f Float) Round() Decimal {
	return Float(math.Round(float64(f)))
}

func (f Float) String() string {
	return strconv.FormatFloat(float64(f), 'f', -1, 64)
}
func (f Float) Float64() float64 {
	return float64(f)
}

func (f Float) IntPart() int64 {
	return int64(f)
}
