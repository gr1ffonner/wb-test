package main

import (
	"context"
	"log/slog"
	"wb-test/pkg/broker"
	"wb-test/pkg/cache"
	"wb-test/pkg/config"
	"wb-test/pkg/db"
	"wb-test/pkg/logger"
	"wb-test/pkg/utils/jwt"
)

func main() {
	// Load config first
	cfg, err := config.Load()
	if err != nil {
		slog.Error("Failed to load config", "error", err)
		panic(err)
	}

	// Initialize unified logger
	logger.InitLogger(cfg.Logger)
	log := slog.Default()

	ctx := context.Background()

	log.Info("Starting application", "config", cfg)

	// Initialize database
	db, err := db.NewPostgres(ctx, cfg.Database.DSN)
	if err != nil {
		log.Error("Failed to connect to database", "error", err)
		panic(err)
	}
	defer db.Close()
	log.Info("Database connected successfully")

	// Initialize cache
	cache, err := cache.NewRedis(ctx, cfg.Redis.Addr, cfg.Redis.Password, cfg.Redis.DB)
	if err != nil {
		log.Error("Failed to connect to Redis", "error", err)
		panic(err)
	}
	defer cache.Close()
	log.Info("Redis connected successfully")

	// Initialize broker
	broker, err := broker.NewNATS(ctx, cfg.NATS.URL)
	if err != nil {
		log.Error("Failed to connect to NATS", "error", err)
		panic(err)
	}
	defer broker.Close()
	log.Info("NATS connected successfully")

	log.Info("All services initialized successfully")

	// JWT Token Generation and Validation Example
	demonstrateJWTFlow(log)
}

// demonstrateJWTFlow shows token generation and validation
func demonstrateJWTFlow(log *slog.Logger) {
	log.Info("=== JWT Token Flow Demo ===")

	// 1. Generate a JWT token
	userID := 123
	username := "john_doe"

	token, err := jwt.GenerateJWT(userID, username)
	if err != nil {
		log.Error("Failed to generate JWT", "error", err)
		return
	}
	log.Info("JWT Token generated", "user_id", userID, "username", username, "token", token[:20]+"...")

	// 2. Validate the token
	err = jwt.ValidateToken(token)
	if err != nil {
		log.Error("Token validation failed", "error", err)
		return
	}
	log.Info("Token validation successful")

	// 3. Parse the token
	claims, err := jwt.ParseJWT(token)
	if err != nil {
		log.Error("Failed to parse JWT", "error", err)
		return
	}
	log.Info("JWT Token validated", "user_id", claims.UserID, "username", claims.Username, "expires_at", claims.ExpiresAt)

	// 4. Demonstrate error handling with invalid token
	invalidToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.invalid.signature"
	_, err = jwt.ParseJWT(invalidToken)
	if err != nil {
		log.Info("Invalid token correctly rejected", "error", err)
	}

	log.Info("=== JWT Demo Complete ===")
}
