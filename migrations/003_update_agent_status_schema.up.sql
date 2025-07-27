-- Migration to update agent_status table schema
-- Change from tracking online/offline/busy/away to logged_in/logged_out
-- Drop existing constraint if exists
ALTER TABLE agent_status
DROP CONSTRAINT IF EXISTS agent_status_status_check;

-- Update status values: convert existing statuses to new schema
UPDATE agent_status
SET
    status = CASE
        WHEN status IN ('online', 'busy', 'away') THEN 'logged_in'
        WHEN status = 'offline' THEN 'logged_out'
        ELSE 'logged_out'
    END;

-- Rename column for clarity
ALTER TABLE agent_status
RENAME COLUMN last_active_at TO last_login_at;

-- Add new constraint
ALTER TABLE agent_status ADD CONSTRAINT agent_status_status_check CHECK (status IN ('logged_in', 'logged_out'));