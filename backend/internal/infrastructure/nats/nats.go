package nats

import (
	"context"
	"fmt"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

type Conn struct {
	conn *nats.Conn
	js   jetstream.JetStream
}

func Connect(url, token string) (*Conn, error) {
	opts := []nats.Option{
		nats.Name("SiGizi"),
		nats.Timeout(10 * time.Second),
		nats.ReconnectWait(2 * time.Second),
		nats.MaxReconnects(-1),
		nats.RetryOnFailedConnect(true),
	}
	if token != "" {
		opts = append(opts, nats.Token(token))
	}

	nc, err := nats.Connect(url, opts...)
	if err != nil {
		return nil, fmt.Errorf("nats connect: %w", err)
	}

	js, err := jetstream.New(nc)
	if err != nil {
		nc.Close()
		return nil, fmt.Errorf("nats jetstream: %w", err)
	}

	return &Conn{conn: nc, js: js}, nil
}

func (c *Conn) Close() {
	if c.conn != nil {
		c.conn.Close()
	}
}

func (c *Conn) JetStream() jetstream.JetStream {
	return c.js
}

func (c *Conn) Conn() *nats.Conn {
	return c.conn
}

func (c *Conn) IsConnected() bool {
	return c.conn != nil && c.conn.IsConnected()
}

func (c *Conn) HealthCheck(ctx context.Context) error {
	if !c.IsConnected() {
		return fmt.Errorf("nats not connected")
	}
	return nil
}

func (c *Conn) CreateKeyValue(ctx context.Context, bucket string, ttl time.Duration) (jetstream.KeyValue, error) {
	kv, err := c.js.CreateKeyValue(ctx, jetstream.KeyValueConfig{
		Bucket: bucket,
		TTL:    ttl,
	})
	if err != nil {
		return nil, fmt.Errorf("create kv bucket %s: %w", bucket, err)
	}
	return kv, nil
}
