package testutils

import (
	"context"
	"time"

	"github.com/stringptr/SiGizi/backend/internal/domain/bannedip"
	"github.com/stringptr/SiGizi/backend/internal/domain/jwtblacklist"
	"github.com/stringptr/SiGizi/backend/internal/domain/notification"
)

type NoopBlacklistRepo struct{}

func (n *NoopBlacklistRepo) Blacklist(_ context.Context, _ string, _ int32, _ string, _ time.Duration) error {
	return nil
}

func (n *NoopBlacklistRepo) IsBlacklisted(_ context.Context, _ string) (bool, error) {
	return false, nil
}

func (n *NoopBlacklistRepo) Remove(_ context.Context, _ string) error {
	return nil
}

var _ jwtblacklist.Repo = (*NoopBlacklistRepo)(nil)

type NoopBanRepo struct{}

func (n *NoopBanRepo) Ban(_ context.Context, _ string, _ string, _ time.Duration) error {
	return nil
}

func (n *NoopBanRepo) IsBanned(_ context.Context, _ string) (bool, error) {
	return false, nil
}

func (n *NoopBanRepo) GetBanInfo(_ context.Context, _ string) (*bannedip.BanInfo, error) {
	return nil, nil
}

func (n *NoopBanRepo) Unban(_ context.Context, _ string) error {
	return nil
}

func (n *NoopBanRepo) IncrementAttempt(_ context.Context, _ string, _ int, _ time.Duration) (int, error) {
	return 0, nil
}

func (n *NoopBanRepo) ClearAttempts(_ context.Context, _ string) error {
	return nil
}

var _ bannedip.Repo = (*NoopBanRepo)(nil)

type NoopNotifPublisher struct{}

func (n *NoopNotifPublisher) PublishToUser(_ int32, _ *notification.Notification) error { return nil }
func (n *NoopNotifPublisher) PublishToRole(_ string, _ *notification.Notification) error { return nil }
func (n *NoopNotifPublisher) PublishBroadcast(_ *notification.Notification) error        { return nil }

var _ notification.Publisher = (*NoopNotifPublisher)(nil)
