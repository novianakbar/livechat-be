# LiveChat API Documentation

## Base URL
```
http://localhost:8080
```

## Authentication
All protected endpoints require JWT token in Authorization header:
```
Authorization: Bearer <your-jwt-token>
```

## Error Response Format
```json
{
  "success": false,
  "message": "Error message",
  "error": "Detailed error information"
}
```

## Success Response Format
```json
{
  "success": true,
  "message": "Success message",
  "data": {...}
}
```

---

## Authentication Endpoints

### POST /api/auth/login
Login user and get JWT token.

**Request Body:**
```json
{
  "email": "admin@livechat.com",
  "password": "password"
}
```

**Response:**
```json
{
  "success": true,
  "message": "Login successful",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": "550e8400-e29b-41d4-a716-446655440010",
      "email": "admin@livechat.com",
      "name": "Administrator",
      "role": "admin",
      "is_active": true,
      "department": null
    },
    "expires_at": "2024-01-02T12:00:00Z"
  }
}
```

### GET /api/auth/profile
Get current user profile.

**Headers:**
```
Authorization: Bearer <jwt-token>
```

**Response:**
```json
{
  "success": true,
  "message": "Profile retrieved successfully",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440010",
    "email": "admin@livechat.com",
    "name": "Administrator",
    "role": "admin",
    "is_active": true,
    "department": null
  }
}
```

### POST /api/auth/register
Register new user (Admin only).

**Headers:**
```
Authorization: Bearer <jwt-token>
```

**Request Body:**
```json
{
  "email": "newagent@livechat.com",
  "password": "password123",
  "name": "New Agent",
  "role": "agent",
  "department_id": "550e8400-e29b-41d4-a716-446655440001"
}
```

**Response:**
```json
{
  "success": true,
  "message": "User registered successfully",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440020",
    "email": "newagent@livechat.com",
    "name": "New Agent",
    "role": "agent",
    "is_active": true,
    "department_id": "550e8400-e29b-41d4-a716-446655440001"
  }
}
```

---

## Public Chat Endpoints

### POST /api/public/chat/start
Start new chat session (public endpoint).

**Request Body:**
```json
{
  "company_name": "PT Maju Terus",
  "person_name": "Budi Santoso",
  "email": "budi@majuterus.com",
  "topic": "Perpanjangan izin usaha UMKM",
  "priority": "normal"
}
```

**Response:**
```json
{
  "success": true,
  "message": "Chat session started successfully",
  "data": {
    "session_id": "550e8400-e29b-41d4-a716-446655440030",
    "status": "waiting",
    "message": "Chat session started. Please wait for an agent to respond."
  }
}
```

### POST /api/public/chat/message
Send message in chat session (public endpoint).

**Request Body:**
```json
{
  "session_id": "550e8400-e29b-41d4-a716-446655440030",
  "message": "Saya ingin mengurus perpanjangan izin usaha UMKM saya. Dokumen apa saja yang diperlukan?",
  "message_type": "text",
  "attachments": []
}
```

**Response:**
```json
{
  "success": true,
  "message": "Message sent successfully",
  "data": {
    "message_id": "550e8400-e29b-41d4-a716-446655440040",
    "timestamp": "2024-01-01T10:00:00Z",
    "status": "sent"
  }
}
```

### GET /api/public/chat/session/{session_id}/messages
Get all messages in a chat session.

**Response:**
```json
{
  "success": true,
  "message": "Messages retrieved successfully",
  "data": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440040",
      "session_id": "550e8400-e29b-41d4-a716-446655440030",
      "sender_id": null,
      "sender_type": "customer",
      "message": "Saya ingin mengurus perpanjangan izin usaha UMKM saya.",
      "message_type": "text",
      "attachments": [],
      "read_at": null,
      "created_at": "2024-01-01T10:00:00Z"
    }
  ]
}
```

---

## Agent Chat Endpoints

### GET /api/chat/agent/sessions
Get all chat sessions assigned to current agent.

**Headers:**
```
Authorization: Bearer <jwt-token>
```

**Response:**
```json
{
  "success": true,
  "message": "Agent sessions retrieved successfully",
  "data": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440030",
      "customer": {
        "id": "550e8400-e29b-41d4-a716-446655440050",
        "company_name": "PT Maju Terus",
        "person_name": "Budi Santoso",
        "email": "budi@majuterus.com",
        "ip_address": "192.168.1.100"
      },
      "agent": {
        "id": "550e8400-e29b-41d4-a716-446655440011",
        "name": "Agent Perizinan 1",
        "email": "agent1@livechat.com"
      },
      "topic": "Perpanjangan izin usaha UMKM",
      "status": "active",
      "priority": "normal",
      "started_at": "2024-01-01T10:00:00Z",
      "ended_at": null
    }
  ]
}
```

### POST /api/chat/agent/message
Send message as agent.

**Headers:**
```
Authorization: Bearer <jwt-token>
```

**Request Body:**
```json
{
  "session_id": "550e8400-e29b-41d4-a716-446655440030",
  "message": "Untuk perpanjangan izin usaha UMKM, diperlukan dokumen: 1) Fotokopi KTP, 2) Fotokopi NPWP, 3) Surat keterangan domisili usaha, 4) Laporan kegiatan usaha.",
  "message_type": "text",
  "attachments": []
}
```

**Response:**
```json
{
  "success": true,
  "message": "Message sent successfully",
  "data": {
    "message_id": "550e8400-e29b-41d4-a716-446655440041",
    "timestamp": "2024-01-01T10:05:00Z",
    "status": "sent"
  }
}
```

### POST /api/chat/agent/assign
Assign agent to chat session.

