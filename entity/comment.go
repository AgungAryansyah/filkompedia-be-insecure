package entity

import (
	"time"

	"github.com/google/uuid"
)

type Comment struct {
	Id        uuid.UUID `json:"id" db:"id"`
	UserId    uuid.UUID `json:"user_id" db:"user_id"`
	BookId    uuid.UUID `json:"book_id" db:"book_id"`
	Comment   string    `json:"comment" db:"comment"`
	Rating    int       `json:"rating" db:"rating"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}
