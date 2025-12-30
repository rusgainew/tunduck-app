package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

// LogrusMiddleware returns a middleware that logs HTTP requests using logrus
func LogrusMiddleware(log *logrus.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		// Continue processing request
		err := c.Next()

		// Log request details
		latency := time.Since(start)
		statusCode := c.Response().StatusCode()
		method := c.Method()
		path := c.OriginalURL()
		ip := c.IP()

		// Determine log level based on status code
		if statusCode >= 400 && statusCode < 500 {
			log.WithFields(logrus.Fields{
				"method":     method,
				"path":       path,
				"status":     statusCode,
				"ip":         ip,
				"latency_ms": latency.Milliseconds(),
			}).Warn("HTTP Request")
		} else if statusCode >= 500 {
			log.WithFields(logrus.Fields{
				"method":     method,
				"path":       path,
				"status":     statusCode,
				"ip":         ip,
				"latency_ms": latency.Milliseconds(),
				"error":      err,
			}).Error("HTTP Request")
		} else {
			log.WithFields(logrus.Fields{
				"method":     method,
				"path":       path,
				"status":     statusCode,
				"ip":         ip,
				"latency_ms": latency.Milliseconds(),
			}).Info("HTTP Request")
		}

		return err
	}
}
