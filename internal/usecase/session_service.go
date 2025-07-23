package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// SessionService handles user session management
// This is a simple in-memory implementation for now
// In production, this should be backed by Redis
type SessionService struct {
	sessions map[string]SessionInfo
}

type SessionInfo struct {
	UserID    uuid.UUID
	Email     string
	Role      string
	LoginTime time.Time
	LastSeen  time.Time
	ClientIP  string
	UserAgent string
}

func NewSessionService() *SessionService {
	return &SessionService{
		sessions: make(map[string]SessionInfo),
	}
}

func (s *SessionService) CreateSession(ctx context.Context, userID uuid.UUID, email, role, clientIP, userAgent string) string {
	sessionID := uuid.New().String()
	s.sessions[sessionID] = SessionInfo{
		UserID:    userID,
		Email:     email,
		Role:      role,
		LoginTime: time.Now(),
		LastSeen:  time.Now(),
		ClientIP:  clientIP,
		UserAgent: userAgent,
	}
	return sessionID
}

func (s *SessionService) ValidateSession(ctx context.Context, sessionID string) (*SessionInfo, error) {
	session, exists := s.sessions[sessionID]
	if !exists {
		return nil, fmt.Errorf("session not found")
	}

	// Update last seen
	session.LastSeen = time.Now()
	s.sessions[sessionID] = session

	return &session, nil
}

func (s *SessionService) InvalidateSession(ctx context.Context, sessionID string) error {
	delete(s.sessions, sessionID)
	return nil
}

func (s *SessionService) InvalidateUserSessions(ctx context.Context, userID uuid.UUID) error {
	for sessionID, session := range s.sessions {
		if session.UserID == userID {
			delete(s.sessions, sessionID)
		}
	}
	return nil
}

func (s *SessionService) GetActiveSessions(ctx context.Context, userID uuid.UUID) ([]SessionInfo, error) {
	var sessions []SessionInfo
	for _, session := range s.sessions {
		if session.UserID == userID {
			sessions = append(sessions, session)
		}
	}
	return sessions, nil
}
