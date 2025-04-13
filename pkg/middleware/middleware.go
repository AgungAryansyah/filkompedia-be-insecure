package middleware

import (
	"github.com/AgungAryansyah/filkompedia-be-unsecure/internal/service"
	"github.com/AgungAryansyah/filkompedia-be-unsecure/pkg/jwt"
	monitoring "github.com/AgungAryansyah/filkompedia-be-unsecure/pkg/prometheus"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type IMiddleware interface {
	Authenticate(ctx *fiber.Ctx) error
	Authorize(roles []int) fiber.Handler
	AuthorizeOrItself(roles []int) fiber.Handler
	PromMiddleware(ctx *fiber.Ctx) error
	LogrusMiddleware(ctx *fiber.Ctx) error
	BookCommentCheck(ctx *fiber.Ctx) error
}

type middleware struct {
	jwtAuth jwt.IJwt
	service *service.Service
	reg     monitoring.Metrics
	logger  *logrus.Logger
}

func Init(jwtAuth jwt.IJwt, service *service.Service, reg monitoring.Metrics, logger *logrus.Logger) IMiddleware {
	return &middleware{
		jwtAuth: jwtAuth,
		service: service,
		reg:     reg,
		logger:  logger,
	}
}
