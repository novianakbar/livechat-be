# Response Models & Mappers Implementation

## Problem Statement

Ketika menggunakan GORM entities langsung sebagai response API, kita menghadapi masalah:

```json
// Response yang tidak clean - sebelum
{
  "success": true,
  "data": {
    "id": "session-123",
    "agent_id": {
      "String": "agent-456",
      "Valid": true
    },
    "ended_at": {
      "Time": "2024-01-01T10:00:00Z",
      "Valid": true
    },
    "deleted_at": {
      "Second": 0,
      "Valid": false
    }
  }
}
```

## Solution

Dengan menggunakan Response Models + Mappers:

```json
// Response yang clean - sesudah
{
  "success": true,
  "data": {
    "id": "session-123",
    "agent_id": "agent-456",
    "ended_at": "2024-01-01T10:00:00Z"
  }
}
```

## Architecture

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Controller    │────│    Usecase      │────│   Repository    │────│     Entity      │
│   (Handler)     │    │                 │    │                 │    │ (livechat-shared)│
└─────────────────┘    └─────────────────┘    └─────────────────┘    └─────────────────┘
         │                                                                      │
         │                                                                      │
         ▼                                                                      │
┌─────────────────┐    ┌─────────────────┐    ◄────────────────────────────────┘
│    Mapper       │────│  Response Model │
│                 │    │                 │
└─────────────────┘    └─────────────────┘
```

## File Structure

```
internal/
├── models/
│   └── response.go              # Clean response structures
├── mappers/
│   └── response_mapper.go       # Entity-to-Response converters
└── delivery/handler/
    ├── chat_handler.go          # Current handlers (to be updated)
    └── example_chat_handler.go  # Example of new implementation
```

## Key Components

### 1. Response Models (`internal/models/response.go`)

Clean structures for API responses:

- `ChatSessionMinimalResponse` - For list endpoints
- `ChatSessionDetailResponse` - For detail endpoints  
- `ChatUserResponse` - Clean user data
- `PaginatedResponse[T]` - Generic pagination wrapper

### 2. Mappers (`internal/mappers/response_mapper.go`)

Conversion functions:

- `ChatSessionToDetailResponse()` - Convert entity to detail response
- `ChatSessionsToMinimalResponse()` - Convert entity slice to minimal responses
- `CreatePaginatedResponse()` - Create paginated wrapper

### 3. Handler Usage

Before:
```go
return c.JSON(domain.ApiResponse{
    Success: true,
    Data:    session, // Raw entity with sql.NullString, etc.
})
```

After:
```go
response := mappers.ChatSessionToDetailResponse(session)
return c.JSON(domain.ApiResponse{
    Success: true,
    Data:    response, // Clean response model
})
```

## Benefits

1. **Clean API**: No more `{"String": "value", "Valid": true}`
2. **Consistent**: All timestamps in RFC3339 format
3. **Type Safe**: Strong typing for all response fields
4. **Maintainable**: Easy to modify response structure without touching entities
5. **Documentation**: Clear response models for API documentation

## Implementation Steps

1. ✅ Create response models in `internal/models/response.go`
2. ✅ Create mappers in `internal/mappers/response_mapper.go` 
3. ✅ Create example handler showing usage
4. 🔄 Update existing handlers to use mappers
5. 🔄 Update API documentation

## Next Steps

1. **Update Existing Handlers**: Replace direct entity returns with mapper usage
2. **Add More Response Models**: For other entities like User, Department, etc.
3. **Add Request Validation Models**: For input validation
4. **Update API Documentation**: Use new response models in Swagger docs

## Example Usage

See `internal/delivery/handler/example_chat_handler.go` for complete examples of:

- Single session retrieval with detailed response
- Session list with pagination and minimal response
- Error handling with clean responses
