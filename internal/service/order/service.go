package order

import (
	"wb-test/internal/models"
)

type OrderRepo interface {
	CreateOrder(order *models.Order) error
	GetOrder(orderUID string) (*models.Order, error)
}

type OrderCache interface {
	GetOrder(orderUID string) (*models.Order, error)
	SetOrder(orderUID string, order *models.Order) error
}

// OrderServiceImpl implements the OrderService interface
type OrderService struct {
	repo  OrderRepo
	cache OrderCache
}

// NewOrderService creates a new order service instance
func NewOrderService(repo OrderRepo, cache OrderCache) *OrderService {
	return &OrderService{
		repo:  repo,
		cache: cache,
	}
}
