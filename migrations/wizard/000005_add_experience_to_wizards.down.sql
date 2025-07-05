-- Remove experience and level columns from wizards table
ALTER TABLE wizards 
DROP COLUMN IF EXISTS experience_points,
DROP COLUMN IF EXISTS level;