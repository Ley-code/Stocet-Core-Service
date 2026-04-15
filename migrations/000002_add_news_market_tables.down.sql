-- ============================================================================
-- Migration 000002 DOWN: Drop News & Market Data Tables
-- ============================================================================

BEGIN;

DROP TABLE IF EXISTS market_prices;
DROP TABLE IF EXISTS articles;

COMMIT;
