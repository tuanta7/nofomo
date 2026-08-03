package spot

import "github.com/tuanta7/nofomo/pkg/binance"

type MarketClient struct {
	*binance.WebsocketClient
}
