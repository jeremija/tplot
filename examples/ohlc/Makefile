.PHONY: all
all: ohlc.json
	go run ./

ohlc.json:
	curl --silent --output ohlc.json 'https://api.cryptowat.ch/markets/kraken/btcusd/ohlc?periods=86400'
