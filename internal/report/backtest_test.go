package report

import (
	"math"
	"testing"
	"time"

	"github.com/tuanta7/nofomo/internal/market/candle"
	"github.com/tuanta7/nofomo/internal/strategy"
)

// bar builds a candle whose open and close differ, so a fill taken from the wrong
// one shows up in the numbers.
func bar(i int, open, close float64) candle.Candle {
	t := time.UnixMilli(int64(i) * 60000).UTC()
	return candle.Candle{
		OpenTime:  t,
		CloseTime: t.Add(time.Minute),
		Open:      open,
		High:      math.Max(open, close),
		Low:       math.Min(open, close),
		Close:     close,
		Volume:    1,
	}
}

// signals replays a fixed script, one entry per bar, ignoring prices.
type signals []strategy.Signal

func (s signals) Evaluate(ctx strategy.Context) strategy.Signal { return s[ctx.Index] }

func TestRunFillsAtNextOpen(t *testing.T) {
	// Buy signalled on bar 0 fills at bar 1's open (110), not bar 0's close (105).
	// Sell signalled on bar 2 fills at bar 3's open (200), not bar 2's close (150).
	cs := []candle.Candle{
		bar(0, 100, 105),
		bar(1, 110, 120),
		bar(2, 130, 150),
		bar(3, 200, 210),
	}
	script := signals{strategy.Buy, strategy.Hold, strategy.Sell, strategy.Hold}

	got := Run(cs, script, 1000, 0)

	// 1000 / 110 units sold at 200 → 1000 * 200/110.
	want := 200.0/110.0 - 1
	if math.Abs(got.Return-want) > 1e-9 {
		t.Errorf("Return = %v, want %v (a fill at the signal bar's close would differ)", got.Return, want)
	}
	if got.Trades != 1 || got.Wins != 1 {
		t.Errorf("Trades/Wins = %d/%d, want 1/1", got.Trades, got.Wins)
	}
	if got.Candles != 4 {
		t.Errorf("Candles = %d, want 4", got.Candles)
	}
}

// A run that ends holding must be closed out, or the return reported is one nobody
// could have realised.
func TestRunLiquidatesOpenPosition(t *testing.T) {
	cs := []candle.Candle{
		bar(0, 100, 100),
		bar(1, 100, 100),
		bar(2, 100, 250),
	}
	got := Run(cs, signals{strategy.Buy, strategy.Hold, strategy.Hold}, 1000, 0)

	if want := 250.0/100.0 - 1; math.Abs(got.Return-want) > 1e-9 {
		t.Errorf("Return = %v, want %v", got.Return, want)
	}
	if got.Trades != 1 {
		t.Errorf("Trades = %d, want 1 (the forced exit counts)", got.Trades)
	}
}

// Fees are charged on both sides, so a flat round trip must lose money.
func TestRunChargesFeesBothSides(t *testing.T) {
	cs := []candle.Candle{
		bar(0, 100, 100),
		bar(1, 100, 100),
		bar(2, 100, 100),
		bar(3, 100, 100),
	}
	script := signals{strategy.Buy, strategy.Hold, strategy.Sell, strategy.Hold}

	if got := Run(cs, script, 1000, 0); math.Abs(got.Return) > 1e-9 {
		t.Errorf("Return with no fee = %v, want 0", got.Return)
	}

	got := Run(cs, script, 1000, 5)
	if want := 0.9995*0.9995 - 1; math.Abs(got.Return-want) > 1e-9 {
		t.Errorf("Return with 5bps = %v, want %v", got.Return, want)
	}
	if got.Wins != 0 {
		t.Errorf("Wins = %d, want 0: a flat round trip loses the fees", got.Wins)
	}
}

func TestRunDrawdownAndBuyHold(t *testing.T) {
	// Hold throughout: equity never moves, so drawdown stays 0 while buy & hold sinks.
	cs := []candle.Candle{
		bar(0, 100, 100),
		bar(1, 100, 50),
		bar(2, 50, 80),
	}
	got := Run(cs, signals{strategy.Hold, strategy.Hold, strategy.Hold}, 1000, 0)

	if got.MaxDrawdown != 0 || got.Return != 0 {
		t.Errorf("flat run: MaxDrawdown = %v, Return = %v, want 0/0", got.MaxDrawdown, got.Return)
	}
	if want := 80.0/100.0 - 1; math.Abs(got.BuyHold-want) > 1e-9 {
		t.Errorf("BuyHold = %v, want %v", got.BuyHold, want)
	}

	// Long through the crash: peak 1000 at bar 0, trough 500 at bar 1.
	long := Run(cs, signals{strategy.Buy, strategy.Hold, strategy.Hold}, 1000, 0)
	if want := -0.5; math.Abs(long.MaxDrawdown-want) > 1e-9 {
		t.Errorf("MaxDrawdown = %v, want %v", long.MaxDrawdown, want)
	}
}

func TestRunEmpty(t *testing.T) {
	if got := Run(nil, signals{}, 1000, 5); got != (BacktestReport{}) {
		t.Errorf("run(nil) = %+v, want zero Result", got)
	}
}
