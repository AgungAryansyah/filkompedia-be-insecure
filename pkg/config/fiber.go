package config

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/yogarn/filkompedia-be/pkg/response"
)

func StartFiber() *fiber.App {
	app := fiber.New(
		fiber.Config{
			ErrorHandler: CustomErrorHandler,
		},
	)

	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:5173, https://filkompedia.yogarn.my.id, https://api.sandbox.midtrans.com",
		AllowCredentials: true,
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS,PATCH",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		ExposeHeaders:    "Set-Cookie",
	}))

	return app
}

func CustomErrorHandler(ctx *fiber.Ctx, err error) error {
	code, message := response.GetErrorInfo(err)

	ctx.Status(code)
	response.Error(ctx, code, message, err)

	return nil
}
