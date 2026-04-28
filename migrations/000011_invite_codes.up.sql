BEGIN;

CREATE TABLE invite_codes (
    code         VARCHAR(128) PRIMARY KEY,
    max_uses     INTEGER NOT NULL,
    current_uses INTEGER NOT NULL DEFAULT 0,
    is_active    BOOLEAN NOT NULL DEFAULT TRUE
);

INSERT INTO invite_codes (code, max_uses, current_uses, is_active)
VALUES ('STOCET-VIP-100', 100, 0, TRUE);

COMMIT;
