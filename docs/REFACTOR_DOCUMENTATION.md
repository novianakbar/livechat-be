# LiveChat OSS System - Dokumentasi Refaktor

## 📋 Overview

Proyek ini merupakan refaktor sistem LiveChat untuk mendukung integrasi dengan sistem OSS (Online Single Submission). Refaktor ini memungkinkan sistem untuk menangani pengguna anonymous dan pengguna yang sudah login dari sistem OSS, serta mendukung transisi dari anonymous ke logged-in user.

## 🎯 Tujuan Refaktor

1. **Mendukung Multiple User Types**: Anonymous users dan logged-in OSS users
2. **Flexible Contact Management**: Informasi kontak terpisah per sesi chat
3. **Seamless User Transition**: Anonymous user dapat di-link ke akun OSS
4. **Scalable Architecture**: Arsitektur yang dapat menangani volume tinggi
5. **Backward Compatibility**: Tetap mendukung sistem existing

## 🔄 Perubahan Arsitektur

### Database Schema Changes

#### 1. **chat_users Table** (menggantikan customers)
```sql
CREATE TABLE chat_users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    browser_uuid UUID UNIQUE,     -- For anonymous users
    oss_user_id VARCHAR(255),     -- For logged-in OSS users  
    email VARCHAR(255),           -- For logged-in users
    is_anonymous BOOLEAN DEFAULT true,
    ip_address VARCHAR(45) NOT NULL,
    user_agent TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    -- Constraints untuk memastikan data valid
    CHECK (
        (is_anonymous = true AND browser_uuid IS NOT NULL) OR
        (is_anonymous = false AND oss_user_id IS NOT NULL AND email IS NOT NULL)
    )
);
```

#### 2. **chat_session_contacts Table** (baru)
```sql
CREATE TABLE chat_session_contacts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    session_id UUID NOT NULL REFERENCES chat_sessions(id),
    contact_name VARCHAR(255) NOT NULL,
    contact_email VARCHAR(255) NOT NULL,
    contact_phone VARCHAR(50),
    position VARCHAR(255),        -- Job position
    company_name VARCHAR(255),    -- Company name
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    UNIQUE(session_id) -- One contact per session
);
```

#### 3. **chat_sessions Table** (update)
```sql
-- Menambahkan chat_user_id, menggantikan customer_id
ALTER TABLE chat_sessions ADD COLUMN chat_user_id UUID REFERENCES chat_users(id);
```

### 📁 Struktur File Baru

```
internal/
├── domain/
│   ├── entities.go              # Updated: ChatUser, ChatSessionContact
│   ├── dto.go                   # Updated: OSS-specific DTOs
│   └── repositories.go          # Updated: New repository interfaces
├── usecase/
│   └── oss_chat_usecase.go      # NEW: OSS chat business logic
├── infrastructure/repository/
│   ├── chat_user_repository.go  # NEW: ChatUser operations
│   └── chat_session_contact_repository.go # NEW: Contact operations
├── delivery/handler/
│   └── oss_chat_handler.go      # NEW: OSS chat API endpoints
└── delivery/routes/
    └── routes.go                # Updated: OSS routes
```

## 🚀 Fitur Utama

### 1. **Multi-Mode User Support**

#### Anonymous User Flow
```javascript
// Frontend generates browser UUID
const browserUUID = crypto.randomUUID();
localStorage.setItem('browser_uuid', browserUUID);

// Start chat as anonymous
const response = await fetch('/api/oss-chat/start', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({
        browser_uuid: browserUUID,
        topic: 'Pertanyaan tentang izin usaha',
        priority: 'normal'
    })
});
```

#### Logged-in User Flow
```javascript
// Start chat as logged-in user
const response = await fetch('/api/oss-chat/start', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({
        oss_user_id: 'USER123',
        email: 'user@example.com',
        topic: 'Pertanyaan tentang izin usaha',
        priority: 'normal'
    })
});
```

#### Anonymous → Login Transition
```javascript
// Link anonymous user to OSS account after login
const linkResponse = await fetch('/api/oss-chat/link-user', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({
        browser_uuid: browserUUID,
        oss_user_id: 'USER123',
        email: 'user@example.com'
    })
});
```

### 2. **Contact Management**

```javascript
// Set contact information for session
const contactResponse = await fetch('/api/oss-chat/contact', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({
        session_id: sessionId,
        contact_name: 'John Doe',
        contact_email: 'john@company.com',
        contact_phone: '+6281234567890',
        position: 'Manager',
        company_name: 'PT. Example'
    })
});
```

### 3. **Chat History**

```javascript
// Get chat history for user
const historyResponse = await fetch('/api/oss-chat/history?oss_user_id=USER123&limit=20&offset=0');
```

## 📊 API Endpoints

### OSS Chat Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/oss-chat/start` | Start new chat session |
| POST | `/api/oss-chat/contact` | Set contact information |
| POST | `/api/oss-chat/link-user` | Link anonymous to OSS user |
| GET | `/api/oss-chat/history` | Get chat history |
| GET | `/api/oss-chat/session/:id` | Get session details |

### Request/Response Examples

#### Start Chat Request
```json
{
  "browser_uuid": "550e8400-e29b-41d4-a716-446655440000", // Optional
  "oss_user_id": "USER123",                              // Optional
  "email": "user@example.com",                           // Optional
  "topic": "Pertanyaan tentang izin usaha",              // Required
  "priority": "normal",                                   // Optional
  "user_agent": "Mozilla/5.0 ..."                        // Optional
}
```

#### Start Chat Response
```json
{
  "session_id": "660e8400-e29b-41d4-a716-446655440000",
  "chat_user_id": "770e8400-e29b-41d4-a716-446655440000",
  "status": "waiting",
  "message": "Chat session started successfully",
  "requires_contact": true
}
```

