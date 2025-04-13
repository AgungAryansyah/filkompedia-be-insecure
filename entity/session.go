package entity

import (
	"time"

	"github.com/google/uuid"
)

type Session struct {
	UserId    uuid.UUID `db:"user_id"`
	Token     string    `db:"token"`
	IPAddress string    `db:"ip_address"`
	ExpiresAt time.Time `db:"expires_at"`
	UserAgent string    `db:"user_agent"`
	DeviceId  string    `db:"device_id"`
}
