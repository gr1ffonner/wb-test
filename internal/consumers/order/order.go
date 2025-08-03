package order

import (
	"encoding/json"
	"fmt"
	"log/slog"

	"wb-test/internal/models"
)

func (oc *OrderConsumer) handleOrder(data []byte) error {
	var order models.Order
	if err := json.Unmarshal(data, &order); err != nil {
		slog.Error("Failed to unmarshal order", "error", err)
		return fmt.Errorf("failed to unmarshal order: %w", err)
	}

	slog.Info("Processing order",
		"order_uid", order.OrderUID,
		"track_number", order.TrackNumber,
		"customer_id", order.CustomerID,
	)

	// Process the order (save to DB, cache, etc.)
	if err := oc.service.ProcessOrder(&order); err != nil {
		slog.Error("Failed to process order", "error", err)
		return fmt.Errorf("failed to process order %s: %w", order.OrderUID, err)
	}

	slog.Info("Order processed successfully", "order_uid", order.OrderUID)
	return nil
}
