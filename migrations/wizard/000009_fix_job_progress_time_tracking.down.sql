-- Rollback for job progress time tracking fixes

-- Drop the function
DROP FUNCTION IF EXISTS calculate_job_progress_percentage(TIMESTAMP WITH TIME ZONE, TIMESTAMP WITH TIME ZONE);

-- Drop the indexes
DROP INDEX IF EXISTS idx_job_progress_expected_end_time;
DROP INDEX IF EXISTS idx_job_progress_actual_start_time;
DROP INDEX IF EXISTS idx_job_progress_last_tick_time;

-- Remove the new columns
ALTER TABLE job_progress DROP COLUMN IF EXISTS expected_end_time;
ALTER TABLE job_progress DROP COLUMN IF EXISTS actual_start_time;
ALTER TABLE job_progress DROP COLUMN IF EXISTS last_tick_time;