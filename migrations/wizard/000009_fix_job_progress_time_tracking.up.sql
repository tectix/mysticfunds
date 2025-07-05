-- Fix job progress table to properly track real time
-- Add fields needed for accurate time-based progress tracking

-- Add expected completion time and actual start time for better tracking
ALTER TABLE job_progress ADD COLUMN expected_end_time TIMESTAMP WITH TIME ZONE;
ALTER TABLE job_progress ADD COLUMN actual_start_time TIMESTAMP WITH TIME ZONE;
ALTER TABLE job_progress ADD COLUMN last_tick_time TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP;

-- Update existing records to have proper start times
UPDATE job_progress 
SET actual_start_time = started_at,
    last_tick_time = CURRENT_TIMESTAMP
WHERE actual_start_time IS NULL;

-- Calculate expected end times for existing active jobs
UPDATE job_progress 
SET expected_end_time = (
    SELECT job_progress.actual_start_time + INTERVAL '1 minute' * j.duration_minutes
    FROM job_assignments ja
    JOIN jobs j ON ja.job_id = j.id
    WHERE ja.id = job_progress.assignment_id
)
WHERE expected_end_time IS NULL 
AND is_active = true;

-- Add indexes for better performance on time queries
CREATE INDEX idx_job_progress_expected_end_time ON job_progress(expected_end_time);
CREATE INDEX idx_job_progress_actual_start_time ON job_progress(actual_start_time);
CREATE INDEX idx_job_progress_last_tick_time ON job_progress(last_tick_time);

-- Add a function to calculate real-time progress percentage
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