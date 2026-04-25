BEGIN;

ALTER TABLE users
    ADD COLUMN IF NOT EXISTS full_name TEXT,
    ADD COLUMN IF NOT EXISTS phone_number VARCHAR(20),
    ADD COLUMN IF NOT EXISTS is_verified BOOLEAN NOT NULL DEFAULT FALSE;

-- Keep existing users able to log in after rollout.
UPDATE users
SET is_verified = TRUE
WHERE is_verified = FALSE;

CREATE TABLE IF NOT EXISTS verification_codes (
    id         uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    email      TEXT NOT NULL,
    code_hash  TEXT NOT NULL,
    expires_at timestamptz NOT NULL,
    used_at    timestamptz,
    created_at timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_verification_codes_email_expires_at
    ON verification_codes (email, expires_at DESC);

CREATE UNIQUE INDEX IF NOT EXISTS idx_verification_codes_email_active
    ON verification_codes (email)
    WHERE used_at IS NULL;

COMMIT;
