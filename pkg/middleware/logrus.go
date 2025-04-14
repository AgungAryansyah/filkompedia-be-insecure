package middleware

import (
	"encoding/json"
	"strings"

	"github.com/AgungAryansyah/filkompedia-be-insecure/pkg/response"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

func GetRealIP(ctx *fiber.Ctx) string {
	if forwardedFor := ctx.Get("X-Forwarded-For"); forwardedFor != "" {
		return strings.Split(forwardedFor, ",")[0]
	}
	return ctx.IP()
}

func (m *middleware) LogrusMiddleware(ctx *fiber.Ctx) error {
	start := ctx.Context().Time()

	contentType := ctx.Get("Content-Type")

	var sanitizedBody string
	if strings.HasPrefix(contentType, "multipart/form-data") {
		sanitizedBody = "[skipped multipart body]"
	} else {
		rawBody := string(ctx.Body())
		sanitizedBody = sanitizeBody(rawBody)
	}

	err := ctx.Next()
	statusCode, _ := response.GetErrorInfo(err)

	clientIp := GetRealIP(ctx)

	entry := m.logger.WithFields(logrus.Fields{
		"method":    ctx.Method(),
		"path":      ctx.Path(),
		"status":    statusCode,
		"latency":   ctx.Context().Time().Sub(start),
		"ip":        clientIp,
		"userAgent": ctx.Get("User-Agent"),
		"body":      sanitizedBody,
	})

	if err != nil {
		entry.Error(err.Error())
	} else {
		entry.Info("incoming request")
	}

	return err
}

func sanitizeBody(body string) string {
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(body), &data); err == nil {
		if _, exists := data["password"]; exists {
			data["password"] = "***REDACTED***"
		}
		sanitized, _ := json.Marshal(data)
		return string(sanitized)
	}

	if strings.Contains(body, `"password"`) {
		return strings.ReplaceAll(body, `"password":"`, `"password":"***REDACTED***"`)
	}

	return body
}
