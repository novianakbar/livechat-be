package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/novianakbar/livechat-be/internal/delivery/handler"
	"github.com/novianakbar/livechat-be/internal/delivery/middleware"
)

func SetupRoutes(
	app *fiber.App,
	authHandler *handler.AuthHandler,
	chatHandler *handler.ChatHandler,
	analyticsHandler *handler.AnalyticsHandler,
	userHandler *handler.UserHandler,
	emailHandler *handler.EmailHandler,
	agentStatusHandler *handler.AgentStatusHandler,
	authMiddleware *middleware.AuthMiddleware,
) {
	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"message": "LiveChat API is running",
		})
	})

	// Handle preflight requests
	app.Options("/*", func(c *fiber.Ctx) error {
		return c.SendStatus(204)
	})

	// API routes
	api := app.Group("/api")

	// Public chat routes (legacy for backward compatibility)
	public := api.Group("/public")
	public.Post("/chat/start", chatHandler.StartChat)
	public.Post("/chat/message", chatHandler.SendMessage)
	public.Get("/chat/session/:session_id/messages", chatHandler.GetSessionMessages)

	// OSS Chat routes (public endpoints for OSS integration)
	ossChat := api.Group("/chat")
	ossChat.Post("/start", chatHandler.StartChat)
	ossChat.Post("/contact", chatHandler.SetSessionContact)
	ossChat.Post("/link-user", chatHandler.LinkOSSUser)
	ossChat.Get("/history", chatHandler.GetChatHistory)
	ossChat.Get("/session/:session_id", chatHandler.GetSession)

	// Authentication routes
	auth := api.Group("/auth")
	auth.Post("/login", authHandler.Login)
	auth.Post("/refresh", authHandler.RefreshToken)
	auth.Post("/logout", authMiddleware.RequireAuth(), authHandler.Logout)
	auth.Get("/validate", authMiddleware.RequireAuth(), authHandler.ValidateSession)
	auth.Get("/profile", authMiddleware.RequireAuth(), authHandler.GetProfile)
	auth.Post("/register", authMiddleware.RequireAuth(), authMiddleware.RequireAdmin(), authHandler.Register)

	// Protected chat management routes
	chatManagement := api.Group("/chat-management")
	chatManagement.Use(authMiddleware.RequireAuth())

	// Agent routes
	agent := chatManagement.Group("/agent")
	agent.Use(authMiddleware.RequireAgent())
	agent.Post("/message", chatHandler.SendMessage)
	agent.Post("/assign", chatHandler.AssignAgent)
	agent.Post("/close", chatHandler.CloseSession)
	agent.Get("/sessions", chatHandler.GetAgentSessions)
	agent.Get("/sessions/:id/connection-status", chatHandler.GetSessionConnectionStatus)
	agent.Get("/sessions/:session_id", chatHandler.GetSession)

	// Admin routes
	admin := chatManagement.Group("/admin")
	admin.Use(authMiddleware.RequireAdmin())
	admin.Get("/waiting", chatHandler.GetWaitingSessions)
	admin.Get("/active", chatHandler.GetActiveSessions)
	admin.Post("/assign", chatHandler.AssignAgent)
	admin.Post("/close", chatHandler.CloseSession)
	admin.Get("/sessions", chatHandler.GetSessions)
	// admin.Get("/sessions/:id/connection-status", chatHandler.GetSessionConnectionStatus)
	// admin.Get("/sessions/:id", chatHandler.GetSession)

	// User routes
	users := api.Group("/users")
	users.Use(authMiddleware.RequireAuth())
	users.Get("/", userHandler.GetUsers)
	users.Get("/agents", userHandler.GetAgents)
	users.Get("/:id", userHandler.GetUser)

	// Analytics routes
	analytics := api.Group("/analytics")
	analytics.Use(authMiddleware.RequireAuth())
	analytics.Get("/dashboard", analyticsHandler.GetDashboardStats)
	analytics.Get("/", analyticsHandler.GetAnalytics)

	// Email routes
	email := api.Group("/email")
	email.Use(authMiddleware.RequireAuth())
	email.Post("/send", emailHandler.SendEmail)
	email.Post("/welcome", emailHandler.SendWelcomeEmail)
	email.Post("/password-reset", emailHandler.SendPasswordResetEmail)
	email.Post("/chat-transcript", emailHandler.SendChatTranscriptEmail)
	email.Post("/custom", emailHandler.SendCustomEmail)

	// Agent status routes
	agentStatus := api.Group("/agent-status")
	agentStatus.Use(authMiddleware.RequireAuth())

	// Agent heartbeat (for agents and admins)
	agentStatus.Post("/heartbeat", authMiddleware.RequireAgent(), agentStatusHandler.AgentHeartbeat)
	agentStatus.Post("/offline", authMiddleware.RequireAgent(), agentStatusHandler.SetAgentOffline)

	// Admin routes for viewing agent status
	agentStatus.Get("/online", agentStatusHandler.GetOnlineAgents)
	agentStatus.Get("/department/:department_id", agentStatusHandler.GetOnlineAgentsByDepartment)
	agentStatus.Get("/agent/:agent_id", agentStatusHandler.GetAgentStatus)
	agentStatus.Get("/stats", agentStatusHandler.GetDepartmentStats)
}
