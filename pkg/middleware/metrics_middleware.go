package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/rusgainew/tunduck-app/pkg/metrics"
)

// MetricsMiddleware записывает HTTP метрики
func MetricsMiddleware(m *metrics.Metrics) fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		// Записываем размер запроса
		m.HTTPRequestsTotal.Inc()
		m.HTTPRequestSize.Observe(float64(len(c.Body())))

		// Выполняем запрос
		err := c.Next()

		// Записываем метрики
		duration := time.Since(start).Seconds()
		m.HTTPRequestDuration.Observe(duration)
		m.HTTPResponseSize.Observe(float64(len(c.Response().Body())))

		return err
	}
}

// CacheMetricsWrapper обёртка для отслеживания cache операций
func CacheMetricsWrapper(m *metrics.Metrics, operation string, fn func() (interface{}, error)) (interface{}, error) {
	start := time.Now()
	result, err := fn()
	duration := time.Since(start).Seconds()

	m.CacheOperationDuration.Observe(duration)

	if err != nil {
		// Cache miss или ошибка
		if operation == "get" {
			m.CacheMissesTotal.Inc()
		}
	} else if operation == "get" && result != nil {
		m.CacheHitsTotal.Inc()
	}

	return result, err
}

// DatabaseMetricsWrapper обёртка для отслеживания DB операций
func DatabaseMetricsWrapper(m *metrics.Metrics, fn func() error) error {
	start := time.Now()
	err := fn()
	duration := time.Since(start).Seconds()

	m.DBQueryDuration.Observe(duration)

	if err != nil {
		m.DBErrorsTotal.Inc()
	}

	return err
}
