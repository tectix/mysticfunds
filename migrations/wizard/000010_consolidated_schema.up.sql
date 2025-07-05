-- Consolidated schema migration for wizard service
-- This replaces migrations 001-009 with a single comprehensive schema

-- Drop existing tables if they exist (for clean reinstall)
DROP TABLE IF EXISTS activity_logs CASCADE;
DROP TABLE IF EXISTS job_progress CASCADE;
DROP TABLE IF EXISTS job_assignments CASCADE;
DROP TABLE IF EXISTS jobs CASCADE;
DROP TABLE IF EXISTS realms CASCADE;
DROP TABLE IF EXISTS wizards CASCADE;
DROP TABLE IF EXISTS guilds CASCADE;

-- Drop functions if they exist
DROP FUNCTION IF EXISTS calculate_job_progress_percentage(TIMESTAMP WITH TIME ZONE, TIMESTAMP WITH TIME ZONE);
DROP FUNCTION IF EXISTS update_updated_at_column();

-- Create guilds table
CREATE TABLE guilds (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create realms table with complete data
CREATE TABLE realms (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    description TEXT,
    element VARCHAR(50),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Insert all realm data
INSERT INTO realms (id, name, description, element) VALUES 
(1, 'Pyrrhian Flame', 'Realm of eternal fire and volcanic titans', 'Fire'),
(2, 'Zepharion Heights', 'Floating islands around an eternal cyclone', 'Air'), 
(3, 'Terravine Hollow', 'Ancient buried realm of stone titans', 'Earth'),
(4, 'Thalorion Depths', 'Submerged empire of the Moonbound Court', 'Water'),
(5, 'Virelya', 'Blinding paradise of pure truth', 'Light'),
(6, 'Umbros', 'Void-split realm where light cannot reach', 'Shadow'),
(7, 'Nyxthar', 'Realm where reality collapses inward', 'Null'),
(8, 'Aetherion', 'Realm between realms of dreaming dead', 'Spirit'),
(9, 'Chronarxis', 'Spiral palace of fractured timelines', 'Time'),
(10, 'Technarok', 'Fusion of steel gods and nano-intelligences', 'Metal');

-- Create comprehensive wizards table
CREATE TABLE wizards (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    name VARCHAR(100) NOT NULL,
    realm VARCHAR(50) NOT NULL,
    element VARCHAR(50) NOT NULL,
    guild_id INTEGER REFERENCES guilds(id) ON DELETE SET NULL,
    mana_balance BIGINT NOT NULL DEFAULT 0,
    experience_points INTEGER NOT NULL DEFAULT 0,
    level INTEGER NOT NULL DEFAULT 1,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    CONSTRAINT unique_user_wizard_name UNIQUE(user_id, name)
);

-- Create comprehensive jobs table
CREATE TABLE jobs (
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
    created_by_wizard_id INTEGER REFERENCES wizards(id) ON DELETE SET NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    is_active BOOLEAN DEFAULT true
);

-- Create job assignments table
CREATE TABLE job_assignments (
    id SERIAL PRIMARY KEY,
    job_id INTEGER NOT NULL REFERENCES jobs(id) ON DELETE CASCADE,
    wizard_id INTEGER NOT NULL REFERENCES wizards(id) ON DELETE CASCADE,
    assigned_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    started_at TIMESTAMP WITH TIME ZONE,
    completed_at TIMESTAMP WITH TIME ZONE,
    status VARCHAR(20) DEFAULT 'assigned' CHECK (status IN ('assigned', 'in_progress', 'completed', 'failed', 'cancelled')),
    mana_earned INTEGER DEFAULT 0,
    exp_earned INTEGER DEFAULT 0,
    notes TEXT
);

-- Create job progress table with real-time tracking
CREATE TABLE job_progress (
    id SERIAL PRIMARY KEY,
    assignment_id INTEGER NOT NULL REFERENCES job_assignments(id) ON DELETE CASCADE UNIQUE,
    started_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    actual_start_time TIMESTAMP WITH TIME ZONE,
    expected_end_time TIMESTAMP WITH TIME ZONE,
    last_updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    last_tick_time TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    progress_percentage INTEGER DEFAULT 0 CHECK (progress_percentage >= 0 AND progress_percentage <= 100),
    time_worked_minutes INTEGER DEFAULT 0 CHECK (time_worked_minutes >= 0),
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create activity logs table
CREATE TABLE activity_logs (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    wizard_id INTEGER REFERENCES wizards(id) ON DELETE SET NULL,
    activity_type VARCHAR(50) NOT NULL,
    activity_description TEXT NOT NULL,
    metadata JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create comprehensive indexes
CREATE INDEX idx_wizards_user_id ON wizards(user_id);
CREATE INDEX idx_wizards_realm ON wizards(realm);
CREATE INDEX idx_wizards_element ON wizards(element);
CREATE INDEX idx_wizards_guild_id ON wizards(guild_id);
CREATE INDEX idx_wizards_level ON wizards(level);

CREATE INDEX idx_jobs_realm_id ON jobs(realm_id);
CREATE INDEX idx_jobs_required_element ON jobs(required_element);
CREATE INDEX idx_jobs_difficulty ON jobs(difficulty);
CREATE INDEX idx_jobs_is_active ON jobs(is_active);
CREATE INDEX idx_jobs_created_by_wizard_id ON jobs(created_by_wizard_id);

CREATE INDEX idx_job_assignments_job_id ON job_assignments(job_id);
CREATE INDEX idx_job_assignments_wizard_id ON job_assignments(wizard_id);
CREATE INDEX idx_job_assignments_status ON job_assignments(status);

CREATE INDEX idx_job_progress_assignment_id ON job_progress(assignment_id);
CREATE INDEX idx_job_progress_is_active ON job_progress(is_active);
CREATE INDEX idx_job_progress_started_at ON job_progress(started_at);
CREATE INDEX idx_job_progress_expected_end_time ON job_progress(expected_end_time);
CREATE INDEX idx_job_progress_actual_start_time ON job_progress(actual_start_time);
CREATE INDEX idx_job_progress_last_tick_time ON job_progress(last_tick_time);

CREATE INDEX idx_activity_logs_user_id ON activity_logs(user_id);
CREATE INDEX idx_activity_logs_wizard_id ON activity_logs(wizard_id);
CREATE INDEX idx_activity_logs_type ON activity_logs(activity_type);
CREATE INDEX idx_activity_logs_created_at ON activity_logs(created_at);

-- Create partial unique index to prevent duplicate active assignments
CREATE UNIQUE INDEX job_assignments_active_unique 
ON job_assignments (job_id, wizard_id) 
WHERE status IN ('assigned', 'in_progress');

-- Create update trigger function
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create triggers for updated_at columns
CREATE TRIGGER update_wizards_updated_at
    BEFORE UPDATE ON wizards
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_jobs_updated_at
    BEFORE UPDATE ON jobs
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Create function for real-time progress calculation
CREATE OR REPLACE FUNCTION calculate_job_progress_percentage(
    actual_start_time TIMESTAMP WITH TIME ZONE,
    expected_end_time TIMESTAMP WITH TIME ZONE
) RETURNS INTEGER AS $$
DECLARE
    now_time TIMESTAMP WITH TIME ZONE := CURRENT_TIMESTAMP;
    total_duration INTERVAL;
    elapsed_duration INTERVAL;
    progress_percentage INTEGER;
BEGIN
    -- If job hasn't started yet
    IF actual_start_time IS NULL OR now_time < actual_start_time THEN
        RETURN 0;
    END IF;
    
    -- If job should be complete
    IF now_time >= expected_end_time THEN
        RETURN 100;
    END IF;
    
    -- Calculate progress
    total_duration := expected_end_time - actual_start_time;
    elapsed_duration := now_time - actual_start_time;
    
    progress_percentage := FLOOR((EXTRACT(EPOCH FROM elapsed_duration) / EXTRACT(EPOCH FROM total_duration)) * 100);
    
    -- Ensure it's within bounds
    IF progress_percentage < 0 THEN
        RETURN 0;
    ELSIF progress_percentage > 100 THEN
        RETURN 100;
    ELSE
        RETURN progress_percentage;
    END IF;
END;
$$ LANGUAGE plpgsql;

-- Insert sample jobs data
INSERT INTO jobs (realm_id, title, description, required_element, required_level, mana_reward_per_hour, exp_reward_per_hour, duration_minutes, max_wizards, difficulty, job_type, location, special_requirements) VALUES
(1, 'Lava Crystal Mining', 'Extract valuable fire crystals from the molten depths', 'Fire', 1, 50, 15, 30, 3, 'Easy', 'Mining', 'Pyrrhian Flame Quarries', 'Heat resistance gear required'),
(1, 'Flame Elemental Pacification', 'Calm aggressive flame elementals threatening settlements', 'Fire', 3, 120, 35, 60, 2, 'Medium', 'Combat', 'Pyrrhian Settlement Borders', 'Combat experience recommended'),
(2, 'Wind Current Mapping', 'Chart the ever-changing air currents around floating islands', 'Air', 1, 40, 12, 45, 2, 'Easy', 'Exploration', 'Zepharion Sky Lanes', 'Flying mount or levitation spell'),
(2, 'Cyclone Core Investigation', 'Study the mysterious energy at the storm''s center', 'Air', 5, 200, 60, 120, 1, 'Expert', 'Research', 'Eternal Cyclone Center', 'Master-level wind magic'),
(3, 'Stone Titan Archaeology', 'Uncover ancient artifacts from buried titan remains', 'Earth', 2, 70, 20, 90, 4, 'Medium', 'Archaeology', 'Terravine Excavation Sites', 'Earth-shaping abilities'),
(4, 'Deep Sea Relic Recovery', 'Retrieve lost magical artifacts from ocean floor', 'Water', 3, 100, 30, 75, 2, 'Hard', 'Recovery', 'Thalorion Abyssal Plains', 'Water breathing enchantment'),
(5, 'Light Prism Maintenance', 'Maintain the realm''s reality-anchoring light prisms', 'Light', 4, 150, 45, 60, 1, 'Hard', 'Maintenance', 'Virelya Crystal Gardens', 'Pure heart and strong will'),
(6, 'Shadow Veil Investigation', 'Investigate anomalies in the realm''s darkness', 'Shadow', 3, 90, 25, 80, 2, 'Medium', 'Investigation', 'Umbros Void Rifts', 'Shadow resistance training'),
(7, 'Null Zone Stabilization', 'Prevent reality collapse in unstable sectors', 'Null', 6, 300, 80, 180, 1, 'Legendary', 'Stabilization', 'Nyxthar Collapse Points', 'Reality anchor certification'),
(8, 'Spirit Guide Escort', 'Guide lost souls to their final destination', 'Spirit', 2, 60, 18, 45, 3, 'Easy', 'Escort', 'Aetherion Dream Bridges', 'Empathic communication ability'),
(9, 'Timeline Repair', 'Fix fractures in the time stream', 'Time', 7, 400, 100, 240, 1, 'Legendary', 'Repair', 'Chronarxis Temporal Nexus', 'Temporal magic mastery'),
(10, 'Nano-Intelligence Debugging', 'Debug malfunctioning AI systems', 'Metal', 4, 180, 50, 90, 2, 'Hard', 'Technical', 'Technarok Core Systems', 'Programming and metal magic skills');