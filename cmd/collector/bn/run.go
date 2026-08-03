package bn

import (
	"context"
	"fmt"

	"github.com/tuanta7/nofomo/pkg/kafka"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.uber.org/zap"

	"github.com/tuanta7/nofomo/pkg/binance"
	"github.com/tuanta7/nofomo/pkg/binance/futures/coinm"
	"github.com/tuanta7/nofomo/pkg/o11y"
)

var streams = []binance.StreamName{
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

func Run(
	ctx context.Context,
	logger *o11y.Logger,
	meterProvider *o11y.MeterProvider,
) error {
	if err := initMetrics(meterProvider.Meter()); err != nil {
		logger.Error("init metrics", zap.Error(err))
		return fmt.Errorf("init metrics: %w", err)
	}

	client, err := coinm.NewMarketClient()
	if err != nil {
		logger.Error("dial market stream", zap.Error(err))
		return fmt.Errorf("dial market stream: %w", err)
	}
	defer client.Close()

	if err := client.Subscribe(streams...); err != nil {
		logger.Error("subscribe", zap.Error(err))
		return fmt.Errorf("subscribe streams: %w", err)
	}

	go func() {
		<-ctx.Done()
		_ = client.Unsubscribe(streams...)
		_ = client.Close()
	}()

	producer, err := kafka.NewSyncProducer([]string{""})
	if err != nil {
		logger.Error("create kafka sync producer", zap.Error(err))
		return err
	}
	defer producer.Close()

	err = client.HandleStreamingEvents(coinm.Handlers{
		AggTrade: func(t coinm.AggTrade) {
			eventTotal.Add(ctx, 1, metric.WithAttributes(
				attribute.String("symbol", t.Symbol),
				attribute.String("event_type", string(t.EventType)),
			))

			producer.Publish("", nil, nil)

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

	return err
}
