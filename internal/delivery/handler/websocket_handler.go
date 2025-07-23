package handler

import (
	"context"
	"log"
	"sync"

	"github.com/gofiber/websocket/v2"
	"github.com/google/uuid"
	"github.com/novianakbar/livechat-be/internal/domain"
	"github.com/novianakbar/livechat-be/internal/usecase"
)

type WebSocketHandler struct {
	chatUsecase *usecase.ChatUsecase
	clients     map[uuid.UUID]*WebSocketClient
	sessions    map[uuid.UUID]map[uuid.UUID]bool // sessionID -> clientID -> true
	clientsMux  sync.RWMutex
}

type WebSocketClient struct {
	conn      *websocket.Conn
	clientID  uuid.UUID
	userType  string // "agent" or "customer"
	userID    *uuid.UUID
	sessionID *uuid.UUID
}

func NewWebSocketHandler(chatUsecase *usecase.ChatUsecase) *WebSocketHandler {
	return &WebSocketHandler{
		chatUsecase: chatUsecase,
		clients:     make(map[uuid.UUID]*WebSocketClient),
		sessions:    make(map[uuid.UUID]map[uuid.UUID]bool),
	}
}

func (h *WebSocketHandler) HandleConnection(c *websocket.Conn) {
	defer c.Close()

	// Generate client ID
	clientID := uuid.New()

	// Create client
	client := &WebSocketClient{
		conn:     c,
		clientID: clientID,
		userType: "customer", // default, will be updated when joining session
	}

	// Register client
	h.clientsMux.Lock()
	h.clients[clientID] = client
	h.clientsMux.Unlock()

	// Remove client on disconnect
	defer func() {
		h.clientsMux.Lock()
		var disconnectedSessionID *uuid.UUID
		// Remove from session if joined
		if client.sessionID != nil {
			disconnectedSessionID = client.sessionID
			if sessionClients, exists := h.sessions[*client.sessionID]; exists {
				delete(sessionClients, clientID)
				if len(sessionClients) == 0 {
					delete(h.sessions, *client.sessionID)
				}
			}
		}
		delete(h.clients, clientID)
		h.clientsMux.Unlock()

		// Broadcast connection status update after disconnect
		if disconnectedSessionID != nil {
			connectionStatus := h.GetSessionConnectedClients(*disconnectedSessionID)
			h.broadcastToSession(*disconnectedSessionID, domain.WebSocketResponse{
				Type:    "connection_status_update",
				Success: true,
				Data: map[string]interface{}{
					"session_id":        *disconnectedSessionID,
					"connection_status": connectionStatus,
				},
			})
		}
	}()

	log.Printf("WebSocket client connected: %s", clientID)

	// Handle messages
	for {
		var msg domain.WebSocketMessage
		if err := c.ReadJSON(&msg); err != nil {
			log.Printf("WebSocket read error: %v", err)
			break
		}

		// Process message based on type
		switch msg.Type {
		case "join_session":
			h.handleJoinSession(c, client, &msg)
		case "send_message":
			h.handleSendMessage(c, client, &msg)
		case "agent_typing":
			h.handleAgentTyping(c, client, &msg)
		case "customer_typing":
			h.handleCustomerTyping(c, client, &msg)
		default:
			h.sendError(c, "Unknown message type", "")
		}
	}
}

