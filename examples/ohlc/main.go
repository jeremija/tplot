package main

import (
	"encoding/json"
	"io"
	"os"
	"time"

	"github.com/jeremija/tplot"
	"github.com/rivo/tview"
)

func main() {
	app := tview.NewApplication()

	var logger io.Writer

	if f := os.Getenv("LOG_FILE"); f != "" {
		log, err := os.OpenFile(f, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
		if err != nil {
			panic(err)
		}

		defer log.Close()

		logger = log
	}

	var factory tplot.FloatFactory

	ohlcPanel := tplot.NewOHLCChart(factory)
	ohlcPanel.SetLogger(logger)

	f, err := os.Open("ohlc.json")
	if err != nil {
		panic(err)
	}

	defer f.Close()

	obj := struct {
		Result map[string][][]json.Number `json:"result"`
	}{}

	if err := json.NewDecoder(f).Decode(&obj); err != nil {
		panic(err)
	}

	f.Close()

	mustDec := func(n json.Number) tplot.Decimal {
		f64, err := n.Float64()
		if err != nil {
			panic(err)
		}

		return tplot.Float(f64)
	}

	mustInt64 := func(n json.Number) int64 {
		v, err := n.Int64()
		if err != nil {
			panic(err)
		}

		return v
	}

	items := make([]tplot.OHLC, len(obj.Result["86400"]))

	for i, r := range obj.Result["86400"] {
		ts := time.Unix(mustInt64(r[0]), 0)
		o := mustDec(r[1])
		h := mustDec(r[2])
		l := mustDec(r[3])
		c := mustDec(r[4])
		v := mustDec(r[5])

		items[i] = tplot.OHLC{
			O:         o,
			H:         h,
			L:         l,
			C:         c,
			V:         v,
			Timestamp: ts,
		}
	}

	ohlcPanel.SetItems(items)
	ohlcPanel.SetBorder(true)

	layout := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(ohlcPanel, 0, 1, true)

	if err := app.SetRoot(layout, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}
