BEGIN;

CREATE TABLE IF NOT EXISTS capital_service_providers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    license_number TEXT NOT NULL,
    category TEXT NOT NULL,
    status TEXT NOT NULL DEFAULT 'Active',
    registration_date TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    services TEXT[],
    paid_up_capital DOUBLE PRECISION NOT NULL DEFAULT 0,
    chief_executive TEXT NOT NULL DEFAULT '',
    phone TEXT NOT NULL DEFAULT '',
    email TEXT NOT NULL DEFAULT '',
    website TEXT NOT NULL DEFAULT '',
    telegram_handle TEXT NOT NULL DEFAULT '',
    head_office TEXT NOT NULL DEFAULT '',
    logo_url TEXT NOT NULL DEFAULT '',
    description TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_csp_license_number ON capital_service_providers (license_number);
CREATE INDEX IF NOT EXISTS idx_csp_name ON capital_service_providers (name);
CREATE INDEX IF NOT EXISTS idx_csp_deleted_at ON capital_service_providers (deleted_at);
CREATE INDEX IF NOT EXISTS idx_csp_category ON capital_service_providers (category);
CREATE INDEX IF NOT EXISTS idx_csp_status ON capital_service_providers (status);

DO $$
BEGIN
  IF EXISTS (
    SELECT 1 FROM information_schema.tables
    WHERE table_schema = 'public' AND table_name = 'brokers'
  ) THEN
    INSERT INTO capital_service_providers (
      id, name, license_number, category, status, registration_date,
      services, paid_up_capital, chief_executive, phone, email, website, telegram_handle,
      head_office, logo_url, description, created_at, updated_at, deleted_at
    )
    SELECT
      b.id,
      b.name,
      -- Replace the old CASE statement with this one
      CASE
        WHEN trim(b.license) IN ('Active', 'Inactive', 'Pending', 'Suspended', 'Revoked', 'Licensed', '') OR b.license IS NULL
          THEN 'MIGRATED-' || b.id::text
        ELSE trim(b.license)
      END,
      COALESCE(NULLIF(trim(b.type), ''), 'Unspecified'),
      CASE
        WHEN trim(b.license) IN ('Active', 'Inactive', 'Pending', 'Suspended', 'Revoked')
          THEN trim(b.license)
        ELSE 'Active'
      END,
      COALESCE(b.created_at, NOW()),
      b.services,
      COALESCE(b.min_account::double precision, 0),
      '',
      COALESCE(b.phone, ''),
      COALESCE(b.email, ''),
      COALESCE(b.website, ''),
      '',
      '',
      '',
      COALESCE(b.description, ''),
      b.created_at,
      b.updated_at,
      b.deleted_at
    FROM brokers b
    ON CONFLICT (id) DO NOTHING;

    DROP TABLE brokers CASCADE;
  END IF;
END $$;

COMMIT;
