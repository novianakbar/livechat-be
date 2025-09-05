package repository

import (
	"context"

	"github.com/novianakbar/livechat-be/internal/domain"
	"gorm.io/gorm"
)

type ticketCategoryRepository struct {
	db *gorm.DB
}

func NewTicketCategoryRepository(db *gorm.DB) domain.TicketCategoryRepository {
	return &ticketCategoryRepository{db: db}
}

func (r *ticketCategoryRepository) Create(ctx context.Context, category *domain.TicketCategory) error {
	return r.db.WithContext(ctx).Create(category).Error
}

func (r *ticketCategoryRepository) GetByID(ctx context.Context, id string) (*domain.TicketCategory, error) {
	var category domain.TicketCategory
	err := r.db.WithContext(ctx).
		Preload("DefaultDepartment").
		First(&category, "id = ?", id).Error

	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *ticketCategoryRepository) GetByCode(ctx context.Context, code string) (*domain.TicketCategory, error) {
	var category domain.TicketCategory
	err := r.db.WithContext(ctx).
		Preload("DefaultDepartment").
		First(&category, "code = ?", code).Error

	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *ticketCategoryRepository) GetAll(ctx context.Context) ([]*domain.TicketCategory, error) {
	var categories []*domain.TicketCategory
	err := r.db.WithContext(ctx).
		Preload("DefaultDepartment").
		Where("is_active = ?", true).
		Order("name ASC").
		Find(&categories).Error

	return categories, err
}

func (r *ticketCategoryRepository) GetByDepartment(ctx context.Context, deptID string) ([]*domain.TicketCategory, error) {
	var categories []*domain.TicketCategory
	err := r.db.WithContext(ctx).
		Preload("DefaultDepartment").
		Where("department_id = ? AND is_active = ?", deptID, true).
		Order("name ASC").
		Find(&categories).Error

	return categories, err
}

func (r *ticketCategoryRepository) GetActive(ctx context.Context) ([]*domain.TicketCategory, error) {
	var categories []*domain.TicketCategory
	err := r.db.WithContext(ctx).
		Preload("DefaultDepartment").
		Where("is_active = ?", true).
		Order("name ASC").
		Find(&categories).Error

	return categories, err
}

func (r *ticketCategoryRepository) Update(ctx context.Context, category *domain.TicketCategory) error {
	return r.db.WithContext(ctx).Save(category).Error
}

func (r *ticketCategoryRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&domain.TicketCategory{}, "id = ?", id).Error
}
