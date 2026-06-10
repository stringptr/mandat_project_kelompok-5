package nats

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/nats-io/nats.go"
)

type PubSub struct {
	conn *nats.Conn
}

func NewPubSub(conn *nats.Conn) *PubSub {
	return &PubSub{conn: conn}
}

func (ps *PubSub) Publish(ctx context.Context, subject string, data any) error {
	msg, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("pubsub marshal: %w", err)
	}
	return ps.conn.Publish(subject, msg)
}

func (ps *PubSub) Subscribe(subject string, handler func(msg []byte)) (*nats.Subscription, error) {
	sub, err := ps.conn.Subscribe(subject, func(m *nats.Msg) {
		handler(m.Data)
	})
	if err != nil {
		return nil, fmt.Errorf("pubsub subscribe %s: %w", subject, err)
	}
	return sub, nil
}

type QueueHandler func(msg []byte)

func (ps *PubSub) QueueSubscribe(subject, queue string, handler QueueHandler) (*nats.Subscription, error) {
	sub, err := ps.conn.QueueSubscribe(subject, queue, func(m *nats.Msg) {
		handler(m.Data)
	})
	if err != nil {
		return nil, fmt.Errorf("pubsub queue subscribe %s/%s: %w", subject, queue, err)
	}
	return sub, nil
}
