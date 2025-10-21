package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/novianakbar/livechat-be/internal/domain"
	"gorm.io/gorm"
)

type departmentRepository struct {
	db *gorm.DB
}

func NewDepartmentRepository(db *gorm.DB) domain.DepartmentRepository {
	return &departmentRepository{db: db}
}

func (r *departmentRepository) Create(ctx context.Context, department *domain.Department) error {
	return r.db.WithContext(ctx).Create(department).Error
}

func (r *departmentRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Department, error) {
	var department domain.Department
	err := r.db.WithContext(ctx).First(&department, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &department, nil
}

func (r *departmentRepository) GetAll(ctx context.Context) ([]*domain.Department, error) {
	var departments []*domain.Department
	err := r.db.WithContext(ctx).Find(&departments).Error
	return departments, err
}

func (r *departmentRepository) GetWithPagination(ctx context.Context, offset, limit int, search string, isActive *bool, supportLevel *int, parentDeptID string) ([]*domain.Department, error) {
	var departments []*domain.Department
	query := r.db.WithContext(ctx)

	// Apply filters
	if search != "" {
		query = query.Where("name ILIKE ? OR description ILIKE ?", "%"+search+"%", "%"+search+"%")
	}

	if isActive != nil {
		query = query.Where("is_active = ?", *isActive)
	}

	if supportLevel != nil {
		query = query.Where("support_level = ?", *supportLevel)
	}

	if parentDeptID != "" {
		query = query.Where("parent_dept_id = ?", parentDeptID)
	}

	// Apply pagination and ordering
	err := query.Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Preload("ParentDept").
		Preload("EscalationDept").
		Find(&departments).Error

	return departments, err
}

func (r *departmentRepository) GetWithRelations(ctx context.Context, id uuid.UUID) (*domain.Department, error) {
	var department domain.Department
	err := r.db.WithContext(ctx).
		Preload("ParentDept").
		Preload("EscalationDept").
		First(&department, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &department, nil
}

func (r *departmentRepository) GetByParent(ctx context.Context, parentDeptID uuid.UUID) ([]*domain.Department, error) {
	var departments []*domain.Department
	err := r.db.WithContext(ctx).
		Where("parent_dept_id = ?", parentDeptID).
		Find(&departments).Error
	return departments, err
}

func (r *departmentRepository) GetBySupportLevel(ctx context.Context, supportLevel int) ([]*domain.Department, error) {
	var departments []*domain.Department
	err := r.db.WithContext(ctx).
		Where("support_level = ? AND is_active = true", supportLevel).
		Find(&departments).Error
	return departments, err
}

func (r *departmentRepository) Update(ctx context.Context, department *domain.Department) error {
	return r.db.WithContext(ctx).
		Model(&domain.Department{}).
		Where("id = ?", department.ID).
		Updates(department).Error
}

func (r *departmentRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&domain.Department{}, "id = ?", id).Error
}

func (r *departmentRepository) Count(ctx context.Context, search string, isActive *bool, supportLevel *int, parentDeptID string) (int64, error) {
	var count int64
	query := r.db.WithContext(ctx).Model(&domain.Department{})

	// Apply same filters as GetWithPagination
	if search != "" {
		query = query.Where("name ILIKE ? OR description ILIKE ?", "%"+search+"%", "%"+search+"%")
	}

	if isActive != nil {
		query = query.Where("is_active = ?", *isActive)
	}

	if supportLevel != nil {
		query = query.Where("support_level = ?", *supportLevel)
	}

	if parentDeptID != "" {
		query = query.Where("parent_dept_id = ?", parentDeptID)
	}

	err := query.Count(&count).Error
	return count, err
}
