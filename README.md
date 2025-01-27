# CryptoSentry

**CryptoSentry** is a Go application that tracks multiple cryptocurrencies via the CoinGecko API, logs real-time price updates (every 10 seconds by default), alerts you when predefined price thresholds are reached, and logs when a coin’s price jumps by a certain percentage within a given timeframe.

Key Features:
- Periodically fetches the latest prices for configured coins from CoinGecko’s free public API (no API key required).
- Supports user-defined thresholds for each coin, triggering log-based “alerts” when prices cross those thresholds.
- Detects significant price jumps (e.g., 10% within 1 minute).
- Simple JSON-based configuration (`config.json`) so you can easily adjust coins, thresholds, and intervals without editing code.
- Logs coin price trends (↑ or ↓) by comparing each new price to the last fetched price.
- Modularized structure, with separate packages for configuration, fetching, and alert logic.