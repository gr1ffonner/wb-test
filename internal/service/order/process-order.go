package order

import (
	"fmt"
	"log/slog"
	"wb-test/internal/models"
)

// ProcessOrder handles the business logic for processing an order
func (s *OrderService) ProcessOrder(order *models.Order) error {
	// Save order to database
	if err := s.repo.CreateOrder(order); err != nil {
		return fmt.Errorf("failed to save order to database: %w", err)
	}

	// Save order to cache
	if err := s.cache.SetOrder(order.OrderUID, order); err != nil {
		// Log cache error but don't fail the process
		slog.Error("Failed to save order to cache", "error", err, "order_uid", order.OrderUID)
	}

	slog.Info("Order processed successfully",
		"order_uid", order.OrderUID,
		"items_count", len(order.Items),
		"total_amount", order.Payment.Amount,
	)

	return nil
}
