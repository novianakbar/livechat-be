# OSS LiveChat API Documentation

## Overview
API endpoints untuk sistem livechat OSS yang mendukung pengguna anonymous dan logged-in dari sistem OSS.

## Base URL
```
/api/chat
```

## Authentication
Tidak ada authentication yang diperlukan untuk endpoints OSS chat, karena sistem mengidentifikasi pengguna melalui:
- `browser_uuid` untuk pengguna anonymous
- `oss_user_id` + `email` untuk pengguna yang sudah login

## Endpoints

### 1. Start Chat Session
Memulai sesi chat baru untuk pengguna OSS.

**Endpoint:** `POST /api/chat/start`

**Request Body:**
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

**Response:**
```json
{
  "session_id": "660e8400-e29b-41d4-a716-446655440000",
  "chat_user_id": "770e8400-e29b-41d4-a716-446655440000",
  "status": "waiting",
  "message": "Chat session started successfully",
  "requires_contact": true
}
```

### 2. Set Session Contact
Mengisi informasi kontak untuk sesi chat.

**Endpoint:** `POST /api/chat/contact`

**Request Body:**
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

**Response:**
```json
{
  "contact_id": "880e8400-e29b-41d4-a716-446655440000",
  "message": "Contact information set successfully"
}
```

### 3. Link OSS User
Menghubungkan user anonymous dengan akun OSS saat login.

**Endpoint:** `POST /api/chat/link-user`

**Request Body:**
```json
{
  "browser_uuid": "550e8400-e29b-41d4-a716-446655440000", // Required
  "oss_user_id": "USER123",                              // Required
  "email": "user@example.com"                            // Required
}
```

**Response:**
```json
{
  "chat_user_id": "770e8400-e29b-41d4-a716-446655440000",
  "message": "Successfully linked to OSS account"
}
```

### 4. Get Chat History
Mengambil histori chat untuk pengguna tertentu.

**Endpoint:** `GET /api/chat/history`

**Query Parameters:**
- `browser_uuid` (string, optional): UUID browser untuk anonymous user
- `oss_user_id` (string, optional): ID pengguna OSS untuk logged-in user
- `limit` (int, optional): Jumlah sesi yang dikembalikan (default: 20)
- `offset` (int, optional): Jumlah sesi yang dilewati (default: 0)

**Example Request:**
```
GET /api/chat/history?oss_user_id=USER123&limit=10&offset=0
```

**Response:**
```json
{
  "sessions": [
    {
      "session_id": "660e8400-e29b-41d4-a716-446655440000",
      "topic": "Pertanyaan tentang izin usaha",
      "status": "closed",
      "priority": "normal",
      "started_at": "2025-01-15T10:00:00Z",
      "ended_at": "2025-01-15T10:30:00Z",
      "agent": {
        "id": "990e8400-e29b-41d4-a716-446655440000",
        "name": "Agent Smith",
        "email": "agent@example.com"
      },
      "department": {
        "id": "aa0e8400-e29b-41d4-a716-446655440000",
        "name": "Customer Support"
      },
      "contact": {
        "contact_name": "John Doe",
        "contact_email": "john@company.com",
        "contact_phone": "+6281234567890",
        "position": "Manager",
        "company_name": "PT. Example"
      },
      "last_message": {
        "id": "bb0e8400-e29b-41d4-a716-446655440000",
        "message": "Terima kasih atas bantuan Anda",
        "sender_type": "customer",
        "created_at": "2025-01-15T10:29:00Z"
      }
    }
  ],
  "total": 5,
  "limit": 10,
  "offset": 0
}
```

### 5. Get Session Details
Mengambil detail sesi chat tertentu (belum diimplementasi).

**Endpoint:** `GET /api/chat/session/{session_id}`

## User Flow Examples

### Flow 1: Anonymous User
1. User membuka website tanpa login
2. Frontend generate `browser_uuid` dan simpan di localStorage
3. User mulai chat dengan `browser_uuid`
4. User mengisi informasi kontak
5. Chat berlangsung
6. (Optional) User login dan link account dengan `browser_uuid`

