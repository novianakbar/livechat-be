package usecase

import (
	"context"
	"database/sql"
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
	uuidV7, _ := uuid.NewV7()
	var senderIDStr sql.NullString
	if senderID != nil {
		senderIDStr = sql.NullString{
			String: senderID.String(),
			Valid:  true,
		}
	}

	message := &domain.ChatMessage{
		ID:          uuidV7.String(),
		SessionID:   req.SessionID.String(),
		SenderID:    senderIDStr,
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
		uuidV7Log, _ := uuid.NewV7()
		var userIDStr sql.NullString
		if senderID != nil {
			userIDStr = sql.NullString{
				String: senderID.String(),
				Valid:  true,
			}
		}

		log := &domain.ChatLog{
			ID:        uuidV7Log.String(),
			SessionID: session.ID,
			Action:    "response",
			Details: sql.NullString{
				String: "Agent responded to customer",
				Valid:  true,
			},
			UserID:    userIDStr,
			CreatedAt: time.Now(),
		}

		if err := uc.logRepo.Create(ctx, log); err != nil {
			return nil, err
		}
	}

	messageUUID, _ := uuid.Parse(message.ID)
	return &domain.SendMessageResponse{
		MessageID: messageUUID,
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
	agent, err := uc.userRepo.GetByID(ctx, req.AgentID.String())
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
	session.AgentID = sql.NullString{
		String: req.AgentID.String(),
		Valid:  true,
	}
	session.DepartmentID = agent.DepartmentID
	session.UpdatedAt = time.Now()

	if err := uc.sessionRepo.Update(ctx, session); err != nil {
		return err
	}

	// Log assignment
	uuidV7Log2, _ := uuid.NewV7()
	log := &domain.ChatLog{
		ID:        uuidV7Log2.String(),
		SessionID: session.ID,
		Action:    "assigned",
		Details: sql.NullString{
			String: "Agent assigned to chat session",
			Valid:  true,
		},
		UserID: sql.NullString{
			String: req.AgentID.String(),
			Valid:  true,
		},
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
	uuidV7Close, _ := uuid.NewV7()
	log := &domain.ChatLog{
		ID:        uuidV7Close.String(),
		SessionID: sessionID.String(),
		Action:    "closed",
		Details: sql.NullString{
			String: reason,
			Valid:  true,
		},
		UserID: func() sql.NullString {
			if userID != nil {
				return sql.NullString{String: userID.String(), Valid: true}
			}
			return sql.NullString{Valid: false}
		}(),
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

	session.AgentID = sql.NullString{
		String: agentID,
		Valid:  true,
	}
	session.Status = "active"
	session.UpdatedAt = time.Now()

	if err := uc.sessionRepo.Update(ctx, session); err != nil {
		return err
	}

	// Log assignment
	uuidV7AutoAssign, _ := uuid.NewV7()
	log := &domain.ChatLog{
		ID:        uuidV7AutoAssign.String(),
		SessionID: session.ID,
		Action:    "auto_assigned",
		Details: sql.NullString{
			String: "Agent automatically assigned to chat session",
			Valid:  true,
		},
		UserID: sql.NullString{
			String: agentID,
			Valid:  true,
		},
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
		uuidV7ChatUser, _ := uuid.NewV7()
		chatUser = &domain.ChatUser{
			ID:        uuidV7ChatUser.String(),
			IPAddress: ipAddress,
			UserAgent: func() sql.NullString {
				if req.UserAgent != nil {
					return sql.NullString{String: *req.UserAgent, Valid: true}
				}
				return sql.NullString{Valid: false}
			}(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		if req.BrowserUUID != nil {
			chatUser.BrowserUUID = sql.NullString{
				String: req.BrowserUUID.String(),
				Valid:  true,
			}
			chatUser.IsAnonymous = true
		}

		if req.OSSUserID != nil && req.Email != nil {
			chatUser.OSSUserID = sql.NullString{
				String: *req.OSSUserID,
				Valid:  true,
			}
			chatUser.Email = sql.NullString{
				String: *req.Email,
				Valid:  true,
			}
			chatUser.IsAnonymous = false
		}

		if err := uc.chatUserRepo.Create(ctx, chatUser); err != nil {
			return nil, err
		}
	}

	// Create chat session
	uuidV7Session, _ := uuid.NewV7()
	session := &domain.ChatSession{
		ID:         uuidV7Session.String(),
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
	uuidV7StartLog, _ := uuid.NewV7()
	log := &domain.ChatLog{
		ID:        uuidV7StartLog.String(),
		SessionID: session.ID,
		Action:    "started",
		Details: sql.NullString{
			String: "Chat session started (OSS mode)",
			Valid:  true,
		},
		CreatedAt: time.Now(),
	}

	if err := uc.logRepo.Create(ctx, log); err != nil {
		return nil, err
	}

	// Try to auto-assign an agent
	sessionUUID, _ := uuid.Parse(session.ID)
	if err := uc.AutoAssignAgent(ctx, sessionUUID); err != nil {
		log.Action = "auto_assignment_failed"
		log.Details = sql.NullString{
			String: "Failed to auto-assign agent: " + err.Error(),
			Valid:  true,
		}
		uc.logRepo.Create(ctx, log)
	} else {
		// Reload session to get updated status
		if updatedSession, err := uc.sessionRepo.GetByID(ctx, sessionUUID); err == nil && updatedSession != nil {
			session = updatedSession
		}
	}

	// Determine if contact information is required
	requiresContact := chatUser.IsAnonymous || (chatUser.OSSUserID.Valid && chatUser.Email.Valid)

	// Parse IDs for DTO compatibility
	sessionUUIDForDTO, _ := uuid.Parse(session.ID)
	chatUserUUIDForDTO, _ := uuid.Parse(chatUser.ID)

	return &domain.StartChatResponse{
		SessionID:       sessionUUIDForDTO,
		ChatUserID:      chatUserUUIDForDTO,
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
		existingContact.ContactPhone = func() sql.NullString {
			if req.ContactPhone != nil {
				return sql.NullString{String: *req.ContactPhone, Valid: true}
			}
			return sql.NullString{Valid: false}
		}()
		existingContact.Position = func() sql.NullString {
			if req.Position != nil {
				return sql.NullString{String: *req.Position, Valid: true}
			}
			return sql.NullString{Valid: false}
		}()
		existingContact.CompanyName = func() sql.NullString {
			if req.CompanyName != nil {
				return sql.NullString{String: *req.CompanyName, Valid: true}
			}
			return sql.NullString{Valid: false}
		}()
		existingContact.UpdatedAt = time.Now()

		if err := uc.contactRepo.Update(ctx, existingContact); err != nil {
			return nil, err
		}

		contactUUIDForDTO, _ := uuid.Parse(existingContact.ID)
		return &domain.SetSessionContactResponse{
			ContactID: contactUUIDForDTO,
			Message:   "Contact information updated successfully",
		}, nil
	}

	// Create new contact
	uuidV7Contact, _ := uuid.NewV7()
	contact := &domain.ChatSessionContact{
		ID:           uuidV7Contact.String(),
		SessionID:    req.SessionID.String(),
		ContactName:  req.ContactName,
		ContactEmail: req.ContactEmail,
		ContactPhone: func() sql.NullString {
			if req.ContactPhone != nil {
				return sql.NullString{String: *req.ContactPhone, Valid: true}
			}
			return sql.NullString{Valid: false}
		}(),
		Position: func() sql.NullString {
			if req.Position != nil {
				return sql.NullString{String: *req.Position, Valid: true}
			}
			return sql.NullString{Valid: false}
		}(),
		CompanyName: func() sql.NullString {
			if req.CompanyName != nil {
				return sql.NullString{String: *req.CompanyName, Valid: true}
			}
			return sql.NullString{Valid: false}
		}(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := uc.contactRepo.Create(ctx, contact); err != nil {
		return nil, err
	}

	// Log the contact information addition
	uuidV7ContactLog, _ := uuid.NewV7()
	log := &domain.ChatLog{
		ID:        uuidV7ContactLog.String(),
		SessionID: req.SessionID.String(),
		Action:    "contact_added",
		Details: sql.NullString{
			String: "Contact information added",
			Valid:  true,
		},
		CreatedAt: time.Now(),
	}
	uc.logRepo.Create(ctx, log)

	contactUUIDForDTO, _ := uuid.Parse(contact.ID)
	return &domain.SetSessionContactResponse{
		ContactID: contactUUIDForDTO,
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

	chatUserUUIDForDTO, _ := uuid.Parse(updatedUser.ID)
	return &domain.LinkOSSUserResponse{
		ChatUserID: chatUserUUIDForDTO,
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
	chatUserUUID, err := uuid.Parse(chatUser.ID)
	if err != nil {
		return nil, errors.New("invalid chat user ID format")
	}
	sessions, err := uc.sessionRepo.GetSessionsWithMessages(ctx, chatUserUUID, req.Limit, req.Offset)
	if err != nil {
		return nil, err
	}

	// Convert to response format
	var sessionHistories []domain.ChatSessionHistory
	for _, session := range sessions {
		sessionUUID, _ := uuid.Parse(session.ID)

		history := domain.ChatSessionHistory{
			SessionID: sessionUUID,
			Topic:     session.Topic,
			Status:    session.Status,
			Priority:  session.Priority,
			StartedAt: session.StartedAt,
			EndedAt: func() *time.Time {
				if session.EndedAt.Valid {
					return &session.EndedAt.Time
				}
				return nil
			}(),
			Agent:      session.Agent,
			Department: session.Department,
		}

		// Get contact information
		contact, err := uc.contactRepo.GetBySessionID(ctx, sessionUUID)
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
