# LiveChat API Routes Documentation

## Overview
Dokumentasi lengkap untuk semua API routes yang tersedia dalam sistem LiveChat OSS.

## Base URL
```
http://localhost:8080/api
```

---

## 1. Health Check & System

### Health Check
- **GET** `/health`
- **Description**: Mengecek status kesehatan aplikasi
- **Auth**: None
- **Response**:
```json
{
  "status": "ok",
  "message": "LiveChat API is running"
}
```

---

## 2. OSS Chat Routes (Public - No Auth Required)

### Base Path: `/api/chat`

#### Start Chat Session
- **POST** `/api/chat/start`
- **Description**: Memulai sesi chat baru untuk pengguna OSS
- **Auth**: None
- **Request Body**:
```json
{
  "browser_uuid": "550e8400-e29b-41d4-a716-446655440000", // Optional: untuk anonymous user
  "oss_user_id": "USER123",                              // Optional: untuk logged-in user
  "email": "user@example.com",                           // Optional: untuk logged-in user  
  "topic": "Pertanyaan tentang izin usaha",              // Required: topik chat
  "priority": "normal",                                   // Optional: low|normal|high|urgent
  "user_agent": "Mozilla/5.0 ..."                        // Optional: browser user agent
}
```

#### Set Session Contact
- **POST** `/api/chat/contact`
- **Description**: Mengisi informasi kontak untuk sesi chat
- **Auth**: None
- **Request Body**:
```json
{
  "session_id": "660e8400-e29b-41d4-a716-446655440000", // Required
  "contact_name": "John Doe",                            // Required
  "contact_email": "john@company.com",                   // Required
  "contact_phone": "+6281234567890",                     // Optional
  "position": "Manager",                                 // Optional
  "company_name": "PT. Example"                          // Optional
}
```

#### Link OSS User
- **POST** `/api/chat/link-user`
- **Description**: Menghubungkan user anonymous dengan akun OSS saat login
- **Auth**: None
- **Request Body**:
```json
{
  "browser_uuid": "550e8400-e29b-41d4-a716-446655440000", // Required
  "oss_user_id": "USER123",                              // Required
  "email": "user@example.com"                            // Required
}
```

#### Get Chat History
- **GET** `/api/chat/history`
- **Description**: Mengambil histori chat untuk pengguna tertentu
- **Auth**: None
- **Query Parameters**:
  - `browser_uuid` (string, optional): UUID browser untuk anonymous user
  - `oss_user_id` (string, optional): ID pengguna OSS untuk logged-in user
  - `limit` (int, optional): Jumlah sesi yang dikembalikan (default: 20)
  - `offset` (int, optional): Jumlah sesi yang dilewati (default: 0)

#### Get Session Details
- **GET** `/api/chat/session/{session_id}`
- **Description**: Mengambil detail sesi chat tertentu
- **Auth**: None

---

## 3. Legacy Public Routes (Backward Compatibility)

### Base Path: `/api/public`

#### Start Chat (Legacy)
- **POST** `/api/public/chat/start`
- **Description**: Legacy endpoint untuk memulai chat (backward compatibility)
- **Auth**: None

#### Send Message (Legacy)
- **POST** `/api/public/chat/message`
- **Description**: Legacy endpoint untuk mengirim pesan (backward compatibility)
- **Auth**: None

#### Get Session Messages (Legacy)
- **GET** `/api/public/chat/session/{session_id}/messages`
- **Description**: Legacy endpoint untuk mengambil pesan sesi (backward compatibility)
- **Auth**: None

---

## 4. Authentication Routes

### Base Path: `/api/auth`

#### Login
- **POST** `/api/auth/login`
- **Description**: Login untuk admin/agent dan mendapatkan JWT token
- **Auth**: None
- **Request Body**:
```json
{
  "email": "admin@livechat.com",
  "password": "password"
}
```

#### Refresh Token
- **POST** `/api/auth/refresh`
- **Description**: Refresh JWT token
- **Auth**: None

#### Logout
- **POST** `/api/auth/logout`
- **Description**: Logout dan invalidate token
- **Auth**: Bearer Token Required

#### Validate Session
- **GET** `/api/auth/validate`
- **Description**: Validasi token yang sedang aktif
- **Auth**: Bearer Token Required

#### Get Profile
- **GET** `/api/auth/profile`
- **Description**: Mendapatkan profil user yang sedang login
- **Auth**: Bearer Token Required

#### Register New User
- **POST** `/api/auth/register`
- **Description**: Registrasi user baru (admin only)
- **Auth**: Bearer Token Required (Admin Only)

---

## 5. Chat Management Routes (Protected)

### Base Path: `/api/chat-management`
**Auth**: Bearer Token Required

#### Agent Routes (`/api/chat-management/agent`)
**Auth**: Bearer Token Required + Agent Role

