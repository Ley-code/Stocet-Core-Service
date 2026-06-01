BEGIN;

CREATE TABLE advisory_reports (
    id         UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id    UUID        NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title      TEXT        NOT NULL DEFAULT '',
    profile    JSONB       NOT NULL DEFAULT '{}',
    result     TEXT        NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_advisory_reports_user_created
    ON advisory_reports (user_id, created_at DESC);

COMMIT;
