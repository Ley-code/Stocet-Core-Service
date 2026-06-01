BEGIN;

CREATE TABLE ai_usage_daily (
    user_id       UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    usage_date    DATE NOT NULL DEFAULT CURRENT_DATE,
    request_count INT NOT NULL DEFAULT 0 CHECK (request_count >= 0),
    PRIMARY KEY (user_id, usage_date)
);

COMMIT;
