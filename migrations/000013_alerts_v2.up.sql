BEGIN;

-- Extend alerts table with new columns
ALTER TABLE alerts
    ADD COLUMN IF NOT EXISTS alert_type text NOT NULL DEFAULT 'price',
    ADD COLUMN IF NOT EXISTS is_recurring boolean NOT NULL DEFAULT false,
    ADD COLUMN IF NOT EXISTS cooldown_hours integer NOT NULL DEFAULT 24,
    ADD COLUMN IF NOT EXISTS last_triggered_at timestamptz;

-- Trigger history log (one row per firing, survives alert deletion via CASCADE)
CREATE TABLE IF NOT EXISTS alert_trigger_history (
    id              uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    alert_id        uuid NOT NULL REFERENCES alerts(id) ON DELETE CASCADE,
    user_id         uuid NOT NULL,
    ticker          text NOT NULL,
    alert_type      text NOT NULL,
    condition       text,
    threshold       float8,
    triggered_price float8,
    percent_change  float8,
    triggered_at    timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_alert_history_alert_id
    ON alert_trigger_history (alert_id);

CREATE INDEX IF NOT EXISTS idx_alert_history_user_triggered
    ON alert_trigger_history (user_id, triggered_at DESC);

COMMIT;
