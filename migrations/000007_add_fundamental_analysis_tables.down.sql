BEGIN;

DROP VIEW IF EXISTS company_latest_eps;

DROP INDEX IF EXISTS idx_corporate_news_company_published_at_desc;
DROP INDEX IF EXISTS idx_corporate_news_company_profile_id;
DROP TABLE IF EXISTS corporate_news;

DROP INDEX IF EXISTS idx_financial_reports_company_period_desc;
DROP INDEX IF EXISTS idx_financial_reports_company_profile_id;
DROP TABLE IF EXISTS financial_reports;

DROP INDEX IF EXISTS idx_board_members_company_profile_id;
DROP TABLE IF EXISTS board_members;

DROP INDEX IF EXISTS idx_company_profiles_ticker_symbol;
DROP INDEX IF EXISTS idx_company_profiles_company_name;
DROP TABLE IF EXISTS company_profiles;

COMMIT;
