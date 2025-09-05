-- Migration untuk Multi-Level Support System
-- Phase 2 - Critical Enhancement
-- Add Multi-Level Support fields to departments table
ALTER TABLE departments
ADD COLUMN support_level INTEGER DEFAULT 0,
ADD COLUMN parent_dept_id VARCHAR(255) REFERENCES departments (id),
ADD COLUMN max_escalation_level INTEGER DEFAULT 3,
ADD COLUMN auto_assignment_rule VARCHAR(50) DEFAULT 'round_robin',
ADD COLUMN escalation_dept_id VARCHAR(255) REFERENCES departments (id);

-- Add Multi-Level Support fields to tickets table
ALTER TABLE tickets
ADD COLUMN current_level INTEGER DEFAULT 0,
ADD COLUMN escalation_path TEXT,
ADD COLUMN can_escalate BOOLEAN DEFAULT true,
ADD COLUMN max_level INTEGER DEFAULT 3,
ADD COLUMN escalation_count INTEGER DEFAULT 0,
ADD COLUMN previous_assigned_to_id VARCHAR(255) REFERENCES users (id),
ADD COLUMN previous_department_id VARCHAR(255) REFERENCES departments (id);

-- Create ticket_escalations table for tracking escalation history
CREATE TABLE
    ticket_escalations (
        id VARCHAR(255) PRIMARY KEY,
        ticket_id VARCHAR(255) NOT NULL REFERENCES tickets (id),
        from_level INTEGER NOT NULL,
        to_level INTEGER NOT NULL,
        from_department_id VARCHAR(255) REFERENCES departments (id),
        to_department_id VARCHAR(255) REFERENCES departments (id),
        from_agent_id VARCHAR(255) REFERENCES users (id),
        to_agent_id VARCHAR(255) REFERENCES users (id),
        reason TEXT NOT NULL,
        escalated_by_id VARCHAR(255) NOT NULL REFERENCES users (id),
        escalated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        is_auto_escalation BOOLEAN DEFAULT false,
        trigger_type VARCHAR(50),
        trigger_data TEXT,
        was_successful BOOLEAN DEFAULT true,
        failure_reason TEXT,
        resolved_at TIMESTAMP NULL,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        deleted_at BIGINT DEFAULT 0
    );

-- Create indexes for performance
CREATE INDEX idx_departments_support_level ON departments (support_level);

CREATE INDEX idx_departments_parent_dept_id ON departments (parent_dept_id);

CREATE INDEX idx_departments_escalation_dept_id ON departments (escalation_dept_id);

CREATE INDEX idx_tickets_current_level ON tickets (current_level);

CREATE INDEX idx_tickets_escalation_count ON tickets (escalation_count);

CREATE INDEX idx_tickets_previous_assigned_to_id ON tickets (previous_assigned_to_id);

CREATE INDEX idx_tickets_previous_department_id ON tickets (previous_department_id);

CREATE INDEX idx_ticket_escalations_ticket_id ON ticket_escalations (ticket_id);

CREATE INDEX idx_ticket_escalations_from_level ON ticket_escalations (from_level);

CREATE INDEX idx_ticket_escalations_to_level ON ticket_escalations (to_level);

CREATE INDEX idx_ticket_escalations_escalated_at ON ticket_escalations (escalated_at);

CREATE INDEX idx_ticket_escalations_escalated_by_id ON ticket_escalations (escalated_by_id);

CREATE INDEX idx_ticket_escalations_trigger_type ON ticket_escalations (trigger_type);

-- Add comments on tables and columns for documentation
COMMENT ON TABLE ticket_escalations IS 'Tracks escalation history and workflow for tickets';

COMMENT ON COLUMN departments.support_level IS '0=L0, 1=L1, 2=L2, 3=L3 - Support tier level';

COMMENT ON COLUMN departments.auto_assignment_rule IS 'Algorithm for auto-assignment: round_robin, least_loaded, skill_based';

COMMENT ON COLUMN tickets.current_level IS 'Current escalation level (0-3)';

COMMENT ON COLUMN tickets.escalation_path IS 'JSON array tracking department path through escalation';

COMMENT ON COLUMN ticket_escalations.trigger_type IS 'What triggered escalation: sla_breach, manual, priority_change, workload';