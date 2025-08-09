# Database Migration Updates

## Overview
Database migration telah diperbarui untuk menyederhanakan data seed dan memastikan kompatibilitas dengan struktur aplikasi saat ini.

## Changes Made

### 1. Initial Schema (001_initial_schema.up.sql) ✓ Verified
Struktur database sudah sesuai dan mencakup semua tabel yang diperlukan:
- `departments` - Departemen untuk mengorganisir agent
- `users` - Admin dan agent users
- `chat_users` - Customer chat users (anonymous & OSS users)
- `chat_sessions` - Chat sessions
- `chat_session_contacts` - Contact information per session
- `chat_messages` - Chat messages
- `chat_logs` - Activity logs
- `chat_tags` - Tags untuk kategorisasi
- `chat_session_tags` - Many-to-many relation untuk session tags
- `agent_status` - Status login agent
- `chat_analytics` - Data analytics

### 2. Seed Data Simplification (002_seed_data.up.sql)

#### Before:
- 4 departments (Perizinan, Investasi, Perpajakan, Teknis)
- 6 users (1 admin + 5 agents)
- 10 chat tags
- 5 agent status records
- Extensive sample data (chat users, sessions, messages, logs)

#### After:
- **2 departments** (General Support, Technical Support)
- **3 users** (1 admin + 2 agents)
- **4 chat tags** (General Question, Technical Issue, Support Request, Urgent)
- **2 agent status records**
- **No sample chat data** (akan dibuat saat runtime)

### 3. Key Benefits

#### Data Minimalis & Focused:
- Hanya data essential untuk menjalankan aplikasi
- Tidak ada sample chat data yang mengotori database
- Setup lebih cepat dan bersih

#### Struktur yang Fleksibel:
- 2 departemen dasar yang bisa disesuaikan
- Tags yang generic tapi berguna
- User minimal tapi functional

#### Production Ready:
- Password menggunakan bcrypt hash yang aman
- Data struktur sesuai dengan kode aplikasi
- Foreign key relationships yang benar

## Default Users

### Admin User:
- **Email**: admin@livechat.com
- **Password**: password (hashed)
- **Role**: admin

### Agent Users:
- **Agent 1**: agent1@livechat.com (General Support)
- **Agent 2**: agent2@livechat.com (Technical Support)
- **Password**: password (hashed)
- **Role**: agent

## Usage

```bash
# Apply migrations
go run main.go migrate up

# Rollback if needed
go run main.go migrate down
```

## Next Steps

1. Sesuaikan departemen sesuai kebutuhan bisnis
2. Update tags sesuai kategori chat yang diinginkan
3. Tambah agent sesuai kebutuhan
4. Konfigurasi email dan sistem lainnya

## Database Schema Verification ✅

All tables dari kode aplikasi sudah ter-cover dalam migration:
- Entity structures match with Go domain models
- Foreign key relationships established
- Indexes created for performance
- Soft delete support (deleted_at field)
- Automatic timestamp updates (triggers)
