package mappers

import (
	"github.com/novianakbar/livechat-be/internal/domain"
	"github.com/novianakbar/livechat-be/internal/models"
)

// TicketMapper handles conversion between domain entities and response models
type TicketMapper struct{}

// NewTicketMapper creates a new ticket mapper
func NewTicketMapper() *TicketMapper {
	return &TicketMapper{}
}

// ToTicketResponse converts domain.Ticket to models.TicketResponse
func (m *TicketMapper) ToTicketResponse(ticket *domain.Ticket) *models.TicketResponse {
	if ticket == nil {
		return nil
	}

	response := &models.TicketResponse{
		ID:            ticket.ID,
		TicketCode:    ticket.TicketCode,
		Subject:       ticket.Subject,
		Description:   ticket.Description,
		CustomerName:  ticket.CustomerName,
		CustomerEmail: ticket.CustomerEmail,
		CustomerPhone: ticket.CustomerPhone,
		Priority:      ticket.Priority,
		Status:        ticket.Status,
		CreatedVia:    ticket.CreatedVia,
		AccessToken:   ticket.AccessToken,
		CreatedAt:     ticket.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:     ticket.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	// Optional fields
	if ticket.CategoryID.Valid {
		response.CategoryID = ticket.CategoryID.String
	}
	if ticket.AssignedToID.Valid {
		response.AssignedToID = ticket.AssignedToID.String
	}
	if ticket.DepartmentID.Valid {
		response.DepartmentID = ticket.DepartmentID.String
	}
	if ticket.CreatedByID.Valid {
		response.CreatedByID = ticket.CreatedByID.String
	}
	if ticket.FirstResponseAt.Valid {
		response.FirstResponseAt = ticket.FirstResponseAt.Time.Format("2006-01-02T15:04:05Z07:00")
	}
	if ticket.ResolvedAt.Valid {
		response.ResolvedAt = ticket.ResolvedAt.Time.Format("2006-01-02T15:04:05Z07:00")
	}
	if ticket.ClosedAt.Valid {
		response.ClosedAt = ticket.ClosedAt.Time.Format("2006-01-02T15:04:05Z07:00")
	}

	return response
}

// ToTicketCategoryResponse converts domain.TicketCategory to models.TicketCategoryResponse
func (m *TicketMapper) ToTicketCategoryResponse(category *domain.TicketCategory) *models.TicketCategoryResponse {
	if category == nil {
		return nil
	}

	response := &models.TicketCategoryResponse{
		ID:               category.ID,
		Name:             category.Name,
		Code:             category.Code,
		Description:      category.Description,
		Color:            category.Color,
		SLAFirstResponse: category.SLAFirstResponse,
		SLAResolution:    category.SLAResolution,
		IsActive:         category.IsActive,
		CreatedAt:        category.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:        category.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	// Optional fields
	if category.DefaultDepartmentID.Valid {
		response.DefaultDepartmentID = category.DefaultDepartmentID.String
	}

	return response
}

// ToTicketCommentResponse converts domain.TicketComment to models.TicketCommentResponse
func (m *TicketMapper) ToTicketCommentResponse(comment *domain.TicketComment) *models.TicketCommentResponse {
	if comment == nil {
		return nil
	}

	response := &models.TicketCommentResponse{
		ID:             comment.ID,
		TicketID:       comment.TicketID,
		Content:        comment.Content,
		IsInternal:     comment.IsInternal,
		IsFromCustomer: comment.IsFromCustomer,
		CreatedAt:      comment.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:      comment.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	// Optional fields
	if comment.UserID.Valid {
		response.UserID = comment.UserID.String
	}

	return response
}

// ToTicketAttachmentResponse converts domain.TicketAttachment to models.TicketAttachmentResponse
func (m *TicketMapper) ToTicketAttachmentResponse(attachment *domain.TicketAttachment) *models.TicketAttachmentResponse {
	if attachment == nil {
		return nil
	}

	response := &models.TicketAttachmentResponse{
		ID:          attachment.ID,
		TicketID:    attachment.TicketID,
		FileName:    attachment.FileName,
		FilePath:    attachment.FilePath,
		FileSize:    attachment.FileSize,
		ContentType: attachment.FileType,
		UploadedBy:  attachment.UploadedBy,
		CreatedAt:   attachment.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	return response
}

// ToTicketHistoryResponse converts domain.TicketHistory to models.TicketHistoryResponse
func (m *TicketMapper) ToTicketHistoryResponse(history *domain.TicketHistory) *models.TicketHistoryResponse {
	if history == nil {
		return nil
	}

	response := &models.TicketHistoryResponse{
		ID:          history.ID,
		TicketID:    history.TicketID,
		Action:      history.Action,
		OldValue:    history.OldValue,
		NewValue:    history.NewValue,
		Description: history.Description,
		CreatedAt:   history.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	// Optional fields
	if history.UserID.Valid {
		response.UserID = history.UserID.String
	}

	return response
}

// ToTicketSLAResponse converts domain.TicketSLA to models.TicketSLAResponse
func (m *TicketMapper) ToTicketSLAResponse(sla *domain.TicketSLA) *models.TicketSLAResponse {
	if sla == nil {
		return nil
	}

	response := &models.TicketSLAResponse{
		ID:                      sla.ID,
		TicketID:                sla.TicketID,
		FirstResponseDue:        sla.FirstResponseDue.Format("2006-01-02T15:04:05Z07:00"),
		ResolutionDue:           sla.ResolutionDue.Format("2006-01-02T15:04:05Z07:00"),
		IsFirstResponseBreached: sla.FirstResponseBreached,
		IsResolutionBreached:    sla.ResolutionBreached,
		CreatedAt:               sla.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:               sla.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	// Optional fields
	if sla.FirstResponseAt.Valid {
		response.FirstResponseAt = sla.FirstResponseAt.Time.Format("2006-01-02T15:04:05Z07:00")
	}
	if sla.ResolvedAt.Valid {
		response.ResolutionAt = sla.ResolvedAt.Time.Format("2006-01-02T15:04:05Z07:00")
	}

	return response
}

// ToTicketListResponse converts slice of domain.Ticket to slice of models.TicketResponse
func (m *TicketMapper) ToTicketListResponse(tickets []*domain.Ticket) []models.TicketResponse {
	if tickets == nil {
		return []models.TicketResponse{}
	}

	responses := make([]models.TicketResponse, len(tickets))
	for i, ticket := range tickets {
		if response := m.ToTicketResponse(ticket); response != nil {
			responses[i] = *response
		}
	}

	return responses
}

// ToTicketCategoryListResponse converts slice of domain.TicketCategory to slice of models.TicketCategoryResponse
func (m *TicketMapper) ToTicketCategoryListResponse(categories []*domain.TicketCategory) []models.TicketCategoryResponse {
	if categories == nil {
		return []models.TicketCategoryResponse{}
	}

	responses := make([]models.TicketCategoryResponse, len(categories))
	for i, category := range categories {
		if response := m.ToTicketCategoryResponse(category); response != nil {
			responses[i] = *response
		}
	}

	return responses
}

// ToTicketDetailResponse converts domain.Ticket with relations to models.TicketDetailResponse
func (m *TicketMapper) ToTicketDetailResponse(ticket *domain.Ticket, comments []*domain.TicketComment, attachments []*domain.TicketAttachment, history []*domain.TicketHistory, sla *domain.TicketSLA) *models.TicketDetailResponse {
	if ticket == nil {
		return nil
	}

	response := &models.TicketDetailResponse{
		TicketResponse: m.ToTicketResponse(ticket),
	}

	// Convert comments
	if comments != nil {
		response.Comments = make([]models.TicketCommentResponse, len(comments))
		for i, comment := range comments {
			if commentResponse := m.ToTicketCommentResponse(comment); commentResponse != nil {
				response.Comments[i] = *commentResponse
			}
		}
	}

	// Convert attachments
	if attachments != nil {
		response.Attachments = make([]models.TicketAttachmentResponse, len(attachments))
		for i, attachment := range attachments {
			if attachmentResponse := m.ToTicketAttachmentResponse(attachment); attachmentResponse != nil {
				response.Attachments[i] = *attachmentResponse
			}
		}
	}

	// Convert history
	if history != nil {
		response.History = make([]models.TicketHistoryResponse, len(history))
		for i, hist := range history {
			if historyResponse := m.ToTicketHistoryResponse(hist); historyResponse != nil {
				response.History[i] = *historyResponse
			}
		}
	}

	// Convert SLA
	response.SLA = m.ToTicketSLAResponse(sla)

	return response
}
