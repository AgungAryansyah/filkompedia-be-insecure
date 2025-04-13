package middleware

import (
	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/yogarn/filkompedia-be/model"
	"github.com/yogarn/filkompedia-be/pkg/response"
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

	validate := validator.New()
	if err := validate.Struct(createReq); err != nil {
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
