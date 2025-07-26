package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/novianakbar/livechat-be/internal/domain"
)

type ChatUsecase struct {
	sessionRepo  domain.ChatSessionRepository
	messageRepo  domain.ChatMessageRepository
	userRepo     domain.UserRepository
	logRepo      domain.ChatLogRepository
	chatUserRepo domain.ChatUserRepository           // Added for OSS support
	contactRepo  domain.ChatSessionContactRepository // Added for OSS support
}

func NewChatUsecase(
	sessionRepo domain.ChatSessionRepository,
	messageRepo domain.ChatMessageRepository,
	userRepo domain.UserRepository,
	logRepo domain.ChatLogRepository,
	chatUserRepo domain.ChatUserRepository, // Added for OSS support
	contactRepo domain.ChatSessionContactRepository, // Added for OSS support
) *ChatUsecase {
	return &ChatUsecase{
		sessionRepo:  sessionRepo,
		messageRepo:  messageRepo,
		userRepo:     userRepo,
		logRepo:      logRepo,
		chatUserRepo: chatUserRepo,
		contactRepo:  contactRepo,
	}
}

func (uc *ChatUsecase) StartChat(ctx context.Context, req *domain.StartChatRequest, ipAddress string) (*domain.StartChatResponse, error) {
	// Handle OSS mode if browser_uuid or oss_user_id is provided
	if req.BrowserUUID != nil || req.OSSUserID != nil {
		return uc.StartOSSChat(ctx, req, ipAddress)
	}

	// For backward compatibility, handle legacy requests
	// This should be removed in future versions
	return nil, errors.New("legacy mode not supported. Please use OSS mode with browser_uuid or oss_user_id")
}

