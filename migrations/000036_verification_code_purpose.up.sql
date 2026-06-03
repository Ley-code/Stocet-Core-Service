BEGIN;

ALTER TABLE verification_codes
    ADD COLUMN IF NOT EXISTS purpose VARCHAR(30) NOT NULL DEFAULT 'email_verification';

DROP INDEX IF EXISTS idx_verification_codes_email_active;

CREATE UNIQUE INDEX IF NOT EXISTS idx_verification_codes_email_purpose_active
    ON verification_codes (email, purpose)
    WHERE used_at IS NULL;

COMMIT;
