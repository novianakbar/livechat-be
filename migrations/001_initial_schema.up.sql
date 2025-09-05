CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- ============================================
-- CORE SYSTEM TABLES
-- ============================================
-- Create departments table with multi-level support
CREATE TABLE
    departments (
        id VARCHAR(255) PRIMARY KEY,
        name VARCHAR(255) NOT NULL,
        description TEXT,
        is_active BOOLEAN DEFAULT true,
        can_handle_tickets BOOLEAN DEFAULT true,
        max_tickets_per_agent INTEGER DEFAULT 10,
        -- Multi-level support fields
        support_level INTEGER DEFAULT 0,
        parent_dept_id VARCHAR(255) REFERENCES departments (id),
        max_escalation_level INTEGER DEFAULT 3,
        auto_assignment_rule VARCHAR(50) DEFAULT 'round_robin',
        escalation_dept_id VARCHAR(255) REFERENCES departments (id),
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        deleted_at BIGINT DEFAULT 0 -- For soft delete support (0 = not deleted, unix timestamp = deleted)
    );

-- Create users table
CREATE TABLE
    users (
        id VARCHAR(255) PRIMARY KEY,
        email VARCHAR(255) UNIQUE NOT NULL,
        password VARCHAR(255) NOT NULL,
        name VARCHAR(255) NOT NULL,
        role VARCHAR(50) NOT NULL CHECK (role IN ('admin', 'agent')),
        is_active BOOLEAN DEFAULT true,
        department_id VARCHAR(255) REFERENCES departments (id),
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        deleted_at BIGINT DEFAULT 0 -- For soft delete support (0 = not deleted, unix timestamp = deleted)
    );

-- ============================================
-- CHAT SYSTEM TABLES
-- ============================================
-- Create chat_users table (refactored from customers)
CREATE TABLE
    chat_users (
        id VARCHAR(255) PRIMARY KEY,
        browser_uuid VARCHAR(255) UNIQUE, -- UUID dari browser untuk anonymous users
        oss_user_id VARCHAR(255), -- ID user dari sistem OSS
        email VARCHAR(255), -- Email untuk logged-in users
        is_anonymous BOOLEAN DEFAULT true,
        ip_address VARCHAR(45) NOT NULL,
        user_agent TEXT, -- Browser user agent
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        deleted_at BIGINT DEFAULT 0, -- For soft delete support (0 = not deleted, unix timestamp = deleted)
        -- Constraints
        CHECK (
            (
                is_anonymous = true
                AND browser_uuid IS NOT NULL
            )
            OR (
                is_anonymous = false
                AND oss_user_id IS NOT NULL
                AND email IS NOT NULL
            )
        )
    );

-- Create chat_sessions table
CREATE TABLE
    chat_sessions (
        id VARCHAR(255) PRIMARY KEY,
        chat_user_id VARCHAR(255) NOT NULL REFERENCES chat_users (id),
        agent_id VARCHAR(255) REFERENCES users (id),
        department_id VARCHAR(255) REFERENCES departments (id),
        topic VARCHAR(255) NOT NULL,
        status VARCHAR(50) NOT NULL DEFAULT 'waiting' CHECK (status IN ('waiting', 'active', 'closed')),
        priority VARCHAR(50) DEFAULT 'normal' CHECK (priority IN ('low', 'normal', 'high', 'urgent')),
        started_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        ended_at TIMESTAMP NULL,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        deleted_at BIGINT DEFAULT 0 -- For soft delete support (0 = not deleted, unix timestamp = deleted)
    );

-- Create chat_session_contacts table
CREATE TABLE
    chat_session_contacts (
        id VARCHAR(255) PRIMARY KEY,
        session_id VARCHAR(255) NOT NULL REFERENCES chat_sessions (id),
        contact_name VARCHAR(255) NOT NULL,
        contact_email VARCHAR(255) NOT NULL,
        contact_phone VARCHAR(50),
        position VARCHAR(255), -- Job position (optional)
        company_name VARCHAR(255), -- Company name if applicable
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        deleted_at BIGINT DEFAULT 0, -- For soft delete support (0 = not deleted, unix timestamp = deleted)
        UNIQUE (session_id) -- One contact per session
    );

