package middleware

import (
	"github.com/AgungAryansyah/filkompedia-be-unsecure/model"
	"github.com/AgungAryansyah/filkompedia-be-unsecure/pkg/response"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func (m *middleware) BookCommentCheck(ctx *fiber.Ctx) error {
	userId, ok := ctx.Locals("userId").(uuid.UUID)
	if !ok {
		return &response.RoleUnauthorized
	}

	createReq := &model.CreateComment{}
	if err := ctx.BodyParser(createReq); err != nil {
		return err
	}

	purchase, err := m.service.PaymentService.CheckUserBookPurchase(userId, createReq.BookId)
	if err != nil {
		return err
	}

	status := *purchase

	if !status {
		return &response.Forbidden
	}

	return ctx.Next()
}
