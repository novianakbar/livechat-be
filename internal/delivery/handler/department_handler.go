package handler

import (
	"math"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/novianakbar/livechat-be/internal/mappers"
	"github.com/novianakbar/livechat-be/internal/models"
	"github.com/novianakbar/livechat-be/internal/usecase"
)

type DepartmentHandler struct {
	departmentUsecase usecase.DepartmentUsecase
}

func NewDepartmentHandler(departmentUsecase usecase.DepartmentUsecase) *DepartmentHandler {
	return &DepartmentHandler{
		departmentUsecase: departmentUsecase,
	}
}

// CreateDepartment creates a new department
// @Summary Create a new department
// @Description Create a new department with optional hierarchy and escalation settings
// @Tags departments
// @Accept json
// @Produce json
// @Param department body models.CreateDepartmentRequest true "Department creation data"
// @Success 201 {object} fiber.Map
// @Failure 400 {object} fiber.Map
// @Failure 401 {object} fiber.Map
// @Failure 500 {object} fiber.Map
// @Router /api/departments [post]
// @Security Bearer
func (h *DepartmentHandler) CreateDepartment(c *fiber.Ctx) error {
	var req models.CreateDepartmentRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Invalid request body",
			"error":   err.Error(),
		})
	}

	// Basic validation
	if req.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Department name is required",
		})
	}

	// Create department
	department, err := h.departmentUsecase.CreateDepartment(c.Context(), &req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to create department",
			"error":   err.Error(),
		})
	}

	// Convert to response
	response := mappers.DepartmentToResponse(department)

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "Department created successfully",
		"data":    response,
	})
}

// GetDepartments gets departments with pagination and filtering
// @Summary Get departments with pagination and filtering
// @Description Get a list of departments with optional filtering by name, status, support level, etc.
// @Tags departments
// @Produce json
// @Param page query int false "Page number (default 1)"
// @Param limit query int false "Items per page (default 20, max 100)"
// @Param search query string false "Search in name and description"
// @Param is_active query bool false "Filter by active status"
// @Param support_level query int false "Filter by support level (0-3)"
// @Param parent_dept_id query string false "Filter by parent department ID"
// @Success 200 {object} fiber.Map
// @Failure 400 {object} fiber.Map
// @Failure 401 {object} fiber.Map
// @Failure 500 {object} fiber.Map
// @Router /api/departments [get]
// @Security Bearer
func (h *DepartmentHandler) GetDepartments(c *fiber.Ctx) error {
	var query models.DepartmentQueryRequest

	// Parse query parameters
	if pageStr := c.Query("page"); pageStr != "" {
		if page, err := strconv.Atoi(pageStr); err == nil && page > 0 {
			query.Page = page
		}
	}

	if limitStr := c.Query("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil && limit > 0 && limit <= 100 {
			query.Limit = limit
		}
	}

	query.Search = c.Query("search")

	if isActiveStr := c.Query("is_active"); isActiveStr != "" {
		if isActive, err := strconv.ParseBool(isActiveStr); err == nil {
			query.IsActive = &isActive
		}
	}

	if supportLevelStr := c.Query("support_level"); supportLevelStr != "" {
		if supportLevel, err := strconv.Atoi(supportLevelStr); err == nil && supportLevel >= 0 && supportLevel <= 3 {
			query.SupportLevel = &supportLevel
		}
	}

	query.ParentDeptID = c.Query("parent_dept_id")

	// Get departments
	departments, total, err := h.departmentUsecase.GetDepartments(c.Context(), &query)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to get departments",
			"error":   err.Error(),
		})
	}

	// Convert to response
	responses := mappers.DepartmentsToResponse(departments)

	// Calculate pagination info
	page := 1
	limit := 20
	if query.Page > 0 {
		page = query.Page
	}
	if query.Limit > 0 {
		limit = query.Limit
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Departments retrieved successfully",
		"data":    responses,
		"meta": fiber.Map{
			"page":        page,
			"limit":       limit,
			"total":       total,
			"total_pages": totalPages,
		},
	})
}

// GetDepartment gets a single department by ID
// @Summary Get department by ID
// @Description Get a single department by its ID with full relations
// @Tags departments
// @Produce json
// @Param id path string true "Department ID"
// @Success 200 {object} fiber.Map
// @Failure 400 {object} fiber.Map
// @Failure 401 {object} fiber.Map
// @Failure 404 {object} fiber.Map
// @Failure 500 {object} fiber.Map
// @Router /api/departments/{id} [get]
// @Security Bearer
func (h *DepartmentHandler) GetDepartment(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Department ID is required",
		})
	}

	// Get department with relations
	department, err := h.departmentUsecase.GetDepartmentWithRelations(c.Context(), id)
	if err != nil {
		if err.Error() == "department not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"success": false,
				"message": "Department not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to get department",
			"error":   err.Error(),
		})
	}

	// Convert to response
	response := mappers.DepartmentToResponse(department)

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Department retrieved successfully",
		"data":    response,
	})
}

