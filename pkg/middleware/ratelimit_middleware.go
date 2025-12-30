package middleware

import (
	"net"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/rusgainew/tunduck-app/pkg/ratelimit"
	"github.com/sirupsen/logrus"
)

// RateLimitMiddleware creates a middleware that enforces rate limits per IP
// category determines which rate limit rules apply to this endpoint
func RateLimitMiddleware(rl *ratelimit.RateLimiter, category string, logger *logrus.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get client IP (consider X-Forwarded-For for proxies)
		clientIP := getClientIP(c)

		// Check rate limit
		allowed, _, resetTime, err := rl.IsAllowed(c.Context(), clientIP, category)
		if err != nil && err.Error() != "redis: nil" {
			logger.WithError(err).Warn("Failed to check rate limit, allowing request")
		}

		if !allowed {
			logger.WithFields(logrus.Fields{
				"ip":       clientIP,
				"category": category,
				"path":     c.Path(),
			}).Warn("Rate limit exceeded")

			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error":   "Too Many Requests",
				"message": "Rate limit exceeded. Please try again later.",
				"reset":   resetTime.Unix(),
			})
		}

		return c.Next()
	}
}

// RateLimitAuthMiddleware for authenticated endpoints
// Uses user ID or token as identifier instead of IP for better granularity
func RateLimitAuthMiddleware(rl *ratelimit.RateLimiter, category string, logger *logrus.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Try to get user ID from context (set by JWT middleware)
		userID := c.Locals("user_id")
		identifier := ""

		if userID != nil {
			identifier = userID.(string)
		} else {
			// Fallback to IP if user not authenticated
			identifier = getClientIP(c)
		}

		// Check rate limit
		allowed, _, resetTime, err := rl.IsAllowed(c.Context(), identifier, category)
		if err != nil && err.Error() != "redis: nil" {
			logger.WithError(err).Warn("Failed to check rate limit, allowing request")
		}

		if !allowed {
			logger.WithFields(logrus.Fields{
				"identifier": identifier,
				"category":   category,
				"path":       c.Path(),
			}).Warn("Rate limit exceeded")

			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error":   "Too Many Requests",
				"message": "Rate limit exceeded. Please try again later.",
				"reset":   resetTime.Unix(),
			})
		}

		return c.Next()
	}
}

// getClientIP extracts the real client IP from request
// Considers X-Forwarded-For header for proxied requests
func getClientIP(c *fiber.Ctx) string {
	// Check X-Forwarded-For header (proxies, load balancers)
	if xForwardedFor := c.Get("X-Forwarded-For"); xForwardedFor != "" {
		// X-Forwarded-For can contain multiple IPs, get the first one
		ips := strings.Split(xForwardedFor, ",")
		return strings.TrimSpace(ips[0])
	}

	// Check X-Real-IP header
	if xRealIP := c.Get("X-Real-IP"); xRealIP != "" {
		return xRealIP
	}

	// Get from connection
	ip, _, err := net.SplitHostPort(c.IP())
	if err != nil {
		return c.IP()
	}
	return ip
}
