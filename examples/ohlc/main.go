package main

import (
	"encoding/json"
	"io"
	"os"
	"time"

	"github.com/jeremija/tplot"
	"github.com/rivo/tview"
	"github.com/shopspring/decimal"
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

	ohlcPanel := tplot.NewOHLC()
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

	mustDec := func(n json.Number) decimal.Decimal {
		return decimal.RequireFromString(string(n))
	}

	mustFloat := func(n json.Number) float64 {
		v, err := n.Float64()
		if err != nil {
			panic(err)
		}

		return v
	}

	items := make(tplot.OHLCItems, len(obj.Result["86400"]))

	for i, r := range obj.Result["86400"] {
		ts := mustFloat(r[0])
		o := mustDec(r[1])
		h := mustDec(r[2])
		l := mustDec(r[3])
		c := mustDec(r[4])
		v := mustDec(r[5])

		items[i] = tplot.OHLCItem{
			O:         o,
			H:         h,
			L:         l,
			C:         c,
			V:         v,
			Timestamp: time.Unix(int64(ts), 0),
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
