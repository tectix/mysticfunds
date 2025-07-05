-- First, remove the existing guild column if it exists
ALTER TABLE wizards DROP COLUMN IF EXISTS guild;

-- Add the guild_id column only if it doesn't exist
DO $$ 
BEGIN 
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='wizards' AND column_name='guild_id') THEN
        ALTER TABLE wizards ADD COLUMN guild_id INTEGER;
    END IF;
END $$;

-- Add a foreign key constraint (only if it doesn't exist)
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM information_schema.table_constraints WHERE constraint_name='fk_wizard_guild') THEN
        ALTER TABLE wizards
        ADD CONSTRAINT fk_wizard_guild
        FOREIGN KEY (guild_id)
        REFERENCES guilds(id)
        ON DELETE SET NULL;
    END IF;
END $$;

-- Create an index on the guild_id column (only if it doesn't exist)
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_indexes WHERE indexname='idx_wizards_guild_id') THEN
        CREATE INDEX idx_wizards_guild_id ON wizards(guild_id);
    END IF;
END $$;