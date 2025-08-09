package mappers

import (
	"math"
	"time"

	"github.com/novianakbar/livechat-be/internal/models"
	"github.com/novianakbar/livechat-shared/entities"
)

// FormatTime formats time to RFC3339 string
func FormatTime(t time.Time) string {
	return t.Format(time.RFC3339)
}

// FormatTimePtr formats time pointer to RFC3339 string, returns empty string if nil
func FormatTimePtr(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.Format(time.RFC3339)
}

// ChatUserToResponse converts ChatUser entity to ChatUserResponse
func ChatUserToResponse(entity *entities.ChatUser) *models.ChatUserResponse {
	if entity == nil {
		return nil
	}

	return &models.ChatUserResponse{
		ID:          entity.ID,
		BrowserUUID: entity.BrowserUUID.String, // Handle sql.NullString
		OSSUserID:   entity.OSSUserID.String,   // Handle sql.NullString
		Email:       entity.Email.String,       // Handle sql.NullString
		IsAnonymous: entity.IsAnonymous,
		IPAddress:   entity.IPAddress,
		UserAgent:   entity.UserAgent.String, // Handle sql.NullString
		CreatedAt:   FormatTime(entity.CreatedAt),
		UpdatedAt:   FormatTime(entity.UpdatedAt),
	}
}

// ChatSessionContactToResponse converts ChatSessionContact entity to ChatSessionContactResponse
func ChatSessionContactToResponse(entity *entities.ChatSessionContact) *models.ChatSessionContactResponse {
	if entity == nil {
		return nil
	}

	return &models.ChatSessionContactResponse{
		ID:           entity.ID,
		SessionID:    entity.SessionID,
		ContactName:  entity.ContactName,
		ContactEmail: entity.ContactEmail,
		ContactPhone: entity.ContactPhone.String, // Handle sql.NullString
		Position:     entity.Position.String,     // Handle sql.NullString
		CompanyName:  entity.CompanyName.String,  // Handle sql.NullString
		CreatedAt:    FormatTime(entity.CreatedAt),
		UpdatedAt:    FormatTime(entity.UpdatedAt),
	}
}

// ChatMessageToResponse converts ChatMessage entity to ChatMessageResponse
func ChatMessageToResponse(entity *entities.ChatMessage) *models.ChatMessageResponse {
	if entity == nil {
		return nil
	}

	response := &models.ChatMessageResponse{
		ID:          entity.ID,
		SessionID:   entity.SessionID,
		SenderType:  entity.SenderType,
		Message:     entity.Message,
		MessageType: entity.MessageType,
		Attachments: entity.Attachments,
		CreatedAt:   FormatTime(entity.CreatedAt),
	}

	// Handle optional sender ID
	if entity.SenderID.Valid {
		response.SenderID = entity.SenderID.String
	}

	// Handle optional read time
	if entity.ReadAt.Valid {
		response.ReadAt = FormatTime(entity.ReadAt.Time)
	}

	return response
}

// ChatMessagesToResponse converts slice of ChatMessage entities to ChatMessageResponse slice
func ChatMessagesToResponse(entities []entities.ChatMessage) []models.ChatMessageResponse {
	if entities == nil {
		return nil
	}

	responses := make([]models.ChatMessageResponse, len(entities))
	for i, entity := range entities {
		response := ChatMessageToResponse(&entity)
		if response != nil {
			responses[i] = *response
		}
	}
	return responses
}

// ChatMessagePointersToResponse converts slice of ChatMessage entity pointers to ChatMessageResponse slice
func ChatMessagePointersToResponse(entities []*entities.ChatMessage) []models.ChatMessageResponse {
	if entities == nil {
		return nil
	}

	responses := make([]models.ChatMessageResponse, len(entities))
	for i, entity := range entities {
		response := ChatMessageToResponse(entity)
		if response != nil {
			responses[i] = *response
		}
	}
	return responses
}

