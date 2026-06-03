BEGIN;

DROP INDEX IF EXISTS idx_verification_codes_email_purpose_active;

CREATE UNIQUE INDEX IF NOT EXISTS idx_verification_codes_email_active
    ON verification_codes (email)
    WHERE used_at IS NULL;

ALTER TABLE verification_codes DROP COLUMN IF EXISTS purpose;

COMMIT;
