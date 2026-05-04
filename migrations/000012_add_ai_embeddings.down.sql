BEGIN;

DROP INDEX IF EXISTS idx_company_profiles_embedding_hnsw;
DROP INDEX IF EXISTS idx_corporate_news_embedding_hnsw;

ALTER TABLE company_profiles DROP COLUMN IF EXISTS embedding;
ALTER TABLE corporate_news   DROP COLUMN IF EXISTS embedding;

-- Note: the `vector` extension itself is left enabled in case other
-- tables/rows still reference it.

COMMIT;
