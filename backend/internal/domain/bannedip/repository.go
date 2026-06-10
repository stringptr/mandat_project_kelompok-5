package bannedip

import (
	"context"
	"time"
)

type BanInfo struct {
	Reason    string    `json:"reason"`
	Attempts  int       `json:"attempts"`
	BannedAt  time.Time `json:"banned_at"`
	ExpiresAt time.Time `json:"expires_at"`
}

type Repo interface {
	Ban(ctx context.Context, ip string, reason string, ttl time.Duration) error
	IsBanned(ctx context.Context, ip string) (bool, error)
	GetBanInfo(ctx context.Context, ip string) (*BanInfo, error)
	Unban(ctx context.Context, ip string) error
	IncrementAttempt(ctx context.Context, ip string, maxAttempts int, banTTL time.Duration) (int, error)
	ClearAttempts(ctx context.Context, ip string) error
}