-- Create chat_messages table
CREATE TABLE
    chat_messages (
        id VARCHAR(255) PRIMARY KEY,
        session_id VARCHAR(255) NOT NULL REFERENCES chat_sessions (id),
        sender_id VARCHAR(255) REFERENCES users (id),
        sender_type VARCHAR(50) NOT NULL CHECK (sender_type IN ('customer', 'agent', 'system')),
        message TEXT NOT NULL,
        message_type VARCHAR(50) DEFAULT 'text' CHECK (
            message_type IN ('text', 'image', 'file', 'system')
        ),
        attachments JSON,
        read_at TIMESTAMP NULL,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        deleted_at BIGINT DEFAULT 0 -- For soft delete support (0 = not deleted, unix timestamp = deleted)
    );

-- Create chat_logs table
CREATE TABLE
    chat_logs (
        id VARCHAR(255) PRIMARY KEY,
        session_id VARCHAR(255) NOT NULL REFERENCES chat_sessions (id),
        action VARCHAR(50) NOT NULL CHECK (
            action IN (
                'started',
                'waiting',
                'response',
                'closed',
                'transferred'
            )
        ),
        details TEXT,
        user_id VARCHAR(255) REFERENCES users (id),
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        deleted_at BIGINT DEFAULT 0 -- For soft delete support (0 = not deleted, unix timestamp = deleted)
    );

-- Create chat_tags table
CREATE TABLE
    chat_tags (
        id VARCHAR(255) PRIMARY KEY,
        name VARCHAR(255) UNIQUE NOT NULL,
        color VARCHAR(7) DEFAULT '#007bff',
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        deleted_at BIGINT DEFAULT 0 -- For soft delete support (0 = not deleted, unix timestamp = deleted)
    );

-- Create chat_session_tags table (many-to-many)
CREATE TABLE
    chat_session_tags (
        id VARCHAR(255) PRIMARY KEY,
        session_id VARCHAR(255) NOT NULL REFERENCES chat_sessions (id),
        tag_id VARCHAR(255) NOT NULL REFERENCES chat_tags (id),
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        deleted_at BIGINT DEFAULT 0, -- For soft delete support (0 = not deleted, unix timestamp = deleted)
        UNIQUE (session_id, tag_id)
    );

-- Create agent_status table (tracks login sessions)
CREATE TABLE
    agent_status (
        id VARCHAR(255) PRIMARY KEY,
        agent_id VARCHAR(255) NOT NULL REFERENCES users (id),
        status VARCHAR(50) NOT NULL CHECK (status IN ('logged_in', 'logged_out')),
        last_login_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        deleted_at BIGINT DEFAULT 0, -- For soft delete support (0 = not deleted, unix timestamp = deleted)
        UNIQUE (agent_id)
    );

-- Create chat_analytics table
CREATE TABLE
    chat_analytics (
        id VARCHAR(255) PRIMARY KEY,
        date DATE NOT NULL,
        total_sessions INTEGER DEFAULT 0,
        completed_sessions INTEGER DEFAULT 0,
        average_response_time FLOAT DEFAULT 0,
        total_messages INTEGER DEFAULT 0,
        department_id VARCHAR(255) REFERENCES departments (id),
        agent_id VARCHAR(255) REFERENCES users (id),
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        deleted_at BIGINT DEFAULT 0 -- For soft delete support (0 = not deleted, unix timestamp = deleted)
    );

-- ============================================
-- TICKETING SYSTEM TABLES
-- ============================================
-- Create ticket_categories table
CREATE TABLE
    ticket_categories (
        id VARCHAR(255) PRIMARY KEY,
        name VARCHAR(100) NOT NULL,
        code VARCHAR(20) UNIQUE NOT NULL,
        description TEXT,
        color VARCHAR(7),
        is_active BOOLEAN DEFAULT true,
        sla_first_response INTEGER DEFAULT 60,
        sla_resolution INTEGER DEFAULT 1440,
        default_department_id VARCHAR(255) REFERENCES departments (id),
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        deleted_at BIGINT DEFAULT 0 -- For soft delete support (0 = not deleted, unix timestamp = deleted)
    );

