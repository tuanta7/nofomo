package strategy

import (
	"testing"
	"time"

	"github.com/tuanta7/nofomo/internal/market/candle"
)

func candles(prices ...float64) []candle.Candle {
	out := make([]candle.Candle, len(prices))
	for i, p := range prices {
		out[i] = candle.Candle{
			OpenTime: time.UnixMilli(int64(i) * 60000).UTC(),
			Open:     p,
			High:     p,
			Low:      p,
			Close:    p,
			Volume:   1,
		}
	}
	return out
}

// run feeds every candle through the strategy and returns the signals that fired.
func run(s Strategy, cs []candle.Candle) []Signal {
	var out []Signal
	for i := range cs {
		if sig := s.Evaluate(Context{Candles: cs, Index: i}); sig != Hold {
			out = append(out, sig)
		}
	}
	return out
}

func ramp(from, to float64, n int) []float64 {
	step := (to - from) / float64(n)
	out := make([]float64, n)
	for i := range out {
		out[i] = from + step*float64(i+1)
	}
	return out
}

func TestEMACross(t *testing.T) {
	flat := make([]float64, 40) // warmup: enough bars for both EMAs to settle
	for i := range flat {
		flat[i] = 100
	}

	tests := []struct {
		name   string
		prices []float64
		want   []Signal
	}{
		{
			name:   "flat prices never cross",
			prices: flat,
			want:   nil,
		},
		{
			name:   "rise then fall crosses both ways exactly once",
			prices: append(append(append([]float64{}, flat...), ramp(100, 200, 40)...), ramp(200, 50, 60)...),
			want:   []Signal{Buy, Sell},
		},
		{
			name:   "sustained rise buys once and holds",
			prices: append(append([]float64{}, flat...), ramp(100, 300, 60)...),
			want:   []Signal{Buy},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, err := NewEMACross(12, 26)
			if err != nil {
				t.Fatal(err)
			}

			got := run(s, candles(tt.prices...))
			if len(got) != len(tt.want) {
				t.Fatalf("got %v, want %v", got, tt.want)
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Fatalf("got %v, want %v", got, tt.want)
				}
			}
		})
	}
}

// The warmup window must stay silent: both EMAs start seeded at the same price, so
// any cross there is an artifact of seeding, not of the trend.
func TestEMACrossSilentDuringWarmup(t *testing.T) {
	s, err := NewEMACross(3, 10)
	if err != nil {
		t.Fatal(err)
	}

	cs := candles(ramp(100, 200, 10)...)
	if got := run(s, cs); got != nil {
		t.Fatalf("signalled during warmup: %v", got)
	}
}

func TestEMACrossRejectsBadPeriods(t *testing.T) {
	for _, tt := range []struct{ fast, slow int }{{26, 12}, {12, 12}, {0, 26}, {12, -1}} {
		if _, err := NewEMACross(tt.fast, tt.slow); err == nil {
			t.Errorf("NewEMACross(%d, %d): want error", tt.fast, tt.slow)
		}
	}
}
