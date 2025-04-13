package model

import "github.com/google/uuid"

type CartParam struct {
	Id     uuid.UUID `json:"id" validate:"uuid"`
	UserId uuid.UUID `json:"user_id" validate:"uuid"`
}

type AddToCart struct {
	BookId uuid.UUID `json:"book_id" validate:"required,uuid"`
	Amount int       `json:"amount" validate:"required,min=1"`
}

type EditCart struct {
	CartId uuid.UUID `json:"cart_id" validate:"required,uuid"`
	Amount int       `json:"amount" validate:"required"`
}
