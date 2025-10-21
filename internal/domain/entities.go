package domain

// Re-export entities from shared package for backward compatibility
import (
	"github.com/novianakbar/livechat-shared/entities"
)

// Re-export all entities
type User = entities.User
type Department = entities.Department
type ChatUser = entities.ChatUser
type ChatSession = entities.ChatSession
type ChatSessionContact = entities.ChatSessionContact
type ChatMessage = entities.ChatMessage
type ChatLog = entities.ChatLog
type ChatTag = entities.ChatTag
type ChatSessionTag = entities.ChatSessionTag
type AgentStatus = entities.AgentStatus
type ChatAnalytics = entities.ChatAnalytics
