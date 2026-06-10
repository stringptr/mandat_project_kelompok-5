package nats

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/nats-io/nats.go/jetstream"
)

type KVEntry struct {
	Key       string
	Value     []byte
	CreatedAt time.Time
	Revision  uint64
	Delta     uint64
	Operation jetstream.KeyValueOp
}

type KV struct {
	kv jetstream.KeyValue
}

func NewKV(kv jetstream.KeyValue) *KV {
	return &KV{kv: kv}
}

func (s *KV) Get(ctx context.Context, key string) ([]byte, error) {
	entry, err := s.kv.Get(ctx, key)
	if err != nil {
		if err == jetstream.ErrKeyNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("kv get %s: %w", key, err)
	}
	return entry.Value(), nil
}

func (s *KV) Put(ctx context.Context, key string, value any) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("kv marshal %s: %w", key, err)
	}
	_, err = s.kv.Put(ctx, key, data)
	if err != nil {
		return fmt.Errorf("kv put %s: %w", key, err)
	}
	return nil
}

func (s *KV) Delete(ctx context.Context, key string) error {
	err := s.kv.Delete(ctx, key)
	if err != nil {
		return fmt.Errorf("kv delete %s: %w", key, err)
	}
	return nil
}

func (s *KV) Keys(ctx context.Context) ([]string, error) {
	keys, err := s.kv.Keys(ctx)
	if err != nil {
		return nil, fmt.Errorf("kv keys: %w", err)
	}
	return keys, nil
}

func (s *KV) Watch(ctx context.Context, key string) (<-chan jetstream.KeyValueEntry, error) {
	watcher, err := s.kv.Watch(ctx, key)
	if err != nil {
		return nil, fmt.Errorf("kv watch: %w", err)
	}
	return watcher.Updates(), nil
}

func (s *KV) WatchAll(ctx context.Context) (<-chan jetstream.KeyValueEntry, error) {
	watcher, err := s.kv.WatchAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("kv watch all: %w", err)
	}
	return watcher.Updates(), nil
}
