package notification

import "time"

type Notification struct {
	ID        int32     `json:"id"`
	UserID    int32     `json:"user_id"`
	Judul     string    `json:"judul"`
	Pesan     string    `json:"pesan"`
	Tipe      string    `json:"tipe"`
	CreatedAt time.Time `json:"created_at"`
}

type Publisher interface {
	PublishToUser(userID int32, notif *Notification) error
	PublishToRole(role string, notif *Notification) error
	PublishBroadcast(notif *Notification) error
}
