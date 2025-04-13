package rest

import (
	"net/http"

	"github.com/AgungAryansyah/filkompedia-be-unsecure/model"
	"github.com/AgungAryansyah/filkompedia-be-unsecure/pkg/response"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func (r *Rest) GetComment(ctx *fiber.Ctx) error {
	commentIdString := ctx.Params("id")
	commentId, err := uuid.Parse(commentIdString)
	if err != nil {
		return err
	}

	comment, err := r.service.CommentService.GetComment(commentId)
	if err != nil {
		return err
	}

	response.Success(ctx, http.StatusOK, "success", comment)
	return nil
}

func (r *Rest) GetCommentByBook(ctx *fiber.Ctx) error {
	bookIdString := ctx.Params("bookId")
	bookId, err := uuid.Parse(bookIdString)
	if err != nil {
		return err
	}

	comment, err := r.service.CommentService.GetCommentByBook(bookId)
	if err != nil {
		return err
	}

	response.Success(ctx, http.StatusOK, "success", comment)
	return nil
}

func (r *Rest) CreateComment(ctx *fiber.Ctx) error {
	createReq := &model.CreateComment{}
	if err := ctx.BodyParser(createReq); err != nil {
		return err
	}

	userId, ok := ctx.Locals("userId").(uuid.UUID)
	if !ok {
		return &response.Unauthorized
	}

	if err := r.service.CommentService.CreateComment(createReq, userId); err != nil {
		return err
	}

	response.Success(ctx, http.StatusOK, "success", nil)
	return nil
}

func (r *Rest) UpdateComment(ctx *fiber.Ctx) error {
	commentIdString := ctx.Params("id")
	commentId, err := uuid.Parse(commentIdString)
	if err != nil {
		return err
	}

	bookIdString := ctx.Params("bookId")
	bookId, err := uuid.Parse(bookIdString)
	if err != nil {
		return err
	}

	updateReq := &model.UpdateComment{}
	if err := ctx.BodyParser(updateReq); err != nil {
		return err
	}

	userId, ok := ctx.Locals("userId").(uuid.UUID)
	if !ok {
		return &response.Unauthorized
	}

	if err := r.service.CommentService.UpdateComment(updateReq, userId, bookId, commentId); err != nil {
		return err
	}

	response.Success(ctx, http.StatusOK, "success", nil)
	return nil
}

func (r *Rest) DeleteComment(ctx *fiber.Ctx) error {
	commentIdString := ctx.Params("id")
	commentId, err := uuid.Parse(commentIdString)
	if err != nil {
		return err
	}

	userId, ok := ctx.Locals("userId").(uuid.UUID)
	if !ok {
		return &response.Unauthorized
	}

	err = r.service.CommentService.DeleteComment(commentId, userId)
	if err != nil {
		return err
	}

	response.Success(ctx, http.StatusOK, "success", nil)
	return nil
}
