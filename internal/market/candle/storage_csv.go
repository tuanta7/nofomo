package candle

import (
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"
	"strconv"
	"time"
)

var header = []string{"open_time", "open", "high", "low", "close", "volume", "close_time"}

var _ Storage = (*CSVStorage)(nil)

type CSVStorage struct {
	basePath string
}

func NewCSVStorage(path string) (*CSVStorage, error) {
	err := os.MkdirAll(path, 0o755)
	if err != nil {
		return nil, fmt.Errorf("storage: %w", err)
	}

	return &CSVStorage{basePath: path}, nil
}

func (c *CSVStorage) Save(ctx context.Context, symbol, interval string, candles []Candle) error {
	fileName := fmt.Sprintf("%s_%s.csv", symbol, interval)
	f, err := os.Create(path.Join(c.basePath, fileName))
	if err != nil {
		return fmt.Errorf("storage: %w", err)
	}
	defer f.Close()

	w := csv.NewWriter(f)
	if err := w.Write(header); err != nil {
		return fmt.Errorf("storage: %w", err)
	}
	for _, c := range candles {
		row := []string{
			strconv.FormatInt(c.OpenTime.UnixMilli(), 10),
			formatPrice(c.Open),
			formatPrice(c.High),
			formatPrice(c.Low),
			formatPrice(c.Close),
			formatPrice(c.Volume),
			strconv.FormatInt(c.CloseTime.UnixMilli(), 10),
		}
		if err := w.Write(row); err != nil {
			return fmt.Errorf("storage: %w", err)
		}
	}

	w.Flush()
	if err := w.Error(); err != nil {
		return fmt.Errorf("storage: %w", err)
	}
	return f.Close()
}

func (c *CSVStorage) Load(ctx context.Context, symbol, interval string) ([]Candle, error) {
	fileName := fmt.Sprintf("%s_%s.csv", symbol, interval)
	f, err := os.Open(path.Join(c.basePath, fileName))
	if errors.Is(err, fs.ErrNotExist) {
		return nil, nil
	} else if err != nil {
		return nil, fmt.Errorf("storage: %w", err)
	}
	defer f.Close()

	r := csv.NewReader(f)
	r.FieldsPerRecord = len(header)

	if _, err := r.Read(); err != nil {
		if errors.Is(err, io.EOF) {
			return nil, nil
		}
		return nil, fmt.Errorf("storage: %w", err)
	}

	var candles []Candle
	for {
		row, err := r.Read()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("storage: %s: %w", c.basePath, err)
		}

		nums := make([]float64, len(row))
		for i, field := range row {
			if nums[i], err = strconv.ParseFloat(field, 64); err != nil {
				return nil, fmt.Errorf("storage: %s: %w", c.basePath, err)
			}
		}

		candles = append(candles, Candle{
			OpenTime:  time.UnixMilli(int64(nums[0])).UTC(),
			Open:      nums[1],
			High:      nums[2],
			Low:       nums[3],
			Close:     nums[4],
			Volume:    nums[5],
			CloseTime: time.UnixMilli(int64(nums[6])).UTC(),
		})
	}
	return candles, nil
}
