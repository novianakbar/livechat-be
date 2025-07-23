-- Insert default departments
INSERT INTO
    departments (id, name, description, is_active)
VALUES
    (
        '550e8400-e29b-41d4-a716-446655440001',
        'Perizinan Usaha',
        'Departemen yang menangani perizinan usaha dan UMKM',
        true
    ),
    (
        '550e8400-e29b-41d4-a716-446655440002',
        'Investasi',
        'Departemen yang menangani investasi dan penanaman modal',
        true
    ),
    (
        '550e8400-e29b-41d4-a716-446655440003',
        'Perpajakan',
        'Departemen yang menangani perpajakan dan retribusi',
        true
    ),
    (
        '550e8400-e29b-41d4-a716-446655440004',
        'Teknis',
        'Departemen yang menangani masalah teknis sistem',
        true
    );

-- Insert default admin user
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
    );

-- Insert default agents
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
        '550e8400-e29b-41d4-a716-446655440011',
        'agent1@livechat.com',
        '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi',
        'Agent Perizinan 1',
        'agent',
        true,
        '550e8400-e29b-41d4-a716-446655440001'
    ),
    (
        '550e8400-e29b-41d4-a716-446655440012',
        'agent2@livechat.com',
        '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi',
        'Agent Perizinan 2',
        'agent',
        true,
        '550e8400-e29b-41d4-a716-446655440001'
    ),
    (
        '550e8400-e29b-41d4-a716-446655440013',
        'agent3@livechat.com',
        '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi',
        'Agent Investasi 1',
        'agent',
        true,
        '550e8400-e29b-41d4-a716-446655440002'
    ),
    (
        '550e8400-e29b-41d4-a716-446655440014',
        'agent4@livechat.com',
        '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi',
        'Agent Perpajakan 1',
        'agent',
        true,
        '550e8400-e29b-41d4-a716-446655440003'
    ),
    (
        '550e8400-e29b-41d4-a716-446655440015',
        'agent5@livechat.com',
        '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi',
        'Agent Teknis 1',
        'agent',
        true,
        '550e8400-e29b-41d4-a716-446655440004'
    );

-- Insert default agent status
INSERT INTO
    agent_status (agent_id, status, last_active_at)
VALUES
    (
        '550e8400-e29b-41d4-a716-446655440011',
        'offline',
        CURRENT_TIMESTAMP
    ),
    (
        '550e8400-e29b-41d4-a716-446655440012',
        'offline',
        CURRENT_TIMESTAMP
    ),
    (
        '550e8400-e29b-41d4-a716-446655440013',
        'offline',
        CURRENT_TIMESTAMP
    ),
    (
        '550e8400-e29b-41d4-a716-446655440014',
        'offline',
        CURRENT_TIMESTAMP
    ),
    (
        '550e8400-e29b-41d4-a716-446655440015',
        'offline',
        CURRENT_TIMESTAMP
    );

-- Insert default chat tags
INSERT INTO
    chat_tags (id, name, color)
VALUES
    (
        '550e8400-e29b-41d4-a716-446655440020',
        'Perizinan Baru',
        '#28a745'
    ),
    (
        '550e8400-e29b-41d4-a716-446655440021',
        'Perpanjangan Izin',
        '#ffc107'
    ),
    (
        '550e8400-e29b-41d4-a716-446655440022',
        'Investasi PMA',
        '#17a2b8'
    ),
    (
        '550e8400-e29b-41d4-a716-446655440023',
        'Investasi PMDN',
        '#007bff'
    ),
    (
        '550e8400-e29b-41d4-a716-446655440024',
        'Pajak Daerah',
        '#6c757d'
    ),
    (
        '550e8400-e29b-41d4-a716-446655440025',
        'Retribusi',
        '#fd7e14'
    ),
    (
        '550e8400-e29b-41d4-a716-446655440026',
        'Masalah Teknis',
        '#dc3545'
    ),
    (
        '550e8400-e29b-41d4-a716-446655440027',
        'Konsultasi',
        '#6f42c1'
    ),
    (
        '550e8400-e29b-41d4-a716-446655440028',
        'Komplain',
        '#e83e8c'
    ),
    (
        '550e8400-e29b-41d4-a716-446655440029',
        'Urgent',
        '#dc3545'
    );