func (uc *ChatUsecase) SendMessage(ctx context.Context, req *domain.SendMessageRequest, senderID *uuid.UUID, senderType string) (*domain.SendMessageResponse, error) {
	// Validate session exists
	session, err := uc.sessionRepo.GetByID(ctx, req.SessionID)
	if err != nil {
		return nil, err
	}

	if session == nil {
		return nil, errors.New("chat session not found")
	}

	// Check if session is active or waiting
	if session.Status == "closed" {
		return nil, errors.New("cannot send message to closed session")
	}

	// Create message
	message := &domain.ChatMessage{
		ID:          uuid.New(),
		SessionID:   req.SessionID,
		SenderID:    senderID,
		SenderType:  senderType,
		Message:     req.Message,
		MessageType: req.MessageType,
		Attachments: req.Attachments,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if message.MessageType == "" {
		message.MessageType = "text"
	}

	if err := uc.messageRepo.Create(ctx, message); err != nil {
		return nil, err
	}

	// If this is an agent response and session is waiting, mark as active
	if senderType == "agent" && session.Status == "waiting" {
		session.Status = "active"
		if err := uc.sessionRepo.Update(ctx, session); err != nil {
			return nil, err
		}

		// Log response
		log := &domain.ChatLog{
			ID:        uuid.New(),
			SessionID: session.ID,
			Action:    "response",
			Details:   "Agent responded to customer",
			UserID:    senderID,
			CreatedAt: time.Now(),
		}

		if err := uc.logRepo.Create(ctx, log); err != nil {
			return nil, err
		}
	}

	return &domain.SendMessageResponse{
		MessageID: message.ID,
		Timestamp: message.CreatedAt,
		Status:    "sent",
	}, nil
}

func (uc *ChatUsecase) AssignAgent(ctx context.Context, req *domain.AssignAgentRequest) error {
	// Validate session exists
	session, err := uc.sessionRepo.GetByID(ctx, req.SessionID)
	if err != nil {
		return err
	}

	if session == nil {
		return errors.New("chat session not found")
	}

	// Validate agent exists
	agent, err := uc.userRepo.GetByID(ctx, req.AgentID)
	if err != nil {
		return err
	}

	if agent == nil {
		return errors.New("agent not found")
	}

	if agent.Role != "agent" {
		return errors.New("user is not an agent")
	}

	// Update session
	session.AgentID = &req.AgentID
	session.DepartmentID = agent.DepartmentID
	session.UpdatedAt = time.Now()

	if err := uc.sessionRepo.Update(ctx, session); err != nil {
		return err
	}

	// Log assignment
	log := &domain.ChatLog{
		ID:        uuid.New(),
		SessionID: session.ID,
		Action:    "assigned",
		Details:   "Agent assigned to chat session",
		UserID:    &req.AgentID,
		CreatedAt: time.Now(),
	}

	if err := uc.logRepo.Create(ctx, log); err != nil {
		return err
	}

	return nil
}

func (uc *ChatUsecase) CloseSession(ctx context.Context, sessionID uuid.UUID, reason string, userID *uuid.UUID) error {
	// Validate session exists
	session, err := uc.sessionRepo.GetByID(ctx, sessionID)
	if err != nil {
		return err
	}

	if session == nil {
		return errors.New("chat session not found")
	}

	if session.Status == "closed" {
		return errors.New("session is already closed")
	}

	// Close session
	if err := uc.sessionRepo.Close(ctx, sessionID); err != nil {
		return err
	}

	// Log closure
	log := &domain.ChatLog{
		ID:        uuid.New(),
		SessionID: sessionID,
		Action:    "closed",
		Details:   reason,
		UserID:    userID,
		CreatedAt: time.Now(),
	}

	if err := uc.logRepo.Create(ctx, log); err != nil {
		return err
	}

	return nil
}

func (uc *ChatUsecase) GetSessionMessages(ctx context.Context, sessionID uuid.UUID) ([]*domain.ChatMessage, error) {
	messages, err := uc.messageRepo.GetBySessionID(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	return messages, nil
}

func (uc *ChatUsecase) GetWaitingSessions(ctx context.Context) ([]*domain.ChatSession, error) {
	sessions, err := uc.sessionRepo.GetWaitingSessions(ctx)
	if err != nil {
		return nil, err
	}

	return sessions, nil
}

func (uc *ChatUsecase) GetActiveSessions(ctx context.Context) ([]*domain.ChatSession, error) {
	sessions, err := uc.sessionRepo.GetActiveSessions(ctx)
	if err != nil {
		return nil, err
	}

	return sessions, nil
}

func (uc *ChatUsecase) GetAgentSessions(ctx context.Context, agentID uuid.UUID) ([]*domain.ChatSession, error) {
	sessions, err := uc.sessionRepo.GetByAgentID(ctx, agentID)
	if err != nil {
		return nil, err
	}

	return sessions, nil
}

func (uc *ChatUsecase) GetAgentSessionsWithPagination(ctx context.Context, agentID uuid.UUID, page, limit int, status string) ([]*domain.ChatSession, int, error) {
	offset := (page - 1) * limit
	sessions, err := uc.sessionRepo.GetWithPagination(ctx, offset, limit, status, &agentID, nil)
	if err != nil {
		return nil, 0, err
	}

	total, err := uc.sessionRepo.Count(ctx, status, &agentID, nil)
	if err != nil {
		return nil, 0, err
	}

	return sessions, total, nil
}

func (uc *ChatUsecase) GetChatUserSessions(ctx context.Context, chatUserID uuid.UUID) ([]*domain.ChatSession, error) {
	sessions, err := uc.sessionRepo.GetByChatUserID(ctx, chatUserID)
	if err != nil {
		return nil, err
	}

	return sessions, nil
}

func (uc *ChatUsecase) GetSessions(ctx context.Context, page, limit int, status string, agentID, departmentID *uuid.UUID) ([]*domain.ChatSession, int, error) {
	offset := (page - 1) * limit
	sessions, err := uc.sessionRepo.GetWithPagination(ctx, offset, limit, status, agentID, departmentID)
	if err != nil {
		return nil, 0, err
	}

	total, err := uc.sessionRepo.Count(ctx, status, agentID, departmentID)
	if err != nil {
		return nil, 0, err
	}

	return sessions, total, nil
}

func (uc *ChatUsecase) GetSession(ctx context.Context, sessionID uuid.UUID) (*domain.ChatSession, error) {
	session, err := uc.sessionRepo.GetByID(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	return session, nil
}

// AutoAssignAgent automatically assigns an available agent to a session
func (uc *ChatUsecase) AutoAssignAgent(ctx context.Context, sessionID uuid.UUID) error {
	// Get available agents (simple round-robin for now)
	agents, err := uc.userRepo.GetAvailableAgents(ctx, nil) // nil for any department
	if err != nil {
		return err
	}

	if len(agents) == 0 {
		// No agents available, session stays in waiting status
		return nil
	}

	// Simple assignment - assign to the first available agent
	// In production, you might want more sophisticated load balancing
	agentID := agents[0].ID

	// Update session
	session, err := uc.sessionRepo.GetByID(ctx, sessionID)
	if err != nil {
		return err
	}

	if session == nil {
		return errors.New("chat session not found")
	}

	session.AgentID = &agentID
	session.Status = "active"
	session.UpdatedAt = time.Now()

	if err := uc.sessionRepo.Update(ctx, session); err != nil {
		return err
	}

	// Log assignment
	log := &domain.ChatLog{
		ID:        uuid.New(),
		SessionID: session.ID,
		Action:    "auto_assigned",
		Details:   "Agent automatically assigned to chat session",
		UserID:    &agentID,
		CreatedAt: time.Now(),
	}

	if err := uc.logRepo.Create(ctx, log); err != nil {
		return err
	}

	return nil
}

func (uc *ChatUsecase) GetMessageByID(ctx context.Context, messageID uuid.UUID) (*domain.ChatMessage, error) {
	return uc.messageRepo.GetByID(ctx, messageID)
}

// StartOSSChat handles OSS-specific chat starting logic
func (uc *ChatUsecase) StartOSSChat(ctx context.Context, req *domain.StartChatRequest, ipAddress string) (*domain.StartChatResponse, error) {
	var chatUser *domain.ChatUser
	var err error

	// Determine if user is anonymous or logged-in
	if req.BrowserUUID != nil {
		// Try to get existing user by browser UUID
		chatUser, err = uc.chatUserRepo.GetByBrowserUUID(ctx, *req.BrowserUUID)
		if err != nil {
			return nil, err
		}
	}

	if chatUser == nil && req.OSSUserID != nil && req.Email != nil {
		// Try to get existing user by OSS user ID
		chatUser, err = uc.chatUserRepo.GetByOSSUserID(ctx, *req.OSSUserID)
		if err != nil {
			return nil, err
		}
	}

	// Create new user if not found
	if chatUser == nil {
		chatUser = &domain.ChatUser{
			ID:        uuid.New(),
			IPAddress: ipAddress,
			UserAgent: req.UserAgent,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		if req.BrowserUUID != nil {
			chatUser.BrowserUUID = req.BrowserUUID
			chatUser.IsAnonymous = true
		}

		if req.OSSUserID != nil && req.Email != nil {
			chatUser.OSSUserID = req.OSSUserID
			chatUser.Email = req.Email
			chatUser.IsAnonymous = false
		}

		if err := uc.chatUserRepo.Create(ctx, chatUser); err != nil {
			return nil, err
		}
	}

	// Create chat session
	session := &domain.ChatSession{
		ID:         uuid.New(),
		ChatUserID: chatUser.ID,
		Topic:      req.Topic,
		Status:     "waiting",
		Priority:   req.Priority,
		StartedAt:  time.Now(),
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	if req.Priority == "" {
		session.Priority = "normal"
	}

	if err := uc.sessionRepo.Create(ctx, session); err != nil {
		return nil, err
	}

	// Log the chat start
	log := &domain.ChatLog{
		ID:        uuid.New(),
		SessionID: session.ID,
		Action:    "started",
		Details:   "Chat session started (OSS mode)",
		CreatedAt: time.Now(),
	}

	if err := uc.logRepo.Create(ctx, log); err != nil {
		return nil, err
	}

	// Try to auto-assign an agent
	if err := uc.AutoAssignAgent(ctx, session.ID); err != nil {
		log.Action = "auto_assignment_failed"
		log.Details = "Failed to auto-assign agent: " + err.Error()
		uc.logRepo.Create(ctx, log)
	} else {
		// Reload session to get updated status
		if updatedSession, err := uc.sessionRepo.GetByID(ctx, session.ID); err == nil && updatedSession != nil {
			session = updatedSession
		}
	}

	// Determine if contact information is required
	requiresContact := chatUser.IsAnonymous || (chatUser.OSSUserID != nil && chatUser.Email != nil)

	return &domain.StartChatResponse{
		SessionID:       session.ID,
		ChatUserID:      chatUser.ID,
		Status:          session.Status,
		Message:         "Chat session started successfully",
		RequiresContact: requiresContact,
	}, nil
}

// SetSessionContact sets contact information for a chat session
func (uc *ChatUsecase) SetSessionContact(ctx context.Context, req *domain.SetSessionContactRequest) (*domain.SetSessionContactResponse, error) {
	// Verify session exists
	session, err := uc.sessionRepo.GetByID(ctx, req.SessionID)
	if err != nil {
		return nil, err
	}
	if session == nil {
		return nil, errors.New("session not found")
	}

	// Check if contact already exists
	existingContact, err := uc.contactRepo.GetBySessionID(ctx, req.SessionID)
	if err != nil {
		return nil, err
	}

	if existingContact != nil {
		// Update existing contact
		existingContact.ContactName = req.ContactName
		existingContact.ContactEmail = req.ContactEmail
		existingContact.ContactPhone = req.ContactPhone
		existingContact.Position = req.Position
		existingContact.CompanyName = req.CompanyName
		existingContact.UpdatedAt = time.Now()

		if err := uc.contactRepo.Update(ctx, existingContact); err != nil {
			return nil, err
		}

		return &domain.SetSessionContactResponse{
			ContactID: existingContact.ID,
			Message:   "Contact information updated successfully",
		}, nil
	}

	// Create new contact
	contact := &domain.ChatSessionContact{
		ID:           uuid.New(),
		SessionID:    req.SessionID,
		ContactName:  req.ContactName,
		ContactEmail: req.ContactEmail,
		ContactPhone: req.ContactPhone,
		Position:     req.Position,
		CompanyName:  req.CompanyName,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := uc.contactRepo.Create(ctx, contact); err != nil {
		return nil, err
	}

	// Log the contact information addition
	log := &domain.ChatLog{
		ID:        uuid.New(),
		SessionID: req.SessionID,
		Action:    "contact_added",
		Details:   "Contact information added",
		CreatedAt: time.Now(),
	}
	uc.logRepo.Create(ctx, log)

	return &domain.SetSessionContactResponse{
		ContactID: contact.ID,
		Message:   "Contact information set successfully",
	}, nil
}

// LinkOSSUser links an anonymous user to an OSS account
func (uc *ChatUsecase) LinkOSSUser(ctx context.Context, req *domain.LinkOSSUserRequest) (*domain.LinkOSSUserResponse, error) {
	// Get the anonymous user by browser UUID
	chatUser, err := uc.chatUserRepo.GetByBrowserUUID(ctx, req.BrowserUUID)
	if err != nil {
		return nil, err
	}
	if chatUser == nil {
		return nil, errors.New("anonymous user not found")
	}

	if !chatUser.IsAnonymous {
		return nil, errors.New("user is already linked to OSS account")
	}

	// Link the user to OSS account
	if err := uc.chatUserRepo.LinkOSSUser(ctx, req.BrowserUUID, req.OSSUserID, req.Email); err != nil {
		return nil, err
	}

	// Get updated user
	updatedUser, err := uc.chatUserRepo.GetByBrowserUUID(ctx, req.BrowserUUID)
	if err != nil {
		return nil, err
	}

	return &domain.LinkOSSUserResponse{
		ChatUserID: updatedUser.ID,
		Message:    "Successfully linked to OSS account",
	}, nil
}

// GetChatHistory gets chat history for a user
func (uc *ChatUsecase) GetChatHistory(ctx context.Context, req *domain.GetChatHistoryRequest) (*domain.GetChatHistoryResponse, error) {
	var chatUser *domain.ChatUser
	var err error

	// Find chat user
	if req.BrowserUUID != nil {
		chatUser, err = uc.chatUserRepo.GetByBrowserUUID(ctx, *req.BrowserUUID)
	} else if req.OSSUserID != nil {
		chatUser, err = uc.chatUserRepo.GetByOSSUserID(ctx, *req.OSSUserID)
	} else {
		return nil, errors.New("either browser_uuid or oss_user_id must be provided")
	}

	if err != nil {
		return nil, err
	}
	if chatUser == nil {
		return nil, errors.New("chat user not found")
	}

	// Set default pagination
	if req.Limit <= 0 {
		req.Limit = 20
	}
	if req.Offset < 0 {
		req.Offset = 0
	}

	// Get sessions with messages
	sessions, err := uc.sessionRepo.GetSessionsWithMessages(ctx, chatUser.ID, req.Limit, req.Offset)
	if err != nil {
		return nil, err
	}

	// Convert to response format
	var sessionHistories []domain.ChatSessionHistory
	for _, session := range sessions {
		history := domain.ChatSessionHistory{
			SessionID:  session.ID,
			Topic:      session.Topic,
			Status:     session.Status,
			Priority:   session.Priority,
			StartedAt:  session.StartedAt,
			EndedAt:    session.EndedAt,
			Agent:      session.Agent,
			Department: session.Department,
		}

		// Get contact information
		contact, err := uc.contactRepo.GetBySessionID(ctx, session.ID)
		if err == nil && contact != nil {
			history.Contact = contact
		}

		// Add messages if available
		if len(session.Messages) > 0 {
			history.Messages = session.Messages
			// Set last message
			history.LastMessage = &session.Messages[len(session.Messages)-1]
		}

		sessionHistories = append(sessionHistories, history)
	}

	// Count total sessions for pagination
	totalSessions, err := uc.sessionRepo.Count(ctx, "", nil, nil)
	if err != nil {
		return nil, err
	}

	return &domain.GetChatHistoryResponse{
		Sessions: sessionHistories,
		Total:    totalSessions,
		Limit:    req.Limit,
		Offset:   req.Offset,
	}, nil
}
