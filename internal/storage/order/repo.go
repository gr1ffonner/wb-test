package order

import (
	"context"
	"fmt"
	"log/slog"

	"wb-test/internal/models"
	"wb-test/pkg/db"

	"github.com/jackc/pgx/v5"
)

type orderRepo struct {
	db *db.PostgresClient
}

func NewOrderRepo(db *db.PostgresClient) *orderRepo {
	return &orderRepo{db: db}
}

func (r *orderRepo) CreateOrder(order *models.Order) error {
	ctx := context.Background()

	// Start a transaction
	tx, err := r.db.Pool().Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Insert order
	orderQuery := `
		INSERT INTO orders (
			order_uid, track_number, entry, locale, internal_signature, 
			customer_id, delivery_service, shard_key, sm_id, date_created, oof_shard
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		ON CONFLICT (order_uid) DO NOTHING
	`
	_, err = tx.Exec(ctx, orderQuery,
		order.OrderUID, order.TrackNumber, order.Entry, order.Locale,
		order.InternalSignature, order.CustomerID, order.DeliveryService,
		order.ShardKey, order.SmID, order.DateCreated, order.OofShard,
	)
	if err != nil {
		return fmt.Errorf("failed to insert order: %w", err)
	}

	// Insert delivery
	deliveryQuery := `
		INSERT INTO deliveries (
			order_uid, name, phone, zip, city, address, region, email
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	_, err = tx.Exec(ctx, deliveryQuery,
		order.OrderUID, order.Delivery.Name, order.Delivery.Phone, order.Delivery.Zip,
		order.Delivery.City, order.Delivery.Address, order.Delivery.Region, order.Delivery.Email,
	)
	if err != nil {
		return fmt.Errorf("failed to insert delivery: %w", err)
	}

	// Insert payment
	paymentQuery := `
		INSERT INTO payments (
			order_uid, transaction, request_id, currency, provider, amount,
			payment_dt, bank, delivery_cost, goods_total, custom_fee
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`
	_, err = tx.Exec(ctx, paymentQuery,
		order.OrderUID, order.Payment.Transaction, order.Payment.RequestID,
		order.Payment.Currency, order.Payment.Provider, order.Payment.Amount,
		order.Payment.PaymentDt, order.Payment.Bank, order.Payment.DeliveryCost,
		order.Payment.GoodsTotal, order.Payment.CustomFee,
	)
	if err != nil {
		return fmt.Errorf("failed to insert payment: %w", err)
	}

	// Insert items
	for _, item := range order.Items {
		itemQuery := `
			INSERT INTO items (
				order_uid, chrt_id, track_number, price, rid, name, sale,
				size, total_price, nm_id, brand, status
			) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		`
		_, err = tx.Exec(ctx, itemQuery,
			order.OrderUID, item.ChrtID, item.TrackNumber, item.Price,
			item.Rid, item.Name, item.Sale, item.Size, item.TotalPrice,
			item.NmID, item.Brand, item.Status,
		)
		if err != nil {
			return fmt.Errorf("failed to insert item: %w", err)
		}
	}

	// Commit transaction
	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	slog.Info("Order saved to database", "order_uid", order.OrderUID)
	return nil
}

func (r *orderRepo) GetOrder(orderUID string) (*models.Order, error) {
	ctx := context.Background()

	// Query order
	orderQuery := `
		SELECT order_uid, track_number, entry, locale, internal_signature,
			   customer_id, delivery_service, shard_key, sm_id, date_created, oof_shard
		FROM orders WHERE order_uid = $1
	`
	var order models.Order
	err := r.db.Pool().QueryRow(ctx, orderQuery, orderUID).Scan(
		&order.OrderUID, &order.TrackNumber, &order.Entry, &order.Locale,
		&order.InternalSignature, &order.CustomerID, &order.DeliveryService,
		&order.ShardKey, &order.SmID, &order.DateCreated, &order.OofShard,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("order not found: %s", orderUID)
		}
		return nil, fmt.Errorf("failed to query order: %w", err)
	}

	// Query delivery
	deliveryQuery := `
		SELECT name, phone, zip, city, address, region, email
		FROM deliveries WHERE order_uid = $1
	`
	err = r.db.Pool().QueryRow(ctx, deliveryQuery, orderUID).Scan(
		&order.Delivery.Name, &order.Delivery.Phone, &order.Delivery.Zip,
		&order.Delivery.City, &order.Delivery.Address, &order.Delivery.Region, &order.Delivery.Email,
	)
	if err != nil && err != pgx.ErrNoRows {
		return nil, fmt.Errorf("failed to query delivery: %w", err)
	}

	// Query payment
	paymentQuery := `
		SELECT transaction, request_id, currency, provider, amount,
			   payment_dt, bank, delivery_cost, goods_total, custom_fee
		FROM payments WHERE order_uid = $1
	`
	err = r.db.Pool().QueryRow(ctx, paymentQuery, orderUID).Scan(
		&order.Payment.Transaction, &order.Payment.RequestID, &order.Payment.Currency,
		&order.Payment.Provider, &order.Payment.Amount, &order.Payment.PaymentDt,
		&order.Payment.Bank, &order.Payment.DeliveryCost, &order.Payment.GoodsTotal, &order.Payment.CustomFee,
	)
	if err != nil && err != pgx.ErrNoRows {
		return nil, fmt.Errorf("failed to query payment: %w", err)
	}

	// Query items
	itemsQuery := `
		SELECT chrt_id, track_number, price, rid, name, sale,
			   size, total_price, nm_id, brand, status
		FROM items WHERE order_uid = $1
	`
	rows, err := r.db.Pool().Query(ctx, itemsQuery, orderUID)
	if err != nil {
		return nil, fmt.Errorf("failed to query items: %w", err)
	}
	defer rows.Close()

	var items []models.Item
	for rows.Next() {
		var item models.Item
		err := rows.Scan(
			&item.ChrtID, &item.TrackNumber, &item.Price, &item.Rid,
			&item.Name, &item.Sale, &item.Size, &item.TotalPrice,
			&item.NmID, &item.Brand, &item.Status,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan item: %w", err)
		}
		items = append(items, item)
	}
	order.Items = items

	return &order, nil
}
