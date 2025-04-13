package response

import "github.com/gofiber/fiber/v2"

type Response struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func Success(ctx *fiber.Ctx, code int, message string, data interface{}) {
	response := Response{
		Message: message,
		Data:    data,
	}

	ctx.Status(code).JSON(response)
}

func Error(ctx *fiber.Ctx, code int, message string, err error) {
	response := Response{
		Message: message,
		Data:    err.Error(),
	}

	ctx.Status(code).JSON(response)
}