**Headers:**
```
Authorization: Bearer <jwt-token>
```

**Request Body:**
```json
{
  "session_id": "550e8400-e29b-41d4-a716-446655440030",
  "agent_id": "550e8400-e29b-41d4-a716-446655440011"
}
```

**Response:**
```json
{
  "success": true,
  "message": "Agent assigned successfully"
}
```

### POST /api/chat/agent/close
Close chat session.

**Headers:**
```
Authorization: Bearer <jwt-token>
```

**Request Body:**
```json
{
  "session_id": "550e8400-e29b-41d4-a716-446655440030",
  "reason": "Pertanyaan customer sudah terjawab dengan lengkap",
  "rating": 5,
  "feedback": "Customer puas dengan layanan yang diberikan"
}
```

**Response:**
```json
{
  "success": true,
  "message": "Session closed successfully"
}
```

---

## Admin Chat Endpoints

### GET /api/chat/admin/waiting
Get all waiting chat sessions (Admin only).

**Headers:**
```
Authorization: Bearer <jwt-token>
```

**Response:**
```json
{
  "success": true,
  "message": "Waiting sessions retrieved successfully",
  "data": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440030",
      "customer": {
        "company_name": "PT Maju Terus",
        "person_name": "Budi Santoso",
        "email": "budi@majuterus.com"
      },
      "topic": "Perpanjangan izin usaha UMKM",
      "status": "waiting",
      "priority": "normal",
      "started_at": "2024-01-01T10:00:00Z"
    }
  ]
}
```

### GET /api/chat/admin/active
Get all active chat sessions (Admin only).

**Headers:**
```
Authorization: Bearer <jwt-token>
```

**Response:**
```json
{
  "success": true,
  "message": "Active sessions retrieved successfully",
  "data": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440030",
      "customer": {
        "company_name": "PT Maju Terus",
        "person_name": "Budi Santoso",
        "email": "budi@majuterus.com"
      },
      "agent": {
        "id": "550e8400-e29b-41d4-a716-446655440011",
        "name": "Agent Perizinan 1",
        "email": "agent1@livechat.com"
      },
      "topic": "Perpanjangan izin usaha UMKM",
      "status": "active",
      "priority": "normal",
      "started_at": "2024-01-01T10:00:00Z"
    }
  ]
}
```

---

## WebSocket Connection

### Connection URL
```
ws://localhost:8080/ws/chat
```

### Message Format
All WebSocket messages follow this format:
```json
{
  "type": "message_type",
  "session_id": "session-uuid",
  "data": {...},
  "timestamp": "2024-01-01T10:00:00Z"
}
```

### Message Types

#### Join Session
```json
{
  "type": "join_session",
  "session_id": "550e8400-e29b-41d4-a716-446655440030",
  "data": {}
}
```

#### Send Message
```json
{
  "type": "send_message",
  "session_id": "550e8400-e29b-41d4-a716-446655440030",
  "data": {
    "message": "Hello, I need help with my business permit"
  }
}
```

#### Typing Indicators
```json
{
  "type": "agent_typing",
  "session_id": "550e8400-e29b-41d4-a716-446655440030",
  "data": {}
}
```

```json
{
  "type": "customer_typing",
  "session_id": "550e8400-e29b-41d4-a716-446655440030",
  "data": {}
}
```

### WebSocket Responses

#### Message Sent
```json
{
  "type": "message_sent",
  "success": true,
  "data": {
    "session_id": "550e8400-e29b-41d4-a716-446655440030",
    "message": "Hello, I need help with my business permit",
    "sender": "customer",
    "timestamp": "2024-01-01T10:00:00Z"
  }
}
```

#### New Message
```json
{
  "type": "new_message",
  "success": true,
  "data": {
    "session_id": "550e8400-e29b-41d4-a716-446655440030",
    "message_id": "550e8400-e29b-41d4-a716-446655440040",
    "message": "I can help you with that. What type of permit do you need?",
    "sender_type": "agent",
    "message_type": "text",
    "timestamp": "2024-01-01T10:05:00Z"
  }
}
```

#### Session Update
```json
{
  "type": "session_update",
  "success": true,
  "data": {
    "session_id": "550e8400-e29b-41d4-a716-446655440030",
    "status": "active"
  }
}
```

#### Error
```json
{
  "type": "error",
  "success": false,
  "data": null,
  "error": "Session not found"
}
```

---

## Status Codes

- `200` - Success
- `201` - Created
- `400` - Bad Request
- `401` - Unauthorized
- `403` - Forbidden
- `404` - Not Found
- `409` - Conflict
- `500` - Internal Server Error

---

## Rate Limiting

Currently no rate limiting is implemented, but it's recommended to add it for production use.

---

## WebSocket Client Example

```javascript
const ws = new WebSocket('ws://localhost:8080/ws/chat');

ws.onopen = function() {
    console.log('Connected to WebSocket');
    
    // Join a session
    ws.send(JSON.stringify({
        type: 'join_session',
        session_id: 'your-session-id',
        data: {}
    }));
};

ws.onmessage = function(event) {
    const response = JSON.parse(event.data);
    console.log('Received:', response);
    
    if (response.type === 'new_message') {
        // Handle new message
        console.log('New message:', response.data.message);
    }
};

ws.onclose = function() {
    console.log('WebSocket connection closed');
};

ws.onerror = function(error) {
    console.log('WebSocket error:', error);
};

// Send a message
function sendMessage(message) {
    ws.send(JSON.stringify({
        type: 'send_message',
        session_id: 'your-session-id',
        data: {
            message: message
        }
    }));
}
```
