package coinm

import (
	"encoding/json"

	"github.com/tuanta7/nofomo/pkg/binance"
)

type MarketClient struct {
	*binance.WebsocketClient
}

func NewMarketClient() (*MarketClient, error) {
	// Market streams for coin-margined futures
	// https://developers.binance.com/en/docs/catalog/core-trading-derivatives-trading-coin-m-futures/api/ws-streams/~
	c, err := binance.NewMarketClient("wss://dstream.binance.com/stream")
	if err != nil {
		return nil, err
	}
	return &MarketClient{WebsocketClient: c}, nil
}

type AggTrade struct {
	EventType    binance.EventType `json:"e"`
	EventTime    int64             `json:"E"`
	AggTradeID   int64             `json:"a"`
	Symbol       string            `json:"s"`
	Price        string            `json:"p"`
	Quantity     string            `json:"q"`
	FirstTradeID int64             `json:"f"`
	LastTradeID  int64             `json:"l"`
	TradeTime    int64             `json:"T"`
	IsBuyerMaker bool              `json:"m"`  // Is the buyer the market maker
	SymbolType   int64             `json:"st"` // 1= USD-M, 2= COIN-M
}

type Kline struct {
	EventType binance.EventType `json:"e"`
	EventTime int64             `json:"E"`
	Symbol    string            `json:"s"`
	Payload   struct {          // Kline (a.k.a Candlestick) data
		OpenTime     int64  `json:"t"`
		CloseTime    int64  `json:"T"`
		Symbol       string `json:"s"`
		Interval     string `json:"i"`
		FirstTradeID int64  `json:"f"`
		LastTradeID  int64  `json:"L"`
		Open         string `json:"o"`
		Close        string `json:"c"`
		High         string `json:"h"`
		Low          string `json:"l"`
		Volume       string `json:"v"` // in contracts
		TradeCount   int64  `json:"n"`
		IsClosed     bool   `json:"x"` // false until the interval ends; the candle is revised in place
		BaseVolume   string `json:"q"`
		TakerVolume  string `json:"V"` // taker buy volume, in contracts
		TakerBase    string `json:"Q"`
	} `json:"k"`
}

type Handlers struct {
	AggTrade func(AggTrade)
	Kline    func(Kline)
}

func (c *MarketClient) HandleStreamingEvents(h Handlers) error {
	for {
		stream, data, err := c.ReadMessage()
		if err != nil {
			return err
		}

		name, err := binance.ParseStreamName(stream)
		if err != nil {
			return err
		}

		switch name.EventType {
		case binance.AggregateTrade:
			err = handle(h.AggTrade, data)
		case binance.Kline:
			err = handle(h.Kline, data)
		}
		if err != nil {
			return err
		}
	}
}

func handle[T any](handler func(T), data json.RawMessage) error {
	if handler == nil {
		return nil
	}
	var event T
	if err := json.Unmarshal(data, &event); err != nil {
		return err
	}
	handler(event)
	return nil
}
