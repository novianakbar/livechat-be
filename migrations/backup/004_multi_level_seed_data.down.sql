-- Rollback Multi-Level Support Seed Data
-- Delete agent status for new agents
DELETE FROM agent_status
WHERE
    agent_id IN (
        '550e8400-e29b-41d4-a716-446655440013',
        '550e8400-e29b-41d4-a716-446655440014',
        '550e8400-e29b-41d4-a716-446655440015'
    );

-- Delete sample escalation data
DELETE FROM ticket_escalations
WHERE
    id = '550e8400-e29b-41d4-a716-446655440201';

-- Delete sample ticket categories
DELETE FROM ticket_categories
WHERE
    id IN (
        '550e8400-e29b-41d4-a716-446655440101',
        '550e8400-e29b-41d4-a716-446655440102',
        '550e8400-e29b-41d4-a716-446655440103'
    );

-- Delete new users
DELETE FROM users
WHERE
    id IN (
        '550e8400-e29b-41d4-a716-446655440013',
        '550e8400-e29b-41d4-a716-446655440014',
        '550e8400-e29b-41d4-a716-446655440015'
    );

-- Delete new departments
DELETE FROM departments
WHERE
    id IN (
        '550e8400-e29b-41d4-a716-446655440003',
        '550e8400-e29b-41d4-a716-446655440004'
    );

-- Reset existing departments to original state
UPDATE departments
SET
    support_level = NULL,
    max_escalation_level = NULL,
    auto_assignment_rule = NULL,
    escalation_dept_id = NULL,
    parent_dept_id = NULL
WHERE
    id IN (
        '550e8400-e29b-41d4-a716-446655440001',
        '550e8400-e29b-41d4-a716-446655440002'
    );