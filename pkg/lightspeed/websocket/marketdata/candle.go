package marketdata

import (
	"context"
	"fmt"
)

func ohlcClosedChannel(resolution, encoding string) string {
	return fmt.Sprintf("ohlc_closed.%s.%s", resolution, encoding)
}

func SubscribeOHLCClosed(ctx context.Context, resolution, encoding, symbol string) {
	channel := ohlcClosedChannel(resolution, encoding)
	fmt.Printf("Subscribing to channel: %s for symbol: %s\n", channel, symbol)
}
