package entity

import "github.com/google/uuid"

type Cart struct {
	Id         uuid.UUID `json:"id" db:"id"`
	UserId     uuid.UUID `json:"user_id" db:"user_id"`
	BookId     uuid.UUID `json:"book_id" db:"book_id"`
	Amount     int       `json:"amount" db:"amount"`
	CheckoutId uuid.UUID `json:"checkout_id" db:"checkout_id"`
}
