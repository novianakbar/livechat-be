CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create departments table
CREATE TABLE departments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
);

-- Create users table
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    role VARCHAR(50) NOT NULL CHECK (role IN ('admin', 'agent')),
    is_active BOOLEAN DEFAULT true,
    department_id UUID REFERENCES departments(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
);

-- Create chat_users table (refactored from customers)
CREATE TABLE chat_users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    browser_uuid UUID UNIQUE, -- UUID dari browser untuk anonymous users
    oss_user_id VARCHAR(255), -- ID user dari sistem OSS
    email VARCHAR(255), -- Email untuk logged-in users
    is_anonymous BOOLEAN DEFAULT true,
    ip_address VARCHAR(45) NOT NULL,
    user_agent TEXT, -- Browser user agent
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    -- Constraints
    CHECK (
        (is_anonymous = true AND browser_uuid IS NOT NULL) OR
        (is_anonymous = false AND oss_user_id IS NOT NULL AND email IS NOT NULL)
    )
);

-- Create chat_sessions table
CREATE TABLE chat_sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    chat_user_id UUID NOT NULL REFERENCES chat_users(id),
    agent_id UUID REFERENCES users(id),
    department_id UUID REFERENCES departments(id),
    topic VARCHAR(255) NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'waiting' CHECK (status IN ('waiting', 'active', 'closed')),
    priority VARCHAR(50) DEFAULT 'normal' CHECK (priority IN ('low', 'normal', 'high', 'urgent')),
    started_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    ended_at TIMESTAMP NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
);

-- Create chat_session_contacts table
CREATE TABLE chat_session_contacts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    session_id UUID NOT NULL REFERENCES chat_sessions(id),
    contact_name VARCHAR(255) NOT NULL,
    contact_email VARCHAR(255) NOT NULL,
    contact_phone VARCHAR(50),
    position VARCHAR(255), -- Job position (optional)
    company_name VARCHAR(255), -- Company name if applicable
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    UNIQUE(session_id) -- One contact per session
);

-- Create chat_messages table
CREATE TABLE chat_messages (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    session_id UUID NOT NULL REFERENCES chat_sessions(id),
    sender_id UUID REFERENCES users(id),
    sender_type VARCHAR(50) NOT NULL CHECK (sender_type IN ('customer', 'agent', 'system')),
    message TEXT NOT NULL,
    message_type VARCHAR(50) DEFAULT 'text' CHECK (message_type IN ('text', 'image', 'file', 'system')),
    attachments JSON,
    read_at TIMESTAMP NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
);

-- Create chat_logs table
CREATE TABLE chat_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    session_id UUID NOT NULL REFERENCES chat_sessions(id),
    action VARCHAR(50) NOT NULL CHECK (action IN ('started', 'waiting', 'response', 'closed', 'transferred')),
    details TEXT,
    user_id UUID REFERENCES users(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
);

-- Create chat_tags table
CREATE TABLE chat_tags (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) UNIQUE NOT NULL,
    color VARCHAR(7) DEFAULT '#007bff',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
);

