package bannedip

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/stringptr/SiGizi/backend/internal/domain/bannedip"
	appnats "github.com/stringptr/SiGizi/backend/internal/infrastructure/nats"
)

type Repo struct {
	kv *appnats.KV
}

func NewRepo(kv *appnats.KV) *Repo {
	return &Repo{kv: kv}
}

func (r *Repo) key(ip string) string {
	return ip
}

func (r *Repo) Ban(ctx context.Context, ip string, reason string, ttl time.Duration) error {
	info := &bannedip.BanInfo{
		Reason:    reason,
		Attempts:  5,
		BannedAt:  time.Now(),
		ExpiresAt: time.Now().Add(ttl),
	}
	return r.kv.Put(ctx, r.key(ip), info)
}

func (r *Repo) IsBanned(ctx context.Context, ip string) (bool, error) {
	data, err := r.kv.Get(ctx, r.key(ip))
	if err != nil {
		return false, fmt.Errorf("check banned ip: %w", err)
	}
	if data == nil {
		return false, nil
	}
	var info bannedip.BanInfo
	if err := json.Unmarshal(data, &info); err != nil {
		return false, nil
	}
	if time.Now().After(info.ExpiresAt) {
		r.kv.Delete(ctx, r.key(ip))
		return false, nil
	}
	return true, nil
}

func (r *Repo) GetBanInfo(ctx context.Context, ip string) (*bannedip.BanInfo, error) {
	data, err := r.kv.Get(ctx, r.key(ip))
	if err != nil {
		return nil, fmt.Errorf("get ban info: %w", err)
	}
	if data == nil {
		return nil, nil
	}
	var info bannedip.BanInfo
	if err := json.Unmarshal(data, &info); err != nil {
		return nil, fmt.Errorf("unmarshal ban info: %w", err)
	}
	if time.Now().After(info.ExpiresAt) {
		r.kv.Delete(ctx, r.key(ip))
		return nil, nil
	}
	return &info, nil
}

func (r *Repo) Unban(ctx context.Context, ip string) error {
	return r.kv.Delete(ctx, r.key(ip))
}

func (r *Repo) IncrementAttempt(ctx context.Context, ip string, maxAttempts int, banTTL time.Duration) (int, error) {
	data, err := r.kv.Get(ctx, r.key(ip))
	attempts := 1
	expiresAt := time.Now().Add(banTTL)
	if err == nil && data != nil {
		var info bannedip.BanInfo
		if err := json.Unmarshal(data, &info); err == nil {
			attempts = info.Attempts + 1
			if info.ExpiresAt.After(time.Now()) {
				expiresAt = info.ExpiresAt
			}
		}
	}
	if attempts >= maxAttempts {
		return attempts, r.Ban(ctx, ip, "max_attempts_exceeded", banTTL)
	}
	info := &bannedip.BanInfo{
		Reason:    "failed_attempts",
		Attempts:  attempts,
		ExpiresAt: expiresAt,
	}
	return attempts, r.kv.Put(ctx, r.key(ip), info)
}

func (r *Repo) ClearAttempts(ctx context.Context, ip string) error {
	return r.kv.Delete(ctx, r.key(ip))
}
