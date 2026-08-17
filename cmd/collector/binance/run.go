package binance

import (
	"context"
	"fmt"

	spot "github.com/binance/binance-connector-go/clients/spot"
	"github.com/binance/binance-connector-go/clients/spot/src/websocketstreams/models"
	"github.com/binance/binance-connector-go/common/v2/common"
	"github.com/tuanta7/nofomo/pkg/kafka"
	"github.com/tuanta7/nofomo/pkg/o11y"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.uber.org/zap"
)

const (
	symbol        = "btcusdt"
	klineInterval = models.KlineIntervalParameterInterval5m
)

func Run(
	ctx context.Context,
	logger *o11y.Logger,
	meterProvider *o11y.MeterProvider,
) error {
	if err := initMetrics(meterProvider.Meter()); err != nil {
		logger.Error("init metrics", zap.Error(err))
		return fmt.Errorf("init metrics: %w", err)
	}

	producer, err := kafka.NewSyncProducer([]string{""})
	if err != nil {
		logger.Error("create kafka sync producer", zap.Error(err))
		return err
	}
	defer producer.Close()

	client := spot.NewBinanceSpotClient(
		spot.WithWebsocketStreams(common.NewConfigurationWebsocketStreams()),
	)

	ws := client.WebsocketStreams
	if err := ws.Connect([]string{}); err != nil {
		logger.Error("dial market stream", zap.Error(err))
		return fmt.Errorf("dial market stream: %w", err)
	}
	defer ws.CloseWebSocketStreamConnection()

	aggTrades, err := ws.DefaultAPI.AggTrade().Symbol(symbol).Execute()
	if err != nil {
		logger.Error("subscribe agg trade", zap.Error(err))
		return fmt.Errorf("subscribe agg trade: %w", err)
	}
	aggTrades.On("message", func(t models.AggTradeResponse) {
		eventTotal.Add(ctx, 1, metric.WithAttributes(
			attribute.String("symbol", t.GetSmalls()),
			attribute.String("event_type", t.GetSmalle()),
		))

		producer.Publish("", nil, nil)

		logger.Info("agg trade", zap.Any("data", t))
	})

	klines, err := ws.DefaultAPI.Kline().Symbol(symbol).Interval(klineInterval).Execute()
	if err != nil {
		logger.Error("subscribe kline", zap.Error(err))
		return fmt.Errorf("subscribe kline: %w", err)
	}
	klines.On("message", func(k models.KlineResponse) {
		eventTotal.Add(ctx, 1, metric.WithAttributes(
			attribute.String("symbol", k.GetSmalls()),
			attribute.String("event_type", k.GetSmalle()),
		))

		logger.Info("kline", zap.Any("data", k))
	})

	<-ctx.Done()
	aggTrades.Unsubscribe()
	klines.Unsubscribe()
	return ctx.Err()
}
