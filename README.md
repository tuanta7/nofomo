# NoFOMO

Algorithmic trading toolkit in Go.

Binance coin-M futures are used for market data, the APIs are free, need no auth, and the volume is high enough for the 
signals to mean something. 

DNSE [LightSpeed](https://entradex.dnse.com.vn/thong-tin-ca-nhan/light-speed) is the intended venue for live VN30 
index futures.

## Quick Start

Backtester replays a year of historical candles  through a strategy and prints how it would have done. 

Self-contained: no Kafka, no database, one CSV cache.

```bash
go run ./cmd/backtester
```

```
     candles  105119  (2025-08-04 → 2026-08-04)
      trades    1927
    win rate   22.0%
      return  -90.0%
  buy & hold  -44.2%
      max DD  -90.3%
```

## Roadmap

- **Core logic and backtesting**: market data → candles → indicators → strategy. Backtester done; indicators beyond 
   EMA pending.
- **Paper trading** — run strategies against the live stream, no real orders.
- **Live trading** — order execution and risk management via LightSpeed.
