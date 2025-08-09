package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JWTClaims struct {
	UserID       string  `json:"user_id"`
	Email        string  `json:"email"`
	Role         string  `json:"role"`
	DepartmentID *string `json:"department_id"`
	TokenType    string  `json:"token_type"` // "access" or "refresh"
	jwt.RegisteredClaims
}

type TokenPair struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresIn    int64     `json:"expires_in"`
	ExpiresAt    time.Time `json:"expires_at"`
}

type JWTUtil struct {
	secret               []byte
	accessTokenDuration  time.Duration
	refreshTokenDuration time.Duration
}

func NewJWTUtil(secret string, accessTokenDuration, refreshTokenDuration time.Duration) *JWTUtil {
	return &JWTUtil{
		secret:               []byte(secret),
		accessTokenDuration:  accessTokenDuration,
		refreshTokenDuration: refreshTokenDuration,
	}
}

func (j *JWTUtil) GenerateTokenPair(userID string, email, role string, departmentID *string) (*TokenPair, error) {
	now := time.Now()

	// Generate Access Token (15 minutes)
	accessExpirationTime := now.Add(j.accessTokenDuration)
	accessClaims := &JWTClaims{
		UserID:       userID,
		Email:        email,
		Role:         role,
		DepartmentID: departmentID,
		TokenType:    "access",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(accessExpirationTime),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Subject:   userID,
			ID:        uuid.New().String(),
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenString, err := accessToken.SignedString(j.secret)
	if err != nil {
		return nil, err
	}

	// Generate Refresh Token (7 days)
	refreshExpirationTime := now.Add(j.refreshTokenDuration)
	refreshClaims := &JWTClaims{
		UserID:       userID,
		Email:        email,
		Role:         role,
		DepartmentID: departmentID,
		TokenType:    "refresh",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(refreshExpirationTime),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Subject:   userID,
			ID:        uuid.New().String(),
		},
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString(j.secret)
	if err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:  accessTokenString,
		RefreshToken: refreshTokenString,
		ExpiresIn:    int64(j.accessTokenDuration.Seconds()),
		ExpiresAt:    accessExpirationTime,
	}, nil
}

func (j *JWTUtil) ValidateAccessToken(tokenString string) (*JWTClaims, error) {
	claims, err := j.validateToken(tokenString)
	if err != nil {
		return nil, err
	}

	if claims.TokenType != "access" {
		return nil, errors.New("invalid token type: expected access token")
	}

	return claims, nil
}

func (j *JWTUtil) ValidateRefreshToken(tokenString string) (*JWTClaims, error) {
	claims, err := j.validateToken(tokenString)
	if err != nil {
		return nil, err
	}

	if claims.TokenType != "refresh" {
		return nil, errors.New("invalid token type: expected refresh token")
	}

	return claims, nil
}

func (j *JWTUtil) validateToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return j.secret, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	return claims, nil
}

func (j *JWTUtil) RefreshTokenPair(refreshTokenString string) (*TokenPair, error) {
	// Validate refresh token
	claims, err := j.ValidateRefreshToken(refreshTokenString)
	if err != nil {
		return nil, err
	}

	// Generate new token pair
	return j.GenerateTokenPair(claims.UserID, claims.Email, claims.Role, claims.DepartmentID)
}

// Legacy support - untuk backward compatibility
func (j *JWTUtil) GenerateToken(userID uuid.UUID, email, role string, departmentID *uuid.UUID) (string, time.Time, error) {
	var deptIDStr *string
	if departmentID != nil {
		deptStr := departmentID.String()
		deptIDStr = &deptStr
	}
	tokenPair, err := j.GenerateTokenPair(userID.String(), email, role, deptIDStr)
	if err != nil {
		return "", time.Time{}, err
	}
	return tokenPair.AccessToken, tokenPair.ExpiresAt, nil
}

func (j *JWTUtil) ValidateToken(tokenString string) (*JWTClaims, error) {
	return j.ValidateAccessToken(tokenString)
}

func (j *JWTUtil) RefreshToken(tokenString string) (string, time.Time, error) {
	tokenPair, err := j.RefreshTokenPair(tokenString)
	if err != nil {
		return "", time.Time{}, err
	}
	return tokenPair.AccessToken, tokenPair.ExpiresAt, nil
}
