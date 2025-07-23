package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/novianakbar/livechat-be/internal/delivery/middleware"
	"github.com/novianakbar/livechat-be/internal/domain"
	"github.com/novianakbar/livechat-be/internal/usecase"
)

type ChatHandler struct {
	chatUsecase *usecase.ChatUsecase
	wsHandler   *WebSocketHandler
}

func NewChatHandler(chatUsecase *usecase.ChatUsecase, wsHandler *WebSocketHandler) *ChatHandler {
	return &ChatHandler{
		chatUsecase: chatUsecase,
		wsHandler:   wsHandler,
	}
}

// StartChat godoc
// @Summary Start new chat session
// @Description Start a new chat session for customer
// @Tags Chat
// @Accept json
// @Produce json
// @Param request body domain.StartChatRequest true "Start chat request"
// @Success 201 {object} domain.ApiResponse{data=domain.StartChatResponse}
// @Failure 400 {object} domain.ApiResponse
// @Router /api/chat/start [post]
func (h *ChatHandler) StartChat(c *fiber.Ctx) error {
	var req domain.StartChatRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(domain.ApiResponse{
			Success: false,
			Message: "Invalid request body",
			Error:   err.Error(),
		})
	}

	// Validate request
	if req.CompanyName == "" || req.PersonName == "" || req.Email == "" || req.Topic == "" {
		return c.Status(fiber.StatusBadRequest).JSON(domain.ApiResponse{
			Success: false,
			Message: "Company name, person name, email, and topic are required",
			Error:   "validation failed",
		})
	}

	// Get client IP
	ipAddress := c.IP()

	response, err := h.chatUsecase.StartChat(c.Context(), &req, ipAddress)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(domain.ApiResponse{
			Success: false,
			Message: "Failed to start chat",
			Error:   err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(domain.ApiResponse{
		Success: true,
		Message: "Chat session started successfully",
		Data:    response,
	})
}

// SendMessage godoc
// @Summary Send message in chat
// @Description Send a message in an existing chat session
// @Tags Chat
// @Accept json
// @Produce json
// @Param request body domain.SendMessageRequest true "Send message request"
// @Success 200 {object} domain.ApiResponse{data=domain.SendMessageResponse}
// @Failure 400 {object} domain.ApiResponse
// @Router /api/chat/message [post]
func (h *ChatHandler) SendMessage(c *fiber.Ctx) error {
	var req domain.SendMessageRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(domain.ApiResponse{
			Success: false,
			Message: "Invalid request body",
			Error:   err.Error(),
		})
	}

	// Validate request
	if req.SessionID == uuid.Nil || req.Message == "" {
		return c.Status(fiber.StatusBadRequest).JSON(domain.ApiResponse{
			Success: false,
			Message: "Session ID and message are required",
			Error:   "validation failed",
		})
	}

	// Determine sender type - if user is in context, it's an agent, otherwise customer
	var senderID *uuid.UUID
	senderType := "customer"

	if user := middleware.GetUserFromContext(c); user != nil {
		senderID = &user.ID
		senderType = "agent"
	}

	response, err := h.chatUsecase.SendMessage(c.Context(), &req, senderID, senderType)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(domain.ApiResponse{
			Success: false,
			Message: "Failed to send message",
			Error:   err.Error(),
		})
	}

	// Get the message from the database to broadcast via WebSocket
	message, err := h.chatUsecase.GetMessageByID(c.Context(), response.MessageID)
	if err == nil && message != nil && h.wsHandler != nil {
		// Broadcast the message to all clients in the session
		h.wsHandler.BroadcastMessage(req.SessionID, message)
	}

	return c.JSON(domain.ApiResponse{
		Success: true,
		Message: "Message sent successfully",
		Data:    response,
	})
}

// AssignAgent godoc
// @Summary Assign agent to chat session
// @Description Assign an agent to a waiting chat session
// @Tags Chat
// @Accept json
// @Produce json
// @Param request body domain.AssignAgentRequest true "Assign agent request"
// @Success 200 {object} domain.ApiResponse
// @Failure 400 {object} domain.ApiResponse
// @Security BearerAuth
// @Router /api/chat/assign [post]
func (h *ChatHandler) AssignAgent(c *fiber.Ctx) error {
	var req domain.AssignAgentRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(domain.ApiResponse{
			Success: false,
			Message: "Invalid request body",
			Error:   err.Error(),
		})
	}

	// Validate request
	if req.SessionID == uuid.Nil || req.AgentID == uuid.Nil {
		return c.Status(fiber.StatusBadRequest).JSON(domain.ApiResponse{
			Success: false,
			Message: "Session ID and agent ID are required",
			Error:   "validation failed",
		})
	}

	err := h.chatUsecase.AssignAgent(c.Context(), &req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(domain.ApiResponse{
			Success: false,
			Message: "Failed to assign agent",
			Error:   err.Error(),
		})
	}

	return c.JSON(domain.ApiResponse{
		Success: true,
		Message: "Agent assigned successfully",
	})
}

