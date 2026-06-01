BEGIN;

CREATE TYPE corporate_action_type AS ENUM (
    'agm',
    'dividend',
    'rights_issue',
    'bonus_shares',
    'listing',
    'other'
);

CREATE TABLE corporate_actions (
    id                 UUID                   PRIMARY KEY DEFAULT gen_random_uuid(),
    company_profile_id UUID                   NOT NULL REFERENCES company_profiles(id) ON DELETE CASCADE,
    action_type        corporate_action_type  NOT NULL,
    title              TEXT                   NOT NULL,
    event_date         DATE                   NOT NULL,
    ex_date            DATE,
    record_date        DATE,
    description        TEXT,
    amount             NUMERIC(24, 6),
    currency           VARCHAR(10)            NOT NULL DEFAULT 'ETB',
    created_at         TIMESTAMPTZ            NOT NULL DEFAULT NOW(),
    updated_at         TIMESTAMPTZ            NOT NULL DEFAULT NOW()
);

-- Fast upcoming-events query: all future events ordered by date
CREATE INDEX idx_corporate_actions_event_date   ON corporate_actions (event_date);
-- Filter by company
CREATE INDEX idx_corporate_actions_company      ON corporate_actions (company_profile_id, event_date);
-- Filter by type
CREATE INDEX idx_corporate_actions_type         ON corporate_actions (action_type, event_date);

COMMIT;
