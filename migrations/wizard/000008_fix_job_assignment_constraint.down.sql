-- Remove the partial unique index
DROP INDEX IF EXISTS job_assignments_active_unique;

-- Restore the basic unique constraint
ALTER TABLE job_assignments ADD CONSTRAINT job_assignments_job_id_wizard_id_key UNIQUE (job_id, wizard_id);