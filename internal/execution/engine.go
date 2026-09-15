package execution

import "github.com/tuanta7/nofomo/internal/strategy"

type Engine interface {
	Execute(strategy.Signal, float64) Fill
	Equity(float64) float64
}

type Fill struct {
	Opened     bool
	Closed     bool
	Profitable bool
}
