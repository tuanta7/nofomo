package market

import (
	"context"
	"time"

	"github.com/tuanta7/nofomo/internal/market/candle"
	"github.com/tuanta7/nofomo/internal/market/tick"
)

type Request struct {
	Symbol   string
	Interval string // Binance kline interval, e.g. "1m", "5m", "1h"
	Start    time.Time
	End      time.Time
}

type DataCollector interface {
	// GetCandleHistory returns the candles covering [Start, End).
	GetCandleHistory(ctx context.Context, request Request) ([]candle.Candle, error)
	// GetCandleStream returns a channel streaming candles for the given symbol.
	GetCandleStream(ctx context.Context, symbol string) (<-chan candle.Candle, error)
	// GetTickStream returns a channel streaming ticks for the given symbol.
	GetTickStream(ctx context.Context, symbol string) (<-chan tick.Tick, error)
}
