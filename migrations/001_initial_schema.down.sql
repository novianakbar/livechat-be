-- Drop triggers for new tables
DROP TRIGGER IF EXISTS update_chat_session_contacts_updated_at ON chat_session_contacts;

DROP TRIGGER IF EXISTS update_chat_users_updated_at ON chat_users;

DROP TRIGGER IF EXISTS update_chat_analytics_updated_at ON chat_analytics;

DROP TRIGGER IF EXISTS update_agent_status_updated_at ON agent_status;

DROP TRIGGER IF EXISTS update_chat_tags_updated_at ON chat_tags;

DROP TRIGGER IF EXISTS update_chat_messages_updated_at ON chat_messages;

DROP TRIGGER IF EXISTS update_chat_sessions_updated_at ON chat_sessions;

DROP TRIGGER IF EXISTS update_users_updated_at ON users;

DROP TRIGGER IF EXISTS update_departments_updated_at ON departments;

-- Drop function
DROP FUNCTION IF EXISTS update_updated_at_column ();

-- Drop indexes
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

DROP INDEX IF EXISTS idx_departments_deleted_at;

DROP INDEX IF EXISTS idx_users_deleted_at;

DROP INDEX IF EXISTS idx_users_role;

DROP INDEX IF EXISTS idx_users_department_id;

DROP INDEX IF EXISTS idx_users_email;

-- Drop tables
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