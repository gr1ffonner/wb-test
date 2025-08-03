package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	ordercache "wb-test/internal/cache/order"
	orderconsumer "wb-test/internal/consumers/order"
	orderservice "wb-test/internal/service/order"
	orderstorage "wb-test/internal/storage/order"
	"wb-test/pkg/broker"
	"wb-test/pkg/cache"
	"wb-test/pkg/config"
	"wb-test/pkg/db"
	"wb-test/pkg/logger"
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

	// Initialize order repo
	orderRepo := orderstorage.NewOrderRepo(db)
	log.Info("Order repo initialized successfully")

	// Initialize order cache
	orderCache := ordercache.NewOrderCache(cache)
	log.Info("Order cache initialized successfully")

	// Initialize order service
	orderService := orderservice.NewOrderService(orderRepo, orderCache)
	log.Info("Order service initialized successfully")

	// Initialize and start order consumer
	orderConsumer := orderconsumer.NewOrderConsumer(broker, orderService)
	log.Info("Order consumer initialized successfully")

	// Create context with cancellation for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigChan
		log.Info("Received shutdown signal", "signal", sig)
		cancel()
	}()

	// Start the consumer
	log.Info("Starting order consumer...")
	if err := orderConsumer.Start(ctx); err != nil {
		log.Error("Consumer failed", "error", err)
		panic(err)
	}

	log.Info("All services initialized successfully")
	log.Info("All consumers initialized successfully")
}
