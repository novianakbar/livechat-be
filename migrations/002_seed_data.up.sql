-- ============================================
-- SEED DATA FOR LIVECHAT SYSTEM
-- Complete seed data for fresh installation
-- ============================================
-- ============================================
-- DEPARTMENTS WITH MULTI-LEVEL SUPPORT
-- ============================================
-- Insert departments with multi-level hierarchy
INSERT INTO
    departments (
        id,
        name,
        description,
        is_active,
        can_handle_tickets,
        max_tickets_per_agent,
        support_level,
        parent_dept_id,
        max_escalation_level,
        auto_assignment_rule,
        escalation_dept_id
    )
VALUES
    -- L0 - General Support (Entry point)
    (
        '550e8400-e29b-41d4-a716-446655440001',
        'General Support',
        'Departemen umum untuk bantuan dan dukungan',
        true,
        true,
        15,
        0,
        NULL,
        3,
        'round_robin',
        '550e8400-e29b-41d4-a716-446655440002'
    ),
    -- L1 - Technical Support
    (
        '550e8400-e29b-41d4-a716-446655440002',
        'Technical Support',
        'Departemen yang menangani masalah teknis sistem',
        true,
        true,
        10,
        1,
        '550e8400-e29b-41d4-a716-446655440001',
        3,
        'least_loaded',
        '550e8400-e29b-41d4-a716-446655440003'
    ),
    -- L2 - Senior Technical Support
    (
        '550e8400-e29b-41d4-a716-446655440003',
        'Senior Technical Support',
        'Level 2 - Senior technical specialists for complex issues',
        true,
        true,
        5,
        2,
        '550e8400-e29b-41d4-a716-446655440002',
        3,
        'skill_based',
        '550e8400-e29b-41d4-a716-446655440004'
    ),
    -- L3 - Management Review
    (
        '550e8400-e29b-41d4-a716-446655440004',
        'Management Review',
        'Level 3 - Management and expert review for critical escalations',
        true,
        true,
        3,
        3,
        '550e8400-e29b-41d4-a716-446655440003',
        3,
        'round_robin',
        NULL -- No further escalation
    );

-- ============================================
-- USERS (ADMIN AND AGENTS)
-- ============================================
-- Insert essential users with proper hierarchy
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
    -- Default admin user
    (
        '550e8400-e29b-41d4-a716-446655440010',
        'admin@livechat.com',
        '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi',
        'Administrator',
        'admin',
        true,
        NULL
    ),
    -- L0 Agent
    (
        '550e8400-e29b-41d4-a716-446655440011',
        'agent1@livechat.com',
        '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi',
        'Agent Support 1',
        'agent',
        true,
        '550e8400-e29b-41d4-a716-446655440001'
    ),
    -- L1 Agent
    (
        '550e8400-e29b-41d4-a716-446655440012',
        'agent2@livechat.com',
        '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi',
        'Agent Technical',
        'agent',
        true,
        '550e8400-e29b-41d4-a716-446655440002'
    ),
    -- L2 Senior Agent
    (
        '550e8400-e29b-41d4-a716-446655440013',
        'senior1@livechat.com',
        '$2a$10$8K9v8K9v8K9v8K9v8K9v8O7Wz7Wz7Wz7Wz7Wz7Wz7Wz7Wz7Wz7W',
        'Sarah Wilson',
        'agent',
        true,
        '550e8400-e29b-41d4-a716-446655440003'
    ),
    -- L2 Senior Agent
    (
        '550e8400-e29b-41d4-a716-446655440014',
        'senior2@livechat.com',
        '$2a$10$8K9v8K9v8K9v8K9v8K9v8O7Wz7Wz7Wz7Wz7Wz7Wz7Wz7Wz7Wz7W',
        'David Chen',
        'agent',
        true,
        '550e8400-e29b-41d4-a716-446655440003'
    ),
    -- L3 Manager
    (
        '550e8400-e29b-41d4-a716-446655440015',
        'manager1@livechat.com',
        '$2a$10$8K9v8K9v8K9v8K9v8K9v8O7Wz7Wz7Wz7Wz7Wz7Wz7Wz7Wz7Wz7W',
        'Emily Rodriguez',
        'agent',
        true,
        '550e8400-e29b-41d4-a716-446655440004'
    );

-- ============================================
-- AGENT STATUS TRACKING
-- ============================================
-- Insert agent status for all agents
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
    ),
    (
        '550e8400-e29b-41d4-a716-446655440113',
        '550e8400-e29b-41d4-a716-446655440013',
        'logged_out',
        CURRENT_TIMESTAMP
    ),
    (
        '550e8400-e29b-41d4-a716-446655440114',
        '550e8400-e29b-41d4-a716-446655440014',
        'logged_out',
        CURRENT_TIMESTAMP
    ),
    (
        '550e8400-e29b-41d4-a716-446655440115',
        '550e8400-e29b-41d4-a716-446655440015',
        'logged_out',
        CURRENT_TIMESTAMP
    );

