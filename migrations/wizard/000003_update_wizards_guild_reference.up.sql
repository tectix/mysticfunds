-- First, remove the existing guild column if it exists
ALTER TABLE wizards DROP COLUMN IF EXISTS guild;

-- Add the guild_id column
ALTER TABLE wizards ADD COLUMN guild_id INTEGER;

-- Add a foreign key constraint
ALTER TABLE wizards
ADD CONSTRAINT fk_wizard_guild
FOREIGN KEY (guild_id)
REFERENCES guilds(id)
ON DELETE SET NULL;

-- Create an index on the guild_id column
CREATE INDEX idx_wizards_guild_id ON wizards(guild_id);