### Flow 2: Logged-in User
1. User sudah login di sistem OSS
2. User mulai chat dengan `oss_user_id` dan `email`
3. User mengisi informasi kontak (nama kontak bisa berbeda dari user OSS)
4. Chat berlangsung

### Flow 3: Anonymous → Login
1. User mulai sebagai anonymous dengan `browser_uuid`
2. User mengisi kontak dan chat berlangsung
3. User login di sistem OSS
4. Sistem call endpoint link-user untuk menghubungkan session anonymous dengan akun OSS
5. User sekarang bisa akses histori chat dari akun OSS

## Error Responses

**400 Bad Request:**
```json
{
  "error": "Invalid request body",
  "details": "missing required field: topic"
}
```

**404 Not Found:**
```json
{
  "error": "Session not found"
}
```

**500 Internal Server Error:**
```json
{
  "error": "Failed to start chat",
  "details": "database connection error"
}
```

## Data Models

### ChatUser
```sql
chat_users (
  id UUID PRIMARY KEY,
  browser_uuid UUID UNIQUE,     -- For anonymous users
  oss_user_id VARCHAR(255),     -- For logged-in OSS users  
  email VARCHAR(255),           -- For logged-in users
  is_anonymous BOOLEAN DEFAULT true,
  ip_address VARCHAR(45) NOT NULL,
  user_agent TEXT,
  created_at TIMESTAMP,
  updated_at TIMESTAMP,
  deleted_at TIMESTAMP
)
```

### ChatSession
```sql
chat_sessions (
  id UUID PRIMARY KEY,
  chat_user_id UUID NOT NULL REFERENCES chat_users(id),
  agent_id UUID REFERENCES users(id),
  department_id UUID REFERENCES departments(id),
  topic VARCHAR(255) NOT NULL,
  status VARCHAR(50) DEFAULT 'waiting',
  priority VARCHAR(50) DEFAULT 'normal',
  started_at TIMESTAMP,
  ended_at TIMESTAMP,
  created_at TIMESTAMP,
  updated_at TIMESTAMP,
  deleted_at TIMESTAMP
)
```

### ChatSessionContact
```sql
chat_session_contacts (
  id UUID PRIMARY KEY,
  session_id UUID NOT NULL REFERENCES chat_sessions(id),
  contact_name VARCHAR(255) NOT NULL,
  contact_email VARCHAR(255) NOT NULL,
  contact_phone VARCHAR(50),
  position VARCHAR(255),
  company_name VARCHAR(255),
  created_at TIMESTAMP,
  updated_at TIMESTAMP,
  deleted_at TIMESTAMP
)
```

## Frontend Integration Tips

1. **Generate Browser UUID:**
```javascript
// Generate and store browser UUID for anonymous users
if (!localStorage.getItem('browser_uuid')) {
  localStorage.setItem('browser_uuid', crypto.randomUUID());
}
```

2. **Start Chat:**
```javascript
const startChat = async (topic, userData = {}) => {
  const payload = {
    topic,
    priority: 'normal',
    user_agent: navigator.userAgent
  };
  
  if (userData.isLoggedIn) {
    payload.oss_user_id = userData.ossUserId;
    payload.email = userData.email;
  } else {
    payload.browser_uuid = localStorage.getItem('browser_uuid');
  }
  
  const response = await fetch('/api/v1/oss-chat/start', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(payload)
  });
  
  return response.json();
};
```

3. **Link User on Login:**
```javascript
const linkUserOnLogin = async (ossUserId, email) => {
  const browserUuid = localStorage.getItem('browser_uuid');
  if (!browserUuid) return;
  
  await fetch('/api/v1/oss-chat/link-user', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({
      browser_uuid: browserUuid,
      oss_user_id: ossUserId,
      email: email
    })
  });
};
```
