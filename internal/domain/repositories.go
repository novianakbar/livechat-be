package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// Usecase interfaces
type OSSChatUsecase interface {
	StartChat(ctx context.Context, req *StartChatRequest, ipAddress string) (*StartChatResponse, error)
	SetSessionContact(ctx context.Context, req *SetSessionContactRequest) (*SetSessionContactResponse, error)
	LinkOSSUser(ctx context.Context, req *LinkOSSUserRequest) (*LinkOSSUserResponse, error)
	GetChatHistory(ctx context.Context, req *GetChatHistoryRequest) (*GetChatHistoryResponse, error)
}

// Repository interfaces
// UserRepository interface for user operations
type UserRepository interface {
	Create(ctx context.Context, user *User) error
	GetByID(ctx context.Context, id uuid.UUID) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetAgentsByDepartment(ctx context.Context, departmentID uuid.UUID) ([]*User, error)
	GetAvailableAgents(ctx context.Context, departmentID *uuid.UUID) ([]*User, error)
	GetWithPagination(ctx context.Context, offset, limit int, role string, departmentID *uuid.UUID) ([]*User, error)
	Count(ctx context.Context, role string, departmentID *uuid.UUID) (int, error)
	GetByRole(ctx context.Context, role string) ([]*User, error)
	// Analytics methods
	CountByRole(ctx context.Context, role string) (int64, error)
	CountOnlineAgents(ctx context.Context) (int64, error)
}

// DepartmentRepository interface for department operations
type DepartmentRepository interface {
	Create(ctx context.Context, department *Department) error
	GetByID(ctx context.Context, id uuid.UUID) (*Department, error)
	GetAll(ctx context.Context) ([]*Department, error)
	Update(ctx context.Context, department *Department) error
	Delete(ctx context.Context, id uuid.UUID) error
}

// ChatUserRepository interface for chat user operations
type ChatUserRepository interface {
	Create(ctx context.Context, user *ChatUser) error
	GetByID(ctx context.Context, id uuid.UUID) (*ChatUser, error)
	GetByBrowserUUID(ctx context.Context, browserUUID uuid.UUID) (*ChatUser, error)
	GetByOSSUserID(ctx context.Context, ossUserID string) (*ChatUser, error)
	GetByEmail(ctx context.Context, email string) (*ChatUser, error)
	Update(ctx context.Context, user *ChatUser) error
	Delete(ctx context.Context, id uuid.UUID) error
	LinkOSSUser(ctx context.Context, browserUUID uuid.UUID, ossUserID string, email string) error
	List(ctx context.Context, limit, offset int) ([]*ChatUser, error)
	Count(ctx context.Context) (int64, error)
}

// ChatSessionContactRepository interface for chat session contact operations
type ChatSessionContactRepository interface {
	Create(ctx context.Context, contact *ChatSessionContact) error
	GetBySessionID(ctx context.Context, sessionID uuid.UUID) (*ChatSessionContact, error)
	Update(ctx context.Context, contact *ChatSessionContact) error
	Delete(ctx context.Context, sessionID uuid.UUID) error
}

// ChatSessionRepository interface for chat session operations
type ChatSessionRepository interface {
	Create(ctx context.Context, session *ChatSession) error
	GetByID(ctx context.Context, id uuid.UUID) (*ChatSession, error)
	GetByChatUserID(ctx context.Context, chatUserID uuid.UUID) ([]*ChatSession, error)
	GetByAgentID(ctx context.Context, agentID uuid.UUID) ([]*ChatSession, error)
	GetActiveSessions(ctx context.Context) ([]*ChatSession, error)
	GetWaitingSessions(ctx context.Context) ([]*ChatSession, error)
	Update(ctx context.Context, session *ChatSession) error
	Close(ctx context.Context, sessionID uuid.UUID) error
	GetSessionsByStatus(ctx context.Context, status string) ([]*ChatSession, error)
	GetSessionsByDateRange(ctx context.Context, start, end time.Time) ([]*ChatSession, error)
	GetWithPagination(ctx context.Context, offset, limit int, status string, agentID, departmentID *uuid.UUID) ([]*ChatSession, error)
	GetSessionsWithMessages(ctx context.Context, chatUserID uuid.UUID, limit, offset int) ([]*ChatSession, error)
	GetSessionHistory(ctx context.Context, chatUserID uuid.UUID, limit, offset int) ([]*ChatSession, error)
	Count(ctx context.Context, status string, agentID, departmentID *uuid.UUID) (int, error)
	// Analytics methods
	CountByStatus(ctx context.Context, status string) (int64, error)
	CountCompletedSince(ctx context.Context, since time.Time) (int64, error)
	GetAverageResponseTime(ctx context.Context) (float64, error)
	GetOSSCategoriesStats(ctx context.Context) ([]CategoryStats, error)
}

