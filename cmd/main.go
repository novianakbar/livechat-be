package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/novianakbar/livechat-be/internal/delivery/handler"
	"github.com/novianakbar/livechat-be/internal/delivery/middleware"
	"github.com/novianakbar/livechat-be/internal/delivery/routes"
	"github.com/novianakbar/livechat-be/internal/infrastructure/database"
	"github.com/novianakbar/livechat-be/internal/infrastructure/email"
	"github.com/novianakbar/livechat-be/internal/infrastructure/repository"
	"github.com/novianakbar/livechat-be/internal/usecase"
	"github.com/novianakbar/livechat-be/pkg/config"
	"github.com/novianakbar/livechat-be/pkg/utils"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Initialize database connections
	db, err := database.NewPostgresConnection(cfg)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	redisClient, err := database.NewRedisConnection(cfg)
	if err != nil {
		log.Fatal("Failed to connect to Redis:", err)
	}

	// Initialize JWT utility
	jwtUtil := utils.NewJWTUtil(cfg.JWT.Secret, cfg.JWT.AccessTokenDuration, cfg.JWT.RefreshTokenDuration)

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	sessionRepo := repository.NewChatSessionRepository(db)
	messageRepo := repository.NewChatMessageRepository(db)
	customerRepo := repository.NewCustomerRepository(db)
	logRepo := repository.NewChatLogRepository(db)

	// Initialize use cases
	authUsecase := usecase.NewAuthUsecase(userRepo, jwtUtil)
	chatUsecase := usecase.NewChatUsecase(sessionRepo, messageRepo, customerRepo, userRepo, logRepo)
	analyticsUsecase := usecase.NewAnalyticsUsecase(sessionRepo, messageRepo, userRepo, customerRepo)
	userUsecase := usecase.NewUserUsecase(userRepo)

	// Initialize email service
	emailService := email.NewSendGridService(&cfg.Email)
	emailUsecase := usecase.NewEmailUseCase(emailService)

	// Initialize handlers
	wsHandler := handler.NewWebSocketHandler(chatUsecase)
	authHandler := handler.NewAuthHandler(authUsecase)
	chatHandler := handler.NewChatHandler(chatUsecase, wsHandler)
	analyticsHandler := handler.NewAnalyticsHandler(analyticsUsecase)
	userHandler := handler.NewUserHandler(userUsecase)
	emailHandler := handler.NewEmailHandler(emailUsecase)

	// Initialize middleware
	authMiddleware := middleware.NewAuthMiddleware(authUsecase)

	// Initialize Fiber app
	app := fiber.New(fiber.Config{
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return ctx.Status(code).JSON(fiber.Map{
				"success": false,
				"message": "Internal server error",
				"error":   err.Error(),
			})
		},
	})

	// Middleware
	app.Use(recover.New())
	app.Use(logger.New())

	// CORS middleware - Allow localhost:3000 for development
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3000,http://127.0.0.1:3000,http://localhost:5173,http://127.0.1:5173",
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
		AllowHeaders:     "Origin,Content-Type,Accept,Authorization,X-Requested-With",
		AllowCredentials: true,
		ExposeHeaders:    "Content-Length,Authorization",
	}))

	// Setup routes
	routes.SetupRoutes(app, authHandler, chatHandler, wsHandler, analyticsHandler, userHandler, emailHandler, authMiddleware)

	// Start server
	serverAddr := cfg.Server.Host + ":" + cfg.Server.Port
	log.Printf("Server starting on %s", serverAddr)
	log.Fatal(app.Listen(serverAddr))

	// Close connections
	if redisClient != nil {
		redisClient.Close()
	}
}
