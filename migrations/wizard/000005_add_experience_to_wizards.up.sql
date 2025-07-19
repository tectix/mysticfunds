-- Add experience and level to wizards table
ALTER TABLE wizards 
ADD COLUMN experience_points INTEGER DEFAULT 0,
ADD COLUMN level INTEGER DEFAULT 1;

-- Create index for level-based queries
CREATE INDEX IF NOT EXISTS idx_wizards_level ON wizards(level);

-- Update existing wizards to have some random levels/experience
UPDATE wizards SET 
    experience_points = (random() * 1000)::INTEGER,
    level = CASE 
        WHEN (random() * 1000)::INTEGER < 100 THEN 1
        WHEN (random() * 1000)::INTEGER < 300 THEN 2
        WHEN (random() * 1000)::INTEGER < 600 THEN 3
        WHEN (random() * 1000)::INTEGER < 900 THEN 4
        ELSE 5
    END;