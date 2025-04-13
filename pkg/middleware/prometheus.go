package middleware

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/yogarn/filkompedia-be/pkg/response"
)

func (m *middleware) PromMiddleware(c *fiber.Ctx) error {
	reg := m.reg
	now := time.Now()

	err := c.Next()

	method := string(append([]byte{}, c.Method()...))
	path := string(append([]byte{}, c.Path()...))

	respStatus := c.Response().StatusCode()
	duration := time.Since(now).Seconds()

	if err != nil {
		respStatus, _ = response.GetErrorInfo(err)
	}

	reg.RequestTotal.With(prometheus.Labels{
		"response_code": strconv.Itoa(respStatus),
		"method":        method,
	}).Inc()

	reg.Duration.WithLabelValues(strconv.Itoa(respStatus), method, path).Observe(duration)
	reg.DurationSummary.Observe(duration)

	if respStatus >= 400 {
		reg.ErrorCount.WithLabelValues(method, strconv.Itoa(respStatus)).Inc()
	}

	return err
}
