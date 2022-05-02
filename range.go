package tplot

type Range struct {
	Min Decimal
	Max Decimal

	isSet bool
}

func NewRange(factory DecimalFactory) Range {
	return Range{
		Min: factory.Zero(),
		Max: factory.Zero(),
	}
}

func (r Range) IsSet() bool {
	return r.isSet
}

func (r Range) Feed(value Decimal) Range {
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
