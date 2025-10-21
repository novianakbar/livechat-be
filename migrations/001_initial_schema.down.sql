-- Drop all indexes in reverse order

DROP INDEX IF EXISTS idx_chat_analytics_deleted_at;

DROP INDEX IF EXISTS idx_chat_analytics_agent_id;

DROP INDEX IF EXISTS idx_chat_analytics_department_id;

DROP INDEX IF EXISTS idx_chat_analytics_date;

DROP INDEX IF EXISTS idx_agent_status_deleted_at;

DROP INDEX IF EXISTS idx_agent_status_status;

DROP INDEX IF EXISTS idx_agent_status_agent_id;

DROP INDEX IF EXISTS idx_chat_session_tags_deleted_at;

DROP INDEX IF EXISTS idx_chat_session_tags_tag_id;

DROP INDEX IF EXISTS idx_chat_session_tags_session_id;

DROP INDEX IF EXISTS idx_chat_tags_deleted_at;

DROP INDEX IF EXISTS idx_chat_logs_deleted_at;

DROP INDEX IF EXISTS idx_chat_logs_created_at;

DROP INDEX IF EXISTS idx_chat_logs_action;

DROP INDEX IF EXISTS idx_chat_logs_session_id;

DROP INDEX IF EXISTS idx_chat_messages_deleted_at;

DROP INDEX IF EXISTS idx_chat_messages_created_at;

DROP INDEX IF EXISTS idx_chat_messages_sender_type;

DROP INDEX IF EXISTS idx_chat_messages_sender_id;

DROP INDEX IF EXISTS idx_chat_messages_session_id;

DROP INDEX IF EXISTS idx_chat_session_contacts_deleted_at;

DROP INDEX IF EXISTS idx_chat_session_contacts_contact_email;

DROP INDEX IF EXISTS idx_chat_session_contacts_session_id;

DROP INDEX IF EXISTS idx_chat_sessions_deleted_at;

DROP INDEX IF EXISTS idx_chat_sessions_started_at;

DROP INDEX IF EXISTS idx_chat_sessions_priority;

DROP INDEX IF EXISTS idx_chat_sessions_status;

DROP INDEX IF EXISTS idx_chat_sessions_department_id;

DROP INDEX IF EXISTS idx_chat_sessions_agent_id;

DROP INDEX IF EXISTS idx_chat_sessions_chat_user_id;

DROP INDEX IF EXISTS idx_chat_users_deleted_at;

DROP INDEX IF EXISTS idx_chat_users_is_anonymous;

DROP INDEX IF EXISTS idx_chat_users_email;

DROP INDEX IF EXISTS idx_chat_users_oss_user_id;

DROP INDEX IF EXISTS idx_chat_users_browser_uuid;

DROP INDEX IF EXISTS idx_departments_escalation_dept_id;

DROP INDEX IF EXISTS idx_departments_parent_dept_id;

DROP INDEX IF EXISTS idx_departments_support_level;

DROP INDEX IF EXISTS idx_departments_deleted_at;

DROP INDEX IF EXISTS idx_users_deleted_at;

DROP INDEX IF EXISTS idx_users_role;

DROP INDEX IF EXISTS idx_users_department_id;

DROP INDEX IF EXISTS idx_users_email;

-- Drop all tables in reverse dependency order

DROP TABLE IF EXISTS chat_analytics;

DROP TABLE IF EXISTS agent_status;

DROP TABLE IF EXISTS chat_session_tags;

DROP TABLE IF EXISTS chat_tags;

DROP TABLE IF EXISTS chat_logs;

DROP TABLE IF EXISTS chat_messages;

DROP TABLE IF EXISTS chat_session_contacts;

DROP TABLE IF EXISTS chat_sessions;

DROP TABLE IF EXISTS chat_users;

DROP TABLE IF EXISTS users;

DROP TABLE IF EXISTS departments;

-- Drop extension if no other tables need it
DROP EXTENSION IF EXISTS "uuid-ossp";