func (h *WebSocketHandler) handleJoinSession(c *websocket.Conn, client *WebSocketClient, msg *domain.WebSocketMessage) {
	// Parse message data to get user info
	if msgData, ok := msg.Data.(map[string]interface{}); ok {
		if agentID, exists := msgData["agent_id"]; exists && agentID != nil {
			// Agent joining
			client.userType = "agent"
			if agentIDStr, ok := agentID.(string); ok {
				if parsedAgentID, err := uuid.Parse(agentIDStr); err == nil {
					client.userID = &parsedAgentID
				}
			}
		} else {
			// Customer joining
			client.userType = "customer"
		}
	}

	// Set session ID
	sessionID := msg.SessionID
	client.sessionID = &sessionID

	// Add client to session
	h.clientsMux.Lock()
	if _, exists := h.sessions[sessionID]; !exists {
		h.sessions[sessionID] = make(map[uuid.UUID]bool)
	}
	h.sessions[sessionID][client.clientID] = true
	h.clientsMux.Unlock()

	log.Printf("Client %s (%s) joined session %s", client.clientID, client.userType, sessionID)

	// Send confirmation
	response := domain.WebSocketResponse{
		Type:    "joined_session",
		Success: true,
		Data: map[string]interface{}{
			"session_id": sessionID,
			"client_id":  client.clientID,
			"user_type":  client.userType,
		},
	}

	if err := c.WriteJSON(response); err != nil {
		log.Printf("WebSocket write error: %v", err)
	}

	// Send initial connection status to the newly joined client
	initialConnectionStatus := h.GetSessionConnectedClients(sessionID)
	if err := c.WriteJSON(domain.WebSocketResponse{
		Type:    "connection_status_update",
		Success: true,
		Data: map[string]interface{}{
			"session_id":        sessionID,
			"connection_status": initialConnectionStatus,
		},
	}); err != nil {
		log.Printf("WebSocket write error for connection status: %v", err)
	}

	// Notify other clients in session about new participant
	h.broadcastToSessionExcept(sessionID, client.clientID, domain.WebSocketResponse{
		Type:    "user_joined",
		Success: true,
		Data: map[string]interface{}{
			"session_id": sessionID,
			"user_type":  client.userType,
			"user_id":    client.userID,
		},
	})

	// Broadcast connection status update
	connectionStatus := h.GetSessionConnectedClients(sessionID)
	h.broadcastToSession(sessionID, domain.WebSocketResponse{
		Type:    "connection_status_update",
		Success: true,
		Data: map[string]interface{}{
			"session_id":        sessionID,
			"connection_status": connectionStatus,
		},
	})
}

func (h *WebSocketHandler) handleSendMessage(c *websocket.Conn, client *WebSocketClient, msg *domain.WebSocketMessage) {
	// Parse message data
	messageData, ok := msg.Data.(map[string]interface{})
	if !ok {
		h.sendError(c, "Invalid message data", "")
		return
	}

	message, ok := messageData["message"].(string)
	if !ok {
		h.sendError(c, "Message text is required", "")
		return
	}

	// Create send message request
	req := &domain.SendMessageRequest{
		SessionID:   msg.SessionID,
		Message:     message,
		MessageType: "text",
	}

	// Send message through usecase
	ctx := context.Background()
	response, err := h.chatUsecase.SendMessage(ctx, req, client.userID, client.userType)
	if err != nil {
		h.sendError(c, "Failed to send message", err.Error())
		return
	}

	// Broadcast message to all clients in the session
	h.broadcastToSession(msg.SessionID, domain.WebSocketResponse{
		Type:    "new_message",
		Success: true,
		Data: map[string]interface{}{
			"session_id":   msg.SessionID,
			"message_id":   response.MessageID,
			"message":      message,
			"sender_type":  client.userType,
			"sender_id":    client.userID,
			"message_type": "text",
			"timestamp":    response.Timestamp,
		},
	})
}

func (h *WebSocketHandler) handleAgentTyping(c *websocket.Conn, client *WebSocketClient, msg *domain.WebSocketMessage) {
	// Parse typing data
	typingData, ok := msg.Data.(map[string]interface{})
	if !ok {
		return
	}

	isTyping, _ := typingData["is_typing"].(bool)

	// Broadcast typing indicator to customers in the session
	h.broadcastToSessionExcept(msg.SessionID, client.clientID, domain.WebSocketResponse{
		Type:    "typing",
		Success: true,
		Data: map[string]interface{}{
			"session_id":  msg.SessionID,
			"sender_type": "agent",
			"is_typing":   isTyping,
		},
	})
}

