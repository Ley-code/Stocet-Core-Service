BEGIN;

CREATE TABLE IF NOT EXISTS company_profiles (
    id                 uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    ticker_symbol      TEXT NOT NULL UNIQUE,
    company_name       TEXT NOT NULL,
    sector             TEXT,
    nature_of_business TEXT,
    security_type      TEXT,
    listing_date       DATE,
    incorporation_date DATE,
    auditor            TEXT,
    shares             BIGINT,
    capital            NUMERIC(24, 2),
    contact_email      TEXT,
    contact_phone      TEXT,
    contact_address    TEXT,
    website            TEXT,
    created_at         timestamptz NOT NULL DEFAULT now(),
    updated_at         timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_company_profiles_ticker_symbol
    ON company_profiles (ticker_symbol);

CREATE INDEX IF NOT EXISTS idx_company_profiles_company_name
    ON company_profiles (company_name);

CREATE TABLE IF NOT EXISTS board_members (
    id                 uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    company_profile_id uuid NOT NULL REFERENCES company_profiles (id) ON DELETE CASCADE,
    name               TEXT NOT NULL,
    role               TEXT NOT NULL,
    created_at         timestamptz NOT NULL DEFAULT now(),
    updated_at         timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_board_members_company_profile_id
    ON board_members (company_profile_id);

CREATE TABLE IF NOT EXISTS financial_reports (
    id                 uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    company_profile_id uuid NOT NULL REFERENCES company_profiles (id) ON DELETE CASCADE,
    period_end_date    DATE NOT NULL,
    net_profit         NUMERIC(24, 2),
    total_operating_income NUMERIC(24, 2),
    profit_before_tax  NUMERIC(24, 2),
    total_assets       NUMERIC(24, 2),
    total_equity       NUMERIC(24, 2),
    total_liabilities  NUMERIC(24, 2),
    eps                NUMERIC(24, 6),
    customer_deposits  NUMERIC(24, 2),
    loans_and_advances NUMERIC(24, 2),
    extended_metrics   JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at         timestamptz NOT NULL DEFAULT now(),
    updated_at         timestamptz NOT NULL DEFAULT now(),
    CONSTRAINT uq_financial_reports_company_period UNIQUE (company_profile_id, period_end_date)
);

CREATE INDEX IF NOT EXISTS idx_financial_reports_company_profile_id
    ON financial_reports (company_profile_id);

CREATE INDEX IF NOT EXISTS idx_financial_reports_company_period_desc
    ON financial_reports (company_profile_id, period_end_date DESC);

CREATE TABLE IF NOT EXISTS corporate_news (
    id                 uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    company_profile_id uuid NOT NULL REFERENCES company_profiles (id) ON DELETE CASCADE,
    category           TEXT NOT NULL CHECK (category IN ('Financial', 'Governance', 'Regulatory', 'Corporate Action', 'General')),
    title              TEXT NOT NULL,
    body               TEXT,
    source_url         TEXT,
    published_at       timestamptz NOT NULL DEFAULT now(),
    created_at         timestamptz NOT NULL DEFAULT now(),
    updated_at         timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_corporate_news_company_profile_id
    ON corporate_news (company_profile_id);

CREATE INDEX IF NOT EXISTS idx_corporate_news_company_published_at_desc
    ON corporate_news (company_profile_id, published_at DESC);

CREATE OR REPLACE VIEW company_latest_eps AS
SELECT DISTINCT ON (cp.id)
    cp.id AS company_profile_id,
    cp.ticker_symbol,
    cp.company_name,
    fr.eps AS latest_eps,
    fr.period_end_date,
    fr.updated_at
FROM company_profiles cp
JOIN financial_reports fr
    ON fr.company_profile_id = cp.id
WHERE fr.eps IS NOT NULL
ORDER BY cp.id, fr.period_end_date DESC, fr.updated_at DESC;

COMMIT;
