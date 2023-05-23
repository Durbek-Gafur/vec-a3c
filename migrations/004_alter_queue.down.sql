-- Add the old columns back
ALTER TABLE queue
ADD COLUMN workflow_parameters JSON AFTER workflow_name,
ADD COLUMN created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP AFTER workflow_parameters;

-- Remove the new columns
ALTER TABLE queue
DROP COLUMN workflow_id,
DROP COLUMN status,
DROP COLUMN enqueued_at;

