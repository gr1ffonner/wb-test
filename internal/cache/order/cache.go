package order

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"wb-test/internal/models"
	"wb-test/pkg/cache"

	"github.com/redis/go-redis/v9"
)

type orderCache struct {
	client *cache.RedisClient
}

func NewOrderCache(client *cache.RedisClient) *orderCache {
	return &orderCache{client: client}
}

func (c *orderCache) GetOrder(orderUID string) (*models.Order, error) {
	ctx := context.Background()
	key := fmt.Sprintf("order:%s", orderUID)

	data, err := c.client.Client().Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			slog.Error("Order not found in cache", "order_uid", orderUID)
			return nil, nil
		}
		slog.Error("Failed to get order from cache", "error", err)
		return nil, err
	}

	var order models.Order
	if err := json.Unmarshal(data, &order); err != nil {
		return nil, fmt.Errorf("failed to unmarshal order from cache: %w", err)
	}

	slog.Info("Order retrieved from cache", "order_uid", orderUID)
	return &order, nil
}

func (c *orderCache) SetOrder(orderUID string, order *models.Order) error {
	ctx := context.Background()
	key := fmt.Sprintf("order:%s", orderUID)

	data, err := json.Marshal(order)
	if err != nil {
		return fmt.Errorf("failed to marshal order for cache: %w", err)
	}

	// Set with TTL (24 hours)
	err = c.client.Client().Set(ctx, key, data, 24*time.Hour).Err()
	if err != nil {
		return fmt.Errorf("failed to set order in cache: %w", err)
	}

	slog.Info("Order saved to cache", "order_uid", orderUID)
	return nil
}

func (c *orderCache) DeleteOrder(orderUID string) error {
	ctx := context.Background()
	key := fmt.Sprintf("order:%s", orderUID)

	err := c.client.Client().Del(ctx, key).Err()
	if err != nil {
		return fmt.Errorf("failed to delete order from cache: %w", err)
	}

	slog.Info("Order deleted from cache", "order_uid", orderUID)
	return nil
}
