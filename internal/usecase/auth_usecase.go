package usecase

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/novianakbar/livechat-be/internal/domain"
	"github.com/novianakbar/livechat-be/pkg/utils"
)

type AgentSessionRepository interface {
	SetAgentLoggedIn(ctx context.Context, agentID string) error
	SetAgentLoggedOut(ctx context.Context, agentID string) error
}

type AuthUsecase struct {
	userRepo         domain.UserRepository
	agentSessionRepo AgentSessionRepository
	jwtUtil          *utils.JWTUtil
}

func NewAuthUsecase(userRepo domain.UserRepository, agentSessionRepo AgentSessionRepository, jwtUtil *utils.JWTUtil) *AuthUsecase {
	return &AuthUsecase{
		userRepo:         userRepo,
		agentSessionRepo: agentSessionRepo,
		jwtUtil:          jwtUtil,
	}
}

func (uc *AuthUsecase) Login(ctx context.Context, req *domain.LoginRequest, clientIP, userAgent string) (*domain.LoginResponse, error) {
	// TODO: Implement rate limiting based on clientIP and userAgent

	// Find user by email
	user, err := uc.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.New("invalid email or password")
	}

	// Check if user is active
	if !user.IsActive {
		return nil, errors.New("user account is inactive")
	}

	// Verify password
	if !utils.CheckPasswordHash(req.Password, user.Password) {
		return nil, errors.New("invalid email or password")
	}

	// Generate JWT token pair
	var departmentID *string
	if user.DepartmentID.Valid {
		departmentID = &user.DepartmentID.String
	}
	tokenPair, err := uc.jwtUtil.GenerateTokenPair(user.ID, user.Email, user.Role, departmentID)
	if err != nil {
		return nil, err
	}

	// Track agent login in database if user is agent or admin
	if user.Role == "agent" || user.Role == "admin" {
		if err := uc.agentSessionRepo.SetAgentLoggedIn(ctx, user.ID); err != nil {
			// Log error but don't fail login process
			// TODO: Add proper logging
		}
	}

	// TODO: Store session info in Redis/database for revocation support

	// Hide password from response
	user.Password = ""

	return &domain.LoginResponse{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		ExpiresIn:    tokenPair.ExpiresIn,
		ExpiresAt:    tokenPair.ExpiresAt,
		User:         user,
	}, nil
}

func (uc *AuthUsecase) Logout(ctx context.Context, userID string, accessToken, refreshToken string) error {
	// Get user info to check role
	user, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	// Track agent logout in database if user is agent or admin
	if user != nil && (user.Role == "agent" || user.Role == "admin") {
		if err := uc.agentSessionRepo.SetAgentLoggedOut(ctx, userID); err != nil {
			// Log error but don't fail logout process
			// TODO: Add proper logging
		}
	}

	// TODO: Implement token blacklisting
	// For now, we'll just validate that the tokens exist and are from the right user

	if accessToken != "" {
		claims, err := uc.jwtUtil.ValidateAccessToken(accessToken)
		if err == nil && claims.UserID == userID {
			// TODO: Add token to blacklist in Redis
		}
	}

	if refreshToken != "" {
		claims, err := uc.jwtUtil.ValidateRefreshToken(refreshToken)
		if err == nil && claims.UserID == userID {
			// TODO: Add token to blacklist in Redis
		}
	}

	return nil
}

func (uc *AuthUsecase) Register(ctx context.Context, req *domain.RegisterRequest) (*domain.User, error) {
	// Check if user already exists
	existingUser, err := uc.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}

	if existingUser != nil {
		return nil, errors.New("user with this email already exists")
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	// Create user
	uuidV7, _ := uuid.NewV7()
	var departmentID sql.NullString
	if req.DepartmentID != nil {
		departmentID = sql.NullString{
			String: req.DepartmentID.String(),
			Valid:  true,
		}
	}

	user := &domain.User{
		ID:           uuidV7.String(),
		Email:        req.Email,
		Password:     hashedPassword,
		Name:         req.Name,
		Role:         req.Role,
		IsActive:     true,
		DepartmentID: departmentID,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := uc.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	// Hide password from response
	user.Password = ""

	return user, nil
}

func (uc *AuthUsecase) GetUserByID(ctx context.Context, userID string) (*domain.User, error) {
	user, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.New("user not found")
	}

	// Hide password from response
	user.Password = ""

	return user, nil
}

func (uc *AuthUsecase) RefreshToken(ctx context.Context, req *domain.RefreshTokenRequest) (*domain.RefreshTokenResponse, error) {
	// Validate refresh token
	claims, err := uc.jwtUtil.ValidateRefreshToken(req.RefreshToken)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	// Check if user still exists and is active
	user, err := uc.userRepo.GetByID(ctx, claims.UserID)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.New("user not found")
	}

	if !user.IsActive {
		return nil, errors.New("user account is inactive")
	}

	// Generate new token pair
	var departmentID *string
	if user.DepartmentID.Valid {
		departmentID = &user.DepartmentID.String
	}
	tokenPair, err := uc.jwtUtil.GenerateTokenPair(user.ID, user.Email, user.Role, departmentID)
	if err != nil {
		return nil, err
	}

	return &domain.RefreshTokenResponse{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		ExpiresIn:    tokenPair.ExpiresIn,
		ExpiresAt:    tokenPair.ExpiresAt,
	}, nil
}

func (uc *AuthUsecase) ValidateToken(ctx context.Context, tokenString string) (*domain.User, error) {
	claims, err := uc.jwtUtil.ValidateAccessToken(tokenString)
	if err != nil {
		return nil, err
	}

	user, err := uc.userRepo.GetByID(ctx, claims.UserID)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.New("user not found")
	}

	if !user.IsActive {
		return nil, errors.New("user account is inactive")
	}

	// Hide password from response
	user.Password = ""

	return user, nil
}
