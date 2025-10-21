-- Migration to allow 'ai' as sender_type in chat_messages table
ALTER TABLE chat_messages 
DROP CONSTRAINT IF EXISTS chat_messages_sender_type_check;

-- Add constraint with 'ai' as a valid sender_type
ALTER TABLE chat_messages 
ADD CONSTRAINT chat_messages_sender_type_check CHECK (sender_type IN ('customer', 'agent', 'system', 'ai'));