BEGIN;

DROP INDEX IF EXISTS idx_invite_codes_code;

ALTER TABLE invite_codes DROP CONSTRAINT invite_codes_pkey;

ALTER TABLE invite_codes
    DROP COLUMN IF EXISTS id,
    DROP COLUMN IF EXISTS description,
    DROP COLUMN IF EXISTS created_by,
    DROP COLUMN IF EXISTS created_at,
    DROP COLUMN IF EXISTS updated_at;

ALTER TABLE invite_codes ADD PRIMARY KEY (code);

COMMIT;