-- Create tickets table with multi-level support
CREATE TABLE
    tickets (
        id VARCHAR(255) PRIMARY KEY,
        ticket_code VARCHAR(20) UNIQUE NOT NULL,
        subject VARCHAR(255) NOT NULL,
        description TEXT,
        customer_name VARCHAR(100),
        customer_email VARCHAR(100),
        customer_phone VARCHAR(20),
        category_id VARCHAR(255) REFERENCES ticket_categories (id),
        priority VARCHAR(20) DEFAULT 'medium' CHECK (priority IN ('low', 'medium', 'high', 'urgent')),
        status VARCHAR(20) DEFAULT 'open' CHECK (
            status IN (
                'open',
                'in_progress',
                'resolved',
                'closed',
                'escalated'
            )
        ),
        assigned_to_id VARCHAR(255) REFERENCES users (id),
        department_id VARCHAR(255) REFERENCES departments (id),
        created_by_id VARCHAR(255) REFERENCES users (id),
        created_via VARCHAR(20) CHECK (created_via IN ('customer', 'agent', 'ai')),
        first_response_at TIMESTAMP NULL,
        resolved_at TIMESTAMP NULL,
        closed_at TIMESTAMP NULL,
        access_token VARCHAR(64),
        -- Multi-level support fields
        current_level INTEGER DEFAULT 0,
        escalation_path TEXT,
        can_escalate BOOLEAN DEFAULT true,
        max_level INTEGER DEFAULT 3,
        escalation_count INTEGER DEFAULT 0,
        previous_assigned_to_id VARCHAR(255) REFERENCES users (id),
        previous_department_id VARCHAR(255) REFERENCES departments (id),
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        deleted_at BIGINT DEFAULT 0 -- For soft delete support (0 = not deleted, unix timestamp = deleted)
    );

-- Create ticket_comments table
CREATE TABLE
    ticket_comments (
        id VARCHAR(255) PRIMARY KEY,
        ticket_id VARCHAR(255) NOT NULL REFERENCES tickets (id),
        user_id VARCHAR(255) REFERENCES users (id),
        content TEXT NOT NULL,
        is_internal BOOLEAN DEFAULT false,
        is_from_customer BOOLEAN DEFAULT false,
        author_name VARCHAR(100),
        author_email VARCHAR(100),
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        deleted_at BIGINT DEFAULT 0 -- For soft delete support (0 = not deleted, unix timestamp = deleted)
    );

-- Create ticket_attachments table
CREATE TABLE
    ticket_attachments (
        id VARCHAR(255) PRIMARY KEY,
        ticket_id VARCHAR(255) NOT NULL REFERENCES tickets (id),
        file_name VARCHAR(255) NOT NULL,
        file_path VARCHAR(500) NOT NULL,
        file_size BIGINT,
        file_type VARCHAR(50),
        uploaded_by VARCHAR(20) CHECK (uploaded_by IN ('customer', 'agent')),
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        deleted_at BIGINT DEFAULT 0 -- For soft delete support (0 = not deleted, unix timestamp = deleted)
    );

-- Create ticket_history table
CREATE TABLE
    ticket_history (
        id VARCHAR(255) PRIMARY KEY,
        ticket_id VARCHAR(255) NOT NULL REFERENCES tickets (id),
        user_id VARCHAR(255) REFERENCES users (id),
        action VARCHAR(50) NOT NULL,
        field_name VARCHAR(50),
        old_value VARCHAR(255),
        new_value VARCHAR(255),
        description TEXT,
        actor_name VARCHAR(100),
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );

-- Create ticket_sla table
CREATE TABLE
    ticket_sla (
        id VARCHAR(255) PRIMARY KEY,
        ticket_id VARCHAR(255) UNIQUE NOT NULL REFERENCES tickets (id),
        first_response_due TIMESTAMP NOT NULL,
        resolution_due TIMESTAMP NOT NULL,
        first_response_at TIMESTAMP NULL,
        resolved_at TIMESTAMP NULL,
        first_response_breached BOOLEAN DEFAULT false,
        resolution_breached BOOLEAN DEFAULT false,
        first_response_time INTEGER DEFAULT 0,
        resolution_time INTEGER DEFAULT 0,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );

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

-- ============================================
-- INDEXES FOR PERFORMANCE
-- ============================================
-- User indexes
CREATE INDEX idx_users_email ON users (email);

CREATE INDEX idx_users_department_id ON users (department_id);

CREATE INDEX idx_users_role ON users (role);

CREATE INDEX idx_users_deleted_at ON users (deleted_at);

-- Department indexes
CREATE INDEX idx_departments_deleted_at ON departments (deleted_at);

