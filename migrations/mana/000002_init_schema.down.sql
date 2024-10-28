

-- Drop triggers first
DROP TRIGGER IF EXISTS update_wizard_investments_updated_at ON wizard_investments;
DROP TRIGGER IF EXISTS update_investment_types_updated_at ON investment_types;

-- Drop trigger function
DROP FUNCTION IF EXISTS update_updated_at_column();

-- Drop tables (in correct order due to foreign key constraints)
DROP TABLE IF EXISTS wizard_investments;
DROP TABLE IF EXISTS investment_types;
DROP TABLE IF EXISTS mana_transactions;

-- Drop indices (they'll be dropped with their tables, but included for completeness)
DROP INDEX IF EXISTS idx_mana_transactions_from_wizard;
DROP INDEX IF EXISTS idx_mana_transactions_to_wizard;
DROP INDEX IF EXISTS idx_mana_transactions_created_at;
DROP INDEX IF EXISTS idx_investment_types_risk_level;
DROP INDEX IF EXISTS idx_wizard_investments_wizard_id;
DROP INDEX IF EXISTS idx_wizard_investments_status;
DROP INDEX IF EXISTS idx_wizard_investments_end_time;