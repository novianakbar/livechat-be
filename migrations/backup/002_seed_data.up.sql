-- Insert essential departments (basic version)
INSERT INTO
    departments (
        id,
        name,
        description,
        is_active,
        can_handle_tickets,
        max_tickets_per_agent
    )
VALUES
    (
        '550e8400-e29b-41d4-a716-446655440001',
        'General Support',
        'Departemen umum untuk bantuan dan dukungan',
        true,
        true,
        15
    ),
    (
        '550e8400-e29b-41d4-a716-446655440002',
        'Technical Support',
        'Departemen yang menangani masalah teknis sistem',
        true,
        true,
        10
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

-- ============================================
-- TICKETING SYSTEM SEED DATA
-- ============================================
-- Insert essential ticket categories
INSERT INTO
    ticket_categories (
        id,
        name,
        code,
        description,
        color,
        is_active,
        sla_first_response,
        sla_resolution,
        default_department_id
    )
VALUES
    (
        '550e8400-e29b-41d4-a716-446655440030',
        'General Inquiry',
        'GENERAL',
        'Pertanyaan umum tentang produk atau layanan',
        '#007bff',
        true,
        60,
        1440,
        '550e8400-e29b-41d4-a716-446655440001'
    ),
    (
        '550e8400-e29b-41d4-a716-446655440031',
        'Technical Issue',
        'TECH',
        'Masalah teknis sistem atau aplikasi',
        '#dc3545',
        true,
        30,
        720,
        '550e8400-e29b-41d4-a716-446655440002'
    ),
    (
        '550e8400-e29b-41d4-a716-446655440032',
        'Billing Support',
        'BILLING',
        'Masalah pembayaran atau tagihan',
        '#28a745',
        true,
        120,
        2880,
        '550e8400-e29b-41d4-a716-446655440001'
    ),
    (
        '550e8400-e29b-41d4-a716-446655440033',
        'Account Support',
        'ACCOUNT',
        'Masalah akun pengguna atau akses',
        '#ffc107',
        true,
        60,
        1440,
        '550e8400-e29b-41d4-a716-446655440001'
    ),
    (
        '550e8400-e29b-41d4-a716-446655440034',
        'Bug Report',
        'BUG',
        'Laporan bug atau error sistem',
        '#6f42c1',
        true,
        15,
        480,
        '550e8400-e29b-41d4-a716-446655440002'
    );