CREATE INDEX idx_departments_support_level ON departments (support_level);

CREATE INDEX idx_departments_parent_dept_id ON departments (parent_dept_id);

CREATE INDEX idx_departments_escalation_dept_id ON departments (escalation_dept_id);

-- Chat user indexes
CREATE INDEX idx_chat_users_browser_uuid ON chat_users (browser_uuid);

CREATE INDEX idx_chat_users_oss_user_id ON chat_users (oss_user_id);

CREATE INDEX idx_chat_users_email ON chat_users (email);

CREATE INDEX idx_chat_users_is_anonymous ON chat_users (is_anonymous);

CREATE INDEX idx_chat_users_deleted_at ON chat_users (deleted_at);

-- Chat session indexes
CREATE INDEX idx_chat_sessions_chat_user_id ON chat_sessions (chat_user_id);

CREATE INDEX idx_chat_sessions_agent_id ON chat_sessions (agent_id);

CREATE INDEX idx_chat_sessions_department_id ON chat_sessions (department_id);

CREATE INDEX idx_chat_sessions_status ON chat_sessions (status);

CREATE INDEX idx_chat_sessions_priority ON chat_sessions (priority);

CREATE INDEX idx_chat_sessions_started_at ON chat_sessions (started_at);

CREATE INDEX idx_chat_sessions_deleted_at ON chat_sessions (deleted_at);

-- Chat session contact indexes
CREATE INDEX idx_chat_session_contacts_session_id ON chat_session_contacts (session_id);

CREATE INDEX idx_chat_session_contacts_contact_email ON chat_session_contacts (contact_email);

CREATE INDEX idx_chat_session_contacts_deleted_at ON chat_session_contacts (deleted_at);

-- Chat message indexes
CREATE INDEX idx_chat_messages_session_id ON chat_messages (session_id);

CREATE INDEX idx_chat_messages_sender_id ON chat_messages (sender_id);

CREATE INDEX idx_chat_messages_sender_type ON chat_messages (sender_type);

CREATE INDEX idx_chat_messages_created_at ON chat_messages (created_at);

CREATE INDEX idx_chat_messages_deleted_at ON chat_messages (deleted_at);

-- Chat log indexes
CREATE INDEX idx_chat_logs_session_id ON chat_logs (session_id);

CREATE INDEX idx_chat_logs_action ON chat_logs (action);

CREATE INDEX idx_chat_logs_created_at ON chat_logs (created_at);

CREATE INDEX idx_chat_logs_deleted_at ON chat_logs (deleted_at);

-- Chat tag indexes
CREATE INDEX idx_chat_tags_deleted_at ON chat_tags (deleted_at);

CREATE INDEX idx_chat_session_tags_session_id ON chat_session_tags (session_id);

CREATE INDEX idx_chat_session_tags_tag_id ON chat_session_tags (tag_id);

CREATE INDEX idx_chat_session_tags_deleted_at ON chat_session_tags (deleted_at);

-- Agent status indexes
CREATE INDEX idx_agent_status_agent_id ON agent_status (agent_id);

CREATE INDEX idx_agent_status_status ON agent_status (status);

CREATE INDEX idx_agent_status_deleted_at ON agent_status (deleted_at);

-- Chat analytics indexes
CREATE INDEX idx_chat_analytics_date ON chat_analytics (date);

CREATE INDEX idx_chat_analytics_department_id ON chat_analytics (department_id);

CREATE INDEX idx_chat_analytics_agent_id ON chat_analytics (agent_id);

CREATE INDEX idx_chat_analytics_deleted_at ON chat_analytics (deleted_at);

-- Ticket category indexes
CREATE INDEX idx_ticket_categories_code ON ticket_categories (code);

CREATE INDEX idx_ticket_categories_is_active ON ticket_categories (is_active);

CREATE INDEX idx_ticket_categories_deleted_at ON ticket_categories (deleted_at);

-- Ticket indexes
CREATE INDEX idx_tickets_ticket_code ON tickets (ticket_code);

CREATE INDEX idx_tickets_status ON tickets (status);

CREATE INDEX idx_tickets_priority ON tickets (priority);

CREATE INDEX idx_tickets_category_id ON tickets (category_id);

CREATE INDEX idx_tickets_assigned_to_id ON tickets (assigned_to_id);

