package bn

import (
	"go.opentelemetry.io/otel/metric"
)

var (
	eventTotal metric.Int64Counter
)

func initMetrics(meter metric.Meter) (err error) {
	eventTotal, err = meter.Int64Counter("events_total",
		metric.WithDescription("Total number of websocket events"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return err
	}

	return nil
}
