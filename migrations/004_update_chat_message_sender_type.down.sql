-- Revert chat_messages sender_type constraint change
ALTER TABLE chat_messages 
DROP CONSTRAINT IF EXISTS chat_messages_sender_type_check;

-- Restore original constraint without 'ai'
ALTER TABLE chat_messages 
ADD CONSTRAINT chat_messages_sender_type_check CHECK (sender_type IN ('customer', 'agent', 'system'));