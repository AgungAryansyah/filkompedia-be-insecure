package entity

import "github.com/google/uuid"

type Checkout struct {
	Id     uuid.UUID `json:"id" db:"id"`
	UserID uuid.UUID `json:"user_id" db:"user_id"`
}
