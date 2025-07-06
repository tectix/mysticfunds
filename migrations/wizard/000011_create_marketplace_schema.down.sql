-- Drop marketplace system tables in reverse dependency order

-- Drop indexes first
DROP INDEX IF EXISTS idx_marketplace_transactions_type;
DROP INDEX IF EXISTS idx_marketplace_transactions_seller;
DROP INDEX IF EXISTS idx_marketplace_transactions_buyer;
DROP INDEX IF EXISTS idx_wizard_spell_teaching_spell_id;
DROP INDEX IF EXISTS idx_wizard_spell_teaching_wizard_id;
DROP INDEX IF EXISTS idx_wizard_spells_learned_from;
DROP INDEX IF EXISTS idx_wizard_spells_wizard_id;
DROP INDEX IF EXISTS idx_wizard_scrolls_wizard_id;
DROP INDEX IF EXISTS idx_wizard_artifacts_wizard_id;
DROP INDEX IF EXISTS idx_spells_element;
DROP INDEX IF EXISTS idx_spells_school;
DROP INDEX IF EXISTS idx_scrolls_available;
DROP INDEX IF EXISTS idx_scrolls_skill_type;
DROP INDEX IF EXISTS idx_artifacts_available;
DROP INDEX IF EXISTS idx_artifacts_rarity;
DROP INDEX IF EXISTS idx_artifacts_realm_id;

-- Drop tables in reverse dependency order
DROP TABLE IF EXISTS marketplace_transactions;
DROP TABLE IF EXISTS wizard_spell_teaching;
DROP TABLE IF EXISTS wizard_spells;
DROP TABLE IF EXISTS wizard_scrolls;
DROP TABLE IF EXISTS wizard_artifacts;
DROP TABLE IF EXISTS spells;
DROP TABLE IF EXISTS scrolls;
DROP TABLE IF EXISTS artifacts;