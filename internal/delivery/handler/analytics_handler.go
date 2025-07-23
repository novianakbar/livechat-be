package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/novianakbar/livechat-be/internal/domain"
	"github.com/novianakbar/livechat-be/internal/usecase"
)

type AnalyticsHandler struct {
	analyticsUsecase *usecase.AnalyticsUsecase
}

func NewAnalyticsHandler(analyticsUsecase *usecase.AnalyticsUsecase) *AnalyticsHandler {
	return &AnalyticsHandler{
		analyticsUsecase: analyticsUsecase,
	}
}

// GetDashboardStats godoc
// @Summary Get dashboard statistics
// @Description Get dashboard statistics for OSS support system
// @Tags Analytics
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} domain.ApiResponse{data=domain.DashboardStats}
// @Failure 401 {object} domain.ApiResponse
// @Failure 500 {object} domain.ApiResponse
// @Router /api/analytics/dashboard [get]
func (h *AnalyticsHandler) GetDashboardStats(c *fiber.Ctx) error {
	stats, err := h.analyticsUsecase.GetDashboardStats()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(domain.ApiResponse{
			Success: false,
			Message: "Failed to get dashboard statistics",
			Error:   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(domain.ApiResponse{
		Success: true,
		Message: "Dashboard statistics retrieved successfully",
		Data:    stats,
	})
}

// GetAnalytics godoc
// @Summary Get analytics data
// @Description Get detailed analytics data with filtering options
// @Tags Analytics
// @Accept json
// @Produce json
// @Security Bearer
// @Param start_date query string false "Start date (YYYY-MM-DD)"
// @Param end_date query string false "End date (YYYY-MM-DD)"
// @Param department_id query string false "Department ID"
// @Param agent_id query string false "Agent ID"
// @Success 200 {object} domain.ApiResponse{data=[]domain.ChatAnalytics}
// @Failure 401 {object} domain.ApiResponse
// @Failure 500 {object} domain.ApiResponse
// @Router /api/analytics [get]
func (h *AnalyticsHandler) GetAnalytics(c *fiber.Ctx) error {
	var req domain.GetAnalyticsRequest
	if err := c.QueryParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(domain.ApiResponse{
			Success: false,
			Message: "Invalid query parameters",
			Error:   err.Error(),
		})
	}

	analytics, err := h.analyticsUsecase.GetAnalytics(&req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(domain.ApiResponse{
			Success: false,
			Message: "Failed to get analytics data",
			Error:   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(domain.ApiResponse{
		Success: true,
		Message: "Analytics data retrieved successfully",
		Data:    analytics,
	})
}
