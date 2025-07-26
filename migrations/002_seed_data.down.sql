-- Delete seed data (in reverse order due to foreign key constraints)
-- Delete chat logs first
DELETE FROM chat_logs
WHERE
    id IN (
        'bb0e8400-e29b-41d4-a716-446655440001',
        'bb0e8400-e29b-41d4-a716-446655440002',
        'bb0e8400-e29b-41d4-a716-446655440003',
        'bb0e8400-e29b-41d4-a716-446655440004',
        'bb0e8400-e29b-41d4-a716-446655440005'
    );

-- Delete chat messages
DELETE FROM chat_messages
WHERE
    id IN (
        'aa0e8400-e29b-41d4-a716-446655440001',
        'aa0e8400-e29b-41d4-a716-446655440002',
        'aa0e8400-e29b-41d4-a716-446655440003',
        'aa0e8400-e29b-41d4-a716-446655440004',
        'aa0e8400-e29b-41d4-a716-446655440005',
        'aa0e8400-e29b-41d4-a716-446655440006'
    );

-- Delete chat session contacts
DELETE FROM chat_session_contacts
WHERE
    id IN (
        '990e8400-e29b-41d4-a716-446655440001',
        '990e8400-e29b-41d4-a716-446655440002',
        '990e8400-e29b-41d4-a716-446655440003'
    );

-- Delete chat sessions
DELETE FROM chat_sessions
WHERE
    id IN (
        '880e8400-e29b-41d4-a716-446655440001',
        '880e8400-e29b-41d4-a716-446655440002',
        '880e8400-e29b-41d4-a716-446655440003'
    );

-- Delete chat users
DELETE FROM chat_users
WHERE
    id IN (
        '660e8400-e29b-41d4-a716-446655440001',
        '660e8400-e29b-41d4-a716-446655440002',
        '660e8400-e29b-41d4-a716-446655440003'
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
        '550e8400-e29b-41d4-a716-446655440025',
        '550e8400-e29b-41d4-a716-446655440026',
        '550e8400-e29b-41d4-a716-446655440027',
        '550e8400-e29b-41d4-a716-446655440028',
        '550e8400-e29b-41d4-a716-446655440029'
    );

-- Delete agent status
DELETE FROM agent_status
WHERE
    agent_id IN (
        '550e8400-e29b-41d4-a716-446655440011',
        '550e8400-e29b-41d4-a716-446655440012',
        '550e8400-e29b-41d4-a716-446655440013',
        '550e8400-e29b-41d4-a716-446655440014',
        '550e8400-e29b-41d4-a716-446655440015'
    );

-- Delete users
DELETE FROM users
WHERE
    id IN (
        '550e8400-e29b-41d4-a716-446655440010',
        '550e8400-e29b-41d4-a716-446655440011',
        '550e8400-e29b-41d4-a716-446655440012',
        '550e8400-e29b-41d4-a716-446655440013',
        '550e8400-e29b-41d4-a716-446655440014',
        '550e8400-e29b-41d4-a716-446655440015'
    );

-- Delete departments
DELETE FROM departments
WHERE
    id IN (
        '550e8400-e29b-41d4-a716-446655440001',
        '550e8400-e29b-41d4-a716-446655440002',
        '550e8400-e29b-41d4-a716-446655440003',
        '550e8400-e29b-41d4-a716-446655440004'
    );