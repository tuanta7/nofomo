package market

import (
	"github.com/tuanta7/nofomo/internal/market/candle"
	"github.com/tuanta7/nofomo/internal/market/tick"
)

type CandleBus interface {
	NextCandle() candle.Candle
}

type TickerBus interface {
	NextTick() tick.Tick
}