## 🔧 Implementation Details

### Repository Layer

#### ChatUserRepository
```go
type ChatUserRepository interface {
    Create(ctx context.Context, user *ChatUser) error
    GetByBrowserUUID(ctx context.Context, browserUUID uuid.UUID) (*ChatUser, error)
    GetByOSSUserID(ctx context.Context, ossUserID string) (*ChatUser, error)
    LinkOSSUser(ctx context.Context, browserUUID uuid.UUID, ossUserID string, email string) error
    // ... other methods
}
```

#### ChatSessionContactRepository
```go
type ChatSessionContactRepository interface {
    Create(ctx context.Context, contact *ChatSessionContact) error
    GetBySessionID(ctx context.Context, sessionID uuid.UUID) (*ChatSessionContact, error)
    Update(ctx context.Context, contact *ChatSessionContact) error
    Delete(ctx context.Context, sessionID uuid.UUID) error
}
```

### Usecase Layer

#### OSSChatUsecase
```go
type OSSChatUsecase interface {
    StartChat(ctx context.Context, req *StartChatRequest, ipAddress string) (*StartChatResponse, error)
    SetSessionContact(ctx context.Context, req *SetSessionContactRequest) (*SetSessionContactResponse, error)
    LinkOSSUser(ctx context.Context, req *LinkOSSUserRequest) (*LinkOSSUserResponse, error)
    GetChatHistory(ctx context.Context, req *GetChatHistoryRequest) (*GetChatHistoryResponse, error)
}
```

## 🧪 Testing Strategy

### Unit Tests
```go
func TestOSSChatUsecase_StartChat_AnonymousUser(t *testing.T) {
    // Test starting chat as anonymous user
}

func TestOSSChatUsecase_StartChat_LoggedInUser(t *testing.T) {
    // Test starting chat as logged-in user
}

func TestOSSChatUsecase_LinkOSSUser(t *testing.T) {
    // Test linking anonymous user to OSS account
}
```

### Integration Tests
```go
func TestOSSChatAPI_CompleteFlow(t *testing.T) {
    // Test complete flow: anonymous → chat → login → link → history
}
```

## 🗃️ Database Migration

### Migration Files
- `001_initial_schema.up.sql` - Updated with new schema
- `003_add_oss_support.up.sql` - OSS-specific additions
- `003_add_oss_support.down.sql` - Rollback migrations

### Migration Commands
```bash
# Apply migrations
migrate -path ./migrations -database "postgres://..." up

# Rollback specific migration
migrate -path ./migrations -database "postgres://..." down 1
```

## 📈 Performance Considerations

### Database Indexes
```sql
-- Performance indexes
CREATE INDEX idx_chat_users_browser_uuid ON chat_users(browser_uuid);
CREATE INDEX idx_chat_users_oss_user_id ON chat_users(oss_user_id);
CREATE INDEX idx_chat_users_email ON chat_users(email);
CREATE INDEX idx_chat_sessions_chat_user_id ON chat_sessions(chat_user_id);
```

### Query Optimization
- Pagination untuk chat history
- Selective preloading untuk associations
- Proper indexing untuk filter queries

## 🔒 Security Considerations

### Data Validation
- Input validation pada semua endpoints
- UUID format validation
- Email format validation
- SQL injection prevention via parameterized queries

### Access Control
- IP address logging untuk tracking
- Rate limiting pada endpoints
- CORS configuration untuk cross-origin requests

## 🚦 Deployment Guide

### Environment Variables
```env
DB_HOST=localhost
DB_PORT=5432
DB_NAME=livechat_oss
DB_USER=postgres
DB_PASSWORD=password
KAFKA_BROKERS=localhost:9092
```

### Docker Deployment
```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o livechat-backend ./cmd/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/livechat-backend .
CMD ["./livechat-backend"]
```

## 📚 Best Practices

### Code Organization
- Repository pattern untuk data access
- Usecase pattern untuk business logic
- Interface-based dependency injection
- Clean architecture principles

### Error Handling
```go
// Consistent error responses
return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
    "error":   "Invalid request",
    "details": err.Error(),
})
```

### Logging
```go
log.Printf("Starting chat for user: %s, session: %s", 
    userID, sessionID)
```

## 🔮 Future Enhancements

### Planned Features
1. **WebSocket Integration** - Real-time messaging
2. **File Upload Support** - Document sharing
3. **Chat Analytics** - Advanced reporting
4. **Multi-language Support** - Internationalization
5. **Bot Integration** - Automated responses

### Scalability Improvements
1. **Redis Caching** - Session and user caching
2. **Horizontal Scaling** - Load balancer support
3. **Message Queue** - Async message processing
4. **CDN Integration** - Static asset delivery

## 📞 Support & Maintenance

### Monitoring
- Database performance monitoring
- API response time tracking
- Error rate monitoring
- User activity analytics

### Troubleshooting
- Comprehensive logging
- Error tracking
- Performance profiling
- Database query analysis

---

## 📝 Conclusion

Refaktor sistem LiveChat OSS ini berhasil mencapai tujuan utama:

✅ **Mendukung multi-mode users** (anonymous & logged-in)  
✅ **Flexible contact management** per session  
✅ **Seamless user transition** anonymous → logged-in  
✅ **Scalable architecture** dengan proper separation of concerns  
✅ **Comprehensive API** dengan dokumentasi lengkap  
✅ **Database optimization** dengan indexes dan constraints  
✅ **Security best practices** dengan validation dan access control  

Sistem ini siap untuk production dan dapat menangani kebutuhan OSS yang kompleks dengan volume tinggi.
