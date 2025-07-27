# Agent Online Status Feature Documentation

## Overview
Fitur Agent Online Status memungkinkan sistem untuk melacak agent mana saja yang sedang online, status mereka, dan department mereka. Data disimpan di Redis untuk performa tinggi dan dapat digunakan untuk auto-assignment dan routing pesan di masa depan.

## Architecture

### Components
1. **AgentStatusRepository** - Mengelola data agent status di Redis
2. **AgentStatusService** - Business logic untuk mengelola agent status
3. **AgentStatusHandler** - HTTP handlers untuk API endpoints
4. **Redis Storage** - Menyimpan data status agent dengan TTL

### Data Structure in Redis

#### Individual Agent Status
**Key**: `agent:online:{agent_id}`
**TTL**: 5 minutes
**Value**: JSON object
```json
{
  "agent_id": "550e8400-e29b-41d4-a716-446655440000",
  "name": "John Doe",
  "email": "john@example.com",
  "department_id": "660e8400-e29b-41d4-a716-446655440000",
  "department": "Customer Support",
  "status": "online",
  "last_heartbeat": "2025-01-27T10:30:00Z"
}
```

#### All Agents Set
**Key**: `agents:all`
**TTL**: 10 minutes
**Value**: Redis Set containing agent IDs

#### Department Agents Set
**Key**: `agents:dept:{department_id}`
**TTL**: 10 minutes
**Value**: Redis Set containing agent IDs for specific department

## API Endpoints

### 1. Agent Heartbeat
**Endpoint**: `POST /api/agent-status/heartbeat`
**Purpose**: Agent mengirim heartbeat untuk menandai dirinya online
**Authentication**: Bearer token (Agent/Admin only)
**Frequency**: Recommended setiap 2-3 menit

**Request**:
```json
{
  "status": "online"  // Optional: "online", "busy", "away"
}
```

**Response**:
```json
{
  "success": true,
  "message": "Heartbeat updated successfully",
  "data": {
    "agent_id": "550e8400-e29b-41d4-a716-446655440000",
    "status": "online"
  }
}
```

### 2. Get Online Agents
**Endpoint**: `GET /api/agent-status/online`
**Purpose**: Mendapatkan daftar semua agent yang online
**Use case**: Dashboard admin, routing decisions

### 3. Get Agents by Department
**Endpoint**: `GET /api/agent-status/department/{department_id}`
**Purpose**: Mendapatkan agent online untuk department tertentu
**Use case**: Department-specific routing

### 4. Get Agent Status
**Endpoint**: `GET /api/agent-status/agent/{agent_id}`
**Purpose**: Cek status spesifik agent
**Use case**: Before assigning chat to agent

### 5. Department Statistics
**Endpoint**: `GET /api/agent-status/stats`
**Purpose**: Statistik agent online per department
**Use case**: Dashboard analytics

### 6. Set Agent Offline
**Endpoint**: `POST /api/agent-status/offline`
**Purpose**: Menandai agent sebagai offline (logout)

## Implementation Details

### Heartbeat Mechanism
- Agent frontend harus mengirim heartbeat setiap 2-3 menit
- Jika tidak ada heartbeat dalam 5 menit, agent dianggap offline
- Status agent otomatis expired dari Redis setelah TTL

### Status Types
- **online**: Agent tersedia untuk menerima chat
- **busy**: Agent sedang handling chat (optional untuk implementasi masa depan)
- **away**: Agent sementara tidak tersedia

### Error Handling
- Invalid agent ID: 400 Bad Request
- Agent not found: 404 Not Found
- Unauthorized access: 401 Unauthorized
- Redis connection error: 500 Internal Server Error

### Security
- Semua endpoint memerlukan authentication
- Agent data diambil dari JWT token, bukan request body
- Validasi role agent/admin untuk heartbeat

## Frontend Implementation Guide

