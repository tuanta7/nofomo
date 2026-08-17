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
	GetCandleHistory(ctx context.Context, request Request) ([]candle.Candle, error)
	GetCandleStream(ctx context.Context, symbol string) (<-chan candle.Candle, error)
	GetTickStream(ctx context.Context, symbol string) (<-chan tick.Tick, error)
}
