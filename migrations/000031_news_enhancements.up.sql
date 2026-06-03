CREATE TABLE IF NOT EXISTS article_bookmarks (
    user_id    uuid        NOT NULL,
    article_id uuid        NOT NULL,
    created_at timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY (user_id, article_id)
);

ALTER TABLE articles
    ADD COLUMN IF NOT EXISTS sentiment_confidence double precision;
