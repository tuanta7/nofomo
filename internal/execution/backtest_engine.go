package execution

import "github.com/tuanta7/nofomo/internal/strategy"

// BacktestEngine simulates an all-in long/flat account
type BacktestEngine struct {
	// cash is the amount of cash available to trade
	cash float64
	// quantity is the number of units held (BTC, Contracts, etc.)
	quantity float64
	// entryCost is the cash spent to open the position
	entryCost float64
	// fee is the fee rate
	fee float64
}

func NewBacktestEngine(cash, feeBasisPoints float64) *BacktestEngine {
	return &BacktestEngine{
		cash: cash,
		fee:  feeBasisPoints / 10000,
	}
}

func (e *BacktestEngine) Equity(price float64) float64 {
	return e.cash + e.quantity*price
}

func (e *BacktestEngine) Execute(signal strategy.Signal, price float64) Fill {
	switch signal {
	case strategy.Buy:
		if e.quantity > 0 {
			return Fill{}
		}

		e.entryCost = e.cash
		e.quantity = e.cash * (1 - e.fee) / price
		e.cash = 0
		return Fill{
			Opened: true,
		}
	case strategy.Sell:
		if e.quantity == 0 {
			return Fill{}
		}

		e.cash = e.quantity * price * (1 - e.fee)
		e.quantity = 0
		return Fill{
			Closed:     true,
			Profitable: e.cash > e.entryCost,
		}
	default:
		return Fill{}
	}
}
