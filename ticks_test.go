package tplot_test

import (
	"fmt"
	"testing"

	"github.com/jeremija/tplot"
	"github.com/jeremija/tplot/test"
	"github.com/stretchr/testify/assert"
)

func TestTicks(t *testing.T) {
	var factory tplot.FloatFactory

	p := tplot.NewTicks(factory)
	scr := test.NewScreen()

	vals := []tplot.Float{
		0,
		1.5,
		2.2,
		3.4,
		4.6,
		5.8,
		6,
		9.2,
		8.5,
		10,
	}

	data := make([]tplot.Decimal, len(vals))

	for i, val := range vals {
		data[i] = val
	}

	p.SetRect(0, 0, 12, 10)

	p.SetData(data)
	p.Draw(scr)

	exp := `
         ⎽ ⎺
          ⎼


       ⎻⎺
      —
     ⎼
    ⎽
   —
  ⎽`

	fmt.Println("== expected ==")
	fmt.Println(exp)
	fmt.Println("==  actual  ==")
	fmt.Println(scr.Content())
	fmt.Println("==============")

	assert.Equal(t, exp, "\n"+scr.Content())
}
