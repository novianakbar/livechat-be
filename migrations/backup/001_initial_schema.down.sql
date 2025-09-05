-- Drop indexes for ticketing tables
DROP INDEX IF EXISTS idx_ticket_sla_resolution_breached;

DROP INDEX IF EXISTS idx_ticket_sla_first_response_breached;

DROP INDEX IF EXISTS idx_ticket_sla_resolution_due;

DROP INDEX IF EXISTS idx_ticket_sla_first_response_due;

DROP INDEX IF EXISTS idx_ticket_sla_ticket_id;

DROP INDEX IF EXISTS idx_ticket_history_created_at;

DROP INDEX IF EXISTS idx_ticket_history_action;

DROP INDEX IF EXISTS idx_ticket_history_user_id;

DROP INDEX IF EXISTS idx_ticket_history_ticket_id;

DROP INDEX IF EXISTS idx_ticket_attachments_deleted_at;

DROP INDEX IF EXISTS idx_ticket_attachments_created_at;

DROP INDEX IF EXISTS idx_ticket_attachments_uploaded_by;

DROP INDEX IF EXISTS idx_ticket_attachments_ticket_id;

DROP INDEX IF EXISTS idx_ticket_comments_deleted_at;

DROP INDEX IF EXISTS idx_ticket_comments_created_at;

DROP INDEX IF EXISTS idx_ticket_comments_is_from_customer;

DROP INDEX IF EXISTS idx_ticket_comments_is_internal;

DROP INDEX IF EXISTS idx_ticket_comments_user_id;

DROP INDEX IF EXISTS idx_ticket_comments_ticket_id;

DROP INDEX IF EXISTS idx_tickets_deleted_at;

DROP INDEX IF EXISTS idx_tickets_created_at;

DROP INDEX IF EXISTS idx_tickets_customer_email;

DROP INDEX IF EXISTS idx_tickets_access_token;

DROP INDEX IF EXISTS idx_tickets_created_by_id;

DROP INDEX IF EXISTS idx_tickets_department_id;

DROP INDEX IF EXISTS idx_tickets_assigned_to_id;

DROP INDEX IF EXISTS idx_tickets_category_id;

DROP INDEX IF EXISTS idx_tickets_priority;

DROP INDEX IF EXISTS idx_tickets_status;

DROP INDEX IF EXISTS idx_tickets_ticket_code;

DROP INDEX IF EXISTS idx_ticket_categories_deleted_at;

DROP INDEX IF EXISTS idx_ticket_categories_is_active;

DROP INDEX IF EXISTS idx_ticket_categories_code;

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
-- Drop ticketing tables first (due to foreign key constraints)
DROP TABLE IF EXISTS ticket_sla;

DROP TABLE IF EXISTS ticket_history;

DROP TABLE IF EXISTS ticket_attachments;

DROP TABLE IF EXISTS ticket_comments;

DROP TABLE IF EXISTS tickets;

DROP TABLE IF EXISTS ticket_categories;

-- Drop chat tables
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