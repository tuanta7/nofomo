package strategy

import (
	"github.com/tuanta7/nofomo/internal/market/candle"
)

type Signal int

const (
	Hold Signal = iota
	Buy
	Sell
)

// Context is the view a strategy gets when a candle closes. Only candles up to
// and including Index have happened; reading past it is lookahead bias.
type Context struct {
	Candles []candle.Candle
	Index   int
}

// Candle is the bar that just closed.
func (c Context) Candle() candle.Candle {
	return c.Candles[c.Index]
}

type Strategy interface {
	// Evaluate is called once per candle, in chronological order, with no gaps.
	Evaluate(Context) Signal
}
