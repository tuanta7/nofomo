package market

import (
	"context"

	"github.com/tuanta7/nofomo/internal/market/candle"
	"github.com/tuanta7/nofomo/internal/market/tick"
	"github.com/tuanta7/nofomo/pkg/lightspeed/websocket"
)

type FuturesDataCollector struct {
	client *websocket.Client
	store  candle.Storage
}

func NewFuturesDataCollector(store candle.Storage) *FuturesDataCollector {
	return &FuturesDataCollector{
		store: store,
	}
}

func (f FuturesDataCollector) GetCandleHistory(ctx context.Context, request Request) ([]candle.Candle, error) {
	//TODO implement me
	panic("implement me")
}

func (f FuturesDataCollector) GetCandleStream(ctx context.Context, symbol string) (<-chan candle.Candle, error) {
	//TODO implement me
	panic("implement me")
}

func (f FuturesDataCollector) GetTickStream(ctx context.Context, symbol string) (<-chan tick.Tick, error) {
	//TODO implement me
	panic("implement me")
}
