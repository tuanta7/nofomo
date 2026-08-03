# Backtester 

Backtesting trading strategies with Binance data.

```mermaid
flowchart LR
    B[Binance WS] --> I[Collector]
    I --> K[(Kafka)]
    K --> T[Backtester]
    K --> P[Price Board]
    T --> C[(ClickHouse)]
    C --> T
```