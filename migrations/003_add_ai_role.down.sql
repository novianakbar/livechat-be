-- Revert AI integration changes
-- Remove AI user
DELETE FROM users WHERE id = '550e8400-e29b-41d4-a716-446655440099';

-- Restore original role constraint
ALTER TABLE users 
DROP CONSTRAINT users_role_check;

ALTER TABLE users 
ADD CONSTRAINT users_role_check CHECK (role IN ('admin', 'agent'));