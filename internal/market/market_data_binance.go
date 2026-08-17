package market

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	spot "github.com/binance/binance-connector-go/clients/spot"
	"github.com/binance/binance-connector-go/clients/spot/src/restapi/models"
	"github.com/binance/binance-connector-go/common/v2/common"
	"github.com/tuanta7/nofomo/internal/market/candle"
	"github.com/tuanta7/nofomo/internal/market/tick"
	"github.com/tuanta7/nofomo/pkg/o11y"
	"go.uber.org/zap"
)

// Spot klines cap at 1000 rows per request. Each costs weight 2 against a
// 6000/min budget, so pagination is nowhere near the limit.
const maxKlinesPerPage = 1000

var _ DataCollector = (*SpotDataCollector)(nil)

type SpotDataCollector struct {
	client *spot.BinanceSpotClient
	store  candle.Storage
	logger *o11y.Logger
}

func NewSpotDataCollector(store candle.Storage, logger *o11y.Logger) *SpotDataCollector {
	return &SpotDataCollector{
		// Market data needs no API key.
		client: spot.NewBinanceSpotClient(
			spot.WithRestAPI(common.NewConfigurationRestAPI()),
		),
		store:  store,
		logger: logger,
	}
}

// GetCandleHistory returns the candles covering [Start, End). It also
// extends the cache at the tail when the window reaches past what's stored.
func (c *SpotDataCollector) GetCandleHistory(ctx context.Context, r Request) ([]candle.Candle, error) {
	stored, err := c.store.Load(ctx, r.Symbol, r.Interval)
	if err != nil {
		return nil, err
	}

	from := r.Start
	if len(stored) > 1 {
		slack := stored[1].OpenTime.Sub(stored[0].OpenTime)
		if stored[0].OpenTime.Before(r.Start.Add(slack)) {
			from = stored[len(stored)-1].CloseTime
		}
	}

	var fetched []candle.Candle
	if from.Before(r.End) {
		c.logger.Info("fetching candles",
			zap.String("symbol", r.Symbol),
			zap.String("interval", r.Interval),
			zap.String("start_from", from.Format(time.DateTime)),
		)

		fetched, err = c.fetchKlines(ctx, r.Symbol, r.Interval, from, r.End)
		if err != nil {
			return nil, err
		}
	}

	merged := candle.Merge(stored, fetched)
	if len(fetched) > 0 {
		err = c.store.Save(ctx, r.Symbol, r.Interval, merged)
		if err != nil {
			return nil, err
		}
	}

	return candle.Window(merged, r.Start, r.End), nil
}

// fetchKlines returns closed candlesticks in [start, end), paginating as needed.
// The still-forming bar is excluded: its OHLC would change after we read it.
func (c *SpotDataCollector) fetchKlines(ctx context.Context, symbol, interval string, start, end time.Time) ([]candle.Candle, error) {
	var bars []candle.Candle

	for cursor := start; cursor.Before(end); {
		resp, err := c.client.RestApi.MarketAPI.Klines(ctx).
			Symbol(symbol).
			Interval(models.KlinesIntervalParameter(interval)).
			StartTime(cursor.UnixMilli()).
			EndTime(end.UnixMilli()).
			Limit(maxKlinesPerPage).
			Execute()
		if err != nil {
			return nil, fmt.Errorf("binance: klines: %w", err)
		}

		page := resp.Data.Items
		if len(page) == 0 {
			break
		}

		var lastOpen time.Time
		for _, k := range page {
			b, err := toCandle(k)
			if err != nil {
				return nil, err
			}
			lastOpen = b.OpenTime

			// Binance includes the in-progress bar; drop it and anything past
			// the requested window.
			if b.CloseTime.After(end) || !b.CloseTime.Before(time.Now()) {
				continue
			}
			bars = append(bars, b)
		}

		// Advance past the last bar Binance returned, closed or not, so a
		// trailing in-progress bar can't stall the loop.
		next := lastOpen.Add(time.Millisecond)
		if !next.After(cursor) {
			break
		}
		cursor = next

		if len(page) < maxKlinesPerPage {
			break
		}
	}

	return bars, nil
}

// toCandle decodes Binance's kline representation, a heterogeneous array with
// numbers quoted as strings:
//
//	[openTime, "open", "high", "low", "close", "volume", closeTime, ...]
func toCandle(k models.KlinesItem) (candle.Candle, error) {
	if len(k.Items) < 7 {
		return candle.Candle{}, fmt.Errorf("binance: kline: got %d fields, want at least 7", len(k.Items))
	}
	if k.Items[0].Int64 == nil || k.Items[6].Int64 == nil {
		return candle.Candle{}, errors.New("binance: kline: missing timestamps")
	}

	b := candle.Candle{
		OpenTime:  time.UnixMilli(*k.Items[0].Int64).UTC(),
		CloseTime: time.UnixMilli(*k.Items[6].Int64).UTC(),
	}

	for i, dst := range []*float64{&b.Open, &b.High, &b.Low, &b.Close, &b.Volume} {
		s := k.Items[i+1].String
		if s == nil {
			return candle.Candle{}, fmt.Errorf("binance: kline: field %d is not a string", i+1)
		}
		v, err := strconv.ParseFloat(*s, 64)
		if err != nil {
			return candle.Candle{}, fmt.Errorf("binance: kline: %w", err)
		}
		*dst = v
	}
	return b, nil
}

func (c *SpotDataCollector) GetCandleStream(ctx context.Context, symbol string) (<-chan candle.Candle, error) {
	return nil, errors.New("candle: streaming not implemented")
}

func (c *SpotDataCollector) GetTickStream(ctx context.Context, symbol string) (<-chan tick.Tick, error) {
	return nil, errors.New("tick: streaming not implemented")
}
