-- Add the SubmittedBy field to the Workflow info table
ALTER TABLE workflow_info
ADD COLUMN submitted_by VARCHAR(255),
MODIFY COLUMN status ENUM('pending', 'processing', 'done') DEFAULT 'pending';
