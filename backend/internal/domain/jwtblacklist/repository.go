package jwtblacklist

import (
	"context"
	"time"
)

type BlacklistEntry struct {
	JTI          string    `json:"jti"`
	UserID       int32     `json:"user_id"`
	Reason       string    `json:"reason"`
	BlacklistedAt time.Time `json:"blacklisted_at"`
	ExpiresAt    time.Time `json:"expires_at"`
}

type Repo interface {
	Blacklist(ctx context.Context, jti string, userID int32, reason string, ttl time.Duration) error
	IsBlacklisted(ctx context.Context, jti string) (bool, error)
	Remove(ctx context.Context, jti string) error
}
