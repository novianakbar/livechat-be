package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/novianakbar/livechat-be/internal/domain"
)

type TicketBulkHandler struct {
	ticketService domain.TicketService
}

func NewTicketBulkHandler(ticketService domain.TicketService) *TicketBulkHandler {
	return &TicketBulkHandler{
		ticketService: ticketService,
	}
}

// BulkAssignTickets assigns multiple tickets to an agent
func (h *TicketBulkHandler) BulkAssignTickets(c *fiber.Ctx) error {
	var req domain.BulkAssignRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if err := h.ticketService.BulkAssignTickets(c.Context(), &req); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Tickets assigned successfully",
	})
}

// BulkUpdateStatus updates status for multiple tickets
func (h *TicketBulkHandler) BulkUpdateStatus(c *fiber.Ctx) error {
	var req domain.BulkUpdateStatusRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if err := h.ticketService.BulkUpdateStatus(c.Context(), &req); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Ticket statuses updated successfully",
	})
}

// BulkCloseTickets closes multiple tickets
func (h *TicketBulkHandler) BulkCloseTickets(c *fiber.Ctx) error {
	var req domain.BulkCloseRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if err := h.ticketService.BulkCloseTickets(c.Context(), &req); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Tickets closed successfully",
	})
}

// AutoAssignTicket automatically assigns a ticket to best available agent
func (h *TicketBulkHandler) AutoAssignTicket(c *fiber.Ctx) error {
	ticketID := c.Params("id")
	if ticketID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Ticket ID is required",
		})
	}

	var req struct {
		DepartmentID string `json:"department_id" validate:"required"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if err := h.ticketService.AutoAssignTicket(c.Context(), ticketID, req.DepartmentID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Ticket auto-assigned successfully",
	})
}
