package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/novianakbar/livechat-be/internal/domain"
	"github.com/novianakbar/livechat-be/internal/mappers"
)

// Updated version of GetSession method to show the difference

// GetSession (BEFORE - returns raw entity)
func (h *ChatHandler) GetSessionOld(c *fiber.Ctx) error {
	sessionIDStr := c.Params("session_id")
	if sessionIDStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(domain.ApiResponse{
			Success: false,
			Message: "Session ID is required",
			Error:   "invalid parameter",
		})
	}

	sessionID, err := uuid.Parse(sessionIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(domain.ApiResponse{
			Success: false,
			Message: "Invalid session ID format",
			Error:   err.Error(),
		})
	}

	session, err := h.chatUsecase.GetSession(c.Context(), sessionID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(domain.ApiResponse{
			Success: false,
			Message: "Failed to get session",
			Error:   err.Error(),
		})
	}

	if session == nil {
		return c.Status(fiber.StatusNotFound).JSON(domain.ApiResponse{
			Success: false,
			Message: "Session not found",
			Error:   "session does not exist",
		})
	}

	// ❌ MASALAH: Mengembalikan raw entity dengan sql.NullString, soft_delete.DeletedAt, etc.
	// Response akan terlihat seperti:
	// {
	//   "agent_id": {"String": "agent-123", "Valid": true},
	//   "ended_at": {"Time": "2024-01-01T10:00:00Z", "Valid": false},
	//   "deleted_at": {"Second": 0, "Valid": false}
	// }
	return c.JSON(domain.ApiResponse{
		Success: true,
		Message: "Session retrieved successfully",
		Data:    session, // Raw entity
	})
}

// GetSession (AFTER - returns clean response)
func (h *ChatHandler) GetSessionNew(c *fiber.Ctx) error {
	sessionIDStr := c.Params("session_id")
	if sessionIDStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(domain.ApiResponse{
			Success: false,
			Message: "Session ID is required",
			Error:   "invalid parameter",
		})
	}

	sessionID, err := uuid.Parse(sessionIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(domain.ApiResponse{
			Success: false,
			Message: "Invalid session ID format",
			Error:   err.Error(),
		})
	}

	session, err := h.chatUsecase.GetSession(c.Context(), sessionID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(domain.ApiResponse{
			Success: false,
			Message: "Failed to get session",
			Error:   err.Error(),
		})
	}

	if session == nil {
		return c.Status(fiber.StatusNotFound).JSON(domain.ApiResponse{
			Success: false,
			Message: "Session not found",
			Error:   "session does not exist",
		})
	}

	// ✅ SOLUSI: Convert entity ke clean response menggunakan mapper
	// Response akan terlihat seperti:
	// {
	//   "id": "session-123",
	//   "agent_id": "agent-123",           // Clean string, bukan object
	//   "ended_at": "2024-01-01T10:00:00Z", // Clean ISO string
	//   "status": "active",
	//   "chat_user": { ... },              // Nested objects juga clean
	//   "messages": [ ... ]                // Array messages clean
	// }
	response := mappers.ChatSessionToDetailResponse(session)

	return c.JSON(domain.ApiResponse{
		Success: true,
		Message: "Session retrieved successfully",
		Data:    response, // Clean response model
	})
}

/*
CONTOH PERBANDINGAN RESPONSE:

❌ SEBELUM (Raw Entity):
{
  "success": true,
  "data": {
    "id": "session-123",
    "agent_id": {
      "String": "agent-456",
      "Valid": true
    },
    "department_id": {
      "String": "",
      "Valid": false
    },
    "ended_at": {
      "Time": "0001-01-01T00:00:00Z",
      "Valid": false
    },
    "deleted_at": {
      "Second": 0,
      "Valid": false
    },
    "chat_user": {
      "browser_uuid": {
        "String": "browser-789",
        "Valid": true
      },
      "email": {
        "String": "",
        "Valid": false
      }
    }
  }
}

✅ SESUDAH (Clean Response):
{
  "success": true,
  "data": {
    "id": "session-123",
    "agent_id": "agent-456",
    "department_id": "",
    "ended_at": "",
    "status": "active",
    "priority": "normal",
    "topic": "Customer Support",
    "started_at": "2024-01-01T08:00:00Z",
    "created_at": "2024-01-01T08:00:00Z",
    "updated_at": "2024-01-01T09:00:00Z",
    "chat_user": {
      "id": "user-123",
      "browser_uuid": "browser-789",
      "email": "",
      "is_anonymous": true,
      "ip_address": "192.168.1.1",
      "created_at": "2024-01-01T08:00:00Z"
    },
    "messages": [
      {
        "id": "msg-1",
        "session_id": "session-123",
        "sender_type": "customer",
        "message": "Hello, I need help",
        "message_type": "text",
        "created_at": "2024-01-01T08:01:00Z"
      }
    ]
  }
}
*/
