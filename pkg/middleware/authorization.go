package middleware

import (
	"slices"

	"github.com/AgungAryansyah/filkompedia-be-unsecure/entity"
	"github.com/AgungAryansyah/filkompedia-be-unsecure/pkg/response"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func (m *middleware) Authorize(roles []int) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		userId, ok := ctx.Locals("userId").(uuid.UUID)
		if !ok {
			return &response.RoleUnauthorized
		}

		var user entity.User

		err := m.service.UserService.GetUserById(&user, userId)
		if err != nil {
			return err
		}

		if slices.Contains(roles, user.RoleId) {
			return ctx.Next()
		}

		return &response.RoleUnauthorized
	}
}

func (m *middleware) AuthorizeOrItself(roles []int) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		userId, ok := ctx.Locals("userId").(uuid.UUID)

		if !ok {
			return &response.RoleUnauthorized
		}

		userIdParamStr := ctx.Params("userId")
		userIdParam, err := uuid.Parse(userIdParamStr)
		if err != nil {
			return &response.RoleUnauthorized
		}

		var user entity.User

		err = m.service.UserService.GetUserById(&user, userId)
		if err != nil {
			return err
		}

		if slices.Contains(roles, user.RoleId) || userId == userIdParam {
			return ctx.Next()
		}

		return &response.RoleUnauthorized
	}
}
