package config

import (
	"github.com/AgungAryansyah/filkompedia-be-insecure/pkg/response"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func StartFiber() *fiber.App {
	app := fiber.New(
		fiber.Config{
			ErrorHandler: CustomErrorHandler,
		},
	)

	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:5173, https://filkompedia.yogarn.my.id, https://api.sandbox.midtrans.com, http://10.34.100.139:5173",
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
