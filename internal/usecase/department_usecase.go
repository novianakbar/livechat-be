package usecase

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/novianakbar/livechat-be/internal/domain"
	"github.com/novianakbar/livechat-be/internal/models"
	"gorm.io/gorm"
)

type DepartmentUsecase interface {
	CreateDepartment(ctx context.Context, req *models.CreateDepartmentRequest) (*domain.Department, error)
	GetDepartment(ctx context.Context, id string) (*domain.Department, error)
	GetDepartmentWithRelations(ctx context.Context, id string) (*domain.Department, error)
	GetDepartments(ctx context.Context, query *models.DepartmentQueryRequest) ([]*domain.Department, int64, error)
	UpdateDepartment(ctx context.Context, id string, req *models.UpdateDepartmentRequest) (*domain.Department, error)
	DeleteDepartment(ctx context.Context, id string) error
	GetDepartmentsByParent(ctx context.Context, parentID string) ([]*domain.Department, error)
	GetDepartmentsBySupportLevel(ctx context.Context, level int) ([]*domain.Department, error)
}

type departmentUsecase struct {
	departmentRepo domain.DepartmentRepository
}

func NewDepartmentUsecase(departmentRepo domain.DepartmentRepository) DepartmentUsecase {
	return &departmentUsecase{
		departmentRepo: departmentRepo,
	}
}

func (u *departmentUsecase) CreateDepartment(ctx context.Context, req *models.CreateDepartmentRequest) (*domain.Department, error) {
	// Validate parent department if provided
	if req.ParentDeptID != nil && *req.ParentDeptID != "" {
		parentID, err := uuid.Parse(*req.ParentDeptID)
		if err != nil {
			return nil, fmt.Errorf("invalid parent department ID: %w", err)
		}

		_, err = u.departmentRepo.GetByID(ctx, parentID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, fmt.Errorf("parent department not found")
			}
			return nil, fmt.Errorf("failed to validate parent department: %w", err)
		}
	}

	// Validate escalation department if provided
	if req.EscalationDeptID != nil && *req.EscalationDeptID != "" {
		escalationID, err := uuid.Parse(*req.EscalationDeptID)
		if err != nil {
			return nil, fmt.Errorf("invalid escalation department ID: %w", err)
		}

		_, err = u.departmentRepo.GetByID(ctx, escalationID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, fmt.Errorf("escalation department not found")
			}
			return nil, fmt.Errorf("failed to validate escalation department: %w", err)
		}
	}

	// Create department entity
	department := &domain.Department{
		Name: req.Name,
	}

	// Set optional fields with defaults
	if req.Description != nil {
		department.Description = sql.NullString{String: *req.Description, Valid: true}
	}

	if req.IsActive != nil {
		department.IsActive = *req.IsActive
	} else {
		department.IsActive = true
	}

	if req.CanHandleTickets != nil {
		department.CanHandleTickets = *req.CanHandleTickets
	} else {
		department.CanHandleTickets = true
	}

	if req.MaxTicketsPerAgent != nil {
		department.MaxTicketsPerAgent = *req.MaxTicketsPerAgent
	} else {
		department.MaxTicketsPerAgent = 10
	}

	// Multi-level support fields
	if req.SupportLevel != nil {
		department.SupportLevel = *req.SupportLevel
	} else {
		department.SupportLevel = 0 // Default to L0
	}

	if req.ParentDeptID != nil && *req.ParentDeptID != "" {
		department.ParentDeptID = sql.NullString{String: *req.ParentDeptID, Valid: true}
	}

	if req.MaxEscalationLevel != nil {
		department.MaxEscalationLevel = *req.MaxEscalationLevel
	} else {
		department.MaxEscalationLevel = 3
	}

	if req.AutoAssignmentRule != nil {
		department.AutoAssignmentRule = *req.AutoAssignmentRule
	} else {
		department.AutoAssignmentRule = "round_robin"
	}

	if req.EscalationDeptID != nil && *req.EscalationDeptID != "" {
		department.EscalationDeptID = sql.NullString{String: *req.EscalationDeptID, Valid: true}
	}

	// Save to database
	err := u.departmentRepo.Create(ctx, department)
	if err != nil {
		return nil, fmt.Errorf("failed to create department: %w", err)
	}

	return department, nil
}

func (u *departmentUsecase) GetDepartment(ctx context.Context, id string) (*domain.Department, error) {
	departmentID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid department ID: %w", err)
	}

	department, err := u.departmentRepo.GetByID(ctx, departmentID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("department not found")
		}
		return nil, fmt.Errorf("failed to get department: %w", err)
	}

	return department, nil
}

func (u *departmentUsecase) GetDepartmentWithRelations(ctx context.Context, id string) (*domain.Department, error) {
	departmentID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid department ID: %w", err)
	}

	department, err := u.departmentRepo.GetWithRelations(ctx, departmentID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("department not found")
		}
		return nil, fmt.Errorf("failed to get department: %w", err)
	}

	return department, nil
}

