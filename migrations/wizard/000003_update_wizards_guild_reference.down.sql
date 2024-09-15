-- Remove the foreign key constraint
ALTER TABLE wizards DROP CONSTRAINT IF EXISTS fk_wizard_guild;

-- Remove the index
DROP INDEX IF EXISTS idx_wizards_guild_id;

-- Remove the guild_id column
ALTER TABLE wizards DROP COLUMN IF EXISTS guild_id;

-- Add back the original guild column (as a VARCHAR)
ALTER TABLE wizards ADD COLUMN guild VARCHAR(100);