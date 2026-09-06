package strategy

import "fmt"

// EMACross goes long when the fast EMA crosses above the slow one and flat when it
// crosses back below. It keeps rolling EMA state, so Evaluate must be called once
// per candle in order, as the Strategy contract requires.
type EMACross struct {
	fastK, slowK float64 // smoothing multipliers, 2/(period+1)
	fast, slow   float64
	seen         int
	warmup       int
	fastAbove    bool
}

func NewEMACross(fastPeriod, slowPeriod int) (*EMACross, error) {
	if fastPeriod < 1 || slowPeriod < 1 {
		return nil, fmt.Errorf("ema cross: periods must be positive, got fast=%d slow=%d", fastPeriod, slowPeriod)
	}
	if fastPeriod >= slowPeriod {
		return nil, fmt.Errorf("ema cross: fast period %d must be shorter than slow %d", fastPeriod, slowPeriod)
	}
	return &EMACross{
		fastK:  2 / (float64(fastPeriod) + 1),
		slowK:  2 / (float64(slowPeriod) + 1),
		warmup: slowPeriod,
	}, nil
}

func (e *EMACross) Evaluate(ctx Context) Signal {
	price := ctx.Candle().Close

	if e.seen == 0 {
		e.fast, e.slow = price, price
	} else {
		e.fast += (price - e.fast) * e.fastK
		e.slow += (price - e.slow) * e.slowK
	}
	e.seen++

	// Both EMAs are seeded from the same price, so early bars say nothing about
	// trend. Track the relationship through warmup, but stay silent.
	above := e.fast > e.slow
	if e.seen <= e.warmup {
		e.fastAbove = above
		return Hold
	}

	crossed := above != e.fastAbove
	e.fastAbove = above
	switch {
	case crossed && above:
		return Buy
	case crossed && !above:
		return Sell
	default:
		return Hold
	}
}
