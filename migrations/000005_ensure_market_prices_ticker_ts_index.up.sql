BEGIN;

CREATE INDEX IF NOT EXISTS idx_market_prices_ticker_ts
ON market_prices (ticker_symbol, timestamp DESC);

COMMIT;
