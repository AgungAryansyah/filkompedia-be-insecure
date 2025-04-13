package model

import (
	"time"

	"github.com/google/uuid"
)

type CommentRes struct {
	Id             uuid.UUID `json:"id"`
	UserId         uuid.UUID `json:"user_id"`
	Username       string    `json:"username"`
	ProfilePicture string    `json:"profilePicture"`
	BookId         uuid.UUID `json:"book_id"`
	Comment        string    `json:"comment"`
	Rating         int       `json:"rating"`
	CreatedAt      time.Time `json:"created_at"`
}

type CreateComment struct {
	BookId  uuid.UUID `json:"book_id" validate:"required"`
	Comment string    `json:"comment" validate:"required,min=5,max=500"`
	Rating  int       `json:"rating" validate:"required,gte=1,lte=5"`
}

type UpdateComment struct {
	Comment string `json:"comment" validate:"required,min=5,max=500"`
	Rating  int    `json:"rating" validate:"required,gte=1,lte=5"`
}
