-- ============================================================================
-- Migration 000002: Add News & Market Data Tables
-- ============================================================================
-- These tables mirror the MongoDB collections in the News Aggregation Service
-- and will be populated via a dual-write strategy during the migration period.
-- ============================================================================

BEGIN;

-- -------------------------------------------------------------------------
-- 1. Articles table (mirrors MongoDB "articles" collection)
-- -------------------------------------------------------------------------
CREATE TABLE articles (
    id              uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    title           text NOT NULL,
    url             text NOT NULL,
    source          text NOT NULL,
    published_date  timestamptz NOT NULL,
    content         text NOT NULL,
    summary         text,
    tags            text[],
    image_url       text,
    created_at      timestamptz NOT NULL DEFAULT now(),

    CONSTRAINT uq_articles_url UNIQUE (url)
);

-- Index for time-based queries (latest articles)
CREATE INDEX idx_articles_published_date ON articles (published_date DESC);

-- Index for filtering by source
CREATE INDEX idx_articles_source ON articles (source);

-- GIN index for array containment queries on tags
CREATE INDEX idx_articles_tags ON articles USING GIN (tags);

-- -------------------------------------------------------------------------
-- 2. Market Prices table (mirrors MongoDB "market_prices" collection)
-- -------------------------------------------------------------------------
CREATE TABLE market_prices (
    id              uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    ticker_symbol   text NOT NULL,
    price           numeric NOT NULL,
    percent_change  numeric NOT NULL,
    variation       numeric NOT NULL,
    timestamp       timestamptz NOT NULL,
    source          text NOT NULL DEFAULT 'ESX',
    created_at      timestamptz NOT NULL DEFAULT now()
);

-- Composite index for fast time-series retrieval by ticker
CREATE INDEX idx_market_prices_ticker_ts ON market_prices (ticker_symbol, timestamp DESC);

-- Index for filtering by source
CREATE INDEX idx_market_prices_source ON market_prices (source);

COMMIT;
