-- Create jobs table for realm-based jobs that wizards can be assigned to
CREATE TABLE IF NOT EXISTS jobs (
    id SERIAL PRIMARY KEY,
    realm_id INTEGER NOT NULL REFERENCES realms(id) ON DELETE CASCADE,
    title VARCHAR(200) NOT NULL,
    description TEXT NOT NULL,
    required_element VARCHAR(50) NOT NULL,
    required_level INTEGER DEFAULT 1,
    mana_reward_per_hour INTEGER NOT NULL,
    exp_reward_per_hour INTEGER DEFAULT 10,
    duration_minutes INTEGER NOT NULL,
    max_wizards INTEGER DEFAULT 1,
    currently_assigned INTEGER DEFAULT 0,
    difficulty VARCHAR(20) NOT NULL CHECK (difficulty IN ('Easy', 'Medium', 'Hard', 'Expert', 'Legendary')),
    job_type VARCHAR(50) NOT NULL,
    location VARCHAR(200),
    special_requirements TEXT,
    created_by_wizard_id INTEGER REFERENCES wizards(id) ON DELETE SET NULL, -- For player-created jobs
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    is_active BOOLEAN DEFAULT true
);

-- Create job assignments table to track which wizards are assigned to which jobs
CREATE TABLE IF NOT EXISTS job_assignments (
    id SERIAL PRIMARY KEY,
    job_id INTEGER NOT NULL REFERENCES jobs(id) ON DELETE CASCADE,
    wizard_id INTEGER NOT NULL REFERENCES wizards(id) ON DELETE CASCADE,
    assigned_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    started_at TIMESTAMP WITH TIME ZONE,
    completed_at TIMESTAMP WITH TIME ZONE,
    status VARCHAR(20) DEFAULT 'assigned' CHECK (status IN ('assigned', 'in_progress', 'completed', 'failed', 'cancelled')),
    mana_earned INTEGER DEFAULT 0,
    exp_earned INTEGER DEFAULT 0,
    notes TEXT,
    UNIQUE(job_id, wizard_id)
);

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_jobs_realm_id ON jobs(realm_id);
CREATE INDEX IF NOT EXISTS idx_jobs_required_element ON jobs(required_element);
CREATE INDEX IF NOT EXISTS idx_jobs_difficulty ON jobs(difficulty);
CREATE INDEX IF NOT EXISTS idx_jobs_is_active ON jobs(is_active);
CREATE INDEX IF NOT EXISTS idx_jobs_created_by_wizard_id ON jobs(created_by_wizard_id);
CREATE INDEX IF NOT EXISTS idx_job_assignments_job_id ON job_assignments(job_id);
CREATE INDEX IF NOT EXISTS idx_job_assignments_wizard_id ON job_assignments(wizard_id);
CREATE INDEX IF NOT EXISTS idx_job_assignments_status ON job_assignments(status);

-- Insert sample mystical jobs for each realm type
INSERT INTO jobs (realm_id, title, description, required_element, required_level, mana_reward_per_hour, exp_reward_per_hour, duration_minutes, max_wizards, difficulty, job_type, location, special_requirements) VALUES 
-- Fire realm jobs (Pyrrhian Flame - realm_id 1)
(1, 'Forge Guardian', 'Tend the Eternal Forge and maintain its blazing flames. Requires constant vigilance against flame wraiths.', 'Fire', 1, 150, 15, 480, 2, 'Medium', 'Guardian', 'The Eternal Forge Chamber', 'Must resist fire damage'),
(1, 'Salamander Whisperer', 'Communicate with the ancient Salamandrine Lords to negotiate territorial agreements.', 'Fire', 3, 300, 35, 720, 1, 'Hard', 'Diplomatic', 'Salamander Council Halls', 'Fire immunity required'),
(1, 'Volcanic Mining Supervisor', 'Oversee extraction of rare fire crystals from active volcanic vents.', 'Fire', 2, 200, 25, 360, 3, 'Medium', 'Industrial', 'Crimson Peaks Mining Site', 'Heat resistance gear provided'),

-- Wind realm jobs (Zepharion Heights - realm_id 2)  
(2, 'Sky Patrol Ranger', 'Monitor the floating islands and report any atmospheric disturbances or sky pirate activity.', 'Air', 1, 120, 12, 600, 2, 'Easy', 'Patrol', 'Floating Island Networks', 'Flight capability preferred'),
(2, 'Wind Current Mapper', 'Study and document the ever-changing wind patterns around the crystal spires.', 'Air', 2, 180, 20, 480, 1, 'Medium', 'Research', 'Crystal Spire Observatory', 'Advanced wind magic knowledge'),
(2, 'Storm Caller', 'Assist in summoning controlled weather patterns for agricultural and defensive purposes.', 'Air', 4, 400, 45, 240, 1, 'Expert', 'Magical', 'Storm Summoning Circle', 'Master-level air magic required'),

