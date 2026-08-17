package candle

import (
	"sort"
	"time"
)

// Candle is a single OHLCV bar. Also known as K-line or Kyosen-line
type Candle struct {
	OpenTime  time.Time `json:"openTime"`
	CloseTime time.Time `json:"closeTime"`
	Open      float64   `json:"o"`
	High      float64   `json:"h"`
	Low       float64   `json:"l"`
	Close     float64   `json:"c"`
	Volume    float64   `json:"v"`
}

// Merge combines candle sets into one chronological run, dropping duplicate bars.
func Merge(sets ...[]Candle) []Candle {
	var all []Candle
	for _, set := range sets {
		all = append(all, set...)
	}
	sort.Slice(all, func(i, j int) bool {
		return all[i].OpenTime.Before(all[j].OpenTime)
	})

	out := all[:0]
	for _, c := range all {
		if len(out) > 0 && out[len(out)-1].OpenTime.Equal(c.OpenTime) {
			continue
		}
		out = append(out, c)
	}
	return out
}

func Window(candles []Candle, start, end time.Time) []Candle {
	var out []Candle
	for _, c := range candles {
		if c.OpenTime.Before(start) || !c.CloseTime.Before(end) {
			continue
		}
		out = append(out, c)
	}
	return out
}
