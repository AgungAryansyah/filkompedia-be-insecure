package rest

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/yogarn/filkompedia-be/model"
	"github.com/yogarn/filkompedia-be/pkg/response"
)

func (r *Rest) GetMe(ctx *fiber.Ctx) (err error) {
	userId, ok := ctx.Locals("userId").(uuid.UUID)
	if !ok {
		return &response.Unauthorized
	}

	var profile model.Profile
	if err := r.service.UserService.GetProfile(&profile, userId); err != nil {
		return err
	}

	response.Success(ctx, http.StatusOK, "success", profile)
	return nil
}

func (r *Rest) GetUserProfile(ctx *fiber.Ctx) (err error) {
	param := ctx.Params("userId")
	userId, err := uuid.Parse(param)
	if err != nil {
		return err
	}

	var profile model.Profile
	if err := r.service.UserService.GetProfile(&profile, userId); err != nil {
		return err
	}

	response.Success(ctx, http.StatusOK, "success", profile)
	return nil
}

func (r *Rest) GetAllUserProfile(ctx *fiber.Ctx) (err error) {
	var userProfile model.ProfilesReq
	userProfile.Page = ctx.QueryInt("page", 1)
	userProfile.PageSize = ctx.QueryInt("size", 9)

	var profiles []model.Profile
	if err := r.service.UserService.GetProfiles(&profiles, userProfile); err != nil {
		return err
	}

	response.Success(ctx, http.StatusOK, "success", profiles)
	return nil
}

func (r *Rest) UpdateRole(ctx *fiber.Ctx) (err error) {
	roleReq := &model.RoleUpdate{}
	if err := ctx.BodyParser(roleReq); err != nil {
		return err
	}

	validate := validator.New()
	if err := validate.Struct(roleReq); err != nil {
		return err
	}

	if err := r.service.UserService.UpdateRole(roleReq); err != nil {
		return err
	}

	if err := r.service.AuthService.ClearToken(roleReq.Id); err != nil {
		return err
	}

	response.Success(ctx, http.StatusOK, "success", nil)
	return nil
}

func (r *Rest) EditProfile(ctx *fiber.Ctx) error {
	var edit model.EditProfile
	if err := ctx.BodyParser(&edit); err != nil {
		return err
	}

	if edit.Id != uuid.Nil {
		return &response.BadRequest
	}

	userId, ok := ctx.Locals("userId").(uuid.UUID)
	if !ok {
		return &response.Unauthorized
	}
	edit.Id = userId

	if err := r.validator.Struct(edit); err != nil {
		return err
	}

	if err := r.service.UserService.EditProfile(&edit); err != nil {
		return err
	}

	response.Success(ctx, http.StatusOK, "success", nil)
	return nil
}

func (r *Rest) DeleteUser(ctx *fiber.Ctx) error {
	param := ctx.Params("userId")
	userId, err := uuid.Parse(param)
	if err != nil {
		return err
	}

	if err := r.service.UserService.DeleteUser(userId); err != nil {
		return err
	}

	response.Success(ctx, http.StatusOK, "success", nil)
	return nil
}

func (r *Rest) UploadProfilePicture(ctx *fiber.Ctx) error {
	file, err := ctx.FormFile("file")
	if err != nil {
		return err
	}

	profilePicture := model.Image{File: file}
	if err := r.validator.Struct(profilePicture); err != nil {
		return &response.BadRequest
	}

	url, err := r.service.UserService.UploadProfilePicture(file)
	if err != nil {
		return err
	}

	response.Success(ctx, http.StatusOK, "success", url)
	return nil
}
