package entity

import (
	"time"

	"github.com/google/uuid"
)

type Payment struct {
	Id         uuid.UUID `json:"id" db:"id"`
	Token      uuid.UUID `json:"token" db:"token"`
	UserId     uuid.UUID `json:"user_id" db:"user_id"`
	CheckoutId uuid.UUID `json:"checkout_id" db:"checkout_id"`
	TotalPrice float64   `json:"total_price" db:"total_price"`
	StatusId   int       `json:"status_id" db:"status_id"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
}

type PaymentStatus struct {
	Id     int    `json:"id" db:"id"`
	Status string `json:"status" db:"status"`
}
