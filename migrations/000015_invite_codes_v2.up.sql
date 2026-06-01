BEGIN;

-- Replace the code-as-PK design with a proper UUID primary key.
-- Existing rows (e.g. STOCET-VIP-100) are preserved; code becomes a unique index.

ALTER TABLE invite_codes DROP CONSTRAINT invite_codes_pkey;

ALTER TABLE invite_codes
    ADD COLUMN id          UUID        NOT NULL DEFAULT gen_random_uuid(),
    ADD COLUMN description TEXT        NOT NULL DEFAULT '',
    ADD COLUMN created_by  UUID        REFERENCES users(id) ON DELETE SET NULL,
    ADD COLUMN created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    ADD COLUMN updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW();

ALTER TABLE invite_codes ADD PRIMARY KEY (id);

CREATE UNIQUE INDEX idx_invite_codes_code ON invite_codes(code);

COMMIT;
