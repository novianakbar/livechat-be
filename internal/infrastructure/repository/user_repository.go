package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/novianakbar/livechat-be/internal/domain"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) domain.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *domain.User) error {
	if err := r.db.WithContext(ctx).Create(user).Error; err != nil {
		return err
	}
	return nil
}

func (r *userRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	var user domain.User
	if err := r.db.WithContext(ctx).Preload("Department").First(&user, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User
	if err := r.db.WithContext(ctx).Preload("Department").First(&user, "email = ?", email).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Update(ctx context.Context, user *domain.User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

func (r *userRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&domain.User{}, "id = ?", id).Error
}

func (r *userRepository) GetAgentsByDepartment(ctx context.Context, departmentID uuid.UUID) ([]*domain.User, error) {
	var users []*domain.User
	if err := r.db.WithContext(ctx).
		Preload("Department").
		Where("department_id = ? AND role = ? AND is_active = ?", departmentID, "agent", true).
		Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (r *userRepository) GetAvailableAgents(ctx context.Context, departmentID *uuid.UUID) ([]*domain.User, error) {
	query := r.db.WithContext(ctx).Where("role = ? AND is_active = ?", "agent", true)

	if departmentID != nil {
		query = query.Where("department_id = ?", *departmentID)
	}

	var users []*domain.User
	if err := query.Preload("Department").Find(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}

// Analytics methods
func (r *userRepository) CountByRole(ctx context.Context, role string) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&domain.User{}).Where("role = ?", role).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *userRepository) CountOnlineAgents(ctx context.Context) (int64, error) {
	var count int64
	// Assuming we have an agent_status table to track online status
	// For now, we'll return a mock count
	if err := r.db.WithContext(ctx).Model(&domain.User{}).Where("role = ? AND is_active = ?", "agent", true).Count(&count).Error; err != nil {
		return 0, err
	}
	// Assume 80% of active agents are online
	return int64(float64(count) * 0.8), nil
}

func (r *userRepository) GetWithPagination(ctx context.Context, offset, limit int, role string, departmentID *uuid.UUID) ([]*domain.User, error) {
	query := r.db.WithContext(ctx).Preload("Department")

	if role != "" {
		query = query.Where("role = ?", role)
	}
	if departmentID != nil {
		query = query.Where("department_id = ?", *departmentID)
	}

	var users []*domain.User
	if err := query.
		Where("is_active = ?", true).
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}

func (r *userRepository) Count(ctx context.Context, role string, departmentID *uuid.UUID) (int, error) {
	query := r.db.WithContext(ctx).Model(&domain.User{})

	if role != "" {
		query = query.Where("role = ?", role)
	}
	if departmentID != nil {
		query = query.Where("department_id = ?", *departmentID)
	}

	var count int64
	if err := query.Where("is_active = ?", true).Count(&count).Error; err != nil {
		return 0, err
	}

	return int(count), nil
}

func (r *userRepository) GetByRole(ctx context.Context, role string) ([]*domain.User, error) {
	var users []*domain.User
	if err := r.db.WithContext(ctx).
		Preload("Department").
		Where("role = ? AND is_active = ?", role, true).
		Order("name ASC").
		Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}
