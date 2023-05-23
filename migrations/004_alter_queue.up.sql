-- Add the new columns
ALTER TABLE queue
ADD COLUMN workflow_id INT AFTER id,
ADD COLUMN status ENUM('pending', 'processing', 'done') DEFAULT 'pending' AFTER workflow_name,
ADD COLUMN enqueued_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP AFTER status;

-- Remove the old columns
ALTER TABLE queue
DROP COLUMN workflow_parameters,
DROP COLUMN created_at;


