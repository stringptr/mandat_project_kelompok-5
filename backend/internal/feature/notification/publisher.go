package notification

import (
	"context"
	"fmt"

	"github.com/stringptr/SiGizi/backend/internal/domain/notification"
	appnats "github.com/stringptr/SiGizi/backend/internal/infrastructure/nats"
)

type Publisher struct {
	ps *appnats.PubSub
}

func NewPublisher(ps *appnats.PubSub) *Publisher {
	return &Publisher{ps: ps}
}

func (p *Publisher) PublishToUser(userID int32, notif *notification.Notification) error {
	subject := fmt.Sprintf("notif.user.%d", userID)
	return p.ps.Publish(context.Background(), subject, notif)
}

func (p *Publisher) PublishToRole(role string, notif *notification.Notification) error {
	subject := fmt.Sprintf("notif.role.%s", role)
	return p.ps.Publish(context.Background(), subject, notif)
}

func (p *Publisher) PublishBroadcast(notif *notification.Notification) error {
	return p.ps.Publish(context.Background(), "notif.broadcast", notif)
}
