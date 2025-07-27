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

-- Insert default agent status (login session tracking)
INSERT INTO
    agent_status (agent_id, status, last_login_at)
VALUES
    (
        '550e8400-e29b-41d4-a716-446655440011',
        'logged_out',
        CURRENT_TIMESTAMP
    ),
    (
        '550e8400-e29b-41d4-a716-446655440012',
        'logged_out',
        CURRENT_TIMESTAMP
    ),
    (
        '550e8400-e29b-41d4-a716-446655440013',
        'logged_out',
        CURRENT_TIMESTAMP
    ),
    (
        '550e8400-e29b-41d4-a716-446655440014',
        'logged_out',
        CURRENT_TIMESTAMP
    ),
    (
        '550e8400-e29b-41d4-a716-446655440015',
        'logged_out',
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

-- Insert sample chat users (both anonymous and logged-in OSS users)
INSERT INTO
    chat_users (
        id,
        browser_uuid,
        oss_user_id,
        email,
        is_anonymous,
        ip_address,
        user_agent
    )
VALUES
    -- Anonymous user
    (
        '660e8400-e29b-41d4-a716-446655440001',
        '770e8400-e29b-41d4-a716-446655440001',
        NULL,
        NULL,
        true,
        '192.168.1.100',
        'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36'
    ),
    -- Logged-in OSS user 1
    (
        '660e8400-e29b-41d4-a716-446655440002',
        '770e8400-e29b-41d4-a716-446655440002',
        'OSS001',
        'user1@company.com',
        false,
        '192.168.1.101',
        'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36'
    ),
    -- Logged-in OSS user 2
    (
        '660e8400-e29b-41d4-a716-446655440003',
        NULL,
        'OSS002',
        'user2@company.com',
        false,
        '192.168.1.102',
        'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36'
    );

-- Insert sample chat sessions
INSERT INTO
    chat_sessions (
        id,
        chat_user_id,
        agent_id,
        department_id,
        topic,
        status,
        priority,
        started_at
    )
VALUES
    -- Anonymous user session
    (
        '880e8400-e29b-41d4-a716-446655440001',
        '660e8400-e29b-41d4-a716-446655440001',
        '550e8400-e29b-41d4-a716-446655440011',
        '550e8400-e29b-41d4-a716-446655440001',
        'Pertanyaan tentang izin usaha mikro',
        'closed',
        'normal',
        CURRENT_TIMESTAMP - INTERVAL '2 hours'
    ),
    -- Logged-in user session (active)
    (
        '880e8400-e29b-41d4-a716-446655440002',
        '660e8400-e29b-41d4-a716-446655440002',
        '550e8400-e29b-41d4-a716-446655440013',
        '550e8400-e29b-41d4-a716-446655440002',
        'Konsultasi investasi PMA bidang teknologi',
        'active',
        'high',
        CURRENT_TIMESTAMP - INTERVAL '30 minutes'
    ),
    -- Waiting session
    (
        '880e8400-e29b-41d4-a716-446655440003',
        '660e8400-e29b-41d4-a716-446655440003',
        NULL,
        '550e8400-e29b-41d4-a716-446655440003',
        'Tanya pajak restoran',
        'waiting',
        'normal',
        CURRENT_TIMESTAMP - INTERVAL '5 minutes'
    );

-- Insert sample chat session contacts
INSERT INTO
    chat_session_contacts (
        id,
        session_id,
        contact_name,
        contact_email,
        contact_phone,
        position,
        company_name
    )
VALUES
    -- Contact for anonymous user session
    (
        '990e8400-e29b-41d4-a716-446655440001',
        '880e8400-e29b-41d4-a716-446655440001',
        'Budi Santoso',
        'budi@warungmakan.com',
        '+6281234567890',
        'Pemilik',
        'Warung Makan Budi'
    ),
    -- Contact for logged-in user session
    (
        '990e8400-e29b-41d4-a716-446655440002',
        '880e8400-e29b-41d4-a716-446655440002',
        'Sarah Johnson',
        'sarah@techstartup.com',
        '+6281234567891',
        'CEO',
        'PT Tech Startup Indonesia'
    ),
    -- Contact for waiting session
    (
        '990e8400-e29b-41d4-a716-446655440003',
        '880e8400-e29b-41d4-a716-446655440003',
        'Ahmad Rizki',
        'ahmad@restoranmewah.com',
        '+6281234567892',
        'Manager',
        'Restoran Mewah Jakarta'
    );

-- Insert sample chat messages
INSERT INTO
    chat_messages (
        id,
        session_id,
        sender_id,
        sender_type,
        message,
        message_type
    )
VALUES
    -- Messages for closed session
    (
        'aa0e8400-e29b-41d4-a716-446655440001',
        '880e8400-e29b-41d4-a716-446655440001',
        NULL,
        'customer',
        'Selamat pagi, saya ingin tanya tentang syarat izin usaha mikro untuk warung makan',
        'text'
    ),
    (
        'aa0e8400-e29b-41d4-a716-446655440002',
        '880e8400-e29b-41d4-a716-446655440001',
        '550e8400-e29b-41d4-a716-446655440011',
        'agent',
        'Selamat pagi! Untuk izin usaha mikro warung makan, Anda memerlukan: 1) KTP, 2) Surat keterangan domisili, 3) NPWP. Apakah ada pertanyaan lain?',
        'text'
    ),
    (
        'aa0e8400-e29b-41d4-a716-446655440003',
        '880e8400-e29b-41d4-a716-446655440001',
        NULL,
        'customer',
        'Terima kasih informasinya, sangat membantu!',
        'text'
    ),
    -- Messages for active session
    (
        'aa0e8400-e29b-41d4-a716-446655440004',
        '880e8400-e29b-41d4-a716-446655440002',
        NULL,
        'customer',
        'Halo, saya sedang merencanakan investasi PMA di bidang teknologi. Bisa bantu informasi prosedurnya?',
        'text'
    ),
    (
        'aa0e8400-e29b-41d4-a716-446655440005',
        '880e8400-e29b-41d4-a716-446655440002',
        '550e8400-e29b-41d4-a716-446655440013',
        'agent',
        'Tentu! Untuk investasi PMA bidang teknologi, ada beberapa tahapan. Bisa tolong jelaskan lebih detail jenis teknologi yang akan dikembangkan?',
        'text'
    ),
    -- Message for waiting session
    (
        'aa0e8400-e29b-41d4-a716-446655440006',
        '880e8400-e29b-41d4-a716-446655440003',
        NULL,
        'customer',
        'Saya mau tanya tentang pajak restoran, berapa persen dan bagaimana cara bayarnya?',
        'text'
    );

-- Insert sample chat logs
INSERT INTO
    chat_logs (id, session_id, action, details, user_id)
VALUES
    (
        'bb0e8400-e29b-41d4-a716-446655440001',
        '880e8400-e29b-41d4-a716-446655440001',
        'started',
        'Chat session started by customer',
        NULL
    ),
    (
        'bb0e8400-e29b-41d4-a716-446655440002',
        '880e8400-e29b-41d4-a716-446655440001',
        'response',
        'Agent responded to customer inquiry',
        '550e8400-e29b-41d4-a716-446655440011'
    ),
    (
        'bb0e8400-e29b-41d4-a716-446655440003',
        '880e8400-e29b-41d4-a716-446655440001',
        'closed',
        'Session closed - customer satisfied',
        '550e8400-e29b-41d4-a716-446655440011'
    ),
    (
        'bb0e8400-e29b-41d4-a716-446655440004',
        '880e8400-e29b-41d4-a716-446655440002',
        'started',
        'Chat session started by OSS user',
        NULL
    ),
    (
        'bb0e8400-e29b-41d4-a716-446655440005',
        '880e8400-e29b-41d4-a716-446655440003',
        'started',
        'Chat session started - waiting for agent',
        NULL
    );