package tplot

// OHLCRunes contains definitions for each OHLC drawing.
type OHLCRunes struct {
	Same          rune
	HighOpenClose rune
	LowOpenClose  rune
	HighOpen      rune
	LowClose      rune
	High          rune
	Low           rune
	OpenClose     rune
	Open          rune
	Close         rune
	Thick         rune
	Thin          rune
}

// DefaultOHLCRunes contains the default runes using the box drawing characters.
// Examples:
//
//     ╽   ╷ ╷     │
//     ┃╻  │ ╵ ┃ ─ ┼ ┼ ┬ ┴
//     ╿╹  │
//     │
//
// See more at: https://en.wikipedia.org/wiki/Box-drawing_character
var DefaultOHLCRunes = OHLCRunes{
	Same:          '─',
	HighOpenClose: '┬',
	LowOpenClose:  '┴',
	HighOpen:      '╻',
	LowClose:      '╹',
	High:          '╷',
	Low:           '╵',
	OpenClose:     '┼',
	Open:          '╽',
	Close:         '╿',
	Thick:         '┃',
	Thin:          '│',
}
