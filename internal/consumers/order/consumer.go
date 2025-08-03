package order

import (
	"context"
	"fmt"
	"log/slog"

	"wb-test/internal/models"
	"wb-test/pkg/broker"
)

const (
	OrderSubject = "orders.new"
	QueueGroup   = "order-processors"
)

type OrderService interface {
	ProcessOrder(order *models.Order) error
}

type OrderConsumer struct {
	broker  *broker.NATSClient
	service OrderService
}

func NewOrderConsumer(broker *broker.NATSClient, service OrderService) *OrderConsumer {
	return &OrderConsumer{
		broker:  broker,
		service: service,
	}
}

// Start starts the order consumer
func (oc *OrderConsumer) Start(ctx context.Context) error {
	slog.Info("Starting order consumer", "subject", OrderSubject, "queue_group", QueueGroup)
	// Subscribe to orders with queue group for load balancing
	sub, err := oc.broker.SubscribeToOrdersWithQueue(OrderSubject, QueueGroup, oc.handleOrder)
	if err != nil {
		return fmt.Errorf("failed to subscribe to orders: %w", err)
	}

	// Wait for context cancellation
	<-ctx.Done()

	// Unsubscribe when context is cancelled
	if err := sub.Unsubscribe(); err != nil {
		slog.Error("Failed to unsubscribe", "error", err)
	}

	slog.Info("Order consumer stopped")
	return nil
}
