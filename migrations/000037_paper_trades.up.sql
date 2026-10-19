BEGIN;

CREATE TABLE IF NOT EXISTS paper_trades (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  ticker_symbol TEXT NOT NULL,
  side TEXT NOT NULL DEFAULT 'buy' CHECK (side IN ('buy', 'sell')),
  status TEXT NOT NULL DEFAULT 'planned' CHECK (status IN ('planned', 'open', 'closed', 'cancelled')),
  source TEXT NOT NULL DEFAULT 'plan' CHECK (source IN ('plan', 'replay', 'backtest', 'neway')),
  quantity NUMERIC NOT NULL CHECK (quantity > 0),
  entry_price NUMERIC NOT NULL CHECK (entry_price > 0),
  exit_price NUMERIC,
  stop_loss NUMERIC,
  take_profit NUMERIC,
  trailing_stop_pct NUMERIC,
  pnl NUMERIC,
  notes TEXT NOT NULL DEFAULT '',
  opened_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  closed_at TIMESTAMPTZ,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_paper_trades_user_created_at
  ON paper_trades (user_id, created_at DESC);

COMMIT;
