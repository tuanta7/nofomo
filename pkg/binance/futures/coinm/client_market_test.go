package coinm

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/websocket"
	"github.com/tuanta7/nofomo/pkg/binance"
)

func TestAggTradeUnmarshal(t *testing.T) {
	const payload = `
		{
			"e":"aggTrade",
			"E":1591261134288,
			"a":424951,
			"s":"BTCUSD_200626",
			"p":"9643.5",
			"q":"2",
			"f":606073,
			"l":606073,
			"T":1591261134199,
			"m":false,
			"st":1
		}`

	var got AggTrade
	if err := json.Unmarshal([]byte(payload), &got); err != nil {
		t.Fatal(err)
	}

	want := AggTrade{
		EventType: binance.AggregateTrade, EventTime: 1591261134288, AggTradeID: 424951,
		Symbol: "BTCUSD_200626", Price: "9643.5", Quantity: "2",
		FirstTradeID: 606073, LastTradeID: 606073, TradeTime: 1591261134199, IsBuyerMaker: false,
		SymbolType: 1,
	}
	if got != want {
		t.Errorf("got %+v, want %+v", got, want)
	}
}

// One connection carries every stream: acks are skipped and each event reaches
// its own handler, whatever the kline interval.
func TestReadDispatch(t *testing.T) {
	const (
		ack      = `{"result":null,"id":1}`
		aggTrade = `{"stream":"btcusd_perp@aggTrade","data":{"e":"aggTrade","p":"9643.5"}}`
		kline    = `{"stream":"btcusd_perp@kline_5m","data":{"e":"kline","s":"BTCUSD_PERP","k":{"i":"5m","c":"9639.8","x":true}}}`
	)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := (&websocket.Upgrader{}).Upgrade(w, r, nil)
		if err != nil {
			t.Error(err)
			return
		}
		defer conn.Close()

		var sub map[string]any
		if err := conn.ReadJSON(&sub); err != nil {
			t.Error(err)
			return
		}
		if sub["method"] != "SUBSCRIBE" {
			t.Errorf("got method %v", sub["method"])
		}

		for _, msg := range []string{ack, aggTrade, kline} {
			if err := conn.WriteMessage(websocket.TextMessage, []byte(msg)); err != nil {
				t.Error(err)
				return
			}
		}
	}))
	defer srv.Close()

	base, err := binance.NewMarketClient("ws" + strings.TrimPrefix(srv.URL, "http"))
	if err != nil {
		t.Fatal(err)
	}
	client := &MarketClient{WebsocketClient: base}
	defer client.Close()

	if err := client.Subscribe(
		binance.StreamName{Symbol: binance.BTCUSDT, EventType: binance.AggregateTrade},
		binance.StreamName{Symbol: binance.BTCUSDT, EventType: binance.Kline, Interval: binance.FiveMinutes},
	); err != nil {
		t.Fatal(err)
	}

	var gotTrade AggTrade
	var gotKline Kline
	// Read only returns once the server hangs up, having delivered both events.
	err = client.HandleStreamingEvents(Handlers{
		AggTrade: func(a AggTrade) { gotTrade = a },
		Kline:    func(k Kline) { gotKline = k },
	})
	if _, ok := err.(*websocket.CloseError); !ok {
		t.Fatalf("want close error, got %v", err)
	}

	if gotTrade.EventType != binance.AggregateTrade || gotTrade.Price != "9643.5" {
		t.Errorf("aggTrade: got %+v", gotTrade)
	}
	if gotKline.Payload.Interval != "5m" || gotKline.Payload.Close != "9639.8" || !gotKline.Payload.IsClosed {
		t.Errorf("kline: got %+v", gotKline.Payload)
	}
}
