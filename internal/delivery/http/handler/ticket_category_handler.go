package handler

import (
	"context"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/novianakbar/livechat-be/internal/domain"
	"github.com/novianakbar/livechat-be/internal/mappers"
)

type TicketCategoryHandler struct {
	categoryRepo domain.TicketCategoryRepository
	mapper       *mappers.TicketMapper
}

func NewTicketCategoryHandler(categoryRepo domain.TicketCategoryRepository) *TicketCategoryHandler {
	return &TicketCategoryHandler{
		categoryRepo: categoryRepo,
		mapper:       mappers.NewTicketMapper(),
	}
}

// GetCategories retrieves all active ticket categories
// @Summary Get all ticket categories
// @Description Get all active ticket categories
// @Tags ticket-categories
// @Produce json
// @Success 200 {array} domain.TicketCategory
// @Failure 500 {object} map[string]string
// @Router /api/ticket-categories [get]
func (h *TicketCategoryHandler) GetCategories(c *fiber.Ctx) error {
	ctx := context.Background()
	categories, err := h.categoryRepo.GetActive(ctx)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	responses := h.mapper.ToTicketCategoryListResponse(categories)
	return c.JSON(fiber.Map{
		"success": true,
		"data":    responses,
	})
}

// GetCategory retrieves a specific ticket category
// @Summary Get ticket category by ID
// @Description Get a specific ticket category by ID
// @Tags ticket-categories
// @Produce json
// @Param id path string true "Category ID"
// @Success 200 {object} domain.TicketCategory
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/ticket-categories/{id} [get]
func (h *TicketCategoryHandler) GetCategory(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Category ID is required",
		})
	}

	ctx := context.Background()
	category, err := h.categoryRepo.GetByID(ctx, id)
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"error": "Category not found",
		})
	}

	return c.JSON(category)
}
