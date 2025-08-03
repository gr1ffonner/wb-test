package broker

import (
	"context"
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
