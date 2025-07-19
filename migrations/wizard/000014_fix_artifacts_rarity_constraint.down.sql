-- Revert artifacts table rarity constraint back to original
ALTER TABLE artifacts DROP CONSTRAINT IF EXISTS artifacts_rarity_check;
ALTER TABLE artifacts ADD CONSTRAINT artifacts_rarity_check 
CHECK (rarity IN ('Common', 'Uncommon', 'Rare', 'Epic', 'Legendary', 'Mythical'));