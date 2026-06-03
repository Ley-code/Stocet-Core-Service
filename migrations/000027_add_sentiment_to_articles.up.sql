ALTER TABLE articles
    ADD COLUMN IF NOT EXISTS sentiment              text,
    ADD COLUMN IF NOT EXISTS sentiment_tickers      text[],
    ADD COLUMN IF NOT EXISTS sentiment_analyzed_at  timestamptz;

CREATE INDEX IF NOT EXISTS idx_articles_sentiment
    ON articles (sentiment, published_date DESC)
    WHERE sentiment IS NOT NULL;

CREATE INDEX IF NOT EXISTS idx_articles_fts
    ON articles USING GIN (to_tsvector('english', title || ' ' || coalesce(content, '')));