func (u *departmentUsecase) GetDepartments(ctx context.Context, query *models.DepartmentQueryRequest) ([]*domain.Department, int64, error) {
	// Set defaults
	page := 1
	limit := 20

	if query.Page > 0 {
		page = query.Page
	}

	if query.Limit > 0 {
		limit = query.Limit
	}

	offset := (page - 1) * limit

	// Get departments with pagination and filters
	departments, err := u.departmentRepo.GetWithPagination(
		ctx,
		offset,
		limit,
		query.Search,
		query.IsActive,
		query.SupportLevel,
		query.ParentDeptID,
	)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get departments: %w", err)
	}

	// Get total count
	total, err := u.departmentRepo.Count(
		ctx,
		query.Search,
		query.IsActive,
		query.SupportLevel,
		query.ParentDeptID,
	)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count departments: %w", err)
	}

	return departments, total, nil
}

func (u *departmentUsecase) UpdateDepartment(ctx context.Context, id string, req *models.UpdateDepartmentRequest) (*domain.Department, error) {
	departmentID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid department ID: %w", err)
	}

	// Get existing department
	department, err := u.departmentRepo.GetByID(ctx, departmentID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("department not found")
		}
		return nil, fmt.Errorf("failed to get department: %w", err)
	}

	// Validate parent department if provided
	if req.ParentDeptID != nil && *req.ParentDeptID != "" {
		parentID, err := uuid.Parse(*req.ParentDeptID)
		if err != nil {
			return nil, fmt.Errorf("invalid parent department ID: %w", err)
		}

		// Prevent circular reference
		if parentID.String() == id {
			return nil, fmt.Errorf("department cannot be its own parent")
		}

		_, err = u.departmentRepo.GetByID(ctx, parentID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, fmt.Errorf("parent department not found")
			}
			return nil, fmt.Errorf("failed to validate parent department: %w", err)
		}
	}

	// Validate escalation department if provided
	if req.EscalationDeptID != nil && *req.EscalationDeptID != "" {
		escalationID, err := uuid.Parse(*req.EscalationDeptID)
		if err != nil {
			return nil, fmt.Errorf("invalid escalation department ID: %w", err)
		}

		_, err = u.departmentRepo.GetByID(ctx, escalationID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, fmt.Errorf("escalation department not found")
			}
			return nil, fmt.Errorf("failed to validate escalation department: %w", err)
		}
	}

	// Update fields
	if req.Name != nil {
		department.Name = *req.Name
	}

	if req.Description != nil {
		department.Description = sql.NullString{String: *req.Description, Valid: true}
	}

	if req.IsActive != nil {
		department.IsActive = *req.IsActive
	}

	if req.CanHandleTickets != nil {
		department.CanHandleTickets = *req.CanHandleTickets
	}

	if req.MaxTicketsPerAgent != nil {
		department.MaxTicketsPerAgent = *req.MaxTicketsPerAgent
	}

	if req.SupportLevel != nil {
		department.SupportLevel = *req.SupportLevel
	}

	if req.ParentDeptID != nil {
		if *req.ParentDeptID == "" {
			department.ParentDeptID = sql.NullString{Valid: false}
		} else {
			department.ParentDeptID = sql.NullString{String: *req.ParentDeptID, Valid: true}
		}
	}

	if req.MaxEscalationLevel != nil {
		department.MaxEscalationLevel = *req.MaxEscalationLevel
	}

	if req.AutoAssignmentRule != nil {
		department.AutoAssignmentRule = *req.AutoAssignmentRule
	}

	if req.EscalationDeptID != nil {
		if *req.EscalationDeptID == "" {
			department.EscalationDeptID = sql.NullString{Valid: false}
		} else {
			department.EscalationDeptID = sql.NullString{String: *req.EscalationDeptID, Valid: true}
		}
	}

	// Save changes
	err = u.departmentRepo.Update(ctx, department)
	if err != nil {
		return nil, fmt.Errorf("failed to update department: %w", err)
	}

	return department, nil
}

func (u *departmentUsecase) DeleteDepartment(ctx context.Context, id string) error {
	departmentID, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid department ID: %w", err)
	}

	// Check if department exists
	_, err = u.departmentRepo.GetByID(ctx, departmentID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("department not found")
		}
		return fmt.Errorf("failed to get department: %w", err)
	}

	// Check if department has child departments
	children, err := u.departmentRepo.GetByParent(ctx, departmentID)
	if err != nil {
		return fmt.Errorf("failed to check child departments: %w", err)
	}

	if len(children) > 0 {
		return fmt.Errorf("cannot delete department with child departments")
	}

	// TODO: Check if department has active tickets
	// This would require access to ticket repository

	// Delete department
	err = u.departmentRepo.Delete(ctx, departmentID)
	if err != nil {
		return fmt.Errorf("failed to delete department: %w", err)
	}

	return nil
}

func (u *departmentUsecase) GetDepartmentsByParent(ctx context.Context, parentID string) ([]*domain.Department, error) {
	parentDeptID, err := uuid.Parse(parentID)
	if err != nil {
		return nil, fmt.Errorf("invalid parent department ID: %w", err)
	}

	departments, err := u.departmentRepo.GetByParent(ctx, parentDeptID)
	if err != nil {
		return nil, fmt.Errorf("failed to get departments by parent: %w", err)
	}

	return departments, nil
}

func (u *departmentUsecase) GetDepartmentsBySupportLevel(ctx context.Context, level int) ([]*domain.Department, error) {
	if level < 0 || level > 3 {
		return nil, fmt.Errorf("invalid support level: must be between 0 and 3")
	}

	departments, err := u.departmentRepo.GetBySupportLevel(ctx, level)
	if err != nil {
		return nil, fmt.Errorf("failed to get departments by support level: %w", err)
	}

	return departments, nil
}
