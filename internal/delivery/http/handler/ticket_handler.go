package handler

import (
	"context"
	"net/http"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/novianakbar/livechat-be/internal/domain"
	"github.com/novianakbar/livechat-be/internal/mappers"
	"github.com/novianakbar/livechat-be/internal/models"
)

type TicketHandler struct {
	ticketService domain.TicketService
	mapper        *mappers.TicketMapper
}

func NewTicketHandler(ticketService domain.TicketService) *TicketHandler {
	return &TicketHandler{
		ticketService: ticketService,
		mapper:        mappers.NewTicketMapper(),
	}
}

// CreateTicket creates a new support ticket
// @Summary Create a new ticket
// @Description Create a new support ticket
// @Tags tickets
// @Accept json
// @Produce json
// @Param ticket body domain.CreateTicketRequest true "Create ticket request"
// @Success 201 {object} domain.Ticket
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/tickets [post]
func (h *TicketHandler) CreateTicket(c *fiber.Ctx) error {
	var req domain.CreateTicketRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	ctx := context.Background()
	ticket, err := h.ticketService.CreateTicket(ctx, &req)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	response := h.mapper.ToTicketResponse(ticket)
	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"success": true,
		"data":    response,
	})
}

// GetTicket retrieves a ticket by ID
// @Summary Get ticket by ID
// @Description Get a specific ticket by its ID
// @Tags tickets
// @Produce json
// @Param id path string true "Ticket ID"
// @Success 200 {object} domain.Ticket
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/tickets/{id} [get]
func (h *TicketHandler) GetTicket(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Ticket ID is required",
		})
	}

	ctx := context.Background()
	ticket, err := h.ticketService.GetTicket(ctx, id)
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"error": "Ticket not found",
		})
	}

	response := h.mapper.ToTicketResponse(ticket)
	return c.JSON(fiber.Map{
		"success": true,
		"data":    response,
	})
}

// GetTicketByCode retrieves a ticket by ticket code
// @Summary Get ticket by code
// @Description Get a specific ticket by its code
// @Tags tickets
// @Produce json
// @Param code path string true "Ticket Code"
// @Success 200 {object} domain.Ticket
// @Failure 404 {object} map[string]string
// @Router /api/tickets/code/{code} [get]
func (h *TicketHandler) GetTicketByCode(c *fiber.Ctx) error {
	code := c.Params("code")
	if code == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Ticket code is required",
		})
	}

	ctx := context.Background()
	ticket, err := h.ticketService.GetTicketByCode(ctx, code)
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"error": "Ticket not found",
		})
	}

	return c.JSON(ticket)
}

// GetTicketByAccessToken retrieves a ticket by access token (for customer portal)
// @Summary Get ticket by access token
// @Description Get a specific ticket by its access token for customer access
// @Tags tickets
// @Produce json
// @Param token path string true "Access Token"
// @Success 200 {object} domain.Ticket
// @Failure 404 {object} map[string]string
// @Router /api/public/tickets/{token} [get]
func (h *TicketHandler) GetTicketByAccessToken(c *fiber.Ctx) error {
	token := c.Params("token")
	if token == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Access token is required",
		})
	}

	ctx := context.Background()
	ticket, err := h.ticketService.GetTicketByAccessToken(ctx, token)
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"error": "Ticket not found",
		})
	}

	return c.JSON(ticket)
}

// UpdateTicket updates an existing ticket
// @Summary Update ticket
// @Description Update an existing ticket
// @Tags tickets
// @Accept json
// @Produce json
// @Param id path string true "Ticket ID"
// @Param ticket body domain.UpdateTicketRequest true "Update ticket request"
// @Success 200 {object} domain.Ticket
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/tickets/{id} [put]
func (h *TicketHandler) UpdateTicket(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Ticket ID is required",
		})
	}

	var req domain.UpdateTicketRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	req.ID = id // Set ID from URL parameter

	ctx := context.Background()
	ticket, err := h.ticketService.UpdateTicket(ctx, &req)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return c.Status(http.StatusNotFound).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(ticket)
}

// AssignTicket assigns a ticket to an agent
// @Summary Assign ticket to agent
// @Description Assign a ticket to a specific agent
// @Tags tickets
// @Accept json
// @Produce json
// @Param id path string true "Ticket ID"
// @Param assignment body domain.AssignTicketRequest true "Assignment request"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/tickets/{id}/assign [post]
func (h *TicketHandler) AssignTicket(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Ticket ID is required",
		})
	}

	var req domain.AssignTicketRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	req.TicketID = id // Set ID from URL parameter

	ctx := context.Background()
	err := h.ticketService.AssignTicket(ctx, &req)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return c.Status(http.StatusNotFound).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Ticket assigned successfully",
	})
}

