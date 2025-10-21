-- ============================================
-- ROLLBACK SEED DATA
-- Remove all seed data in reverse order
-- ============================================
-- Delete sample data first

DELETE FROM chat_sessions
WHERE
    id = '550e8400-e29b-41d4-a716-446655440600';

DELETE FROM chat_users
WHERE
    id IN (
        '550e8400-e29b-41d4-a716-446655440500',
        '550e8400-e29b-41d4-a716-446655440501'
    );

-- Delete chat tags
DELETE FROM chat_tags
WHERE
    id IN (
        '550e8400-e29b-41d4-a716-446655440020',
        '550e8400-e29b-41d4-a716-446655440021',
        '550e8400-e29b-41d4-a716-446655440022',
        '550e8400-e29b-41d4-a716-446655440023',
        '550e8400-e29b-41d4-a716-446655440024',
        '550e8400-e29b-41d4-a716-446655440025'
    );

-- Delete agent status
DELETE FROM agent_status
WHERE
    id IN (
        '550e8400-e29b-41d4-a716-446655440111',
        '550e8400-e29b-41d4-a716-446655440112',
        '550e8400-e29b-41d4-a716-446655440113',
        '550e8400-e29b-41d4-a716-446655440114',
        '550e8400-e29b-41d4-a716-446655440115'
    );

-- Delete users (agents first, then admin)
DELETE FROM users
WHERE
    id IN (
        '550e8400-e29b-41d4-a716-446655440011',
        '550e8400-e29b-41d4-a716-446655440012',
        '550e8400-e29b-41d4-a716-446655440013',
        '550e8400-e29b-41d4-a716-446655440014',
        '550e8400-e29b-41d4-a716-446655440015',
        '550e8400-e29b-41d4-a716-446655440010'
    );

-- Delete departments (reverse hierarchy order - children first)
DELETE FROM departments
WHERE
    id = '550e8400-e29b-41d4-a716-446655440004';

-- L3
DELETE FROM departments
WHERE
    id = '550e8400-e29b-41d4-a716-446655440003';

-- L2
DELETE FROM departments
WHERE
    id = '550e8400-e29b-41d4-a716-446655440002';

-- L1
DELETE FROM departments
WHERE
    id = '550e8400-e29b-41d4-a716-446655440001';

-- L0