package lightspeed

import "time"

// OHLC is a completed DNSE candle. Time and LastUpdated are Unix seconds.
type OHLC struct {
	Symbol      string     `json:"symbol" msgpack:"symbol"`
	Resolution  Resolution `json:"resolution" msgpack:"resolution"`
	Open        float64    `json:"open" msgpack:"open"`
	High        float64    `json:"high" msgpack:"high"`
	Low         float64    `json:"low" msgpack:"low"`
	Close       float64    `json:"close" msgpack:"close"`
	Volume      float64    `json:"volume" msgpack:"volume"`
	Time        int64      `json:"time" msgpack:"time"`
	LastUpdated int64      `json:"lastUpdated" msgpack:"lastUpdated"`
	Type        string     `json:"type" msgpack:"type"`
	ReceivedAt  time.Time  `json:"-" msgpack:"-"`
}
