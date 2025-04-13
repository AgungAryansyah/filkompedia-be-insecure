package rest

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/yogarn/filkompedia-be/entity"
	"github.com/yogarn/filkompedia-be/model"
	"github.com/yogarn/filkompedia-be/pkg/response"
)

func (r *Rest) GetUserCart(ctx *fiber.Ctx) error {
	userId, ok := ctx.Locals("userId").(uuid.UUID)
	if !ok {
		return &response.Unauthorized
	}

	var carts []entity.Cart
	if err := r.service.CartService.GetUserCart(&carts, userId); err != nil {
		return err
	}

	response.Success(ctx, http.StatusOK, "success", carts)
	return nil
}

func (r *Rest) GetUserCartAdmin(ctx *fiber.Ctx) error {
	param := ctx.Params("userId")
	userId, err := uuid.Parse(param)
	if err != nil {
		return err
	}

	var carts []entity.Cart
	if err := r.service.CartService.GetUserCart(&carts, userId); err != nil {
		return err
	}

	response.Success(ctx, http.StatusOK, "success", carts)
	return nil
}

func (r *Rest) GetCart(ctx *fiber.Ctx) error {
	param := ctx.Params("cartId")
	cartId, err := uuid.Parse(param)
	if err != nil {
		return err
	}

	var cart entity.Cart
	if err := r.service.CartService.GetCart(&cart, cartId); err != nil {
		return err
	}

	response.Success(ctx, http.StatusOK, "success", cart)
	return nil
}

func (r *Rest) AddToCart(ctx *fiber.Ctx) error {
	var add model.AddToCart
	if err := ctx.BodyParser(&add); err != nil {
		return err
	}

	userId, ok := ctx.Locals("userId").(uuid.UUID)
	if !ok {
		return &response.Unauthorized
	}

	if err := r.validator.Struct(add); err != nil {
		return err
	}

	if err := r.service.CartService.AddToCart(add, userId); err != nil {
		return err
	}

	response.Success(ctx, http.StatusOK, "success", nil)
	return nil
}

func (r *Rest) EditCart(ctx *fiber.Ctx) error {
	var edit model.EditCart
	if err := ctx.BodyParser(&edit); err != nil {
		return err
	}

	userId, ok := ctx.Locals("userId").(uuid.UUID)
	if !ok {
		return &response.Unauthorized
	}

	if err := r.validator.Struct(edit); err != nil {
		return err
	}

	if err := r.service.CartService.EditCart(edit, userId); err != nil {
		return err
	}

	response.Success(ctx, http.StatusOK, "success", nil)
	return nil
}

func (r *Rest) RemoveFromCart(ctx *fiber.Ctx) error {
	param := ctx.Params("cartId")
	cartId, err := uuid.Parse(param)
	if err != nil {
		return err
	}

	if err := r.service.CartService.RemoveFromCart(cartId); err != nil {
		return err
	}

	response.Success(ctx, http.StatusOK, "success", nil)
	return nil
}
