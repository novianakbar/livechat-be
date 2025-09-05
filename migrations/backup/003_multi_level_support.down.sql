-- Rollback migration untuk Multi-Level Support System
-- Drop indexes
DROP INDEX IF EXISTS idx_ticket_escalations_trigger_type;

DROP INDEX IF EXISTS idx_ticket_escalations_escalated_by_id;

DROP INDEX IF EXISTS idx_ticket_escalations_escalated_at;

DROP INDEX IF EXISTS idx_ticket_escalations_to_level;

DROP INDEX IF EXISTS idx_ticket_escalations_from_level;

DROP INDEX IF EXISTS idx_ticket_escalations_ticket_id;

DROP INDEX IF EXISTS idx_tickets_previous_department_id;

DROP INDEX IF EXISTS idx_tickets_previous_assigned_to_id;

DROP INDEX IF EXISTS idx_tickets_escalation_count;

DROP INDEX IF EXISTS idx_tickets_current_level;

DROP INDEX IF EXISTS idx_departments_escalation_dept_id;

DROP INDEX IF EXISTS idx_departments_parent_dept_id;

DROP INDEX IF EXISTS idx_departments_support_level;

-- Drop table
DROP TABLE IF EXISTS ticket_escalations;

-- Remove columns from tickets table
ALTER TABLE tickets
DROP COLUMN IF EXISTS previous_department_id,
DROP COLUMN IF EXISTS previous_assigned_to_id,
DROP COLUMN IF EXISTS escalation_count,
DROP COLUMN IF EXISTS max_level,
DROP COLUMN IF EXISTS can_escalate,
DROP COLUMN IF EXISTS escalation_path,
DROP COLUMN IF EXISTS current_level;

-- Remove columns from departments table
ALTER TABLE departments
DROP COLUMN IF EXISTS escalation_dept_id,
DROP COLUMN IF EXISTS auto_assignment_rule,
DROP COLUMN IF EXISTS max_escalation_level,
DROP COLUMN IF EXISTS parent_dept_id,
DROP COLUMN IF EXISTS support_level;