// CloseSession godoc
// @Summary Close chat session
// @Description Close an active chat session
// @Tags Chat
// @Accept json
// @Produce json
// @Param request body domain.CloseSessionRequest true "Close session request"
// @Success 200 {object} domain.ApiResponse
// @Failure 400 {object} domain.ApiResponse
// @Security BearerAuth
// @Router /api/chat/close [post]
func (h *ChatHandler) CloseSession(c *fiber.Ctx) error {
	var req domain.CloseSessionRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(domain.ApiResponse{
			Success: false,
			Message: "Invalid request body",
			Error:   err.Error(),
		})
	}

	// Validate request
	if req.SessionID == uuid.Nil {
		return c.Status(fiber.StatusBadRequest).JSON(domain.ApiResponse{
			Success: false,
			Message: "Session ID is required",
			Error:   "validation failed",
		})
	}

	userID := middleware.GetUserIDFromContext(c)
	err := h.chatUsecase.CloseSession(c.Context(), req.SessionID, req.Reason, userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(domain.ApiResponse{
			Success: false,
			Message: "Failed to close session",
			Error:   err.Error(),
		})
	}

	return c.JSON(domain.ApiResponse{
		Success: true,
		Message: "Session closed successfully",
	})
}

// GetSessionMessages godoc
// @Summary Get chat session messages
// @Description Get all messages for a chat session
// @Tags Chat
// @Produce json
// @Param session_id path string true "Session ID"
// @Success 200 {object} domain.ApiResponse{data=[]domain.ChatMessage}
// @Failure 400 {object} domain.ApiResponse
// @Router /api/chat/session/{session_id}/messages [get]
func (h *ChatHandler) GetSessionMessages(c *fiber.Ctx) error {
	sessionIDStr := c.Params("session_id")
	sessionID, err := uuid.Parse(sessionIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(domain.ApiResponse{
			Success: false,
			Message: "Invalid session ID",
			Error:   err.Error(),
		})
	}

	messages, err := h.chatUsecase.GetSessionMessages(c.Context(), sessionID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(domain.ApiResponse{
			Success: false,
			Message: "Failed to get messages",
			Error:   err.Error(),
		})
	}

	return c.JSON(domain.ApiResponse{
		Success: true,
		Message: "Messages retrieved successfully",
		Data:    messages,
	})
}

// GetWaitingSessions godoc
// @Summary Get waiting chat sessions
// @Description Get all chat sessions waiting for agent assignment
// @Tags Chat
// @Produce json
// @Success 200 {object} domain.ApiResponse{data=[]domain.ChatSession}
// @Failure 500 {object} domain.ApiResponse
// @Security BearerAuth
// @Router /api/chat/waiting [get]
func (h *ChatHandler) GetWaitingSessions(c *fiber.Ctx) error {
	sessions, err := h.chatUsecase.GetWaitingSessions(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(domain.ApiResponse{
			Success: false,
			Message: "Failed to get waiting sessions",
			Error:   err.Error(),
		})
	}

	return c.JSON(domain.ApiResponse{
		Success: true,
		Message: "Waiting sessions retrieved successfully",
		Data:    sessions,
	})
}

// GetActiveSessions godoc
// @Summary Get active chat sessions
// @Description Get all active chat sessions
// @Tags Chat
// @Produce json
// @Success 200 {object} domain.ApiResponse{data=[]domain.ChatSession}
// @Failure 500 {object} domain.ApiResponse
// @Security BearerAuth
// @Router /api/chat/active [get]
func (h *ChatHandler) GetActiveSessions(c *fiber.Ctx) error {
	sessions, err := h.chatUsecase.GetActiveSessions(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(domain.ApiResponse{
			Success: false,
			Message: "Failed to get active sessions",
			Error:   err.Error(),
		})
	}

	return c.JSON(domain.ApiResponse{
		Success: true,
		Message: "Active sessions retrieved successfully",
		Data:    sessions,
	})
}