-- Create chat_session_tags table (many-to-many)
CREATE TABLE chat_session_tags (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    session_id UUID NOT NULL REFERENCES chat_sessions(id),
    tag_id UUID NOT NULL REFERENCES chat_tags(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    UNIQUE(session_id, tag_id)
);

-- Create agent_status table (tracks login sessions)
CREATE TABLE agent_status (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    agent_id UUID NOT NULL REFERENCES users(id),
    status VARCHAR(50) NOT NULL CHECK (status IN ('logged_in', 'logged_out')),
    last_login_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    UNIQUE(agent_id)
);

-- Create chat_analytics table
CREATE TABLE chat_analytics (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    date DATE NOT NULL,
    total_sessions INTEGER DEFAULT 0,
    completed_sessions INTEGER DEFAULT 0,
    average_response_time FLOAT DEFAULT 0,
    total_messages INTEGER DEFAULT 0,
    department_id UUID REFERENCES departments(id),
    agent_id UUID REFERENCES users(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
);

-- Create indexes for better performance
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_department_id ON users(department_id);
CREATE INDEX idx_users_role ON users(role);
CREATE INDEX idx_users_deleted_at ON users(deleted_at);

CREATE INDEX idx_chat_users_browser_uuid ON chat_users(browser_uuid);
CREATE INDEX idx_chat_users_oss_user_id ON chat_users(oss_user_id);
CREATE INDEX idx_chat_users_email ON chat_users(email);
CREATE INDEX idx_chat_users_is_anonymous ON chat_users(is_anonymous);
CREATE INDEX idx_chat_users_deleted_at ON chat_users(deleted_at);

CREATE INDEX idx_chat_sessions_chat_user_id ON chat_sessions(chat_user_id);
CREATE INDEX idx_chat_sessions_agent_id ON chat_sessions(agent_id);
CREATE INDEX idx_chat_sessions_department_id ON chat_sessions(department_id);
CREATE INDEX idx_chat_sessions_status ON chat_sessions(status);
CREATE INDEX idx_chat_sessions_priority ON chat_sessions(priority);
CREATE INDEX idx_chat_sessions_started_at ON chat_sessions(started_at);
CREATE INDEX idx_chat_sessions_deleted_at ON chat_sessions(deleted_at);

CREATE INDEX idx_chat_session_contacts_session_id ON chat_session_contacts(session_id);
CREATE INDEX idx_chat_session_contacts_contact_email ON chat_session_contacts(contact_email);
CREATE INDEX idx_chat_session_contacts_deleted_at ON chat_session_contacts(deleted_at);

CREATE INDEX idx_chat_messages_session_id ON chat_messages(session_id);
CREATE INDEX idx_chat_messages_sender_id ON chat_messages(sender_id);
CREATE INDEX idx_chat_messages_sender_type ON chat_messages(sender_type);
CREATE INDEX idx_chat_messages_created_at ON chat_messages(created_at);
CREATE INDEX idx_chat_messages_deleted_at ON chat_messages(deleted_at);

CREATE INDEX idx_chat_logs_session_id ON chat_logs(session_id);
CREATE INDEX idx_chat_logs_action ON chat_logs(action);
CREATE INDEX idx_chat_logs_created_at ON chat_logs(created_at);
CREATE INDEX idx_chat_logs_deleted_at ON chat_logs(deleted_at);

CREATE INDEX idx_chat_session_tags_session_id ON chat_session_tags(session_id);
CREATE INDEX idx_chat_session_tags_tag_id ON chat_session_tags(tag_id);
CREATE INDEX idx_chat_session_tags_deleted_at ON chat_session_tags(deleted_at);

CREATE INDEX idx_agent_status_agent_id ON agent_status(agent_id);
CREATE INDEX idx_agent_status_status ON agent_status(status);
CREATE INDEX idx_agent_status_deleted_at ON agent_status(deleted_at);

CREATE INDEX idx_chat_analytics_date ON chat_analytics(date);
CREATE INDEX idx_chat_analytics_department_id ON chat_analytics(department_id);
CREATE INDEX idx_chat_analytics_agent_id ON chat_analytics(agent_id);
CREATE INDEX idx_chat_analytics_deleted_at ON chat_analytics(deleted_at);

-- Create triggers for updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_departments_updated_at BEFORE UPDATE ON departments FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON users FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_chat_users_updated_at BEFORE UPDATE ON chat_users FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_chat_sessions_updated_at BEFORE UPDATE ON chat_sessions FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_chat_session_contacts_updated_at BEFORE UPDATE ON chat_session_contacts FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_chat_messages_updated_at BEFORE UPDATE ON chat_messages FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_chat_tags_updated_at BEFORE UPDATE ON chat_tags FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_agent_status_updated_at BEFORE UPDATE ON agent_status FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_chat_analytics_updated_at BEFORE UPDATE ON chat_analytics FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
