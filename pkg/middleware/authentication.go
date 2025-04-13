package middleware

import (
	"github.com/AgungAryansyah/filkompedia-be-unsecure/pkg/response"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func (m *middleware) Authenticate(ctx *fiber.Ctx) error {
	tokenString := ctx.Cookies("token")
	if tokenString == "" {
		return &response.InvalidToken
	}

	userId, err := m.jwtAuth.ValidateToken(tokenString)
	if err != nil {
		return &response.InvalidToken
	}

	if userId != uuid.Nil {
		ctx.Locals("userId", userId)
	}

	return ctx.Next()
}