### JavaScript Example for Heartbeat
```javascript
// Setup heartbeat interval
const HEARTBEAT_INTERVAL = 180000; // 3 minutes
let heartbeatInterval;

function startHeartbeat() {
  heartbeatInterval = setInterval(async () => {
    try {
      const response = await fetch('/api/agent-status/heartbeat', {
        method: 'POST',
        headers: {
          'Authorization': `Bearer ${authToken}`,
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({
          status: getCurrentAgentStatus() // "online", "busy", "away"
        })
      });
      
      if (!response.ok) {
        console.error('Heartbeat failed:', response.status);
        // Handle heartbeat failure (maybe redirect to login)
      }
    } catch (error) {
      console.error('Heartbeat error:', error);
    }
  }, HEARTBEAT_INTERVAL);
}

function stopHeartbeat() {
  if (heartbeatInterval) {
    clearInterval(heartbeatInterval);
    heartbeatInterval = null;
  }
}

// Start heartbeat when agent logs in
function onAgentLogin() {
  startHeartbeat();
}

// Stop heartbeat and mark offline when agent logs out
async function onAgentLogout() {
  stopHeartbeat();
  
  try {
    await fetch('/api/agent-status/offline', {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${authToken}`
      }
    });
  } catch (error) {
    console.error('Failed to set agent offline:', error);
  }
}

// Handle page visibility changes
document.addEventListener('visibilitychange', () => {
  if (document.hidden) {
    // Page is hidden, maybe pause heartbeat or set status to "away"
  } else {
    // Page is visible, resume normal heartbeat
  }
});
```

## Future Enhancements

### Auto-Assignment Algorithm
```javascript
// Example algorithm untuk auto-assignment berdasarkan agent online
async function autoAssignAgent(sessionData) {
  // 1. Get agents online for relevant department
  const agents = await getOnlineAgentsByDepartment(sessionData.department_id);
  
  // 2. Filter by workload (implement workload tracking)
  const availableAgents = agents.filter(agent => 
    agent.status === 'online' && agent.current_chats < MAX_CHATS_PER_AGENT
  );
  
  // 3. Round-robin or least-loaded assignment
  const selectedAgent = availableAgents.sort((a, b) => 
    a.current_chats - b.current_chats
  )[0];
  
  return selectedAgent;
}
```

### Workload Tracking
Tambahkan tracking jumlah chat aktif per agent untuk assignment yang lebih baik:
```json
{
  "agent_id": "...",
  "name": "...",
  "status": "online",
  "current_chats": 3,
  "max_chats": 5,
  "last_heartbeat": "..."
}
```

### Real-time Notifications
Integrate dengan WebSocket untuk real-time updates:
- Notify admin ketika agent online/offline
- Update agent list di dashboard secara real-time
- Alert jika tidak ada agent online untuk department tertentu

## Testing

### Manual Testing
1. Login sebagai agent
2. Kirim heartbeat via API
3. Verify agent muncul di list online agents
4. Wait for TTL expiry, verify agent hilang dari list
5. Test dengan multiple agents dan departments

### Load Testing
- Test dengan banyak agents mengirim heartbeat bersamaan
- Monitor Redis performance
- Test cleanup mechanism saat TTL expired

## Monitoring & Metrics

### Key Metrics to Track
- Number of online agents per department
- Heartbeat success/failure rate
- Average response time for status queries
- Redis memory usage for agent status data

### Alerting
- Alert jika tidak ada agent online untuk department critical
- Alert jika heartbeat failure rate tinggi
- Monitor Redis connection health

## Configuration

### Environment Variables
```bash
# Redis configuration for agent status
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0

# Agent status settings
AGENT_STATUS_TTL=300  # 5 minutes in seconds
HEARTBEAT_INTERVAL=180  # 3 minutes in seconds
```

## Troubleshooting

### Common Issues
1. **Agent not appearing online after heartbeat**
   - Check JWT token validity
   - Verify agent role in database
   - Check Redis connectivity

2. **Agent showing online but actually offline**
   - Normal behavior due to TTL mechanism
   - Agent will auto-expire after 5 minutes

3. **Performance issues with many agents**
   - Consider Redis clustering
   - Optimize Redis commands
   - Implement caching for frequently accessed data
