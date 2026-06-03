BEGIN;

CREATE TABLE IF NOT EXISTS company_filings (
    id                 uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    company_profile_id uuid NOT NULL REFERENCES company_profiles (id) ON DELETE CASCADE,
    filing_type        TEXT NOT NULL CHECK (filing_type IN (
                           'Annual Report', 'Interim Report', 'Prospectus',
                           'Rights Issue', 'AGM Notice', 'Material Disclosure', 'Other'
                       )),
    title              TEXT NOT NULL,
    description        TEXT,
    document_url       TEXT,
    filed_at           TIMESTAMPTZ NOT NULL DEFAULT now(),
    created_at         TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at         TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_company_filings_company_profile_id
    ON company_filings (company_profile_id);

CREATE INDEX IF NOT EXISTS idx_company_filings_company_filed_at_desc
    ON company_filings (company_profile_id, filed_at DESC);

COMMIT;
