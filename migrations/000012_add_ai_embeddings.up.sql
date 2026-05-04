-- ============================================================================
-- Migration 000012: AI Embeddings (pgvector)
-- ============================================================================
-- Enables the pgvector extension and adds 768-dimensional embedding columns
-- to corporate_news and company_profiles for the RAG Data Assistant.
-- 768 dims matches Google Gemini text-embedding-004.
-- HNSW indexes use cosine distance (vector_cosine_ops) for fast similarity
-- search. The application performs UPDATE-after-INSERT embedding writes from
-- a goroutine (fire-and-forget) plus a periodic backfill for NULL rows.
-- ============================================================================

BEGIN;

CREATE EXTENSION IF NOT EXISTS vector;

ALTER TABLE corporate_news
    ADD COLUMN IF NOT EXISTS embedding vector(768);

ALTER TABLE company_profiles
    ADD COLUMN IF NOT EXISTS embedding vector(768);

-- HNSW indexes for cosine similarity (pgvector >= 0.5).
CREATE INDEX IF NOT EXISTS idx_corporate_news_embedding_hnsw
    ON corporate_news
    USING hnsw (embedding vector_cosine_ops);

CREATE INDEX IF NOT EXISTS idx_company_profiles_embedding_hnsw
    ON company_profiles
    USING hnsw (embedding vector_cosine_ops);

COMMIT;
