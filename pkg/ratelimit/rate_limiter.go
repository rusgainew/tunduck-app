package ratelimit

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// RateLimiter handles request rate limiting using Redis as a backend
type RateLimiter struct {
	redisClient *redis.Client
}

// LimitConfig contains configuration for different rate limit scenarios
type LimitConfig struct {
	RequestsPerMinute int
	Window            time.Duration
}

// DefaultLimits defines rate limits for different endpoint categories
var DefaultLimits = map[string]LimitConfig{
	"public":    {RequestsPerMinute: 30, Window: time.Minute},  // register, login
	"protected": {RequestsPerMinute: 60, Window: time.Minute},  // authenticated endpoints
	"health":    {RequestsPerMinute: 120, Window: time.Minute}, // /health check
	"metrics":   {RequestsPerMinute: 180, Window: time.Minute}, // /metrics scraping
	"sensitive": {RequestsPerMinute: 5, Window: time.Minute},   // logout, password change
}

// NewRateLimiter creates a new rate limiter instance
func NewRateLimiter(redisClient *redis.Client) *RateLimiter {
	return &RateLimiter{
		redisClient: redisClient,
	}
}

// IsAllowed checks if a request should be allowed based on rate limit
// identifier: unique key (IP, user ID, API key, etc.)
// category: "public", "protected", "health", "metrics", or "sensitive"
// Returns (allowed, remaining, resetTime, error)
func (rl *RateLimiter) IsAllowed(ctx context.Context, identifier string, category string) (bool, int, time.Time, error) {
	config, exists := DefaultLimits[category]
	if !exists {
		config = DefaultLimits["protected"] // fallback to protected
	}

	key := fmt.Sprintf("ratelimit:%s:%s", category, identifier)
	now := time.Now()
	resetTime := now.Add(config.Window)

	// Use Redis INCR with expiration for rate limiting
	pipe := rl.redisClient.Pipeline()

	incrCmd := pipe.Incr(ctx, key)
	pipe.Expire(ctx, key, config.Window)
	_, err := pipe.Exec(ctx)
	if err != nil && err != redis.Nil {
		// On Redis error, allow request (graceful degradation)
		return true, config.RequestsPerMinute, resetTime, nil
	}

	count := incrCmd.Val()
	allowed := count <= int64(config.RequestsPerMinute)
	remaining := int(config.RequestsPerMinute) - int(count)
	if remaining < 0 {
		remaining = 0
	}

	return allowed, remaining, resetTime, nil
}

// Reset clears the rate limit counter for an identifier
func (rl *RateLimiter) Reset(ctx context.Context, identifier string, category string) error {
	key := fmt.Sprintf("ratelimit:%s:%s", category, identifier)
	return rl.redisClient.Del(ctx, key).Err()
}

// GetCount returns current request count for an identifier
func (rl *RateLimiter) GetCount(ctx context.Context, identifier string, category string) (int, error) {
	key := fmt.Sprintf("ratelimit:%s:%s", category, identifier)
	val, err := rl.redisClient.Get(ctx, key).Int()
	if err == redis.Nil {
		return 0, nil
	}
	return val, err
}
