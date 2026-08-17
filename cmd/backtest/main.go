package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"path"
	"time"

	"github.com/tuanta7/nofomo/internal/market"
	"github.com/tuanta7/nofomo/internal/market/candle"
	"github.com/tuanta7/nofomo/internal/report"
	"github.com/tuanta7/nofomo/internal/strategy"
	"github.com/tuanta7/nofomo/pkg/o11y"
)

func main() {
	// A cold fetch of a year of 5m bars takes ~19s; make Ctrl-C work through it.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	const (
		symbol   = "BTCUSDT"
		interval = "5m"
		days     = 365
		fast     = 12
		slow     = 26
		fee      = 5
		cash     = 10000
	)

	logger, err := o11y.NewLogger(ctx, "")
	if err != nil {
		log.Fatal(err)
	}

	var c market.DataCollector
	var s candle.Storage

	s, err = candle.NewCSVStorage(path.Join(os.TempDir(), "nofomo"))
	if err != nil {
		log.Fatal(err)
	}

	c = market.NewSpotDataCollector(s, logger)

	end := time.Now().UTC()
	candles, err := c.GetCandleHistory(ctx, market.Request{
		Symbol:   symbol,
		Interval: interval,
		Start:    end.AddDate(0, 0, -days),
		End:      end,
	})
	if err != nil {
		log.Fatal(err)
	}
	if len(candles) == 0 {
		log.Fatalf("no candles for %s %s over the last %d days", symbol, interval, days)
	}

	st, err := strategy.NewEMACross(fast, slow)
	if err != nil {
		log.Fatal(err)
	}

	report.Run(candles, st, cash, fee).Result(os.Stdout)
}
