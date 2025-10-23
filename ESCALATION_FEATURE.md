# AI to Human Escalation Feature

## Overview
Fitur ini memungkinkan sistem AI untuk mengalihkan (escalate) percakapan ke agent manusia ketika AI mendeteksi bahwa pertanyaan atau situasi memerlukan penanganan manusia.

## How It Works

### 1. AI Detection
AI mendeteksi situasi yang memerlukan eskalasi berdasarkan:
- Kompleksitas pertanyaan yang tidak bisa dijawab AI
- Permintaan eksplisit dari customer untuk berbicara dengan agent
- Situasi yang memerlukan keputusan atau empati manusia
- Masalah teknis yang kompleks

### 2. Escalation Request
Ketika AI memutuskan untuk melakukan eskalasi, AI mengirim response dengan parameter:
```json
{
  "session_id": "uuid-session",
  "message": "Baik, saya akan menghubungkan Anda dengan agent kami...",
  "message_type": "text",
  "escalate_to_human": true,
  "escalation_reason": "Customer request for human agent"
}
```

### 3. Backend Processing
Backend (`ai_handler.go`) menerima request dan:
1. Memeriksa flag `escalate_to_human`
2. Memanggil `EscalateToHuman()` untuk:
   - Mencari agent yang tersedia
   - Jika ada agent: assign langsung dan ubah status ke `active`
   - Jika tidak ada agent: ubah status ke `queued` (antri)
3. Mengirim system message ke customer
4. Menghentikan AI typing indicator
5. Mem-broadcast perubahan status melalui Kafka/WebSocket

### 4. Session Status Flow
```
waiting (new chat) 
    ↓
active (AI responding)
    ↓
queued (escalated, waiting for agent) or active (escalated, agent assigned)
    ↓
active (agent handling)
    ↓
closed
```

## API Endpoint

### AI Response with Escalation
**POST** `/api/ai/response`

**Headers:**
```
Authorization: Bearer <AI_TOKEN>
Content-Type: application/json
```

**Request Body:**
```json
{
  "session_id": "550e8400-e29b-41d4-a716-446655440000",
  "message": "Saya akan menghubungkan Anda dengan agent kami yang dapat membantu lebih lanjut.",
  "message_type": "text",
  "escalate_to_human": true,
  "escalation_reason": "Complex technical issue requiring human expertise"
}
```

**Response (Success):**
```json
{
  "success": true,
  "message": "Chat successfully escalated to human agent",
  "data": {
    "escalated": true,
    "reason": "Complex technical issue requiring human expertise"
  }
}
```

## Database Changes

### Migration: 005_add_queued_status
Menambahkan status `queued` ke tabel `chat_sessions`:
- `waiting`: Menunggu response (initial state)
- `queued`: Di-eskalasi dari AI, menunggu agent tersedia
- `active`: Sedang ditangani oleh AI atau agent
- `closed`: Chat selesai

## New Methods

### ChatUsecase
```go
func (uc *ChatUsecase) EscalateToHuman(ctx context.Context, sessionID uuid.UUID, reason string) error
func (uc *ChatUsecase) GetQueuedSessions(ctx context.Context) ([]*domain.ChatSession, error)
```

### ChatSessionRepository
```go
func (r *chatSessionRepository) GetQueuedSessions(ctx context.Context) ([]*domain.ChatSession, error)
```

### ChatHandler
```go
func (h *ChatHandler) GetQueuedSessions(c *fiber.Ctx) error
```

## Admin Endpoints

### Get Queued Sessions
**GET** `/api/chat-management/admin/queued`

**Response:**
```json
{
  "success": true,
  "message": "Queued sessions retrieved successfully",
  "data": [
    {
      "id": "session-uuid",
      "chat_user": {...},
      "topic": "Bantuan Perizinan",
      "status": "queued",
      "priority": "normal",
      "started_at": "2025-10-23T10:00:00Z",
      "created_at": "2025-10-23T10:00:00Z"
    }
  ]
}
```

## WebSocket Notifications

### Session Status Update
Ketika terjadi eskalasi, WebSocket akan mengirim notifikasi:
```json
{
  "type": "session_status_update",
  "data": {
    "session_id": "uuid",
    "status": "queued",
    "message": "Your chat has been escalated to a human agent"
  }
}
```

### System Message
Customer akan menerima system message:
```json
{
  "type": "new_message",
  "data": {
    "message_id": "uuid",
    "session_id": "uuid",
    "sender_type": "system",
    "message": "Percakapan Anda sedang dialihkan ke agent kami. Mohon tunggu sebentar...",
    "message_type": "system",
    "timestamp": "2025-10-23T10:00:00Z"
  }
}
```

## Frontend Integration

### Handling Escalation on Widget
Widget harus:
1. Mendengarkan WebSocket message dengan type `new_message` dan `sender_type: "system"`
2. Menampilkan pesan sistem ke customer
3. Update UI status dari "AI" ke "Menunggu Agent" atau "Terhubung dengan Agent"
4. Menampilkan indikator jika chat sedang dalam antrian (status: queued)

### Example Widget Code
```typescript
// Listen for system messages
websocket.onmessage = (event) => {
  const data = JSON.parse(event.data);
  
  if (data.type === 'new_message' && data.data.sender_type === 'system') {
    // Show system message to user
    displaySystemMessage(data.data.message);
  }
  
  if (data.type === 'session_status_update') {
    // Update session status
    if (data.data.status === 'queued') {
      showQueueStatus('Menunggu agent tersedia...');
    } else if (data.data.status === 'active') {
      showAgentConnected();
    }
  }
};
```

## Monitoring & Analytics

Admin dashboard dapat menampilkan:
- Jumlah chat yang di-eskalasi per hari
- Alasan eskalasi paling umum
- Waktu tunggu rata-rata untuk queued sessions
- Agent response time setelah eskalasi

## Best Practices

1. **AI Should Provide Context**: Ketika melakukan eskalasi, AI sebaiknya mengirim pesan yang informatif kepada customer
2. **Graceful Handoff**: Transisi dari AI ke agent harus smooth dan transparan
3. **Queue Management**: Implementasikan notifikasi untuk agent ketika ada session dalam queue
4. **SLA Monitoring**: Set SLA untuk waktu maksimal session dalam status queued
5. **Fallback Mechanism**: Jika tidak ada agent dalam waktu tertentu, berikan opsi alternatif (email, callback, dll)

## Testing

### Test Escalation Flow
```bash
# 1. Start a chat session
POST /api/chat/start

# 2. Send AI response with escalation
POST /api/ai/response
{
  "session_id": "...",
  "message": "Let me connect you with an agent",
  "escalate_to_human": true,
  "escalation_reason": "Customer requested human assistance"
}

# 3. Check queued sessions
GET /api/chat-management/admin/queued

# 4. Assign agent (if needed)
POST /api/chat-management/admin/assign
{
  "session_id": "...",
  "agent_id": "..."
}
```

## Future Enhancements

1. **Priority Queue**: Prioritas berdasarkan urgency atau customer type
2. **Smart Routing**: Assign ke agent berdasarkan expertise
3. **Callback Option**: Jika wait time terlalu lama
4. **AI Context Sharing**: Transfer konteks percakapan AI ke agent
5. **Auto-escalation**: Eskalasi otomatis jika AI confidence rendah
