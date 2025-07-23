package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Database  DatabaseConfig
	Redis     RedisConfig
	Server    ServerConfig
	JWT       JWTConfig
	WebSocket WebSocketConfig
	Email     EmailConfig
	App       AppConfig
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

type ServerConfig struct {
	Host string
	Port string
}

type JWTConfig struct {
	Secret               string
	AccessTokenDuration  time.Duration
	RefreshTokenDuration time.Duration
}

type WebSocketConfig struct {
	ReadBufferSize  int
	WriteBufferSize int
}

type EmailConfig struct {
	SendGridAPIKey string
	FromEmail      string
	FromName       string
}

type AppConfig struct {
	Environment string
}

func LoadConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found")
	}

	// JWT Token Durations
	accessTokenDuration, err := time.ParseDuration(getEnv("JWT_ACCESS_TOKEN_DURATION", "15m"))
	if err != nil {
		accessTokenDuration = 15 * time.Minute
	}

	refreshTokenDuration, err := time.ParseDuration(getEnv("JWT_REFRESH_TOKEN_DURATION", "168h"))
	if err != nil {
		refreshTokenDuration = 168 * time.Hour // 7 days
	}

	readBufferSize, err := strconv.Atoi(getEnv("WS_READ_BUFFER_SIZE", "1024"))
	if err != nil {
		readBufferSize = 1024
	}

	writeBufferSize, err := strconv.Atoi(getEnv("WS_WRITE_BUFFER_SIZE", "1024"))
	if err != nil {
		writeBufferSize = 1024
	}

	redisDB, err := strconv.Atoi(getEnv("REDIS_DB", "0"))
	if err != nil {
		redisDB = 0
	}

	return &Config{
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "password"),
			Name:     getEnv("DB_NAME", "livechat_db"),
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnv("REDIS_PORT", "6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       redisDB,
		},
		Server: ServerConfig{
			Host: getEnv("SERVER_HOST", "0.0.0.0"),
			Port: getEnv("SERVER_PORT", "8080"),
		},
		JWT: JWTConfig{
			Secret:               getEnv("JWT_SECRET", "your-secret-key-here"),
			AccessTokenDuration:  accessTokenDuration,
			RefreshTokenDuration: refreshTokenDuration,
		},
		WebSocket: WebSocketConfig{
			ReadBufferSize:  readBufferSize,
			WriteBufferSize: writeBufferSize,
		},
		Email: EmailConfig{
			SendGridAPIKey: getEnv("SENDGRID_API_KEY", ""),
			FromEmail:      getEnv("SENDGRID_FROM_EMAIL", "noreply@yourcompany.com"),
			FromName:       getEnv("SENDGRID_FROM_NAME", "LiveChat System"),
		},
		App: AppConfig{
			Environment: getEnv("APP_ENV", "development"),
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
