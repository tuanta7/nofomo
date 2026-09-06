package strategy

import "fmt"

type EMACross struct {
	// fastK and slowK are the smoothing factors for the fast and slow EMAs
	fastK, slowK float64
	// fast and slow are the current values of the fast and slow EMAs
	fast, slow float64
	// seen counts how many bars have been processed
	seen int
	// warmup is the number of bars to process before the first signal can be generated
	warmup int
	// fastAbove tracks whether the fast EMA is above the slow one, so we can detect crossings
	fastAbove bool
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
