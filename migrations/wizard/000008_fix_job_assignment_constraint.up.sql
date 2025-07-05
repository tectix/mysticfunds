-- Remove the basic UNIQUE constraint and add a partial unique index
-- This allows wizards to be reassigned to jobs after completion
ALTER TABLE job_assignments DROP CONSTRAINT IF EXISTS job_assignments_job_id_wizard_id_key;

-- Create partial unique index to prevent duplicate active assignments
CREATE UNIQUE INDEX job_assignments_active_unique 
ON job_assignments (job_id, wizard_id) 
WHERE status IN ('assigned', 'in_progress');