- **POST** `/agent/message` - Mengirim pesan sebagai agent
- **POST** `/agent/assign` - Assign sesi ke agent
- **POST** `/agent/close` - Menutup sesi chat
- **GET** `/agent/sessions` - Mendapatkan sesi yang ditangani agent
- **GET** `/agent/sessions/{id}/connection-status` - Status koneksi sesi
- **GET** `/agent/sessions/{id}` - Detail sesi tertentu

#### Admin Routes (`/api/chat-management/admin`)
**Auth**: Bearer Token Required + Admin Role

- **GET** `/admin/waiting` - Mendapatkan sesi yang menunggu
- **GET** `/admin/active` - Mendapatkan sesi yang aktif
- **POST** `/admin/assign` - Assign sesi ke agent
- **POST** `/admin/close` - Menutup sesi chat
- **GET** `/admin/sessions` - Mendapatkan semua sesi

---

## 6. User Management Routes

### Base Path: `/api/users`
**Auth**: Bearer Token Required

- **GET** `/` - Mendapatkan daftar semua user
- **GET** `/agents` - Mendapatkan daftar agent
- **GET** `/{id}` - Mendapatkan detail user tertentu

---

## 7. Analytics Routes

### Base Path: `/api/analytics`
**Auth**: Bearer Token Required

- **GET** `/dashboard` - Mendapatkan statistik dashboard
- **GET** `/` - Mendapatkan data analytics umum

---

## 8. Email Routes

### Base Path: `/api/email`
**Auth**: Bearer Token Required

- **POST** `/send` - Mengirim email umum
- **POST** `/welcome` - Mengirim email welcome
- **POST** `/password-reset` - Mengirim email reset password
- **POST** `/chat-transcript` - Mengirim transkrip chat via email
- **POST** `/custom` - Mengirim email custom

---

## Chat Flow Documentation

### Flow 1: Anonymous User Chat
```
1. Frontend generates browser_uuid → localStorage
2. POST /api/chat/start (with browser_uuid, topic)
   ← Response: {session_id, chat_user_id, requires_contact: true}
3. POST /api/chat/contact (with session_id, contact info)
   ← Response: {contact_id, message: "success"}
4. Chat session starts, user can send messages
5. (Optional) User login → POST /api/chat/link-user
```

### Flow 2: Logged-in OSS User Chat
```
1. User already logged in OSS system
2. POST /api/chat/start (with oss_user_id, email, topic)
   ← Response: {session_id, chat_user_id, requires_contact: true}
3. POST /api/chat/contact (with session_id, contact info)
   ← Response: {contact_id, message: "success"}
4. Chat session starts, user can send messages
5. GET /api/chat/history (to see previous chats)
```

### Flow 3: Anonymous → Login Transition
```
1. Start as anonymous (browser_uuid)
2. POST /api/chat/start → POST /api/chat/contact
3. Chat berlangsung
4. User login di OSS system
5. POST /api/chat/link-user (browser_uuid + oss_user_id + email)
6. Now user can access history via oss_user_id
```

### Flow 4: Agent/Admin Management
```
1. Agent/Admin login: POST /api/auth/login
2. Get waiting sessions: GET /api/chat-management/admin/waiting
3. Assign to agent: POST /api/chat-management/admin/assign
4. Agent handle session: POST /api/chat-management/agent/message
5. Close session: POST /api/chat-management/agent/close
```

---

## Key Differences from Previous Implementation

### 1. Route Structure Changes
- **Old**: `/api/v1/oss-chat/*` 
- **New**: `/api/chat/*` (cleaner, unified)

### 2. Authentication Separation
- **Public OSS Routes**: `/api/chat/*` (no auth)
- **Management Routes**: `/api/chat-management/*` (auth required)

### 3. Backward Compatibility
- Legacy routes maintained at `/api/public/chat/*`
- Existing integrations won't break

### 4. Enhanced Features
- Anonymous ↔ Login user transition
- Contact info per session (not per user)
- Unified chat history across anonymous/login states
- Better route organization for different user types

---

## Response Format

### Success Response
```json
{
  "success": true,
  "message": "Operation successful",
  "data": {...}
}
```

### Error Response
```json
{
  "success": false,
  "message": "Error message",
  "error": "Detailed error information"
}
```

---

## Notes

1. **OSS Chat Routes** (`/api/chat/*`) tidak memerlukan authentication karena digunakan oleh sistem eksternal
2. **Chat Management Routes** (`/api/chat-management/*`) memerlukan authentication untuk admin/agent
3. **Legacy Routes** (`/api/public/*`) dipertahankan untuk backward compatibility
4. Semua endpoint menggunakan JSON untuk request/response
5. CORS sudah dikonfigurasi untuk cross-origin requests
