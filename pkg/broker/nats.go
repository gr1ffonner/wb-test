package broker

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/nats-io/nats.go"
)

type NATSClient struct {
	conn *nats.Conn
}

func NewNATS(ctx context.Context, url string) (*NATSClient, error) {
	opts := []nats.Option{
		nats.Name("wb-app"),
		nats.Timeout(10 * time.Second),
		nats.ReconnectWait(2 * time.Second),
		nats.MaxReconnects(5),
		nats.ReconnectJitter(100*time.Millisecond, 1*time.Second),
	}

	conn, err := nats.Connect(url, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to NATS: %w", err)
	}

	// Test the connection
	if !conn.IsConnected() {
		conn.Close()
		return nil, fmt.Errorf("failed to establish NATS connection")
	}

	return &NATSClient{conn: conn}, nil
}

func (n *NATSClient) Close() {
	if n.conn != nil {
		n.conn.Close()
	}
}

// PublishOrder publishes an order to the specified subject
func (n *NATSClient) PublishOrder(subject string, order interface{}) error {
	data, err := json.Marshal(order)
	if err != nil {
		return fmt.Errorf("failed to marshal order: %w", err)
	}

	err = n.conn.Publish(subject, data)
	if err != nil {
		return fmt.Errorf("failed to publish order: %w", err)
	}

	return nil
}

// SubscribeToOrders subscribes to orders on the specified subject
func (n *NATSClient) SubscribeToOrders(subject string, handler func([]byte) error) (*nats.Subscription, error) {
	sub, err := n.conn.Subscribe(subject, func(msg *nats.Msg) {
		if err := handler(msg.Data); err != nil {
			// Log error but don't ack the message to allow retry
			fmt.Printf("Error processing message: %v\n", err)
			return
		}
		// Acknowledge the message
		msg.Ack()
	})
	if err != nil {
		return nil, fmt.Errorf("failed to subscribe: %w", err)
	}

	return sub, nil
}

// SubscribeToOrdersWithQueue subscribes to orders with queue group for load balancing
func (n *NATSClient) SubscribeToOrdersWithQueue(subject, queueGroup string, handler func([]byte) error) (*nats.Subscription, error) {
	sub, err := n.conn.QueueSubscribe(subject, queueGroup, func(msg *nats.Msg) {
		if err := handler(msg.Data); err != nil {
			// Log error but don't ack the message to allow retry
			fmt.Printf("Error processing message: %v\n", err)
			return
		}
		// Acknowledge the message
		msg.Ack()
	})
	if err != nil {
		return nil, fmt.Errorf("failed to subscribe with queue: %w", err)
	}

	return sub, nil
}
