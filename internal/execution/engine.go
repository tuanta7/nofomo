package execution

import "github.com/tuanta7/nofomo/internal/strategy"

type Engine interface {
	Execute(strategy.Signal, float64) Fill
	Equity(float64) float64
}

// Fill describes the outcome of an execution attempt.
type Fill struct {
	Opened     bool
	Closed     bool
	Profitable bool
}
