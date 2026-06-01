BEGIN;

-- One row per user; upserted on first GET/PUT.
CREATE TABLE notification_preferences (
    user_id            UUID        PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    price_alert_push   BOOLEAN     NOT NULL DEFAULT true,
    price_alert_email  BOOLEAN     NOT NULL DEFAULT false,
    news_alert_push    BOOLEAN     NOT NULL DEFAULT true,
    news_alert_email   BOOLEAN     NOT NULL DEFAULT false,
    in_app_enabled     BOOLEAN     NOT NULL DEFAULT true,
    updated_at         TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

COMMIT;
