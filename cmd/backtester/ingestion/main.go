package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.uber.org/zap"

	"github.com/tuanta7/nofomo/pkg/binance"
	"github.com/tuanta7/nofomo/pkg/binance/futures/coinm"
	"github.com/tuanta7/nofomo/pkg/o11y"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	logger, err := o11y.NewLogger(ctx, "ingestion-service")
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = logger.Sync()
		_ = logger.Shutdown(shutdownCtx)
	}()

	meterProvider, err := o11y.NewMeterProvider(ctx, "ingestion-service")
	if err != nil {
		logger.Fatal("new meter provider", zap.Error(err))
	}
	defer func() {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = meterProvider.Shutdown(shutdownCtx)
	}()

	err = initMetrics(meterProvider.Meter())
	if err != nil {
		logger.Fatal("init metrics", zap.Error(err))
	}

	client, err := coinm.NewMarketClient()
	if err != nil {
		logger.Fatal("dial market stream", zap.Error(err))
	}

	streams := []binance.StreamName{
		{
			Symbol:    binance.BTCUSDPerp,
			EventType: binance.AggregateTrade,
		},
		{
			Symbol:    binance.BTCUSDPerp,
			EventType: binance.Kline,
			Interval:  binance.FiveMinutes,
		},
	}

	if err := client.Subscribe(streams...); err != nil {
		logger.Fatal("subscribe", zap.Error(err))
	}

	go func() {
		<-ctx.Done()
		_ = client.Unsubscribe(streams...)
		_ = client.Close()
	}()

	err = client.HandleStreamingEvents(coinm.Handlers{
		AggTrade: func(t coinm.AggTrade) {
			eventTotal.Add(ctx, 1, metric.WithAttributes(
				attribute.String("symbol", t.Symbol),
				attribute.String("event_type", string(t.EventType)),
			))
			logger.Info("agg trade",
				zap.Any("data", t),
			)
		},
		Kline: func(k coinm.Kline) {
			eventTotal.Add(ctx, 1, metric.WithAttributes(
				attribute.String("symbol", k.Symbol),
				attribute.String("event_type", string(k.EventType)),
			))
			logger.Info("kline",
				zap.Any("data", k),
			)
		},
	})
	if err != nil && ctx.Err() == nil {
		logger.Error("market stream ended", zap.Error(err))
	}
}
