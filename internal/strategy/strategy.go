package strategy

import (
	"github.com/tuanta7/nofomo/internal/market/candle"
)

type Strategy interface {
	// Evaluate is called once per candle, in chronological order, with no gaps.
	// That contract lets implementations carry rolling state instead of
	// recomputing over the whole history on every bar.
	Evaluate(Context) Signal
}

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
