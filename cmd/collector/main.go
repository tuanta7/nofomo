package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"
	"time"

	"github.com/tuanta7/nofomo/cmd/collector/binance"
	"github.com/tuanta7/nofomo/pkg/o11y"
	"go.uber.org/zap"

	_ "github.com/joho/godotenv/autoload"
)

const (
	serviceName = "collector-service"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	logger, err := o11y.NewLogger(ctx, serviceName)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = logger.Sync()
		_ = logger.Shutdown(shutdownCtx)
	}()

	meterProvider, err := o11y.NewMeterProvider(ctx, serviceName)
	if err != nil {
		logger.Fatal("new meter provider", zap.Error(err))
	}
	defer func() {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = meterProvider.Shutdown(shutdownCtx)
	}()

	if err = binance.Run(ctx, logger, meterProvider); err != nil {
		logger.Fatal("run collector", zap.Error(err))
	}
}
