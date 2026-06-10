package jwtblacklist

import (
	"context"
	"fmt"
	"time"

	"github.com/stringptr/SiGizi/backend/internal/domain/jwtblacklist"
	appnats "github.com/stringptr/SiGizi/backend/internal/infrastructure/nats"
)

type Repo struct {
	kv *appnats.KV
}

func NewRepo(kv *appnats.KV) *Repo {
	return &Repo{kv: kv}
}

func (r *Repo) key(jti string) string {
	return jti
}

func (r *Repo) Blacklist(ctx context.Context, jti string, userID int32, reason string, ttl time.Duration) error {
	entry := &jwtblacklist.BlacklistEntry{
		JTI:           jti,
		UserID:        userID,
		Reason:        reason,
		BlacklistedAt: time.Now(),
		ExpiresAt:     time.Now().Add(ttl),
	}
	return r.kv.Put(ctx, r.key(jti), entry)
}

func (r *Repo) IsBlacklisted(ctx context.Context, jti string) (bool, error) {
	data, err := r.kv.Get(ctx, r.key(jti))
	if err != nil {
		return false, fmt.Errorf("check blacklisted jti: %w", err)
	}
	return data != nil, nil
}

func (r *Repo) Remove(ctx context.Context, jti string) error {
	return r.kv.Delete(ctx, r.key(jti))
}
