BEGIN;

DROP INDEX IF EXISTS idx_users_org_id;
DROP INDEX IF EXISTS idx_invite_codes_org_id;

ALTER TABLE users DROP COLUMN IF EXISTS org_id;
ALTER TABLE invite_codes DROP COLUMN IF EXISTS org_id;

COMMIT;