// ChatMessageRepository interface for chat message operations
type ChatMessageRepository interface {
	Create(ctx context.Context, message *ChatMessage) error
	GetByID(ctx context.Context, id uuid.UUID) (*ChatMessage, error)
	GetBySessionID(ctx context.Context, sessionID uuid.UUID) ([]*ChatMessage, error)
	Update(ctx context.Context, message *ChatMessage) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetMessagesByDateRange(ctx context.Context, start, end time.Time) ([]*ChatMessage, error)
	// Analytics methods
	GetTopQuestions(ctx context.Context, limit int) ([]QuestionStats, error)
}

// ChatLogRepository interface for chat log operations
type ChatLogRepository interface {
	Create(ctx context.Context, log *ChatLog) error
	GetBySessionID(ctx context.Context, sessionID uuid.UUID) ([]*ChatLog, error)
	GetByDateRange(ctx context.Context, start, end time.Time) ([]*ChatLog, error)
}

// ChatTagRepository interface for chat tag operations
type ChatTagRepository interface {
	Create(ctx context.Context, tag *ChatTag) error
	GetByID(ctx context.Context, id uuid.UUID) (*ChatTag, error)
	GetAll(ctx context.Context) ([]*ChatTag, error)
	Update(ctx context.Context, tag *ChatTag) error
	Delete(ctx context.Context, id uuid.UUID) error
}

// ChatSessionTagRepository interface for chat session tag operations
type ChatSessionTagRepository interface {
	Create(ctx context.Context, sessionTag *ChatSessionTag) error
	GetBySessionID(ctx context.Context, sessionID uuid.UUID) ([]*ChatSessionTag, error)
	DeleteBySessionID(ctx context.Context, sessionID uuid.UUID) error
	DeleteByTagID(ctx context.Context, tagID uuid.UUID) error
}

// AgentStatusRepository interface for agent status operations
type AgentStatusRepository interface {
	Create(ctx context.Context, status *AgentStatus) error
	GetByAgentID(ctx context.Context, agentID uuid.UUID) (*AgentStatus, error)
	Update(ctx context.Context, status *AgentStatus) error
	GetOnlineAgents(ctx context.Context) ([]*AgentStatus, error)
	UpdateLastActive(ctx context.Context, agentID uuid.UUID) error
}

// ChatAnalyticsRepository interface for chat analytics operations
type ChatAnalyticsRepository interface {
	Create(ctx context.Context, analytics *ChatAnalytics) error
	GetByDateRange(ctx context.Context, start, end time.Time) ([]*ChatAnalytics, error)
	GetByAgentAndDateRange(ctx context.Context, agentID uuid.UUID, start, end time.Time) ([]*ChatAnalytics, error)
	GetByDepartmentAndDateRange(ctx context.Context, departmentID uuid.UUID, start, end time.Time) ([]*ChatAnalytics, error)
	UpdateOrCreate(ctx context.Context, analytics *ChatAnalytics) error
}

// EmailService interface for email operations
type EmailService interface {
	SendEmail(ctx context.Context, request *SendEmailRequest) (*EmailResponse, error)
	SendTemplatedEmail(ctx context.Context, template *EmailTemplate, to []string, variables map[string]string) (*EmailResponse, error)
	SendWelcomeEmail(ctx context.Context, to string, name string) (*EmailResponse, error)
	SendPasswordResetEmail(ctx context.Context, to string, resetToken string) (*EmailResponse, error)
	SendChatTranscriptEmail(ctx context.Context, to string, transcript string, sessionID uuid.UUID) (*EmailResponse, error)
}
