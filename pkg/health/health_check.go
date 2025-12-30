package health

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// Status представляет статус компонента
type Status string

const (
	StatusUp   Status = "UP"
	StatusDown Status = "DOWN"
)

// ComponentHealth содержит статус компонента системы
type ComponentHealth struct {
	Name         string    `json:"name"`
	Status       Status    `json:"status"`
	ResponseTime string    `json:"response_time"`
	Message      string    `json:"message,omitempty"`
	LastChecked  time.Time `json:"last_checked"`
}

// HealthCheck содержит информацию о здоровье всей системы
type HealthCheck struct {
	Status     Status            `json:"status"`
	Timestamp  time.Time         `json:"timestamp"`
	Components []ComponentHealth `json:"components"`
	Uptime     string            `json:"uptime,omitempty"`
}

// HealthChecker проверяет здоровье системы
type HealthChecker struct {
	db          *gorm.DB
	redisClient *redis.Client
	logger      *logrus.Logger
	startTime   time.Time
}

// NewHealthChecker создает новый health checker
func NewHealthChecker(db *gorm.DB, redisClient *redis.Client, logger *logrus.Logger) *HealthChecker {
	return &HealthChecker{
		db:          db,
		redisClient: redisClient,
		logger:      logger,
		startTime:   time.Now(),
	}
}

// Check проверяет здоровье всей системы
func (hc *HealthChecker) Check(ctx context.Context) *HealthCheck {
	components := []ComponentHealth{}

	// Проверяем PostgreSQL
	components = append(components, hc.checkDatabase(ctx))

	// Проверяем Redis
	components = append(components, hc.checkRedis(ctx))

	// Определяем общий статус
	overallStatus := StatusUp
	for _, comp := range components {
		if comp.Status == StatusDown {
			overallStatus = StatusDown
			break
		}
	}

	return &HealthCheck{
		Status:     overallStatus,
		Timestamp:  time.Now(),
		Components: components,
		Uptime:     time.Since(hc.startTime).String(),
	}
}

// checkDatabase проверяет статус PostgreSQL
func (hc *HealthChecker) checkDatabase(ctx context.Context) ComponentHealth {
	start := time.Now()
	component := ComponentHealth{
		Name:        "PostgreSQL",
		LastChecked: time.Now(),
	}

	// Добавляем timeout для проверки
	checkCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Пытаемся выполнить простой запрос
	var version string
	if err := hc.db.WithContext(checkCtx).Raw("SELECT version()").Scan(&version).Error; err != nil {
		component.Status = StatusDown
		component.Message = err.Error()
		hc.logger.WithError(err).Warn("Database health check failed")
	} else {
		component.Status = StatusUp
		component.Message = "Database connected successfully"
	}

	component.ResponseTime = time.Since(start).String()
	return component
}

// checkRedis проверяет статус Redis
func (hc *HealthChecker) checkRedis(ctx context.Context) ComponentHealth {
	start := time.Now()
	component := ComponentHealth{
		Name:        "Redis",
		LastChecked: time.Now(),
	}

	// Если redisClient nil, значит Redis не настроен
	if hc.redisClient == nil {
		component.Status = StatusDown
		component.Message = "Redis client not configured"
		component.ResponseTime = time.Since(start).String()
		return component
	}

	// Добавляем timeout для проверки
	checkCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Пытаемся выполнить PING
	if err := hc.redisClient.Ping(checkCtx).Err(); err != nil {
		component.Status = StatusDown
		component.Message = err.Error()
		hc.logger.WithError(err).Warn("Redis health check failed")
	} else {
		component.Status = StatusUp
		component.Message = "Redis connected successfully"
	}

	component.ResponseTime = time.Since(start).String()
	return component
}
