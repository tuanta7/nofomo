package candle

import (
	"context"
	"strconv"
)

type Storage interface {
	Save(ctx context.Context, symbol, interval string, candles []Candle) error
	Load(ctx context.Context, symbol, interval string) ([]Candle, error)
}

// formatPrice round-trips exactly, so a save/load cycle never shifts a price.
func formatPrice(f float64) string {
	return strconv.FormatFloat(f, 'g', -1, 64)
}