func (h *WebSocketHandler) handleCustomerTyping(c *websocket.Conn, client *WebSocketClient, msg *domain.WebSocketMessage) {
	// Parse typing data
	typingData, ok := msg.Data.(map[string]interface{})
	if !ok {
		return
	}

	isTyping, _ := typingData["is_typing"].(bool)

	// Broadcast typing indicator to agents in the session
	h.broadcastToSessionExcept(msg.SessionID, client.clientID, domain.WebSocketResponse{
		Type:    "typing",
		Success: true,
		Data: map[string]interface{}{
			"session_id":  msg.SessionID,
			"sender_type": "customer",
			"is_typing":   isTyping,
		},
	})
}

func (h *WebSocketHandler) sendError(c *websocket.Conn, message, error string) {
	response := domain.WebSocketResponse{
		Type:    "error",
		Success: false,
		Data:    nil,
		Error:   error,
	}

	if err := c.WriteJSON(response); err != nil {
		log.Printf("WebSocket write error: %v", err)
	}
}

func (h *WebSocketHandler) broadcastToSession(sessionID uuid.UUID, response domain.WebSocketResponse) {
	h.clientsMux.RLock()
	defer h.clientsMux.RUnlock()

	if sessionClients, exists := h.sessions[sessionID]; exists {
		for clientID := range sessionClients {
			if client, exists := h.clients[clientID]; exists {
				if err := client.conn.WriteJSON(response); err != nil {
					log.Printf("WebSocket broadcast error: %v", err)
				}
			}
		}
	}
}

func (h *WebSocketHandler) broadcastToSessionExcept(sessionID uuid.UUID, exceptClientID uuid.UUID, response domain.WebSocketResponse) {
	h.clientsMux.RLock()
	defer h.clientsMux.RUnlock()

	if sessionClients, exists := h.sessions[sessionID]; exists {
		for clientID := range sessionClients {
			if clientID == exceptClientID {
				continue
			}
			if client, exists := h.clients[clientID]; exists {
				if err := client.conn.WriteJSON(response); err != nil {
					log.Printf("WebSocket broadcast error: %v", err)
				}
			}
		}
	}
}

func (h *WebSocketHandler) BroadcastMessage(sessionID uuid.UUID, message *domain.ChatMessage) {
	response := domain.WebSocketResponse{
		Type:    "new_message",
		Success: true,
		Data: map[string]interface{}{
			"session_id":   sessionID,
			"message_id":   message.ID,
			"message":      message.Message,
			"sender_type":  message.SenderType,
			"message_type": message.MessageType,
			"timestamp":    message.CreatedAt,
		},
	}

	h.broadcastToSession(sessionID, response)
}

func (h *WebSocketHandler) BroadcastSessionUpdate(sessionID uuid.UUID, status string) {
	response := domain.WebSocketResponse{
		Type:    "session_update",
		Success: true,
		Data: map[string]interface{}{
			"session_id": sessionID,
			"status":     status,
		},
	}

	h.broadcastToSession(sessionID, response)
}

func (h *WebSocketHandler) GetConnectedClients() int {
	h.clientsMux.RLock()
	defer h.clientsMux.RUnlock()
	return len(h.clients)
}

// GetSessionConnectedClients returns connected clients for a specific session
func (h *WebSocketHandler) GetSessionConnectedClients(sessionID uuid.UUID) map[string]interface{} {
	h.clientsMux.RLock()
	defer h.clientsMux.RUnlock()

	result := map[string]interface{}{
		"customer_connected": false,
		"agent_connected":    false,
		"total_customer":     0,
		"total_agent":        0,
	}

	if sessionClients, exists := h.sessions[sessionID]; exists {
		for clientID := range sessionClients {
			if client, exists := h.clients[clientID]; exists {
				if client.userType == "customer" {
					result["total_customer"] = result["total_customer"].(int) + 1
					result["customer_connected"] = true
				} else if client.userType == "agent" {
					result["total_agent"] = result["total_agent"].(int) + 1
					result["agent_connected"] = true
				}
			}
		}
	}

	return result
}
