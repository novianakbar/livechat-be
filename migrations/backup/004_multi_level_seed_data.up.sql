-- Multi-Level Support Seed Data
-- Updates existing departments and adds new levels
-- Update existing departments with support levels
UPDATE departments
SET
    support_level = 0,
    max_escalation_level = 3,
    auto_assignment_rule = 'round_robin',
    escalation_dept_id = '550e8400-e29b-41d4-a716-446655440002'
WHERE
    id = '550e8400-e29b-41d4-a716-446655440001';

-- General Support (L0)
UPDATE departments
SET
    support_level = 1,
    max_escalation_level = 3,
    auto_assignment_rule = 'least_loaded',
    parent_dept_id = '550e8400-e29b-41d4-a716-446655440001',
    escalation_dept_id = '550e8400-e29b-41d4-a716-446655440003'
WHERE
    id = '550e8400-e29b-41d4-a716-446655440002';

-- Technical Support (L1)
-- Insert new higher-level departments
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

-- Insert additional users for higher levels
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

-- Insert sample ticket categories with level-appropriate SLAs
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

-- Sample escalation data for testing
INSERT INTO
    ticket_escalations (
        id,
        ticket_id,
        from_level,
        to_level,
        from_department_id,
        to_department_id,
        reason,
        escalated_by_id,
        is_auto_escalation,
        trigger_type,
        was_successful
    )
VALUES
    (
        '550e8400-e29b-41d4-a716-446655440201',
        -- We'll need to insert this after tickets are created
        NULL, -- Will be updated when tickets exist
        0, -- From L0
        1, -- To L1
        '550e8400-e29b-41d4-a716-446655440001', -- From General Support
        '550e8400-e29b-41d4-a716-446655440002', -- To Technical Support
        'Customer issue requires technical expertise beyond L0 capabilities',
        '550e8400-e29b-41d4-a716-446655440010', -- Escalated by admin
        false, -- Manual escalation
        'manual',
        true
    );

-- Comments for documentation
COMMENT ON TABLE ticket_escalations IS 'Sample escalation data showing L0 to L1 escalation pattern';

-- Update existing agent status records if they exist
UPDATE agent_status
SET
    last_seen = CURRENT_TIMESTAMP
WHERE
    agent_id IN (
        '550e8400-e29b-41d4-a716-446655440011',
        '550e8400-e29b-41d4-a716-446655440012'
    );

-- Insert agent status for new agents
INSERT INTO
    agent_status (id, agent_id, status, last_seen)
VALUES
    (
        '550e8400-e29b-41d4-a716-446655440033',
        '550e8400-e29b-41d4-a716-446655440013',
        'online',
        CURRENT_TIMESTAMP
    ),
    (
        '550e8400-e29b-41d4-a716-446655440034',
        '550e8400-e29b-41d4-a716-446655440014',
        'online',
        CURRENT_TIMESTAMP
    ),
    (
        '550e8400-e29b-41d4-a716-446655440035',
        '550e8400-e29b-41d4-a716-446655440015',
        'away',
        CURRENT_TIMESTAMP
    );