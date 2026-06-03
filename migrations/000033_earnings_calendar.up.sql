CREATE TABLE IF NOT EXISTS earnings_dates (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_profile_id UUID NOT NULL REFERENCES company_profiles(id) ON DELETE CASCADE,
    report_period VARCHAR(20) NOT NULL,
    report_type VARCHAR(10) NOT NULL DEFAULT 'annual',
    expected_date DATE NOT NULL,
    actual_date DATE,
    status VARCHAR(20) NOT NULL DEFAULT 'upcoming',
    notes TEXT,
    eps_estimate NUMERIC(18,4),
    eps_actual NUMERIC(18,4),
    revenue_estimate NUMERIC(24,2),
    revenue_actual NUMERIC(24,2),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_earnings_company ON earnings_dates(company_profile_id);
CREATE INDEX IF NOT EXISTS idx_earnings_expected_date ON earnings_dates(expected_date);
CREATE INDEX IF NOT EXISTS idx_earnings_status ON earnings_dates(status);
