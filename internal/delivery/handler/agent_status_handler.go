package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/novianakbar/livechat-be/internal/delivery/middleware"
	"github.com/novianakbar/livechat-be/internal/domain"
	"github.com/novianakbar/livechat-be/internal/service"
)

type AgentStatusHandler struct {
	agentStatusService *service.AgentStatusService
}

func NewAgentStatusHandler(agentStatusService *service.AgentStatusService) *AgentStatusHandler {
	return &AgentStatusHandler{
		agentStatusService: agentStatusService,
	}
}

// HeartbeatRequest represents the payload for agent heartbeat
type HeartbeatRequest struct {
	Status string `json:"status,omitempty"` // online, busy, away - optional, defaults to "online"
}

// AgentHeartbeat handles agent heartbeat requests
func (h *AgentStatusHandler) AgentHeartbeat(c *fiber.Ctx) error {
	// Get user from token
	user := middleware.GetUserFromContext(c)
	if user == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(domain.ApiResponse{
			Success: false,
			Message: "User not found in context",
			Error:   "authentication required",
		})
	}

	// Verify user is an agent or admin
	if user.Role != "agent" && user.Role != "admin" {
		return c.Status(fiber.StatusForbidden).JSON(domain.ApiResponse{
			Success: false,
			Message: "Access denied",
			Error:   "only agents and admins can send heartbeat",
		})
	}

	// Parse request body
	var req HeartbeatRequest
	if err := c.BodyParser(&req); err != nil {
		// If body parsing fails, just use default status
		req.Status = "online"
	}

	// Default status if not provided
	if req.Status == "" {
		req.Status = "online"
	}

	// Update agent heartbeat
	userUUID, err := uuid.Parse(user.ID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(domain.ApiResponse{
			Success: false,
			Message: "Invalid user ID format",
			Error:   err.Error(),
		})
	}

	if err := h.agentStatusService.UpdateAgentHeartbeat(c.Context(), userUUID, req.Status); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(domain.ApiResponse{
			Success: false,
			Message: "Failed to update agent heartbeat",
			Error:   err.Error(),
		})
	}

	return c.JSON(domain.ApiResponse{
		Success: true,
		Message: "Heartbeat updated successfully",
		Data: fiber.Map{
			"agent_id": user.ID,
			"status":   req.Status,
		},
	})
}

// GetOnlineAgents returns all currently online agents
func (h *AgentStatusHandler) GetOnlineAgents(c *fiber.Ctx) error {
	agents, err := h.agentStatusService.GetOnlineAgents(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(domain.ApiResponse{
			Success: false,
			Message: "Failed to get online agents",
			Error:   err.Error(),
		})
	}

	return c.JSON(domain.ApiResponse{
		Success: true,
		Message: "Online agents retrieved successfully",
		Data: fiber.Map{
			"agents": agents,
			"count":  len(agents),
		},
	})
}

// GetOnlineAgentsByDepartment returns online agents by department
func (h *AgentStatusHandler) GetOnlineAgentsByDepartment(c *fiber.Ctx) error {
	departmentIDStr := c.Params("department_id")
	if departmentIDStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(domain.ApiResponse{
			Success: false,
			Message: "Department ID is required",
			Error:   "missing department_id parameter",
		})
	}

	departmentID, err := uuid.Parse(departmentIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(domain.ApiResponse{
			Success: false,
			Message: "Invalid department ID format",
			Error:   "department_id must be a valid UUID",
		})
	}

	agents, err := h.agentStatusService.GetOnlineAgentsByDepartment(c.Context(), departmentID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(domain.ApiResponse{
			Success: false,
			Message: "Failed to get online agents by department",
			Error:   err.Error(),
		})
	}

	return c.JSON(domain.ApiResponse{
		Success: true,
		Message: "Online agents by department retrieved successfully",
		Data: fiber.Map{
			"department_id": departmentID,
			"agents":        agents,
			"count":         len(agents),
		},
	})
}

// GetAgentStatus returns specific agent status
func (h *AgentStatusHandler) GetAgentStatus(c *fiber.Ctx) error {
	agentIDStr := c.Params("agent_id")
	if agentIDStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(domain.ApiResponse{
			Success: false,
			Message: "Agent ID is required",
			Error:   "missing agent_id parameter",
		})
	}

	agentID, err := uuid.Parse(agentIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(domain.ApiResponse{
			Success: false,
			Message: "Invalid agent ID format",
			Error:   "agent_id must be a valid UUID",
		})
	}

	status, err := h.agentStatusService.GetAgentStatus(c.Context(), agentID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(domain.ApiResponse{
			Success: false,
			Message: "Failed to get agent status",
			Error:   err.Error(),
		})
	}

	if status == nil {
		return c.Status(fiber.StatusNotFound).JSON(domain.ApiResponse{
			Success: false,
			Message: "Agent is not online",
			Data: fiber.Map{
				"agent_id": agentID,
				"online":   false,
			},
		})
	}

	return c.JSON(domain.ApiResponse{
		Success: true,
		Message: "Agent status retrieved successfully",
		Data: fiber.Map{
			"agent":  status,
			"online": true,
		},
	})
}

// GetDepartmentStats returns online agents count by department
func (h *AgentStatusHandler) GetDepartmentStats(c *fiber.Ctx) error {
	stats, err := h.agentStatusService.GetDepartmentStats(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(domain.ApiResponse{
			Success: false,
			Message: "Failed to get department statistics",
			Error:   err.Error(),
		})
	}

	return c.JSON(domain.ApiResponse{
		Success: true,
		Message: "Department statistics retrieved successfully",
		Data: fiber.Map{
			"departments": stats,
		},
	})
}

// SetAgentOffline sets agent as offline (useful for logout)
func (h *AgentStatusHandler) SetAgentOffline(c *fiber.Ctx) error {
	user := middleware.GetUserFromContext(c)
	if user == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(domain.ApiResponse{
			Success: false,
			Message: "User not found in context",
			Error:   "authentication required",
		})
	}

	userUUID, err := uuid.Parse(user.ID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(domain.ApiResponse{
			Success: false,
			Message: "Invalid user ID format",
			Error:   err.Error(),
		})
	}

	if err := h.agentStatusService.SetAgentOffline(c.Context(), userUUID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(domain.ApiResponse{
			Success: false,
			Message: "Failed to set agent offline",
			Error:   err.Error(),
		})
	}

	return c.JSON(domain.ApiResponse{
		Success: true,
		Message: "Agent set offline successfully",
		Data: fiber.Map{
			"agent_id": user.ID,
			"status":   "offline",
		},
	})
}