CREATE INDEX idx_tickets_department_id ON tickets (department_id);

CREATE INDEX idx_tickets_created_by_id ON tickets (created_by_id);

CREATE INDEX idx_tickets_access_token ON tickets (access_token);

CREATE INDEX idx_tickets_customer_email ON tickets (customer_email);

CREATE INDEX idx_tickets_created_at ON tickets (created_at);

CREATE INDEX idx_tickets_deleted_at ON tickets (deleted_at);

CREATE INDEX idx_tickets_current_level ON tickets (current_level);

CREATE INDEX idx_tickets_escalation_count ON tickets (escalation_count);

CREATE INDEX idx_tickets_previous_assigned_to_id ON tickets (previous_assigned_to_id);

CREATE INDEX idx_tickets_previous_department_id ON tickets (previous_department_id);

-- Ticket comment indexes
CREATE INDEX idx_ticket_comments_ticket_id ON ticket_comments (ticket_id);

CREATE INDEX idx_ticket_comments_user_id ON ticket_comments (user_id);

CREATE INDEX idx_ticket_comments_is_internal ON ticket_comments (is_internal);

CREATE INDEX idx_ticket_comments_is_from_customer ON ticket_comments (is_from_customer);

CREATE INDEX idx_ticket_comments_created_at ON ticket_comments (created_at);

CREATE INDEX idx_ticket_comments_deleted_at ON ticket_comments (deleted_at);

-- Ticket attachment indexes
CREATE INDEX idx_ticket_attachments_ticket_id ON ticket_attachments (ticket_id);

CREATE INDEX idx_ticket_attachments_uploaded_by ON ticket_attachments (uploaded_by);

CREATE INDEX idx_ticket_attachments_created_at ON ticket_attachments (created_at);

CREATE INDEX idx_ticket_attachments_deleted_at ON ticket_attachments (deleted_at);

-- Ticket history indexes
CREATE INDEX idx_ticket_history_ticket_id ON ticket_history (ticket_id);

CREATE INDEX idx_ticket_history_user_id ON ticket_history (user_id);

CREATE INDEX idx_ticket_history_action ON ticket_history (action);

CREATE INDEX idx_ticket_history_created_at ON ticket_history (created_at);

-- Ticket SLA indexes
CREATE INDEX idx_ticket_sla_ticket_id ON ticket_sla (ticket_id);

CREATE INDEX idx_ticket_sla_first_response_due ON ticket_sla (first_response_due);

CREATE INDEX idx_ticket_sla_resolution_due ON ticket_sla (resolution_due);

CREATE INDEX idx_ticket_sla_first_response_breached ON ticket_sla (first_response_breached);

CREATE INDEX idx_ticket_sla_resolution_breached ON ticket_sla (resolution_breached);

-- Ticket escalation indexes
CREATE INDEX idx_ticket_escalations_ticket_id ON ticket_escalations (ticket_id);

CREATE INDEX idx_ticket_escalations_from_level ON ticket_escalations (from_level);

CREATE INDEX idx_ticket_escalations_to_level ON ticket_escalations (to_level);

CREATE INDEX idx_ticket_escalations_escalated_at ON ticket_escalations (escalated_at);

CREATE INDEX idx_ticket_escalations_escalated_by_id ON ticket_escalations (escalated_by_id);

CREATE INDEX idx_ticket_escalations_trigger_type ON ticket_escalations (trigger_type);

-- ============================================
-- COMMENTS FOR DOCUMENTATION
-- ============================================
COMMENT ON TABLE departments IS 'Departments with multi-level support structure';

COMMENT ON COLUMN departments.support_level IS '0=L0, 1=L1, 2=L2, 3=L3 - Support tier level';

COMMENT ON COLUMN departments.auto_assignment_rule IS 'Algorithm for auto-assignment: round_robin, least_loaded, skill_based';

COMMENT ON TABLE tickets IS 'Tickets with multi-level escalation support';

COMMENT ON COLUMN tickets.current_level IS 'Current escalation level (0-3)';

COMMENT ON COLUMN tickets.escalation_path IS 'JSON array tracking department path through escalation';

COMMENT ON TABLE ticket_escalations IS 'Tracks escalation history and workflow for tickets';

COMMENT ON COLUMN ticket_escalations.trigger_type IS 'What triggered escalation: sla_breach, manual, priority_change, workload';