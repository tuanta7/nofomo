package market

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
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

func TestMergeSortsAndDedupes(t *testing.T) {
	// Overlapping sets, out of order: the cache tail and a refetch covering it.
	got := candle.Merge(bars(3, 1, 2), bars(2, 4), nil)
	assert.Equal(t, []time.Time{at(1), at(2), at(3), at(4)}, opens(got))
}

func TestMergeEmpty(t *testing.T) {
	assert.Empty(t, candle.Merge(nil, nil))
}

func TestWindowExcludesBarsClosingAtOrAfterEnd(t *testing.T) {
	cs := bars(0, 1, 2, 3, 4)

	// [1, 3): bar 1 closes at 2 (in), bar 2 closes at 3 (out — closes at end).
	assert.Equal(t, []time.Time{at(1)}, opens(candle.Window(cs, at(1), at(3))))
	assert.Equal(t, []time.Time{at(1), at(2)}, opens(candle.Window(cs, at(1), at(4))))
	assert.Empty(t, candle.Window(cs, at(10), at(20)))
}