// EscalateTicket escalates a ticket
// @Summary Escalate ticket
// @Description Escalate a ticket to higher level support
// @Tags tickets
// @Accept json
// @Produce json
// @Param id path string true "Ticket ID"
// @Param escalation body domain.EscalateTicketRequest true "Escalation request"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/tickets/{id}/escalate [post]
func (h *TicketHandler) EscalateTicket(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Ticket ID is required",
		})
	}

	var req domain.EscalateTicketRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	req.TicketID = id // Set ID from URL parameter

	ctx := context.Background()
	err := h.ticketService.EscalateTicket(ctx, &req)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return c.Status(http.StatusNotFound).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Ticket escalated successfully",
	})
}

// AddComment adds a comment to a ticket
// @Summary Add comment to ticket
// @Description Add a comment or response to a ticket
// @Tags tickets
// @Accept json
// @Produce json
// @Param id path string true "Ticket ID"
// @Param comment body domain.AddCommentRequest true "Comment request"
// @Success 201 {object} domain.TicketComment
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/tickets/{id}/comments [post]
func (h *TicketHandler) AddComment(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Ticket ID is required",
		})
	}

	var req domain.AddCommentRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	req.TicketID = id // Set ID from URL parameter

	ctx := context.Background()
	comment, err := h.ticketService.AddComment(ctx, &req)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return c.Status(http.StatusNotFound).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(http.StatusCreated).JSON(comment)
}

// GetTicketList retrieves a list of tickets with filtering
// @Summary Get tickets list
// @Description Get a list of tickets with optional filtering
// @Tags tickets
// @Produce json
// @Param status query string false "Filter by status"
// @Param priority query string false "Filter by priority"
// @Param category_id query string false "Filter by category ID"
// @Param department_id query string false "Filter by department ID"
// @Param assigned_to query string false "Filter by assigned agent ID"
// @Param limit query int false "Number of results per page" default(20)
// @Param offset query int false "Number of results to skip" default(0)
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]string
// @Router /api/tickets [get]
func (h *TicketHandler) GetTicketList(c *fiber.Ctx) error {
	// Parse query parameters
	filter := &domain.TicketFilter{}

	if status := c.Query("status"); status != "" {
		filter.Status = &status
	}
	if priority := c.Query("priority"); priority != "" {
		filter.Priority = &priority
	}
	if categoryID := c.Query("category_id"); categoryID != "" {
		filter.CategoryID = &categoryID
	}
	if departmentID := c.Query("department_id"); departmentID != "" {
		filter.DepartmentID = &departmentID
	}
	if assignedTo := c.Query("assigned_to"); assignedTo != "" {
		filter.AssignedTo = &assignedTo
	}
	if customerInfo := c.Query("customer_info"); customerInfo != "" {
		filter.CustomerInfo = &customerInfo
	}
	if subject := c.Query("subject"); subject != "" {
		filter.Subject = &subject
	}

	// Parse pagination
	limit, _ := strconv.Atoi(c.Query("limit", "20"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))
	filter.Limit = limit
	filter.Offset = offset

	// Parse sorting
	filter.SortBy = c.Query("sort_by", "created_at")
	filter.SortOrder = c.Query("sort_order", "desc")

	ctx := context.Background()
	tickets, total, err := h.ticketService.GetTicketList(ctx, filter)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	responses := h.mapper.ToTicketListResponse(tickets)
	page := (offset / limit) + 1
	totalPages := int((total + int64(limit) - 1) / int64(limit))

	return c.JSON(models.PaginatedResponse[models.TicketResponse]{
		Data:       responses,
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
	})
}

// GetAgentTickets retrieves tickets assigned to a specific agent
// @Summary Get agent tickets
// @Description Get tickets assigned to a specific agent
// @Tags tickets
// @Produce json
// @Param agent_id path string true "Agent ID"
// @Param status query string false "Filter by status (comma-separated)"
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]string
// @Router /api/agents/{agent_id}/tickets [get]
func (h *TicketHandler) GetAgentTickets(c *fiber.Ctx) error {
	agentID := c.Params("agent_id")
	if agentID == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Agent ID is required",
		})
	}

	// Parse status filter
	var status []string
	if statusParam := c.Query("status"); statusParam != "" {
		status = strings.Split(statusParam, ",")
	}

	ctx := context.Background()
	tickets, err := h.ticketService.GetAgentTickets(ctx, agentID, status)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"tickets": tickets,
		"total":   len(tickets),
	})
}

// GetDepartmentTickets retrieves tickets for a specific department
// @Summary Get department tickets
// @Description Get tickets for a specific department
// @Tags tickets
// @Produce json
// @Param department_id path string true "Department ID"
// @Param status query string false "Filter by status (comma-separated)"
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]string
// @Router /api/departments/{department_id}/tickets [get]
func (h *TicketHandler) GetDepartmentTickets(c *fiber.Ctx) error {
	deptID := c.Params("department_id")
	if deptID == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Department ID is required",
		})
	}

	// Parse status filter
	var status []string
	if statusParam := c.Query("status"); statusParam != "" {
		status = strings.Split(statusParam, ",")
	}

	ctx := context.Background()
	tickets, err := h.ticketService.GetDepartmentTickets(ctx, deptID, status)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"tickets": tickets,
		"total":   len(tickets),
	})
}
