# Panduan Implementasi Response Models & Mappers

## Overview

Implementasi ini memecahkan masalah response API yang tidak clean akibat penggunaan GORM entities langsung. Dengan menggunakan Response Models + Mappers, kita mendapatkan response API yang bersih dan konsisten.

## Quick Start

### 1. Import yang diperlukan

```go
import (
    "github.com/novianakbar/livechat-be/internal/mappers"
    "github.com/novianakbar/livechat-be/internal/models"
)
```

### 2. Gunakan Mapper di Handler

```go
// Sebelum (❌ Raw Entity)
func (h *ChatHandler) GetSession(c *fiber.Ctx) error {
    session, err := h.chatUsecase.GetSession(ctx, sessionID)
    if err != nil {
        return handleError(c, err)
    }
    
    return c.JSON(domain.ApiResponse{
        Success: true,
        Data:    session, // Raw entity dengan sql.NullString
    })
}

// Sesudah (✅ Clean Response)
func (h *ChatHandler) GetSession(c *fiber.Ctx) error {
    session, err := h.chatUsecase.GetSession(ctx, sessionID)
    if err != nil {
        return handleError(c, err)
    }
    
    // Convert ke clean response
    response := mappers.ChatSessionToDetailResponse(session)
    
    return c.JSON(domain.ApiResponse{
        Success: true,
        Data:    response, // Clean response model
    })
}
```

## Available Mappers

### Single Entity Mappers
- `ChatSessionToDetailResponse(entity)` - Detail lengkap session
- `ChatSessionToMinimalResponse(entity)` - Minimal session untuk list
- `ChatUserToResponse(entity)` - Clean user data
- `UserToResponse(entity)` - Agent/admin user data
- `DepartmentToResponse(entity)` - Department data
- `ChatMessageToResponse(entity)` - Single message

### Slice Mappers
- `ChatSessionsToMinimalResponse(entities)` - Multiple sessions
- `ChatMessagesToResponse(entities)` - Multiple messages
- `UsersToResponse(entities)` - Multiple users
- `DepartmentsToResponse(entities)` - Multiple departments

### Pagination
- `CreatePaginatedResponse(data, page, limit, total)` - Wrap dengan pagination

## Response Models

### ChatSessionDetailResponse
Digunakan untuk endpoint detail session (`GET /sessions/{id}`):

```go
type ChatSessionDetailResponse struct {
    ID           string                          `json:"id"`
    ChatUserID   string                          `json:"chat_user_id"`
    AgentID      string                          `json:"agent_id,omitempty"`
    DepartmentID string                          `json:"department_id,omitempty"`
    Topic        string                          `json:"topic"`
    Status       string                          `json:"status"`
    Priority     string                          `json:"priority"`
    StartedAt    string                          `json:"started_at"`
    EndedAt      string                          `json:"ended_at,omitempty"`
    ChatUser     *ChatUserResponse               `json:"chat_user,omitempty"`
    Agent        *UserResponse                   `json:"agent,omitempty"`
    Department   *DepartmentResponse             `json:"department,omitempty"`
    Messages     []ChatMessageResponse           `json:"messages,omitempty"`
    Contact      *ChatSessionContactResponse     `json:"contact,omitempty"`
    CreatedAt    string                          `json:"created_at"`
    UpdatedAt    string                          `json:"updated_at"`
}
```

### ChatSessionMinimalResponse  
Digunakan untuk endpoint list sessions (`GET /sessions`):

```go
type ChatSessionMinimalResponse struct {
    ID         string            `json:"id"`
    ChatUserID string            `json:"chat_user_id"`
    AgentID    string            `json:"agent_id,omitempty"`
    Topic      string            `json:"topic"`
    Status     string            `json:"status"`
    Priority   string            `json:"priority"`
    StartedAt  string            `json:"started_at"`
    EndedAt    string            `json:"ended_at,omitempty"`
    ChatUser   *ChatUserResponse `json:"chat_user,omitempty"`
    Agent      *UserResponse     `json:"agent,omitempty"`
    CreatedAt  string            `json:"created_at"`
    UpdatedAt  string            `json:"updated_at"`
}
```