// ChatSessionToMinimalResponse converts ChatSession entity to ChatSessionMinimalResponse
func ChatSessionToMinimalResponse(entity *entities.ChatSession) *models.ChatSessionMinimalResponse {
	if entity == nil {
		return nil
	}

	response := &models.ChatSessionMinimalResponse{
		ID:         entity.ID,
		ChatUserID: entity.ChatUserID,
		Topic:      entity.Topic,
		Status:     entity.Status,
		Priority:   entity.Priority,
		StartedAt:  FormatTime(entity.StartedAt),
		CreatedAt:  FormatTime(entity.CreatedAt),
		UpdatedAt:  FormatTime(entity.UpdatedAt),
	}

	// Handle optional agent ID
	if entity.AgentID.Valid {
		response.AgentID = entity.AgentID.String
	}

	// Handle optional end time
	if entity.EndedAt.Valid {
		response.EndedAt = FormatTime(entity.EndedAt.Time)
	}

	// Handle optional relations
	if entity.ChatUser.ID != "" {
		response.ChatUser = ChatUserToResponse(&entity.ChatUser)
	}

	if entity.Agent != nil {
		response.Agent = UserToResponse(entity.Agent)
	}

	return response
}

// ChatSessionToDetailResponse converts ChatSession entity to ChatSessionDetailResponse
func ChatSessionToDetailResponse(entity *entities.ChatSession) *models.ChatSessionDetailResponse {
	if entity == nil {
		return nil
	}

	response := &models.ChatSessionDetailResponse{
		ID:         entity.ID,
		ChatUserID: entity.ChatUserID,
		Topic:      entity.Topic,
		Status:     entity.Status,
		Priority:   entity.Priority,
		StartedAt:  FormatTime(entity.StartedAt),
		CreatedAt:  FormatTime(entity.CreatedAt),
		UpdatedAt:  FormatTime(entity.UpdatedAt),
	}

	// Handle optional agent ID
	if entity.AgentID.Valid {
		response.AgentID = entity.AgentID.String
	}

	// Handle optional department ID
	if entity.DepartmentID.Valid {
		response.DepartmentID = entity.DepartmentID.String
	}

	// Handle optional end time
	if entity.EndedAt.Valid {
		response.EndedAt = FormatTime(entity.EndedAt.Time)
	}

	// Handle optional relations
	if entity.ChatUser.ID != "" {
		response.ChatUser = ChatUserToResponse(&entity.ChatUser)
	}

	if entity.Agent != nil {
		response.Agent = UserToResponse(entity.Agent)
	}

	if entity.Department != nil {
		response.Department = DepartmentToResponse(entity.Department)
	}

	if entity.Messages != nil && len(entity.Messages) > 0 {
		response.Messages = ChatMessagesToResponse(entity.Messages)
	}

	if entity.Contact != nil {
		response.Contact = ChatSessionContactToResponse(entity.Contact)
	}

	return response
}

// ChatSessionsToMinimalResponse converts slice of ChatSession entities (or pointers) to ChatSessionMinimalResponse slice
func ChatSessionsToMinimalResponse(entities []*entities.ChatSession) []models.ChatSessionMinimalResponse {
	if entities == nil {
		return nil
	}

	responses := make([]models.ChatSessionMinimalResponse, len(entities))
	for i, entity := range entities {
		response := ChatSessionToMinimalResponse(entity)
		if response != nil {
			responses[i] = *response
		}
	}
	return responses
}

// ChatSessionsToDetailResponse converts slice of ChatSession entities (or pointers) to ChatSessionDetailResponse slice
func ChatSessionsToDetailResponse(entities []*entities.ChatSession) []models.ChatSessionDetailResponse {
	if entities == nil {
		return nil
	}

	responses := make([]models.ChatSessionDetailResponse, len(entities))
	for i, entity := range entities {
		response := ChatSessionToDetailResponse(entity)
		if response != nil {
			responses[i] = *response
		}
	}

	return responses
}

// CreatePaginatedResponse creates a paginated response
func CreatePaginatedResponse[T any](data []T, page, limit int, total int64) *models.PaginatedResponse[T] {
	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	return &models.PaginatedResponse[T]{
		Data:       data,
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
	}
}
