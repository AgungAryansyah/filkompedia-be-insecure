package model

import "github.com/google/uuid"

type CheckoutRequest struct {
	CartsId []uuid.UUID `json:"carts_id" validate:"required"`
}
