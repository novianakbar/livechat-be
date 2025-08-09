-- Insert essential departments
INSERT INTO
    departments (id, name, description, is_active)
VALUES
    (
        '550e8400-e29b-41d4-a716-446655440001',
        'General Support',
        'Departemen umum untuk bantuan dan dukungan',
        true
    ),
    (
        '550e8400-e29b-41d4-a716-446655440002',
        'Technical Support',
        'Departemen yang menangani masalah teknis sistem',
        true
    );

-- Insert essential users
-- Default admin user
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
        '550e8400-e29b-41d4-a716-446655440010',
        'admin@livechat.com',
        '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi',
        'Administrator',
        'admin',
        true,
        NULL
    ),
    -- Sample agents
    (
        '550e8400-e29b-41d4-a716-446655440011',
        'agent1@livechat.com',
        '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi',
        'Agent Support 1',
        'agent',
        true,
        '550e8400-e29b-41d4-a716-446655440001'
    ),
    (
        '550e8400-e29b-41d4-a716-446655440012',
        'agent2@livechat.com',
        '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi',
        'Agent Technical',
        'agent',
        true,
        '550e8400-e29b-41d4-a716-446655440002'
    );

-- Insert essential agent status
INSERT INTO
    agent_status (id, agent_id, status, last_login_at)
VALUES
    (
        '550e8400-e29b-41d4-a716-446655440111',
        '550e8400-e29b-41d4-a716-446655440011',
        'logged_out',
        CURRENT_TIMESTAMP
    ),
    (
        '550e8400-e29b-41d4-a716-446655440112',
        '550e8400-e29b-41d4-a716-446655440012',
        'logged_out',
        CURRENT_TIMESTAMP
    );

-- Insert essential chat tags
INSERT INTO
    chat_tags (id, name, color)
VALUES
    (
        '550e8400-e29b-41d4-a716-446655440020',
        'General Question',
        '#007bff'
    ),
    (
        '550e8400-e29b-41d4-a716-446655440021',
        'Technical Issue',
        '#dc3545'
    ),
    (
        '550e8400-e29b-41d4-a716-446655440022',
        'Support Request',
        '#28a745'
    ),
    (
        '550e8400-e29b-41d4-a716-446655440023',
        'Urgent',
        '#fd7e14'
    );