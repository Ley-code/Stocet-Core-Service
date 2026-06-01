BEGIN;
DROP TABLE IF EXISTS alert_trigger_history;
ALTER TABLE alerts
    DROP COLUMN IF EXISTS last_triggered_at,
    DROP COLUMN IF EXISTS cooldown_hours,
    DROP COLUMN IF EXISTS is_recurring,
    DROP COLUMN IF EXISTS alert_type;
COMMIT;
