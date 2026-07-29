package dnse

import "context"

const (
	baseURL = "wss://ws-openapi.dnse.com.vn"
)

type WebSocketClient struct {
}

// SubscribeOHLC
// https://developers.dnse.com.vn/docs/guide/market-data/connect#ohlc
func (w *WebSocketClient) SubscribeOHLC(ctx context.Context, symbol string) {

}
