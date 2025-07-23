package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/novianakbar/livechat-be/internal/domain"
)

type UserUsecase struct {
	userRepo domain.UserRepository
}

func NewUserUsecase(userRepo domain.UserRepository) *UserUsecase {
	return &UserUsecase{
		userRepo: userRepo,
	}
}

func (uc *UserUsecase) GetUsers(ctx context.Context, page, limit int, role string, departmentID *uuid.UUID) ([]*domain.User, int, error) {
	offset := (page - 1) * limit
	users, err := uc.userRepo.GetWithPagination(ctx, offset, limit, role, departmentID)
	if err != nil {
		return nil, 0, err
	}

	total, err := uc.userRepo.Count(ctx, role, departmentID)
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func (uc *UserUsecase) GetUser(ctx context.Context, userID uuid.UUID) (*domain.User, error) {
	user, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (uc *UserUsecase) GetAgents(ctx context.Context) ([]*domain.User, error) {
	agents, err := uc.userRepo.GetByRole(ctx, "agent")
	if err != nil {
		return nil, err
	}

	return agents, nil
}
