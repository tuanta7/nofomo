package domain

// Tick is the record of a single trade in the system.
// All data from external sources will eventually be mapped to this.
type Tick struct {
	TradeID   string `json:"trade_id"`
	Symbol    string `json:"symbol"`
	Price     string `json:"price"`
	Quantity  string `json:"quantity"`
	Timestamp int64  `json:"timestamp"`
}
