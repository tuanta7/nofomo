package report

import (
	"fmt"
	"io"
	"time"

	"github.com/tuanta7/nofomo/internal/execution"
	"github.com/tuanta7/nofomo/internal/market/candle"
	"github.com/tuanta7/nofomo/internal/strategy"
)

type BacktestReport struct {
	Candles     int
	Start, End  time.Time
	Trades      int
	Wins        int
	Return      float64
	BuyHold     float64
	MaxDrawdown float64
}

func (r BacktestReport) PrintResult(w io.Writer) {
	var winRate float64
	if r.Trades > 0 {
		winRate = float64(r.Wins) / float64(r.Trades)
	}

	fmt.Fprintf(w, "%12s  %6d  (%s → %s)\n", "candles", r.Candles, r.Start.Format(time.DateOnly), r.End.Format(time.DateOnly))
	fmt.Fprintf(w, "%12s  %6d\n", "trades", r.Trades)
	fmt.Fprintf(w, "%12s  %5.1f%%\n", "win rate", winRate*100)
	fmt.Fprintf(w, "%12s  %5.1f%%\n", "return", r.Return*100)
	fmt.Fprintf(w, "%12s  %5.1f%%\n", "buy & hold", r.BuyHold*100)
	fmt.Fprintf(w, "%12s  %5.1f%%\n", "max DD", r.MaxDrawdown*100)
}

// RunBacktest replays candles through a strategy, long/flat and all-in.
func RunBacktest(
	candles []candle.Candle,
	backtestStrategy strategy.Strategy,
	cash, feeBasisPoint float64,
) BacktestReport {
	if len(candles) == 0 {
		return BacktestReport{}
	}

	var (
		start   = cash
		pending strategy.Signal  // order placed last bar, filled at this bar's open
		engine  execution.Engine = execution.NewBacktestEngine(cash, feeBasisPoint)
		peak                     = cash
		res                      = BacktestReport{
			Candles: len(candles),
			Start:   candles[0].OpenTime,
			End:     candles[len(candles)-1].CloseTime,
			BuyHold: candles[len(candles)-1].Close/candles[0].Open - 1,
		}
	)

	for i, c := range candles {
		// Fill what the previous bar's close asked for, at the earliest price actually
		// reachable: this bar's open. Filling on the signal bar's own close would be
		// lookahead bias and would flatter every result.
		fill := engine.Execute(pending, c.Open)
		if fill.Closed {
			res.Trades++
			if fill.Profitable {
				res.Wins++
			}
		}
		pending = strategy.Hold

		equity := engine.Equity(c.Close)
		peak = max(peak, equity)
		res.MaxDrawdown = min(res.MaxDrawdown, equity/peak-1)

		switch sig := backtestStrategy.Evaluate(strategy.Context{Candles: candles, Index: i}); {
		case sig == strategy.Buy:
			pending = strategy.Buy
		case sig == strategy.Sell:
			pending = strategy.Sell
		}
		// A signal on the final bar never fills: there is no next open.
	}

	// Ending long is a position, not a result. Close it at the last close so the
	// number reported is one a trader could have realised.
	if fill := engine.Execute(strategy.Sell, candles[len(candles)-1].Close); fill.Closed {
		res.Trades++
		if fill.Profitable {
			res.Wins++
		}
	}

	res.Return = engine.Equity(candles[len(candles)-1].Close)/start - 1
	return res
}
