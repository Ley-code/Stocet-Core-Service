BEGIN;

ALTER TABLE invite_codes
    ADD COLUMN org_id UUID REFERENCES organizations(id) ON DELETE SET NULL;

ALTER TABLE users
    ADD COLUMN org_id UUID REFERENCES organizations(id) ON DELETE SET NULL;

CREATE INDEX idx_invite_codes_org_id ON invite_codes(org_id);
CREATE INDEX idx_users_org_id ON users(org_id);

COMMIT;
