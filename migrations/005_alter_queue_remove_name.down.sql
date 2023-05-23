-- Add the old columns back
ALTER TABLE queue
ADD COLUMN workflow_name JSON AFTER id,