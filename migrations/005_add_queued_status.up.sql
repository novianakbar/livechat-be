-- Add 'queued' status to chat_sessions status enum
-- This status is used when AI escalates to human but no agent is available

ALTER TABLE chat_sessions 
DROP CONSTRAINT IF EXISTS chat_sessions_status_check;

ALTER TABLE chat_sessions 
ADD CONSTRAINT chat_sessions_status_check 
CHECK (status IN ('waiting', 'queued', 'active', 'closed'));
