package binance

import (
	"fmt"
	"strings"
)

type Symbol string

const (
	BTCUSDT    Symbol = "btcusdt"
	ETHUSDT    Symbol = "ethusdt"
	BTCUSDPerp Symbol = "btcusd_perp"
	ETHUSDPerp Symbol = "ethusd_perp"
)

type EventType string

const (
	AggregateTrade EventType = "aggTrade"  // aggregate trade (price and quantity)
	MarkPrice      EventType = "markPrice" // futures estimated price of index price
	Kline          EventType = "kline"     // candlestick (OHLCV), suffixed with the interval
)

type Duration string

const (
	OneSecond   Duration = "1s"
	TenSeconds  Duration = "10s"
	OneMinute   Duration = "1m"
	FiveMinutes Duration = "5m"
)

type StreamName struct {
	Symbol      Symbol
	EventType   EventType
	Interval    Duration
	UpdateSpeed Duration
}

func ParseStreamName(name string) (StreamName, error) {
	symbol, rest, ok := strings.Cut(name, "@")
	if !ok || symbol == "" || rest == "" {
		return StreamName{}, fmt.Errorf("invalid stream name: %q", name)
	}
	s := StreamName{Symbol: Symbol(symbol)}

	if event, speed, ok := strings.Cut(rest, "@"); ok {
		if strings.Contains(speed, "_") {
			return StreamName{}, fmt.Errorf("invalid stream name: %q", name)
		}

		s.EventType, s.UpdateSpeed = EventType(event), Duration(speed)
		return s, nil
	}

	if event, interval, ok := strings.Cut(rest, "_"); ok {
		s.EventType, s.Interval = EventType(event), Duration(interval)
		return s, nil
	}

	s.EventType = EventType(rest)
	return s, nil
}

func (s StreamName) String() string {
	switch {
	case s.UpdateSpeed != "":
		return fmt.Sprintf("%s@%s@%s", s.Symbol, s.EventType, s.UpdateSpeed)
	case s.Interval != "":
		return fmt.Sprintf("%s@%s_%s", s.Symbol, s.EventType, s.Interval)
	default:
		return fmt.Sprintf("%s@%s", s.Symbol, s.EventType)
	}
}
