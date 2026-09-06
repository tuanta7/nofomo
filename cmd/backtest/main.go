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
	"go.uber.org/zap"
)

const (
	symbol   = "BTCUSDT"
	interval = "5m"
	days     = 365
	fast     = 9
	slow     = 21
	cash     = 1000
	fee      = 10 // 0.10%
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	logger, err := o11y.NewLogger(ctx, "")
	if err != nil {
		log.Fatal(err)
	}

	storagePath := path.Join(os.TempDir(), "nofomo")
	candleStorage, err := candle.NewCSVStorage(storagePath)
	if err != nil {
		logger.Fatal("failed to create candle storage", zap.Error(err))
	}
	logger.Info("candle storage", zap.String("path", storagePath))

	spotCollector := market.NewSpotDataCollector(candleStorage, logger)
	end := time.Now().UTC()
	candles, err := spotCollector.GetCandleHistory(ctx, market.Request{
		Symbol:   symbol,
		Interval: interval,
		Start:    end.AddDate(0, 0, -days),
		End:      end,
	})
	if err != nil {
		logger.Fatal("failed to get candle history", zap.Error(err))
	} else if len(candles) == 0 {
		logger.Fatal("no candles", zap.String("symbol", symbol), zap.String("interval", interval), zap.Int("days", days))
	}

	ema, err := strategy.NewEMACross(fast, slow)
	if err != nil {
		logger.Fatal("failed to create EMA strategy", zap.Error(err))
	}

	report.RunBacktest(candles, ema, cash, fee).PrintResult(os.Stdout)
}
