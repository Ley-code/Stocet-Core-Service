BEGIN;

CREATE TABLE audit_logs (
    id           UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    actor_id     UUID        NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    actor_email  TEXT        NOT NULL,
    method       VARCHAR(10) NOT NULL,
    path         TEXT        NOT NULL,
    resource_id  TEXT,
    status_code  INT         NOT NULL,
    ip_address   TEXT,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Admin views: most recent first, filter by actor
CREATE INDEX idx_audit_logs_created_at  ON audit_logs (created_at DESC);
CREATE INDEX idx_audit_logs_actor_id    ON audit_logs (actor_id, created_at DESC);

COMMIT;
