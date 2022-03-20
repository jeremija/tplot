package ohlc

// Runes contains definitions for each OHLC drawing.
type Runes struct {
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

// DefaultRunes contains the default runes using the box drawing characters.
// Examples:
//
//     ╽   ╷ ╷     │
//     ┃╻  │ ╵ ┃ ─ ┼ ┼ ┬ ┴
//     ╿╹  │
//     │
//
// See more at: https://en.wikipedia.org/wiki/Box-drawing_character
var DefaultRunes = Runes{
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
