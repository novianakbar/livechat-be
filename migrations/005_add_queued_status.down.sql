-- Rollback: Remove 'queued' status from chat_sessions

-- First, update any sessions with 'queued' status to 'waiting'
UPDATE chat_sessions SET status = 'waiting' WHERE status = 'queued';

-- Then restore original constraint
ALTER TABLE chat_sessions 
DROP CONSTRAINT IF EXISTS chat_sessions_status_check;

ALTER TABLE chat_sessions 
ADD CONSTRAINT chat_sessions_status_check 
CHECK (status IN ('waiting', 'active', 'closed'));
