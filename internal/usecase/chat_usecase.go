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
	customerRepo domain.CustomerRepository
	userRepo     domain.UserRepository
	logRepo      domain.ChatLogRepository
}

func NewChatUsecase(
	sessionRepo domain.ChatSessionRepository,
	messageRepo domain.ChatMessageRepository,
	customerRepo domain.CustomerRepository,
	userRepo domain.UserRepository,
	logRepo domain.ChatLogRepository,
) *ChatUsecase {
	return &ChatUsecase{
		sessionRepo:  sessionRepo,
		messageRepo:  messageRepo,
		customerRepo: customerRepo,
		userRepo:     userRepo,
		logRepo:      logRepo,
	}
}

func (uc *ChatUsecase) StartChat(ctx context.Context, req *domain.StartChatRequest, ipAddress string) (*domain.StartChatResponse, error) {
	// Create or get customer
	customer := &domain.Customer{
		ID:          uuid.New(),
		CompanyName: req.CompanyName,
		PersonName:  req.PersonName,
		Email:       req.Email,
		IPAddress:   ipAddress,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	customer, err := uc.customerRepo.GetOrCreate(ctx, customer)
	if err != nil {
		return nil, err
	}

	// Create chat session
	session := &domain.ChatSession{
		ID:         uuid.New(),
		CustomerID: customer.ID,
		Topic:      req.Topic,
		Priority:   req.Priority,
		Status:     "waiting",
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

	// Log chat started
	log := &domain.ChatLog{
		ID:        uuid.New(),
		SessionID: session.ID,
		Action:    "started",
		Details:   "Chat session started by customer",
		CreatedAt: time.Now(),
	}

	if err := uc.logRepo.Create(ctx, log); err != nil {
		return nil, err
	}

	// Try to auto-assign an agent
	if err := uc.AutoAssignAgent(ctx, session.ID); err != nil {
		// Log the error but don't fail the session creation
		// The session will remain in "waiting" status
		log.Action = "auto_assignment_failed"
		log.Details = "Failed to auto-assign agent: " + err.Error()
		uc.logRepo.Create(ctx, log)
	} else {
		// Reload session to get updated status
		if updatedSession, err := uc.sessionRepo.GetByID(ctx, session.ID); err == nil && updatedSession != nil {
			session = updatedSession
		}
	}

	return &domain.StartChatResponse{
		SessionID: session.ID,
		Status:    session.Status,
		Message:   "Chat session started. Please wait for an agent to respond.",
	}, nil
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

func (uc *ChatUsecase) GetCustomerSessions(ctx context.Context, customerID uuid.UUID) ([]*domain.ChatSession, error) {
	sessions, err := uc.sessionRepo.GetByCustomerID(ctx, customerID)
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
