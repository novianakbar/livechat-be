# Migration Consolidation Summary

## Overview
Migration files have been consolidated from 4 separate migrations into 2 comprehensive migrations for easier deployment and management.

## New Migration Structure

### 001_initial_schema.up.sql
**Complete database schema including:**
- All core system tables (departments, users)
- Chat system tables (chat_users, chat_sessions, chat_messages, etc.)
- Ticketing system tables (tickets, ticket_categories, ticket_comments, etc.)
- Multi-level support system (integrated into core tables)
- All necessary indexes for performance
- Table comments for documentation

**Key Features Included:**
- Multi-level department hierarchy (L0-L3)
- Ticket escalation system
- Chat system with anonymous and authenticated users
- Comprehensive indexing strategy
- Soft delete support across all tables

### 002_seed_data.up.sql
**Complete seed data including:**
- 4-level department hierarchy (General Support → Technical → Senior Technical → Management)
- Admin user and agents for each level
- Essential chat tags
- Ticket categories with appropriate SLAs
- Sample data for testing (demo chat session, ticket, etc.)

## Migration History Backup
Original migration files are preserved in `/migrations/backup/` directory:
- 001_initial_schema.{up,down}.sql
- 002_seed_data.{up,down}.sql  
- 003_multi_level_support.{up,down}.sql
- 004_multi_level_seed_data.{up,down}.sql

## Benefits
1. **Simplified deployment** - Only 2 migrations to run instead of 4
2. **Consistent state** - All features available from initial deployment
3. **Better maintenance** - Easier to understand complete system structure
4. **Production ready** - Includes all necessary indexes and constraints

## Usage
For fresh installation:
```bash
# Run schema migration
migrate -path migrations -database "postgres://..." up 1

# Run seed data migration  
migrate -path migrations -database "postgres://..." up 2
```

## Default Credentials
- Admin: admin@livechat.com / password
- Agent L0: agent1@livechat.com / password
- Agent L1: agent2@livechat.com / password
- Agent L2: senior1@livechat.com, senior2@livechat.com / password
- Agent L3: manager1@livechat.com / password

*Note: Change default passwords before production deployment*
