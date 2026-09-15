package websocket

import "context"

const (
	baseURL        = "wss://ws-openapi.dnse.com.vn"
	sandboxBaseURL = "wss://ws-sb-openapi.dnse.com.vn"
)

type Client struct {
}

// SubscribeOHLC
// https://developers.dnse.com.vn/docs/guide/market-data/connect#ohlc
func (w *Client) SubscribeOHLC(ctx context.Context, symbol string) {

}
