package tplot

import "github.com/shopspring/decimal"

type Range struct {
	Min decimal.Decimal
	Max decimal.Decimal

	isSet bool
}

func (r Range) IsSet() bool {
	return r.isSet
}

func (r Range) Feed(value decimal.Decimal) Range {
	if !r.isSet {
		r.Min = value
		r.Max = value
		r.isSet = true

		return r
	}

	if value.LessThan(r.Min) {
		r.Min = value
	}

	if value.GreaterThan(r.Max) {
		r.Max = value
	}

	return r
}
