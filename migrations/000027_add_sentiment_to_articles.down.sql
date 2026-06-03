DROP INDEX IF EXISTS idx_articles_fts;
DROP INDEX IF EXISTS idx_articles_sentiment;

ALTER TABLE articles
    DROP COLUMN IF EXISTS sentiment_analyzed_at,
    DROP COLUMN IF EXISTS sentiment_tickers,
    DROP COLUMN IF EXISTS sentiment;
