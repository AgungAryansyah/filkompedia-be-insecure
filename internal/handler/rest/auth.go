package rest

import (
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/yogarn/filkompedia-be/model"
	"github.com/yogarn/filkompedia-be/pkg/response"
)

func (r *Rest) Register(ctx *fiber.Ctx) (err error) {
	registerReq := &model.RegisterReq{}
	if err := ctx.BodyParser(registerReq); err != nil {
		return err
	}

	validate := validator.New()
	if err := validate.Struct(registerReq); err != nil {
		return err
	}

	user, err := r.service.AuthService.Register(registerReq)
	if err != nil {
		return err
	}

	response.Success(ctx, http.StatusCreated, "success", user)
	return nil
}

func (r *Rest) SendOtp(ctx *fiber.Ctx) (err error) {
	otpReq := &model.OtpReq{}
	if err := ctx.BodyParser(otpReq); err != nil {
		return err
	}

	validate := validator.New()
	if err := validate.Struct(otpReq); err != nil {
		return err
	}

	err = r.service.AuthService.SendOTP(otpReq.Email)
	if err != nil {
		return err
	}

	response.Success(ctx, http.StatusCreated, "success", nil)
	return nil
}

func (r *Rest) VerifyOtp(ctx *fiber.Ctx) (err error) {
	OtpVerifyReq := &model.OtpVerifyReq{}
	if err := ctx.BodyParser(OtpVerifyReq); err != nil {
		return err
	}

	validate := validator.New()
	if err := validate.Struct(OtpVerifyReq); err != nil {
		return err
	}

	err = r.service.AuthService.VerifyOTP(OtpVerifyReq.Email, OtpVerifyReq.Otp)
	if err != nil {
		return err
	}

	response.Success(ctx, http.StatusOK, "success", nil)
	return nil
}

func (r *Rest) Login(ctx *fiber.Ctx) (err error) {
	loginReq := &model.LoginReq{}
	if err := ctx.BodyParser(loginReq); err != nil {
		return err
	}

	validate := validator.New()
	if err := validate.Struct(loginReq); err != nil {
		return err
	}

	ipAddress := ctx.IP()
	userAgent := ctx.Get("User-Agent")

	refreshTokenExpiresIn, err := strconv.Atoi(os.Getenv("REFRESH_EXPIRED_TIME"))
	if err != nil {
		return err
	}

	loginRes, err := r.service.AuthService.Login(loginReq, ipAddress, userAgent, refreshTokenExpiresIn)
	if err != nil {
		return err
	}

	expiresIn, err := strconv.Atoi(os.Getenv("JWT_EXPIRED_TIME"))
	if err != nil {
		return err
	}

	ctx.Cookie(&fiber.Cookie{
		Name:     "token",
		Value:    loginRes.JwtToken,
		Expires:  time.Now().Add(time.Duration(expiresIn) * time.Second),
		HTTPOnly: true,
		Secure:   false, // should set true in prod
		Path:     "/",
		SameSite: "None",
	})

	ctx.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    loginRes.RefreshToken,
		Expires:  time.Now().Add(time.Duration(refreshTokenExpiresIn) * time.Second),
		HTTPOnly: true,
		Secure:   false, // this one too
		Path:     "/",
		SameSite: "None",
	})

	response.Success(ctx, http.StatusOK, "success", nil)
	return nil
}

func (r *Rest) GetSessions(ctx *fiber.Ctx) (err error) {
	userId, ok := ctx.Locals("userId").(uuid.UUID)
	if !ok {
		return &response.Unauthorized
	}

	sessions, err := r.service.AuthService.GetSessions(userId)
	if err != nil {
		return err
	}

	response.Success(ctx, http.StatusOK, "success", sessions)
	return nil
}

func (r *Rest) ExchangeToken(ctx *fiber.Ctx) (err error) {
	token := ctx.Cookies("refresh_token")
	if token == "" {
		return &response.InvalidToken
	}

	refreshTokenExpiresIn, err := strconv.Atoi(os.Getenv("REFRESH_EXPIRED_TIME"))
	if err != nil {
		return err
	}

	jwtToken, newToken, err := r.service.AuthService.ExchangeToken(token, refreshTokenExpiresIn)
	if err != nil {
		return err
	}

	expiresIn, err := strconv.Atoi(os.Getenv("JWT_EXPIRED_TIME"))
	if err != nil {
		return err
	}

	ctx.Cookie(&fiber.Cookie{
		Name:     "token",
		Value:    jwtToken,
		Expires:  time.Now().Add(time.Duration(expiresIn) * time.Second),
		HTTPOnly: true,
		Secure:   true,
		Path:     "/",
	})

	ctx.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    newToken,
		Expires:  time.Now().Add(time.Duration(refreshTokenExpiresIn) * time.Second),
		HTTPOnly: true,
		Secure:   true,
		Path:     "/",
	})

	response.Success(ctx, http.StatusOK, "success", nil)
	return nil
}

func (r *Rest) Logout(ctx *fiber.Ctx) (err error) {
	userId, ok := ctx.Locals("userId").(uuid.UUID)
	if !ok {
		return &response.Unauthorized
	}

	tokenString := ctx.Cookies("refresh_token")
	if tokenString == "" {
		return &response.InvalidToken
	}

	deleteTokenReq := &model.DeleteToken{}
	deleteTokenReq.UserId = userId
	deleteTokenReq.Token = tokenString

	err = r.service.AuthService.DeleteToken(deleteTokenReq)
	if err != nil {
		return err
	}

	ctx.Cookie(&fiber.Cookie{
		Name:     "token",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
		Secure:   false,
		Path:     "/",
	})

	ctx.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
		Secure:   false,
		Path:     "/",
	})

	response.Success(ctx, http.StatusOK, "Logged out successfully", nil)
	return nil
}

func (r *Rest) ChangePassword(ctx *fiber.Ctx) error {
	var req model.ChangePassword
	if err := ctx.BodyParser(&req); err != nil {
		return err
	}

	if err := r.service.AuthService.ChangePassword(&req); err != nil {
		return err
	}

	response.Success(ctx, http.StatusOK, "success", nil)
	return nil
}
