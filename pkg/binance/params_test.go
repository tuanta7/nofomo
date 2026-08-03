package binance

import (
	"testing"
)

func TestStreamNameRoundTrip(t *testing.T) {
	for _, want := range []StreamName{
		{Symbol: ETHUSDT, EventType: AggregateTrade},
		{Symbol: ETHUSDT, EventType: Kline, Interval: FiveMinutes},
		{Symbol: ETHUSDT, EventType: MarkPrice, UpdateSpeed: FiveMinutes},
	} {
		got, err := ParseStreamName(want.String())
		if err != nil {
			t.Fatalf("%s: %v", want, err)
		}
		if got != want {
			t.Errorf("%s: got %+v, want %+v", want, got, want)
		}
	}
}
