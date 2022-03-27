package scale_test

import (
	"testing"

	"github.com/jeremija/tplot/scale"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestLiner(t *testing.T) {
	l := scale.NewLinear()

	min := decimal.NewFromInt(5)
	max := decimal.NewFromInt(10)

	l.SetRange(min, max)
	l.SetSize(7)

	min2, max2 := l.Range()
	assert.True(t, min2.Equal(min))
	assert.True(t, max2.Equal(max))

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
