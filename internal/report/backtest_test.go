package report

import (
	"math"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
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

	got := RunBacktest(cs, script, 1000, 0)

	// 1000 / 110 units sold at 200 → 1000 * 200/110.
	want := 200.0/110.0 - 1
	assert.InDelta(t, want, got.Return, 1e-9, "a fill at the signal bar's close would differ")
	assert.Equal(t, 1, got.Trades)
	assert.Equal(t, 1, got.Wins)
	assert.Equal(t, 4, got.Candles)
}

// A run that ends holding must be closed out, or the return reported is one nobody
// could have realised.
func TestRunLiquidatesOpenPosition(t *testing.T) {
	cs := []candle.Candle{
		bar(0, 100, 100),
		bar(1, 100, 100),
		bar(2, 100, 250),
	}
	got := RunBacktest(cs, signals{strategy.Buy, strategy.Hold, strategy.Hold}, 1000, 0)

	assert.InDelta(t, 250.0/100.0-1, got.Return, 1e-9)
	assert.Equal(t, 1, got.Trades, "the forced exit counts")
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

	assert.InDelta(t, 0, RunBacktest(cs, script, 1000, 0).Return, 1e-9)

	got := RunBacktest(cs, script, 1000, 5)
	assert.InDelta(t, 0.9995*0.9995-1, got.Return, 1e-9)
	assert.Equal(t, 0, got.Wins, "a flat round trip loses the fees")
}

func TestRunDrawdownAndBuyHold(t *testing.T) {
	// Hold throughout: equity never moves, so drawdown stays 0 while buy & hold sinks.
	cs := []candle.Candle{
		bar(0, 100, 100),
		bar(1, 100, 50),
		bar(2, 50, 80),
	}
	got := RunBacktest(cs, signals{strategy.Hold, strategy.Hold, strategy.Hold}, 1000, 0)

	assert.Zero(t, got.MaxDrawdown)
	assert.Zero(t, got.Return)
	assert.InDelta(t, 80.0/100.0-1, got.BuyHold, 1e-9)

	// Long through the crash: peak 1000 at bar 0, trough 500 at bar 1.
	long := RunBacktest(cs, signals{strategy.Buy, strategy.Hold, strategy.Hold}, 1000, 0)
	assert.InDelta(t, -0.5, long.MaxDrawdown, 1e-9)
}

func TestRunEmpty(t *testing.T) {
	assert.Equal(t, BacktestReport{}, RunBacktest(nil, signals{}, 1000, 5))
}