-- ============================================
-- CHAT TAGS
-- ============================================
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
    ),
    (
        '550e8400-e29b-41d4-a716-446655440024',
        'Bug Report',
        '#6f42c1'
    ),
    (
        '550e8400-e29b-41d4-a716-446655440025',
        'Feature Request',
        '#20c997'
    );

-- ============================================
-- TICKET CATEGORIES WITH MULTI-LEVEL SUPPORT
-- ============================================
-- Insert ticket categories with appropriate SLAs and default departments
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
    -- Basic categories
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
    ),
    -- Advanced categories for multi-level escalation
    (
        '550e8400-e29b-41d4-a716-446655440101',
        'Critical System Issue',
        'CRITICAL',
        'Critical system failures requiring immediate escalation',
        '#dc2626',
        true,
        300, -- 5 minutes first response
        3600, -- 1 hour resolution
        '550e8400-e29b-41d4-a716-446655440002' -- Start at L1
    ),
    (
        '550e8400-e29b-41d4-a716-446655440102',
        'Complex Technical',
        'COMPLEX',
        'Complex technical issues that may require L2+ support',
        '#f59e0b',
        true,
        900, -- 15 minutes first response
        7200, -- 2 hours resolution
        '550e8400-e29b-41d4-a716-446655440001' -- Start at L0
    ),
    (
        '550e8400-e29b-41d4-a716-446655440103',
        'Escalation Review',
        'ESCALATION',
        'Issues requiring management review and approval',
        '#7c3aed',
        true,
        1800, -- 30 minutes first response
        86400, -- 24 hours resolution
        '550e8400-e29b-41d4-a716-446655440004' -- Start at L3
    );

-- ============================================
-- SAMPLE DATA FOR DEMONSTRATION
-- ============================================
-- Insert sample chat user
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
    (
        '550e8400-e29b-41d4-a716-446655440500',
        'demo-browser-uuid-12345',
        NULL,
        NULL,
        true,
        '192.168.1.100',
        'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36'
    ),
    (
        '550e8400-e29b-41d4-a716-446655440501',
        NULL,
        'oss-user-123',
        'user@example.com',
        false,
        '192.168.1.101',
        'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36'
    );

-- Insert sample chat session
INSERT INTO
    chat_sessions (
        id,
        chat_user_id,
        agent_id,
        department_id,
        topic,
        status,
        priority
    )
VALUES
    (
        '550e8400-e29b-41d4-a716-446655440600',
        '550e8400-e29b-41d4-a716-446655440500',
        '550e8400-e29b-41d4-a716-446655440011',
        '550e8400-e29b-41d4-a716-446655440001',
        'Demo chat session',
        'active',
        'normal'
    );

-- Insert sample ticket for demonstration
INSERT INTO
    tickets (
        id,
        ticket_code,
        subject,
        description,
        customer_name,
        customer_email,
        customer_phone,
        category_id,
        priority,
        status,
        assigned_to_id,
        department_id,
        created_via,
        access_token,
        current_level,
        can_escalate,
        max_level
    )
VALUES
    (
        '550e8400-e29b-41d4-a716-446655440700',
        'DEMO-001',
        'Demo Ticket - General Inquiry',
        'This is a demo ticket for testing purposes',
        'John Doe',
        'john.doe@example.com',
        '+1234567890',
        '550e8400-e29b-41d4-a716-446655440030',
        'medium',
        'open',
        '550e8400-e29b-41d4-a716-446655440011',
        '550e8400-e29b-41d4-a716-446655440001',
        'customer',
        'demo-access-token-123',
        0,
        true,
        3
    );

-- Insert sample ticket SLA
INSERT INTO
    ticket_sla (id, ticket_id, first_response_due, resolution_due)
VALUES
    (
        '550e8400-e29b-41d4-a716-446655440800',
        '550e8400-e29b-41d4-a716-446655440700',
        CURRENT_TIMESTAMP + INTERVAL '1 hour',
        CURRENT_TIMESTAMP + INTERVAL '24 hours'
    );

-- ============================================
-- COMMENTS AND NOTES
-- ============================================
-- Add helpful comments for administrators
INSERT INTO
    ticket_comments (
        id,
        ticket_id,
        user_id,
        content,
        is_internal,
        is_from_customer,
        author_name,
        author_email
    )
VALUES
    (
        '550e8400-e29b-41d4-a716-446655440900',
        '550e8400-e29b-41d4-a716-446655440700',
        '550e8400-e29b-41d4-a716-446655440010',
        'This ticket has been created for demonstration purposes. The multi-level support system is now active.',
        true,
        false,
        'Administrator',
        'admin@livechat.com'
    );

-- ============================================
-- FINAL NOTES
-- ============================================
-- Update any timestamp-dependent records
UPDATE departments
SET
    updated_at = CURRENT_TIMESTAMP
WHERE
    id IS NOT NULL;

UPDATE users
SET
    updated_at = CURRENT_TIMESTAMP
WHERE
    id IS NOT NULL;

UPDATE chat_tags
SET
    updated_at = CURRENT_TIMESTAMP
WHERE
    id IS NOT NULL;

UPDATE ticket_categories
SET
    updated_at = CURRENT_TIMESTAMP
WHERE
    id IS NOT NULL;