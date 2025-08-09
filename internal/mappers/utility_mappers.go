package mappers

import (
	"database/sql"

	"github.com/novianakbar/livechat-be/internal/models"
	"github.com/novianakbar/livechat-shared/entities"
)

// UserToResponse converts User entity to UserResponse
func UserToResponse(entity *entities.User) *models.UserResponse {
	if entity == nil {
		return nil
	}

	response := &models.UserResponse{
		ID:        entity.ID,
		Email:     entity.Email,
		Name:      entity.Name,
		Role:      entity.Role,
		IsActive:  entity.IsActive,
		CreatedAt: FormatTime(entity.CreatedAt),
		UpdatedAt: FormatTime(entity.UpdatedAt),
	}

	// Handle optional department ID
	if entity.DepartmentID.Valid {
		response.DepartmentID = entity.DepartmentID.String
	}

	// Handle nested department if loaded
	if entity.Department != nil {
		response.Department = DepartmentToResponse(entity.Department)
	}

	return response
}

// DepartmentToResponse converts Department entity to DepartmentResponse
func DepartmentToResponse(entity *entities.Department) *models.DepartmentResponse {
	if entity == nil {
		return nil
	}

	response := &models.DepartmentResponse{
		ID:        entity.ID,
		Name:      entity.Name,
		CreatedAt: FormatTime(entity.CreatedAt),
		UpdatedAt: FormatTime(entity.UpdatedAt),
	}

	// Handle optional description
	if entity.Description.Valid {
		response.Description = entity.Description.String
	}

	return response
}

// UsersToResponse converts slice of User entity pointers to UserResponse slice
func UsersToResponse(entities []*entities.User) []models.UserResponse {
	if entities == nil {
		return nil
	}

	responses := make([]models.UserResponse, len(entities))
	for i, entity := range entities {
		response := UserToResponse(entity)
		if response != nil {
			responses[i] = *response
		}
	}
	return responses
}

// DepartmentsToResponse converts slice of Department entities to DepartmentResponse slice
func DepartmentsToResponse(entities []*entities.Department) []models.DepartmentResponse {
	if entities == nil {
		return nil
	}

	responses := make([]models.DepartmentResponse, len(entities))
	for i, entity := range entities {
		response := DepartmentToResponse(entity)
		if response != nil {
			responses[i] = *response
		}
	}
	return responses
}

// Helper function to safely get string from sql.NullString
func SafeStringFromNull(ns sql.NullString) string {
	if ns.Valid {
		return ns.String
	}
	return ""
}

// Helper function to safely get bool from sql.NullBool
func SafeBoolFromNull(nb sql.NullBool) bool {
	if nb.Valid {
		return nb.Bool
	}
	return false
}
