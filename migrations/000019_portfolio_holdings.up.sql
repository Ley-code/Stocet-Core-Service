BEGIN;

CREATE TABLE portfolio_holdings (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id    UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    ticker     VARCHAR(16) NOT NULL,
    shares     NUMERIC(18,4) NOT NULL CHECK (shares > 0),
    buy_price  NUMERIC(18,4) NOT NULL CHECK (buy_price > 0),
    buy_date   DATE,
    notes      TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (user_id, ticker)
);

CREATE INDEX idx_portfolio_holdings_user_id ON portfolio_holdings(user_id);

COMMIT;
