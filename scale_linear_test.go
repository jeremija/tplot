package tplot_test

import (
	"testing"

	"github.com/jeremija/tplot"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestLiner(t *testing.T) {
	l := tplot.NewScaleLinear()

	rng := tplot.Range{
		Min: decimal.NewFromInt(5),
		Max: decimal.NewFromInt(10),
	}

	l.SetRange(rng)
	l.SetSize(7)

	assert.Equal(t, rng, l.Range())

	assert.Equal(t, l.Size(), 7)

	assert.Equal(t, 0, l.NumDecimals())

	assert.Equal(t, 0, l.Value(decimal.NewFromInt(5)))

	rev0 := l.Reverse(0)
	rev1 := l.Reverse(1)
	rev2 := l.Reverse(2)
	rev3 := l.Reverse(3)

	assert.True(t, rev0.Equal(decimal.NewFromInt(5)), rev0.String())
	assert.True(t, rev1.Equal(decimal.RequireFromString("5.8333333333333333")), rev1.String())
	assert.True(t, rev2.Equal(decimal.RequireFromString("6.6666666666666666")), rev2.String())
	assert.True(t, rev3.Equal(decimal.RequireFromString("7.4999999999999999")), rev3.String())
}
