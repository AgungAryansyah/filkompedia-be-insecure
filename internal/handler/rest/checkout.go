package rest

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/yogarn/filkompedia-be/model"
	"github.com/yogarn/filkompedia-be/pkg/response"
)

func (r *Rest) GetUserCheckouts(ctx *fiber.Ctx) error {
	userId, ok := ctx.Locals("userId").(uuid.UUID)
	if !ok {
		return &response.Unauthorized
	}

	checkouts, err := r.service.CheckoutService.GetUserCheckouts(userId)
	if err != nil {
		return err
	}

	response.Success(ctx, http.StatusOK, "success", checkouts)
	return nil
}

func (r *Rest) GetUserCheckoutsAdmin(ctx *fiber.Ctx) error {
	param := ctx.Params("userId")
	userId, err := uuid.Parse(param)
	if err != nil {
		return err
	}

	checkouts, err := r.service.CheckoutService.GetUserCheckouts(userId)
	if err != nil {
		return err
	}

	response.Success(ctx, http.StatusOK, "success", checkouts)
	return nil
}

func (r *Rest) GetCheckoutCarts(ctx *fiber.Ctx) error {
	param := ctx.Params("checkoutId")
	checkoutId, err := uuid.Parse(param)
	if err != nil {
		return err
	}

	carts, err := r.service.CheckoutService.GetCheckoutCarts(checkoutId)
	if err != nil {
		return err
	}

	response.Success(ctx, http.StatusOK, "success", carts)
	return nil
}

func (r *Rest) Checkout(ctx *fiber.Ctx) error {
	var request model.CheckoutRequest
	if err := ctx.BodyParser(&request); err != nil {
		return err
	}

	userId, ok := ctx.Locals("userId").(uuid.UUID)
	if !ok {
		return &response.Unauthorized
	}

	checkoutId := uuid.New()
	totalPrice, err := r.service.CheckoutService.Checkout(request, userId, checkoutId)
	if err != nil {
		return err
	}

	snapRes, err := r.service.PaymentService.CreatePayment(userId, checkoutId, totalPrice)
	if err != nil {
		return err
	}

	response.Success(ctx, http.StatusOK, "success", snapRes)
	return nil
}