// UpdateDepartment updates an existing department
// @Summary Update department
// @Description Update an existing department's information
// @Tags departments
// @Accept json
// @Produce json
// @Param id path string true "Department ID"
// @Param department body models.UpdateDepartmentRequest true "Department update data"
// @Success 200 {object} fiber.Map
// @Failure 400 {object} fiber.Map
// @Failure 401 {object} fiber.Map
// @Failure 404 {object} fiber.Map
// @Failure 500 {object} fiber.Map
// @Router /api/departments/{id} [put]
// @Security Bearer
func (h *DepartmentHandler) UpdateDepartment(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Department ID is required",
		})
	}

	var req models.UpdateDepartmentRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Invalid request body",
			"error":   err.Error(),
		})
	}

	// Update department
	department, err := h.departmentUsecase.UpdateDepartment(c.Context(), id, &req)
	if err != nil {
		if err.Error() == "department not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"success": false,
				"message": "Department not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to update department",
			"error":   err.Error(),
		})
	}

	// Convert to response
	response := mappers.DepartmentToResponse(department)

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Department updated successfully",
		"data":    response,
	})
}

// DeleteDepartment deletes a department
// @Summary Delete department
// @Description Delete a department (soft delete)
// @Tags departments
// @Produce json
// @Param id path string true "Department ID"
// @Success 200 {object} fiber.Map
// @Failure 400 {object} fiber.Map
// @Failure 401 {object} fiber.Map
// @Failure 404 {object} fiber.Map
// @Failure 409 {object} fiber.Map
// @Failure 500 {object} fiber.Map
// @Router /api/departments/{id} [delete]
// @Security Bearer
func (h *DepartmentHandler) DeleteDepartment(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Department ID is required",
		})
	}

	// Delete department
	err := h.departmentUsecase.DeleteDepartment(c.Context(), id)
	if err != nil {
		if err.Error() == "department not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"success": false,
				"message": "Department not found",
			})
		}
		if err.Error() == "cannot delete department with child departments" {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"success": false,
				"message": "Cannot delete department with child departments",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to delete department",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Department deleted successfully",
	})
}

// GetDepartmentsByParent gets child departments of a parent department
// @Summary Get child departments
// @Description Get all child departments of a specific parent department
// @Tags departments
// @Produce json
// @Param parent_id path string true "Parent Department ID"
// @Success 200 {object} fiber.Map
// @Failure 400 {object} fiber.Map
// @Failure 401 {object} fiber.Map
// @Failure 500 {object} fiber.Map
// @Router /api/departments/parent/{parent_id} [get]
// @Security Bearer
func (h *DepartmentHandler) GetDepartmentsByParent(c *fiber.Ctx) error {
	parentID := c.Params("parent_id")
	if parentID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Parent department ID is required",
		})
	}

	// Get child departments
	departments, err := h.departmentUsecase.GetDepartmentsByParent(c.Context(), parentID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to get child departments",
			"error":   err.Error(),
		})
	}

	// Convert to response
	responses := mappers.DepartmentsToResponse(departments)

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Child departments retrieved successfully",
		"data":    responses,
	})
}

// GetDepartmentsBySupportLevel gets departments by support level
// @Summary Get departments by support level
// @Description Get all departments by their support level (L0, L1, L2, L3)
// @Tags departments
// @Produce json
// @Param level path int true "Support Level (0-3)"
// @Success 200 {object} fiber.Map
// @Failure 400 {object} fiber.Map
// @Failure 401 {object} fiber.Map
// @Failure 500 {object} fiber.Map
// @Router /api/departments/level/{level} [get]
// @Security Bearer
func (h *DepartmentHandler) GetDepartmentsBySupportLevel(c *fiber.Ctx) error {
	levelStr := c.Params("level")
	if levelStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Support level is required",
		})
	}

	level, err := strconv.Atoi(levelStr)
	if err != nil || level < 0 || level > 3 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Invalid support level: must be between 0 and 3",
		})
	}

	// Get departments by support level
	departments, err := h.departmentUsecase.GetDepartmentsBySupportLevel(c.Context(), level)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to get departments by support level",
			"error":   err.Error(),
		})
	}

	// Convert to response
	responses := mappers.DepartmentsToResponse(departments)

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Departments retrieved successfully",
		"data":    responses,
	})
}
