-- Remove the SubmittedBy field from the Workflow info table
ALTER TABLE workflow_info
DROP COLUMN submitted_by,
MODIFY COLUMN status VARCHAR(50) DEFAULT NULL;
