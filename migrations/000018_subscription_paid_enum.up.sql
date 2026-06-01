BEGIN;

ALTER TYPE subscription_status_enum ADD VALUE IF NOT EXISTS 'paid' AFTER 'free';

COMMIT;
