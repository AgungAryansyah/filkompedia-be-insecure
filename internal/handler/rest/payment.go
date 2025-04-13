package rest

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/yogarn/filkompedia-be/model"
	"github.com/yogarn/filkompedia-be/pkg/response"
)

func (r *Rest) HandleMidtransWebhook(ctx *fiber.Ctx) error {
	var PaymentDetails map[string]any
	if err := ctx.BodyParser(&PaymentDetails); err != nil {
		return err
	}

	if err := r.service.PaymentService.UpdatePaymentStatus(PaymentDetails); err != nil {
		return err
	}

	response.Success(ctx, http.StatusOK, "success", nil)
	return nil
}

func (r *Rest) GetPayment(ctx *fiber.Ctx) error {
	paymentIdString := ctx.Params("id")
	paymentId, err := uuid.Parse(paymentIdString)
	if err != nil {
		return err
	}

	payment, err := r.service.PaymentService.GetPayment(paymentId)
	if err != nil {
		return err
	}

	response.Success(ctx, http.StatusOK, "success", payment)
	return nil
}

func (r *Rest) CheckUserBookPurchase(ctx *fiber.Ctx) error {
	userId, ok := ctx.Locals("userId").(uuid.UUID)
	if !ok {
		return &response.Unauthorized
	}

	bookIdString := ctx.Params("id")
	bookId, err := uuid.Parse(bookIdString)
	if err != nil {
		return err
	}

	purchase, err := r.service.PaymentService.CheckUserBookPurchase(userId, bookId)
	if err != nil {
		return err
	}

	response.Success(ctx, http.StatusOK, "success", purchase)
	return nil
}

func (r *Rest) GetPayments(ctx *fiber.Ctx) error {
	var req model.PaymentReq
	req.Page = ctx.QueryInt("page", 1)
	req.PageSize = ctx.QueryInt("size", 9)

	payments, err := r.service.PaymentService.GetPayments(req)
	if err != nil {
		return err
	}

	response.Success(ctx, http.StatusOK, "success", payments)
	return nil
}

func (r *Rest) GetPaymentByCheckout(ctx *fiber.Ctx) error {
	checkoutIdString := ctx.Params("id")
	checkoutId, err := uuid.Parse(checkoutIdString)
	if err != nil {
		return err
	}

	payment, err := r.service.PaymentService.GetPaymentByCheckout(checkoutId)
	if err != nil {
		return nil
	}

	response.Success(ctx, http.StatusOK, "success", payment)
	return nil
}

func (r *Rest) GetPaymentByUser(ctx *fiber.Ctx) error {
	userId, ok := ctx.Locals("userId").(uuid.UUID)
	if !ok {
		return &response.Unauthorized
	}

	payments, err := r.service.PaymentService.GetPaymentByUser(userId)
	if err != nil {
		return nil
	}

	response.Success(ctx, http.StatusOK, "success", payments)
	return nil
}
