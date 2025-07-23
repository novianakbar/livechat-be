package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User represents system user (admin, agent, etc.)
type User struct {
	ID           uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Email        string         `json:"email" gorm:"uniqueIndex;not null"`
	Password     string         `json:"-" gorm:"not null"`
	Name         string         `json:"name" gorm:"not null"`
	Role         string         `json:"role" gorm:"not null"` // admin, agent
	IsActive     bool           `json:"is_active" gorm:"default:true"`
	DepartmentID *uuid.UUID     `json:"department_id" gorm:"type:uuid"`
	Department   *Department    `json:"department,omitempty" gorm:"foreignKey:DepartmentID"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

// Department represents department for agents
type Department struct {
	ID          uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Name        string         `json:"name" gorm:"not null"`
	Description string         `json:"description"`
	IsActive    bool           `json:"is_active" gorm:"default:true"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

// Customer represents chat customer
type Customer struct {
	ID          uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	CompanyName string         `json:"company_name" gorm:"not null"`
	PersonName  string         `json:"person_name" gorm:"not null"`
	Email       string         `json:"email" gorm:"not null"`
	IPAddress   string         `json:"ip_address" gorm:"not null"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

// ChatSession represents chat session
type ChatSession struct {
	ID           uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	CustomerID   uuid.UUID      `json:"customer_id" gorm:"type:uuid;not null"`
	Customer     Customer       `json:"customer" gorm:"foreignKey:CustomerID"`
	AgentID      *uuid.UUID     `json:"agent_id" gorm:"type:uuid"`
	Agent        *User          `json:"agent,omitempty" gorm:"foreignKey:AgentID"`
	DepartmentID *uuid.UUID     `json:"department_id" gorm:"type:uuid"`
	Department   *Department    `json:"department,omitempty" gorm:"foreignKey:DepartmentID"`
	Topic        string         `json:"topic" gorm:"not null"`
	Status       string         `json:"status" gorm:"not null;default:'waiting'"` // waiting, active, closed
	Priority     string         `json:"priority" gorm:"default:'normal'"`         // low, normal, high, urgent
	StartedAt    time.Time      `json:"started_at"`
	EndedAt      *time.Time     `json:"ended_at"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

// ChatMessage represents chat message
type ChatMessage struct {
	ID          uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	SessionID   uuid.UUID      `json:"session_id" gorm:"type:uuid;not null"`
	Session     ChatSession    `json:"session" gorm:"foreignKey:SessionID"`
	SenderID    *uuid.UUID     `json:"sender_id" gorm:"type:uuid"`
	Sender      *User          `json:"sender,omitempty" gorm:"foreignKey:SenderID"`
	SenderType  string         `json:"sender_type" gorm:"not null"` // customer, agent, system
	Message     string         `json:"message" gorm:"not null"`
	MessageType string         `json:"message_type" gorm:"default:'text'"` // text, image, file, system
	Attachments []string       `json:"attachments" gorm:"type:json"`
	ReadAt      *time.Time     `json:"read_at"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

// ChatLog represents chat activity log
type ChatLog struct {
	ID        uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	SessionID uuid.UUID      `json:"session_id" gorm:"type:uuid;not null"`
	Session   ChatSession    `json:"session" gorm:"foreignKey:SessionID"`
	Action    string         `json:"action" gorm:"not null"` // started, waiting, response, closed, transferred
	Details   string         `json:"details"`
	UserID    *uuid.UUID     `json:"user_id" gorm:"type:uuid"`
	User      *User          `json:"user,omitempty" gorm:"foreignKey:UserID"`
	CreatedAt time.Time      `json:"created_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

// ChatTag represents chat tags
type ChatTag struct {
	ID        uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Name      string         `json:"name" gorm:"uniqueIndex;not null"`
	Color     string         `json:"color" gorm:"default:'#007bff'"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

// ChatSessionTag represents many-to-many relationship between sessions and tags
type ChatSessionTag struct {
	ID        uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	SessionID uuid.UUID      `json:"session_id" gorm:"type:uuid;not null"`
	Session   ChatSession    `json:"session" gorm:"foreignKey:SessionID"`
	TagID     uuid.UUID      `json:"tag_id" gorm:"type:uuid;not null"`
	Tag       ChatTag        `json:"tag" gorm:"foreignKey:TagID"`
	CreatedAt time.Time      `json:"created_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

// AgentStatus represents agent online status
type AgentStatus struct {
	ID           uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	AgentID      uuid.UUID      `json:"agent_id" gorm:"type:uuid;not null"`
	Agent        User           `json:"agent" gorm:"foreignKey:AgentID"`
	Status       string         `json:"status" gorm:"not null"` // online, offline, busy, away
	LastActiveAt time.Time      `json:"last_active_at"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

// ChatAnalytics represents chat analytics
type ChatAnalytics struct {
	ID                  uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Date                time.Time      `json:"date" gorm:"not null"`
	TotalSessions       int            `json:"total_sessions" gorm:"default:0"`
	CompletedSessions   int            `json:"completed_sessions" gorm:"default:0"`
	AverageResponseTime float64        `json:"average_response_time" gorm:"default:0"` // in seconds
	TotalMessages       int            `json:"total_messages" gorm:"default:0"`
	DepartmentID        *uuid.UUID     `json:"department_id" gorm:"type:uuid"`
	Department          *Department    `json:"department,omitempty" gorm:"foreignKey:DepartmentID"`
	AgentID             *uuid.UUID     `json:"agent_id" gorm:"type:uuid"`
	Agent               *User          `json:"agent,omitempty" gorm:"foreignKey:AgentID"`
	CreatedAt           time.Time      `json:"created_at"`
	UpdatedAt           time.Time      `json:"updated_at"`
	DeletedAt           gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}
