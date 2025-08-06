package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	ordercache "wb-test/internal/cache/order"
	orderconsumer "wb-test/internal/consumers/order"
	handler "wb-test/internal/handlers"
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

	slog.Info("Config and logger initialized")

	ctx := context.Background()

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

	handlers := handler.NewHandler()
	router := handler.InitRouter(handlers)

	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.Server.Port),
		Handler: router,
	}

	// Start the consumer in a goroutine
	go func() {
		if err := orderConsumer.Start(ctx); err != nil {
			log.Error("Consumer failed", "error", err)
			cancel()
		}
	}()

	// Start HTTP server in a goroutine
	go func() {
		log.Info("Starting HTTP server", "port", cfg.Server.Port, "addr", httpServer.Addr)
		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("HTTP server error", "error", err)
			cancel()
		} else if err != nil {
			log.Info("HTTP server closed normally", "error", err)
		}
	}()

	log.Info("All servers are ready to handle requests")

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	select {
	case sig := <-stop:
		log.Info("Interrupt signal received", "signal", sig)
	case <-ctx.Done():
		log.Info("Context canceled")
	}

	log.Info("Shutting down servers...")

	// Cancel context to stop consumer
	cancel()
	log.Info("Order consumer shutdown initiated")

	// Create shutdown context with timeout for HTTP server
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	// Shutdown HTTP server gracefully
	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		log.Error("HTTP server shutdown failed", "error", err)
	} else {
		log.Info("HTTP server stopped gracefully")
	}

	log.Info("Application shutdown complete")
}
