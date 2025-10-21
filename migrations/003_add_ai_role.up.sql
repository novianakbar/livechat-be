-- Migration to support AI integration
-- Add 'ai' as a valid role
ALTER TABLE users 
DROP CONSTRAINT users_role_check;

ALTER TABLE users 
ADD CONSTRAINT users_role_check CHECK (role IN ('admin', 'agent', 'ai'));

-- Create AI user for system integration
INSERT INTO
    users (
        id,
        email,
        password,
        name,
        role,
        is_active,
        department_id
    )
VALUES
    (
        '550e8400-e29b-41d4-a716-446655440099',
        'ai@livechat.com',
        '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', -- password: password
        'AI Assistant',
        'ai',
        true,
        NULL
    );