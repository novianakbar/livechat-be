package handler

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/novianakbar/livechat-be/internal/domain"
)

type TicketStatsHandler struct {
	ticketService domain.TicketService
}

func NewTicketStatsHandler(ticketService domain.TicketService) *TicketStatsHandler {
	return &TicketStatsHandler{
		ticketService: ticketService,
	}
}

// GetDashboardStats returns comprehensive dashboard statistics
func (h *TicketStatsHandler) GetDashboardStats(c *fiber.Ctx) error {
	stats, err := h.ticketService.GetDashboardStats(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Dashboard stats retrieved successfully",
		"data":    stats,
	})
}

// GetTicketStats returns detailed ticket statistics
func (h *TicketStatsHandler) GetTicketStats(c *fiber.Ctx) error {
	// Parse query parameters for date range
	startDateStr := c.Query("start_date", "")
	endDateStr := c.Query("end_date", "")
	departmentID := c.Query("department_id", "")
	agentID := c.Query("agent_id", "")

	var startDate, endDate *time.Time
	if startDateStr != "" {
		if parsed, err := time.Parse("2006-01-02", startDateStr); err == nil {
			startDate = &parsed
		}
	}
	if endDateStr != "" {
		if parsed, err := time.Parse("2006-01-02", endDateStr); err == nil {
			endDate = &parsed
		}
	}

	req := &domain.TicketStatsRequest{
		StartDate:    startDate,
		EndDate:      endDate,
		DepartmentID: departmentID,
		AgentID:      agentID,
	}

	stats, err := h.ticketService.GetTicketStats(c.Context(), req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Ticket stats retrieved successfully",
		"data":    stats,
	})
}

// GetPerformanceMetrics returns agent and department performance metrics
func (h *TicketStatsHandler) GetPerformanceMetrics(c *fiber.Ctx) error {
	// Parse query parameters
	periodStr := c.Query("period", "30") // days
	departmentID := c.Query("department_id", "")

	period, err := strconv.Atoi(periodStr)
	if err != nil {
		period = 30
	}

	req := &domain.TicketPerformanceRequest{
		Period:       period,
		DepartmentID: departmentID,
	}

	metrics, err := h.ticketService.GetPerformanceMetrics(c.Context(), req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Performance metrics retrieved successfully",
		"data":    metrics,
	})
}

// GetSLAReport returns SLA compliance reports
func (h *TicketStatsHandler) GetSLAReport(c *fiber.Ctx) error {
	// Parse query parameters
	periodStr := c.Query("period", "30") // days
	departmentID := c.Query("department_id", "")

	period, err := strconv.Atoi(periodStr)
	if err != nil {
		period = 30
	}

	req := &domain.TicketSLAReportRequest{
		Period:       period,
		DepartmentID: departmentID,
	}

	report, err := h.ticketService.GetSLAReport(c.Context(), req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "SLA report retrieved successfully",
		"data":    report,
	})
}

// GetEscalationStats returns escalation statistics and trends
func (h *TicketStatsHandler) GetEscalationStats(c *fiber.Ctx) error {
	// Parse query parameters
	periodStr := c.Query("period", "30") // days
	departmentID := c.Query("department_id", "")

	period, err := strconv.Atoi(periodStr)
	if err != nil {
		period = 30
	}

	req := &domain.TicketEscalationStatsRequest{
		Period:       period,
		DepartmentID: departmentID,
	}

	stats, err := h.ticketService.GetEscalationStats(c.Context(), req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Escalation stats retrieved successfully",
		"data":    stats,
	})
}
