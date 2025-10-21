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
