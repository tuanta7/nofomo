package market

import (
	"testing"
	"time"

	"github.com/tuanta7/nofomo/internal/market/candle"
)

func at(min int) time.Time { return time.UnixMilli(int64(min) * 60000).UTC() }

// bars builds one-minute candles opening at the given minute offsets, priced so each
// bar is identifiable.
func bars(mins ...int) []candle.Candle {
	out := make([]candle.Candle, len(mins))
	for i, m := range mins {
		out[i] = candle.Candle{OpenTime: at(m), CloseTime: at(m + 1), Close: float64(m)}
	}
	return out
}

func opens(cs []candle.Candle) []time.Time {
	out := make([]time.Time, len(cs))
	for i, c := range cs {
		out[i] = c.OpenTime
	}
	return out
}

func equal(got []candle.Candle, wantMins ...int) bool {
	if len(got) != len(wantMins) {
		return false
	}
	for i, m := range wantMins {
		if !got[i].OpenTime.Equal(at(m)) {
			return false
		}
	}
	return true
}

func TestMergeSortsAndDedupes(t *testing.T) {
	// Overlapping sets, out of order: the cache tail and a refetch covering it.
	got := candle.Merge(bars(3, 1, 2), bars(2, 4), nil)
	if !equal(got, 1, 2, 3, 4) {
		t.Errorf("merge = %v, want opens at minutes 1,2,3,4", opens(got))
	}
}

func TestMergeEmpty(t *testing.T) {
	if got := candle.Merge(nil, nil); len(got) != 0 {
		t.Errorf("merge(nil, nil) = %v, want empty", opens(got))
	}
}

func TestWindowExcludesBarsClosingAtOrAfterEnd(t *testing.T) {
	cs := bars(0, 1, 2, 3, 4)

	// [1, 3): bar 1 closes at 2 (in), bar 2 closes at 3 (out — closes at end).
	if got := candle.Window(cs, at(1), at(3)); !equal(got, 1) {
		t.Errorf("window = %v, want only the bar opening at minute 1", opens(got))
	}
	if got := candle.Window(cs, at(1), at(4)); !equal(got, 1, 2) {
		t.Errorf("window = %v, want bars opening at minutes 1,2", opens(got))
	}
	if got := candle.Window(cs, at(10), at(20)); len(got) != 0 {
		t.Errorf("window outside range = %v, want empty", opens(got))
	}
}
