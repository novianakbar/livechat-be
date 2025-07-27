-- Rollback migration: revert agent_status table changes
-- Drop new constraint
ALTER TABLE agent_status
DROP CONSTRAINT IF EXISTS agent_status_status_check;

-- Rename column back
ALTER TABLE agent_status
RENAME COLUMN last_login_at TO last_active_at;

-- Update status values back to original schema
UPDATE agent_status
SET
    status = CASE
        WHEN status = 'logged_in' THEN 'online'
        WHEN status = 'logged_out' THEN 'offline'
        ELSE 'offline'
    END;

-- Add original constraint
ALTER TABLE agent_status ADD CONSTRAINT agent_status_status_check CHECK (status IN ('online', 'offline', 'busy', 'away'));