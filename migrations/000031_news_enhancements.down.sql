ALTER TABLE articles
    DROP COLUMN IF EXISTS sentiment_confidence;

DROP TABLE IF EXISTS article_bookmarks;