-- Earth realm jobs (Terravine Hollow - realm_id 3)
(3, 'Grove Protector', 'Defend the sacred groves from stone blight and maintain the magical root networks.', 'Earth', 1, 140, 14, 720, 2, 'Medium', 'Guardian', 'Sacred Grove Network', 'Plant communication ability helpful'),
(3, 'Crystal Cave Archaeologist', 'Explore and document ancient ruins within the crystal cave systems.', 'Earth', 2, 220, 28, 360, 1, 'Hard', 'Research', 'Deep Crystal Caverns', 'Cave navigation expertise required'),
(3, 'Terravine Cultivator', 'Tend to the magical vines that form the realm''s living architecture.', 'Earth', 1, 100, 10, 480, 3, 'Easy', 'Agricultural', 'Living Architecture Districts', 'Basic plant magic knowledge'),

-- Water realm jobs (Thalorion Depths - realm_id 4)
(4, 'Depth Sentinel', 'Patrol the deepest ocean trenches and monitor for ancient sea creature activity.', 'Water', 2, 250, 30, 600, 1, 'Hard', 'Patrol', 'Abyssal Trenches', 'Deep sea pressure resistance'),
(4, 'Coral City Architect', 'Design and grow new living coral structures for expanding underwater cities.', 'Water', 3, 280, 32, 480, 2, 'Hard', 'Construction', 'Coral City Planning District', 'Advanced water shaping skills'),
(4, 'Tidal Oracle Assistant', 'Help interpret prophetic visions seen in the shifting tidal patterns.', 'Water', 1, 160, 18, 360, 1, 'Medium', 'Spiritual', 'Oracle Chambers', 'Divination sensitivity required'),

-- Light realm jobs (Virelya - realm_id 5)
(5, 'Dawn Herald', 'Announce the daily awakening of the realm and maintain the eternal dawn cycle.', 'Light', 1, 130, 13, 240, 1, 'Easy', 'Ceremonial', 'Dawn Proclamation Tower', 'Strong voice and presence required'),
(5, 'Shadow Cleanser', 'Purify areas corrupted by shadow magic and maintain the realm''s luminous purity.', 'Light', 2, 200, 22, 480, 2, 'Medium', 'Purification', 'Corrupted Zones', 'Shadow resistance essential'),
(5, 'Prism Master', 'Operate the great light prisms that focus and direct the realm''s magical illumination.', 'Light', 4, 350, 40, 360, 1, 'Expert', 'Magical', 'Central Prism Array', 'Master-level light magic control'),

-- Short silly jobs (2-10 minutes) for quick testing and mana rewards
(1, 'Dragon Sneeze Collector', 'Collect magical dragon sneezes in tiny crystal vials. Dragons are surprisingly ticklish when you mention dusty caves.', 'Fire', 1, 300, 25, 3, 1, 'Easy', 'Silly', 'Dragon Napping Quarters', 'Handkerchief provided'),
(1, 'Candle Lighter Apprentice', 'Light exactly 47 candles in the ceremonial chamber without burning your robes.', 'Fire', 1, 180, 20, 5, 1, 'Easy', 'Silly', 'Ceremonial Candle Chamber', 'Eyebrow insurance recommended'),
(2, 'Cloud Sheep Herder', 'Herd fluffy cloud sheep that keep floating away. Use your magic lasso and a lot of patience.', 'Air', 1, 220, 25, 4, 1, 'Easy', 'Silly', 'Sky Pastures', 'Magic lasso included'),
(2, 'Flying Carpet Traffic Controller', 'Direct flying carpet traffic during rush hour. Wizards are surprisingly bad at following air traffic rules.', 'Air', 1, 320, 35, 2, 1, 'Easy', 'Silly', 'Sky Traffic Control Tower', 'Loud whistle and strong voice needed'),
(3, 'Mushroom Whisperer', 'Convince stubborn magical mushrooms to grow in the right direction. They have strong opinions about soil pH.', 'Earth', 1, 190, 22, 6, 1, 'Easy', 'Silly', 'Mushroom Philosophy Gardens', 'Mushroom language phrasebook provided'),
(3, 'Flower Compliment Giver', 'Give sincere compliments to flowers to boost their self-esteem. They bloom brighter when appreciated.', 'Earth', 1, 200, 23, 3, 1, 'Easy', 'Silly', 'Self-Esteem Flower Garden', 'Genuine enthusiasm required'),
(4, 'Bubble Wrap Popper', 'Pop magical bubble wrap to release stored water magic. Surprisingly therapeutic but highly addictive.', 'Water', 1, 160, 18, 7, 1, 'Easy', 'Silly', 'Bubble Storage Facility', 'Strong restraint required'),
(4, 'Seahorse Racing Commentator', 'Provide exciting commentary for seahorse races. They move very slowly, so you need to be creative.', 'Water', 1, 240, 27, 8, 1, 'Easy', 'Silly', 'Seahorse Racing Track', 'Dramatic voice and vivid imagination'),
(5, 'Unicorn Mane Stylist', 'Style unicorn manes for important celestial events. They are surprisingly vain about their appearance.', 'Light', 1, 300, 35, 10, 1, 'Easy', 'Silly', 'Celestial Beauty Parlor', 'Magic-resistant hair spray provided'),
(5, 'Glitter Cleanup Crew', 'Clean up after unicorn parties. Warning: glitter gets everywhere and is magically permanent.', 'Light', 1, 180, 20, 9, 1, 'Easy', 'Silly', 'Post-Party Cleanup Sites', 'Resignation to eternal sparkles required');