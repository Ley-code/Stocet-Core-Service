BEGIN;

DROP INDEX IF EXISTS idx_verification_codes_email_active;
DROP INDEX IF EXISTS idx_verification_codes_email_expires_at;
DROP TABLE IF EXISTS verification_codes;

ALTER TABLE users
    DROP COLUMN IF EXISTS is_verified,
    DROP COLUMN IF EXISTS phone_number,
    DROP COLUMN IF EXISTS full_name;

COMMIT;
