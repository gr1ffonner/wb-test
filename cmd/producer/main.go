package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"wb-test/internal/models"
	"wb-test/pkg/broker"
	"wb-test/pkg/config"
	"wb-test/pkg/logger"
)

const (
	OrderSubject = "orders.new"
)

func main() {
	// Parse command line flags
	var (
		count    = flag.Int("count", 1, "Number of orders to publish")
		interval = flag.Duration("interval", 1*time.Second, "Interval between orders")
		subject  = flag.String("subject", OrderSubject, "NATS subject to publish to")
	)
	flag.Parse()

	// Load config
	cfg, err := config.Load()
	if err != nil {
		slog.Error("Failed to load config", "error", err)
		os.Exit(1)
	}

	// Initialize logger
	logger.InitLogger(cfg.Logger)
	log := slog.Default()

	ctx := context.Background()

	// Initialize NATS broker
	broker, err := broker.NewNATS(ctx, cfg.NATS.URL)
	if err != nil {
		log.Error("Failed to connect to NATS", "error", err)
		os.Exit(1)
	}
	defer broker.Close()

	log.Info("Connected to NATS", "url", cfg.NATS.URL)

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

	// Start publishing orders
	log.Info("Starting to publish orders",
		"count", *count,
		"interval", *interval,
		"subject", *subject,
	)

	published := 0
	ticker := time.NewTicker(*interval)
	defer ticker.Stop()

	for i := 0; i < *count; i++ {
		select {
		case <-ctx.Done():
			log.Info("Shutting down producer")
			return
		case <-ticker.C:
			order := generateSampleOrder(i)

			if err := broker.PublishOrder(*subject, order); err != nil {
				log.Error("Failed to publish order", "error", err, "order_uid", order.OrderUID)
				continue
			}

			published++
			log.Info("Published order",
				"order_uid", order.OrderUID,
				"track_number", order.TrackNumber,
				"published_count", published,
			)
		}
	}

	log.Info("Finished publishing orders", "total_published", published)
}

// generateSampleOrder creates a sample order with unique identifiers
func generateSampleOrder(index int) *models.Order {
	now := time.Now()
	orderUID := fmt.Sprintf("b563feb7b2b84b6test%d", index)
	trackNumber := fmt.Sprintf("WBILMTESTTRACK%d", index)

	return &models.Order{
		OrderUID:    orderUID,
		TrackNumber: trackNumber,
		Entry:       "WBIL",
		Delivery: models.Delivery{
			Name:    "Test Testov",
			Phone:   "+9720000000",
			Zip:     "2639809",
			City:    "Kiryat Mozkin",
			Address: "Ploshad Mira 15",
			Region:  "Kraiot",
			Email:   "test@gmail.com",
		},
		Payment: models.Payment{
			Transaction:  orderUID,
			RequestID:    "",
			Currency:     "USD",
			Provider:     "wbpay",
			Amount:       1817,
			PaymentDt:    now.Unix(),
			Bank:         "alpha",
			DeliveryCost: 1500,
			GoodsTotal:   317,
			CustomFee:    0,
		},
		Items: []models.Item{
			{
				ChrtID:      9934930 + index,
				TrackNumber: trackNumber,
				Price:       453,
				Rid:         fmt.Sprintf("ab4219087a764ae0btest%d", index),
				Name:        "Mascaras",
				Sale:        30,
				Size:        "0",
				TotalPrice:  317,
				NmID:        2389212 + index,
				Brand:       "Vivienne Sabo",
				Status:      202,
			},
		},
		Locale:            "en",
		InternalSignature: "",
		CustomerID:        "test",
		DeliveryService:   "meest",
		ShardKey:          "9",
		SmID:              99,
		DateCreated:       now,
		OofShard:          "1",
	}
}
