package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/novianakbar/livechat-be/internal/domain"
	"gorm.io/gorm"
)

type chatUserRepository struct {
	db *gorm.DB
}

func NewChatUserRepository(db *gorm.DB) domain.ChatUserRepository {
	return &chatUserRepository{db: db}
}

func (r *chatUserRepository) Create(ctx context.Context, user *domain.ChatUser) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *chatUserRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.ChatUser, error) {
	var user domain.ChatUser
	if err := r.db.WithContext(ctx).First(&user, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("chat user not found")
		}
		return nil, err
	}
	return &user, nil
}

func (r *chatUserRepository) GetByBrowserUUID(ctx context.Context, browserUUID uuid.UUID) (*domain.ChatUser, error) {
	var user domain.ChatUser
	if err := r.db.WithContext(ctx).First(&user, "browser_uuid = ?", browserUUID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *chatUserRepository) GetByOSSUserID(ctx context.Context, ossUserID string) (*domain.ChatUser, error) {
	var user domain.ChatUser
	if err := r.db.WithContext(ctx).First(&user, "oss_user_id = ?", ossUserID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *chatUserRepository) GetByEmail(ctx context.Context, email string) (*domain.ChatUser, error) {
	var user domain.ChatUser
	if err := r.db.WithContext(ctx).First(&user, "email = ?", email).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *chatUserRepository) Update(ctx context.Context, user *domain.ChatUser) error {
	return r.db.WithContext(ctx).Save(user).Error
}

func (r *chatUserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&domain.ChatUser{}, "id = ?", id).Error
}

func (r *chatUserRepository) LinkOSSUser(ctx context.Context, browserUUID uuid.UUID, ossUserID string, email string) error {
	return r.db.WithContext(ctx).Model(&domain.ChatUser{}).
		Where("browser_uuid = ?", browserUUID).
		Updates(map[string]interface{}{
			"oss_user_id":  ossUserID,
			"email":        email,
			"is_anonymous": false,
		}).Error
}

func (r *chatUserRepository) List(ctx context.Context, limit, offset int) ([]*domain.ChatUser, error) {
	var users []*domain.ChatUser
	if err := r.db.WithContext(ctx).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (r *chatUserRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).
		Model(&domain.ChatUser{}).
		Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
