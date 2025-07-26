# LiveChat Backend - OSS Integration

[![Go Version](https://img.shields.io/badge/Go-1.21+-blue.svg)](https://golang.org)
[![Fiber](https://img.shields.io/badge/Fiber-v2-blue.svg)](https://gofiber.io)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-15+-blue.svg)](https://postgresql.org)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

Backend service untuk sistem LiveChat yang terintegrasi dengan OSS (Online Single Submission). Mendukung pengguna anonymous dan logged-in dengan transisi yang seamless.

## 🚀 Features

### Core Features
- ✅ **Multi-Mode Users**: Anonymous dan logged-in OSS users
- ✅ **Contact Management**: Informasi kontak terpisah per sesi
- ✅ **User Transition**: Anonymous → logged-in seamlessly
- ✅ **Real-time Messaging**: WebSocket support via Kafka
- ✅ **Session Management**: Multiple sessions per user
- ✅ **Chat History**: Pagination dan filtering
- ✅ **Agent Assignment**: Manual dan automatic assignment
- ✅ **Analytics**: Dashboard dan reporting
- ✅ **Email Notifications**: Automated email system

### OSS Integration Features
- 🆔 **Browser UUID Tracking**: Anonymous user identification
- 🔗 **Account Linking**: Connect anonymous sessions to OSS accounts
- 📧 **Email Integration**: OSS email-based user management
- 📊 **Session Analytics**: Per-user dan department analytics
- 🏢 **Multi-Contact**: One OSS account, multiple contacts per session

## 🏗️ Architecture

```
┌─────────────────┬─────────────────┬─────────────────┐
│   Presentation  │    Business     │       Data      │
│     Layer       │     Layer       │      Layer      │
├─────────────────┼─────────────────┼─────────────────┤
│ • HTTP Handlers │ • Use Cases     │ • Repositories  │
│ • Middleware    │ • Domain Logic  │ • Database      │
│ • Routes        │ • Validation    │ • External APIs │
│ • WebSocket     │ • Authorization │ • Kafka         │
└─────────────────┴─────────────────┴─────────────────┘
```

## 📋 Prerequisites

- **Go 1.21+**
- **PostgreSQL 15+**
- **Redis 6+** (optional, for caching)
- **Apache Kafka** (for real-time messaging)
- **Docker & Docker Compose** (for development)

## ⚡ Quick Start

### 1. Clone Repository
```bash
git clone https://github.com/your-org/livechat-be.git
cd livechat-be
```

### 2. Environment Setup
```bash
cp .env.example .env
# Edit .env with your configuration
```

### 3. Database Setup
```bash
# Start PostgreSQL
docker-compose up -d postgres

# Run migrations
make migrate-up
```

### 4. Start Services
```bash
# Development mode
make dev

# Production mode
make build && ./bin/livechat-backend
```

## 🔧 Configuration

### Environment Variables
```env
# Database
DB_HOST=localhost
DB_PORT=5432
DB_NAME=livechat_oss
DB_USER=postgres
DB_PASSWORD=your_password

# Server
PORT=8080
GIN_MODE=release

# JWT
JWT_SECRET=your_jwt_secret
JWT_EXPIRES_IN=24h

# Kafka
KAFKA_BROKERS=localhost:9092
KAFKA_TOPIC_MESSAGES=chat.messages

# Email (SendGrid)
SENDGRID_API_KEY=your_sendgrid_key
FROM_EMAIL=noreply@yourcompany.com

# CORS
CORS_ORIGINS=http://localhost:3000,https://yourapp.com
```

## 📡 API Endpoints

### OSS Chat API
```
POST   /api/oss-chat/start              # Start chat session
POST   /api/oss-chat/contact            # Set contact info
POST   /api/oss-chat/link-user          # Link anonymous to OSS
GET    /api/oss-chat/history            # Get chat history
GET    /api/oss-chat/session/:id        # Get session details
```

### Admin API
```
GET    /api/chat/waiting                # Get waiting sessions
GET    /api/chat/active                 # Get active sessions
POST   /api/chat/assign                 # Assign agent
POST   /api/chat/close                  # Close session
```

### Analytics API
```
GET    /api/analytics/dashboard         # Dashboard stats
GET    /api/analytics                   # Detailed analytics
```

For complete API documentation, see [API_Documentation.md](docs/API_Documentation.md)

## 🗄️ Database Schema

### Key Tables

#### chat_users (OSS Users)
```sql
CREATE TABLE chat_users (
    id UUID PRIMARY KEY,
    browser_uuid UUID UNIQUE,        -- Anonymous users
    oss_user_id VARCHAR(255),        -- OSS logged-in users  
    email VARCHAR(255),              -- User email
    is_anonymous BOOLEAN DEFAULT true,
    ip_address VARCHAR(45),
    user_agent TEXT,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP
);
```

#### chat_sessions
```sql
CREATE TABLE chat_sessions (
    id UUID PRIMARY KEY,
    chat_user_id UUID REFERENCES chat_users(id),
    agent_id UUID REFERENCES users(id),
    department_id UUID REFERENCES departments(id),
    topic VARCHAR(255),
    status VARCHAR(50) DEFAULT 'waiting',
    priority VARCHAR(50) DEFAULT 'normal',
    started_at TIMESTAMP,
    ended_at TIMESTAMP
);
```

#### chat_session_contacts
```sql
CREATE TABLE chat_session_contacts (
    id UUID PRIMARY KEY,
    session_id UUID REFERENCES chat_sessions(id),
    contact_name VARCHAR(255),
    contact_email VARCHAR(255),
    contact_phone VARCHAR(50),
    position VARCHAR(255),
    company_name VARCHAR(255)
);
```

## 🎯 Usage Examples

### Frontend Integration

#### Initialize Chat for Anonymous User
```javascript
// Generate browser UUID
const browserUUID = crypto.randomUUID();
localStorage.setItem('browser_uuid', browserUUID);

// Start chat
const response = await fetch('/api/oss-chat/start', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({
        browser_uuid: browserUUID,
        topic: 'Pertanyaan tentang izin usaha',
        priority: 'normal'
    })
});

const { session_id, requires_contact } = await response.json();
```

#### Set Contact Information
```javascript
if (requires_contact) {
    await fetch('/api/oss-chat/contact', {
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
}
```

#### Link to OSS Account (After Login)
```javascript
// After user logs in to OSS
await fetch('/api/oss-chat/link-user', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({
        browser_uuid: browserUUID,
        oss_user_id: user.ossId,
        email: user.email
    })
});
```

#### Get Chat History
```javascript
// For logged-in user
const history = await fetch('/api/oss-chat/history?oss_user_id=USER123&limit=20');

// For anonymous user
const history = await fetch(`/api/oss-chat/history?browser_uuid=${browserUUID}&limit=20`);
```

## 🔄 User Flows

### Flow 1: Anonymous User
1. User opens website (no login)
2. Frontend generates `browser_uuid`
3. User starts chat with `browser_uuid`
4. User fills contact information
5. Chat proceeds normally
6. (Optional) User logs in and links account

### Flow 2: Logged-in User
1. User already logged in to OSS
2. User starts chat with `oss_user_id` + `email`
3. User fills contact information
4. Chat proceeds normally

### Flow 3: Anonymous → Login Transition
1. User starts as anonymous
2. Chat session active
3. User logs in to OSS system
4. System calls link-user endpoint
5. Chat history now accessible from OSS account

## 🧪 Testing

### Run Tests
```bash
# Unit tests
make test

# Integration tests
make test-integration

# Coverage report
make test-coverage
```

### Test Data
```bash
# Seed test data
make seed-test-data

# Clean test data
make clean-test-data
```

## 📦 Deployment

### Docker Deployment
```bash
# Build image
docker build -t livechat-backend .

# Run with docker-compose
docker-compose up -d
```

### Kubernetes Deployment
```yaml
# See deployment/k8s/ directory
kubectl apply -f deployment/k8s/
```

### Production Checklist
- [ ] Environment variables configured
- [ ] Database migrations applied
- [ ] SSL certificates installed
- [ ] Monitoring configured
- [ ] Backup strategy implemented
- [ ] Load balancer configured

## 📊 Monitoring & Logging

### Health Checks
```bash
# Application health
curl http://localhost:8080/health

# Database health
curl http://localhost:8080/health/db
```

### Metrics Endpoints
- `/metrics` - Prometheus metrics
- `/health` - Health check
- `/debug/pprof` - Performance profiling

### Log Levels
```env
LOG_LEVEL=info  # debug, info, warn, error
LOG_FORMAT=json # json, text
```

## 🤝 Contributing

1. Fork the repository
2. Create feature branch (`git checkout -b feature/amazing-feature`)
3. Commit changes (`git commit -m 'Add amazing feature'`)
4. Push to branch (`git push origin feature/amazing-feature`)
5. Open Pull Request

### Development Guidelines
- Follow Go conventions and best practices
- Write tests for new features
- Update documentation
- Run `make lint` before committing

## 📄 Documentation

- [API Documentation](docs/API_Documentation.md)
- [OSS Chat API](docs/OSS_Chat_API.md)
- [Refactor Documentation](docs/REFACTOR_DOCUMENTATION.md)
- [Database Schema](docs/database_schema.md)
- [Deployment Guide](docs/deployment.md)

## 🔧 Troubleshooting

### Common Issues

#### Database Connection
```bash
# Check PostgreSQL status
docker-compose ps postgres

# View logs
docker-compose logs postgres
```

#### Migration Issues
```bash
# Check migration status
migrate -path ./migrations -database $DB_URL version

# Force migration version
migrate -path ./migrations -database $DB_URL force VERSION
```

#### Kafka Connection
```bash
# Check Kafka status
docker-compose ps kafka

# Test Kafka connection
kafka-console-consumer --bootstrap-server localhost:9092 --topic chat.messages
```

## 📞 Support

- **Email**: support@yourcompany.com
- **Slack**: #livechat-support
- **Documentation**: [Wiki](https://github.com/your-org/livechat-be/wiki)
- **Issues**: [GitHub Issues](https://github.com/your-org/livechat-be/issues)

## 📜 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- [Fiber](https://gofiber.io/) - Web framework
- [GORM](https://gorm.io/) - ORM library
- [golang-migrate](https://github.com/golang-migrate/migrate) - Database migrations
- [UUID](https://github.com/google/uuid) - UUID generation

---

**Made with ❤️ for OSS Integration**
