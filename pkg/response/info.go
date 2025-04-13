package response

import (
	"errors"

	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
)

func GetErrorInfo(err error) (int, string) {
	if err == nil {
		return fiber.StatusOK, ""
	}

	code := fiber.StatusInternalServerError
	message := "Internal Server Error"

	var errorRequest *ErrorResponse
	if errors.As(err, &errorRequest) {
		return errorRequest.Code, errorRequest.Error()
	}

	var fiberError *fiber.Error
	if errors.As(err, &fiberError) {
		return fiberError.Code, fiberError.Message
	}

	var validationError validator.ValidationErrors
	if errors.As(err, &validationError) {
		return fiber.StatusBadRequest, "Bad Request"
	}

	return code, message
}
