package timestamp

import "time"

type TimestampEntity struct {
	CreatedAt time.Time  `json:"created_at" db:"created_at" validate:"required"`
	UpdatedAt time.Time  `json:"updated_at" db:"updated_at" validate:"required"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
}

type SessionTimestampEntity struct {
	TimestampEntity
	ExpiredAt time.Time `json:"expired_at" db:"expired_at" validate:"required"`
}