// GetAgentSessions godoc
// @Summary Get agent's chat sessions
// @Description Get chat sessions assigned to current agent with pagination
// @Tags Chat
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param limit query int false "Items per page"
// @Param status query string false "Session status filter"
// @Success 200 {object} domain.PaginatedResponse{data=[]domain.ChatSession}
// @Failure 500 {object} domain.ApiResponse
// @Security BearerAuth
// @Router /api/chat/agent/sessions [get]
func (h *ChatHandler) GetAgentSessions(c *fiber.Ctx) error {
	user := middleware.GetUserFromContext(c)
	if user == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(domain.ApiResponse{
			Success: false,
			Message: "User not found in context",
			Error:   "authentication required",
		})
	}

	// Parse query parameters
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 20)
	status := c.Query("status")

	sessions, total, err := h.chatUsecase.GetAgentSessionsWithPagination(c.Context(), user.ID, page, limit, status)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(domain.ApiResponse{
			Success: false,
			Message: "Failed to get agent sessions",
			Error:   err.Error(),
		})
	}

	totalPages := (total + limit - 1) / limit
	pagination := domain.PaginationInfo{
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
	}

	return c.JSON(domain.PaginatedResponse{
		Success:    true,
		Message:    "Agent sessions retrieved successfully",
		Data:       sessions,
		Pagination: pagination,
	})
}

// GetSessions godoc
// @Summary Get chat sessions
// @Description Get chat sessions with pagination and filters
// @Tags Chat
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param limit query int false "Items per page"
// @Param status query string false "Session status filter"
// @Param agent_id query string false "Agent ID filter"
// @Param department_id query string false "Department ID filter"
// @Success 200 {object} domain.ApiResponse{data=[]domain.ChatSession}
// @Failure 500 {object} domain.ApiResponse
// @Router /api/chat/sessions [get]
func (h *ChatHandler) GetSessions(c *fiber.Ctx) error {
	// Parse query parameters
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 20)
	status := c.Query("status")
	agentIDStr := c.Query("agent_id")
	departmentIDStr := c.Query("department_id")

	var agentID *uuid.UUID
	if agentIDStr != "" {
		id, err := uuid.Parse(agentIDStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(domain.ApiResponse{
				Success: false,
				Message: "Invalid agent ID format",
				Error:   err.Error(),
			})
		}
		agentID = &id
	}

	var departmentID *uuid.UUID
	if departmentIDStr != "" {
		id, err := uuid.Parse(departmentIDStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(domain.ApiResponse{
				Success: false,
				Message: "Invalid department ID format",
				Error:   err.Error(),
			})
		}
		departmentID = &id
	}

	sessions, total, err := h.chatUsecase.GetSessions(c.Context(), page, limit, status, agentID, departmentID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(domain.ApiResponse{
			Success: false,
			Message: "Failed to get sessions",
			Error:   err.Error(),
		})
	}

	totalPages := (total + limit - 1) / limit
	pagination := domain.PaginationInfo{
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
	}

	return c.JSON(domain.PaginatedResponse{
		Success:    true,
		Message:    "Sessions retrieved successfully",
		Data:       sessions,
		Pagination: pagination,
	})
}

// GetSession godoc
// @Summary Get single chat session
// @Description Get a single chat session by ID
// @Tags Chat
// @Accept json
// @Produce json
// @Param id path string true "Session ID"
// @Success 200 {object} domain.ApiResponse{data=domain.ChatSession}
// @Failure 400 {object} domain.ApiResponse
// @Failure 404 {object} domain.ApiResponse
// @Router /api/chat/sessions/{id} [get]
func (h *ChatHandler) GetSession(c *fiber.Ctx) error {
	sessionIDStr := c.Params("id")
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

	return c.JSON(domain.ApiResponse{
		Success: true,
		Message: "Session retrieved successfully",
		Data:    session,
	})
}

// GetSessionConnectionStatus godoc
// @Summary Get session connection status
// @Description Get connection status of clients in a chat session
// @Tags Chat
// @Accept json
// @Produce json
// @Param id path string true "Session ID"
// @Success 200 {object} domain.ApiResponse
// @Failure 400 {object} domain.ApiResponse
// @Router /api/chat/sessions/{id}/connection-status [get]
func (h *ChatHandler) GetSessionConnectionStatus(c *fiber.Ctx) error {
	sessionIDStr := c.Params("id")
	sessionID, err := uuid.Parse(sessionIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(domain.ApiResponse{
			Success: false,
			Message: "Invalid session ID",
			Error:   err.Error(),
		})
	}

	status := h.wsHandler.GetSessionConnectedClients(sessionID)

	return c.JSON(domain.ApiResponse{
		Success: true,
		Message: "Connection status retrieved successfully",
		Data:    status,
	})
}
