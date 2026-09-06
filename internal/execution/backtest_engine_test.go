package execution

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tuanta7/nofomo/internal/strategy"
)

func TestBacktestEngineChargesFeesOnBothFills(t *testing.T) {
	engine := NewBacktestEngine(1000, 5)

	fill := engine.Execute(strategy.Buy, 100)
	assert.True(t, fill.Opened)
	assert.False(t, fill.Closed)
	assert.InDelta(t, 999.5, engine.Equity(100), 1e-9)

	fill = engine.Execute(strategy.Sell, 100)
	assert.True(t, fill.Closed)
	assert.False(t, fill.Profitable)
	assert.InDelta(t, 0.9995*0.9995*1000, engine.Equity(100), 1e-9)
}