## Contoh Penggunaan

### 1. Single Session Detail

```go
func (h *ChatHandler) GetSession(c *fiber.Ctx) error {
    sessionID, _ := uuid.Parse(c.Params("session_id"))
    
    session, err := h.chatUsecase.GetSession(c.Context(), sessionID)
    if err != nil {
        return c.Status(500).JSON(domain.ApiResponse{
            Success: false,
            Error:   err.Error(),
        })
    }
    
    response := mappers.ChatSessionToDetailResponse(session)
    
    return c.JSON(domain.ApiResponse{
        Success: true,
        Data:    response,
    })
}
```

### 2. Sessions List dengan Pagination

```go
func (h *ChatHandler) GetSessions(c *fiber.Ctx) error {
    page := c.QueryInt("page", 1)
    limit := c.QueryInt("limit", 10)
    
    sessions, total, err := h.chatUsecase.GetSessions(c.Context(), page, limit, "", nil, nil)
    if err != nil {
        return c.Status(500).JSON(domain.ApiResponse{
            Success: false,
            Error:   err.Error(),
        })
    }
    
    sessionResponses := mappers.ChatSessionsToMinimalResponse(sessions)
    paginatedResponse := mappers.CreatePaginatedResponse(sessionResponses, page, limit, int64(total))
    
    return c.JSON(domain.ApiResponse{
        Success: true,
        Data:    paginatedResponse,
    })
}
```

### 3. Messages dalam Session

```go
func (h *ChatHandler) GetSessionMessages(c *fiber.Ctx) error {
    sessionID, _ := uuid.Parse(c.Params("session_id"))
    
    messages, err := h.chatUsecase.GetSessionMessages(c.Context(), sessionID)
    if err != nil {
        return c.Status(500).JSON(domain.ApiResponse{
            Success: false,
            Error:   err.Error(),
        })
    }
    
    messageResponses := mappers.ChatMessagesToResponse(messages)
    
    return c.JSON(domain.ApiResponse{
        Success: true,
        Data:    messageResponses,
    })
}
```

## Best Practices

### 1. Pilih Response Model yang Tepat
- **Detail Response**: Untuk endpoint single resource (`GET /sessions/{id}`)
- **Minimal Response**: Untuk endpoint list (`GET /sessions`)

### 2. Handle Null Values dengan Benar
```go
// Mappers sudah handle sql.NullString secara otomatis
// AgentID akan jadi string kosong jika null, bukan {"String": "", "Valid": false}
if entity.AgentID.Valid {
    response.AgentID = entity.AgentID.String
}
```

### 3. Consistent Time Format
Semua timestamp menggunakan RFC3339 format (`2024-01-01T08:00:00Z`)

### 4. Optional Fields
Gunakan `omitempty` tag untuk field yang optional:
```go
AgentID      string `json:"agent_id,omitempty"`
EndedAt      string `json:"ended_at,omitempty"`
```

### 5. Pagination
Selalu gunakan `PaginatedResponse` untuk endpoint list:
```go
paginatedResponse := mappers.CreatePaginatedResponse(data, page, limit, total)
```

## Error Handling

Mappers sudah handle case dimana entity adalah `nil`:

```go
func ChatSessionToDetailResponse(entity *entities.ChatSession) *models.ChatSessionDetailResponse {
    if entity == nil {
        return nil  // Safe handling
    }
    // ... mapping logic
}
```

## Migration Guide

### Step 1: Update Handler
```go
// Ganti ini:
Data: session,

// Dengan ini:
Data: mappers.ChatSessionToDetailResponse(session),
```

### Step 2: Update Tests
Update unit tests untuk expect clean response structure

### Step 3: Update Documentation  
Update Swagger/OpenAPI docs untuk menggunakan response models

## Performance Notes

- Mappers melakukan shallow copy, tidak expensive
- Pagination di-handle di level usecase, bukan mapper
- Time formatting menggunakan built-in Go time package yang efficient
