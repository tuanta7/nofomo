package report

import (
	"fmt"
	"io"
	"time"

	"github.com/tuanta7/nofomo/internal/market/candle"
	"github.com/tuanta7/nofomo/internal/strategy"
)

type BacktestReport struct {
	Candles     int
	Start, End  time.Time
	Trades      int // completed round trips
	Wins        int
	Return      float64 // final cash / starting cash - 1
	BuyHold     float64
	MaxDrawdown float64
}

// Run replays candles through a strategy, long/flat and all-in.
//
// ponytail: P&L is computed in USD terms. BTCUSD_PERP is an inverse contract that
// really settles in BTC, so this understates convexity — fine for judging whether an
// edge exists, wrong for sizing. Switch to BTC-denominated P&L if that starts to matter.
func Run(cs []candle.Candle, s strategy.Strategy, cash, feeBps float64) BacktestReport {
	if len(cs) == 0 {
		return BacktestReport{}
	}

	var (
		fee       = feeBps / 10000
		start     = cash
		qty       float64         // units held; 0 means flat
		entryCost float64         // cash spent entering the open position
		pending   strategy.Signal // order placed last bar, filled at this bar's open
		peak      = cash
		res       = BacktestReport{
			Candles: len(cs),
			Start:   cs[0].OpenTime,
			End:     cs[len(cs)-1].CloseTime,
			BuyHold: cs[len(cs)-1].Close/cs[0].Open - 1,
		}
	)

	for i, c := range cs {
		// Fill what the previous bar's close asked for, at the earliest price actually
		// reachable: this bar's open. Filling on the signal bar's own close would be
		// lookahead bias and would flatter every result.
		switch pending {
		case strategy.Buy:
			entryCost = cash
			qty = cash * (1 - fee) / c.Open
			cash = 0
		case strategy.Sell:
			cash = qty * c.Open * (1 - fee)
			qty = 0
			res.Trades++
			if cash > entryCost {
				res.Wins++
			}
		}
		pending = strategy.Hold

		equity := cash + qty*c.Close
		peak = max(peak, equity)
		res.MaxDrawdown = min(res.MaxDrawdown, equity/peak-1)

		switch sig := s.Evaluate(strategy.Context{Candles: cs, Index: i}); {
		case sig == strategy.Buy && qty == 0:
			pending = strategy.Buy
		case sig == strategy.Sell && qty > 0:
			pending = strategy.Sell
		}
		// A signal on the final bar never fills: there is no next open.
	}

	// Ending long is a position, not a result. Close it at the last close so the
	// number reported is one a trader could have realised.
	if qty > 0 {
		cash = qty * cs[len(cs)-1].Close * (1 - fee)
		res.Trades++
		if cash > entryCost {
			res.Wins++
		}
	}

	res.Return = cash/start - 1
	return res
}

func (r BacktestReport) Result(w io.Writer) {
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
