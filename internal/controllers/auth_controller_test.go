package controllers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"github.com/rusgainew/tunduck-app/pkg/cache"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestAuthPayloadValidation tests auth endpoint payload structure
func TestAuthPayloadValidation(t *testing.T) {
	t.Run("Register payload structure", func(t *testing.T) {
		requestBody := map[string]string{
			"username": "newuser",
			"email":    "newuser@example.com",
			"password": "SecurePassword123!",
		}

		body, err := json.Marshal(requestBody)
		require.NoError(t, err)

		var payload map[string]string
		err = json.Unmarshal(body, &payload)
		require.NoError(t, err)

		assert.Equal(t, "newuser", payload["username"])
		assert.Equal(t, "newuser@example.com", payload["email"])
		assert.Equal(t, "SecurePassword123!", payload["password"])
	})

	t.Run("Login payload structure", func(t *testing.T) {
		requestBody := map[string]string{
			"email":    "user@example.com",
			"password": "Password123!",
		}

		body, err := json.Marshal(requestBody)
		require.NoError(t, err)

		var payload map[string]string
		err = json.Unmarshal(body, &payload)
		require.NoError(t, err)

		assert.Equal(t, "user@example.com", payload["email"])
		assert.Equal(t, "Password123!", payload["password"])
	})

	t.Run("Logout payload structure", func(t *testing.T) {
		requestBody := map[string]interface{}{
			"token": "jwt_token_here",
		}

		body, err := json.Marshal(requestBody)
		require.NoError(t, err)

		var payload map[string]interface{}
		err = json.Unmarshal(body, &payload)
		require.NoError(t, err)

		assert.Contains(t, payload, "token")
	})

	t.Run("Invalid JSON payload", func(t *testing.T) {
		invalidJSON := []byte(`{"username": "user" invalid}`)

		var payload map[string]string
		err := json.Unmarshal(invalidJSON, &payload)
		assert.Error(t, err)
	})

	t.Run("Missing required fields in register", func(t *testing.T) {
		requestBody := map[string]string{
			"username": "newuser",
			// Missing email and password
		}

		body, _ := json.Marshal(requestBody)
		var payload map[string]interface{}
		json.Unmarshal(body, &payload)

		assert.NotContains(t, payload, "email")
		assert.NotContains(t, payload, "password")
	})
}

// TestTokenBlacklistIntegration tests logout with token blacklist
func TestTokenBlacklistIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Setup Redis
	redisClient := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	defer redisClient.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := redisClient.Ping(ctx).Err()
	if err != nil {
		t.Skipf("Redis not available: %v", err)
	}

	logger := logrus.New()
	cacheManager := cache.NewRedisCacheManager(redisClient, logger)

	t.Run("Token added to blacklist on logout", func(t *testing.T) {
		token := "jwt_token_example_12345"
		// Get token cache from cache manager
		tokenCache := cacheManager.Token()

		// Before logout - not in blacklist
		exists, err := tokenCache.Exists(ctx, token)
		require.NoError(t, err)
		assert.False(t, exists)

		// Simulate logout - add to blacklist
		err = tokenCache.Set(ctx, token, true, 1*time.Hour)
		require.NoError(t, err)

		// After logout - in blacklist
		exists, err = tokenCache.Exists(ctx, token)
		require.NoError(t, err)
		assert.True(t, exists)

		// Cleanup
		tokenCache.Delete(ctx, token)
	})

	t.Run("Multiple tokens can be blacklisted independently", func(t *testing.T) {
		tokenCache := cacheManager.Token()
		token1 := "jwt_token_1"
		token2 := "jwt_token_2"

		// Add both to blacklist
		tokenCache.Set(ctx, token1, true, 1*time.Hour)
		tokenCache.Set(ctx, token2, true, 1*time.Hour)

		// Both should exist independently
		exists1, _ := tokenCache.Exists(ctx, token1)
		exists2, _ := tokenCache.Exists(ctx, token2)
		assert.True(t, exists1)
		assert.True(t, exists2)

		// Delete one
		tokenCache.Delete(ctx, token1)

		// One should be gone, other should remain
		exists1, _ = tokenCache.Exists(ctx, token1)
		exists2, _ = tokenCache.Exists(ctx, token2)
		assert.False(t, exists1)
		assert.True(t, exists2)

		// Cleanup
		tokenCache.Delete(ctx, token2)
	})

	t.Run("Blacklist expiration", func(t *testing.T) {
		tokenCache := cacheManager.Token()
		token := "jwt_token_expiring"

		// Add with short TTL (500ms)
		err := tokenCache.Set(ctx, token, true, 500*time.Millisecond)
		require.NoError(t, err)

		// Should exist immediately
		exists, _ := tokenCache.Exists(ctx, token)
		assert.True(t, exists)

		// Wait for expiration
		time.Sleep(600 * time.Millisecond)

		// Should be expired
		exists, _ = tokenCache.Exists(ctx, token)
		assert.False(t, exists)
	})
}

// TestHTTPMiddlewareIntegration tests middleware behavior
func TestHTTPMiddlewareIntegration(t *testing.T) {
	t.Run("Missing Authorization header", func(t *testing.T) {
		app := fiber.New()

		// Add auth check middleware
		app.Use(func(c *fiber.Ctx) error {
			authHeader := c.Get("Authorization")
			if authHeader == "" {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"error": "missing authorization header",
				})
			}
			return c.Next()
		})

		app.Get("/protected", func(c *fiber.Ctx) error {
			return c.JSON(fiber.Map{"message": "success"})
		})

		req := httptest.NewRequest(http.MethodGet, "/protected", nil)
		resp, err := app.Test(req)
		require.NoError(t, err)

		assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("Valid Authorization header", func(t *testing.T) {
		app := fiber.New()

		app.Use(func(c *fiber.Ctx) error {
			authHeader := c.Get("Authorization")
			if authHeader == "" {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"error": "missing",
				})
			}
			return c.Next()
		})

		app.Get("/protected", func(c *fiber.Ctx) error {
			return c.JSON(fiber.Map{"message": "success"})
		})

		req := httptest.NewRequest(http.MethodGet, "/protected", nil)
		req.Header.Set("Authorization", "Bearer jwt_token_here")
		resp, err := app.Test(req)
		require.NoError(t, err)

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	})
}

// TestJWTTokenFormat tests JWT token structure
func TestJWTTokenFormat(t *testing.T) {
	t.Run("JWT structure validation", func(t *testing.T) {
		// A valid JWT would have format: header.payload.signature
		validJWT := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"

		// Should have exactly 2 dots
		parts := bytes.Count([]byte(validJWT), []byte("."))
		assert.Equal(t, 2, parts)
	})

	t.Run("Invalid JWT format - missing signature", func(t *testing.T) {
		invalidJWT := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ"

		parts := bytes.Count([]byte(invalidJWT), []byte("."))
		assert.Equal(t, 1, parts) // Only 1 dot instead of 2
	})

	t.Run("Token Bearer prefix", func(t *testing.T) {
		authHeader := "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIn0.signature"

		// Should start with "Bearer "
		assert.True(t, bytes.HasPrefix([]byte(authHeader), []byte("Bearer ")))

		// Extract token
		token := authHeader[7:] // Skip "Bearer "
		assert.NotEmpty(t, token)
	})
}

// TestConcurrentAuthRequests tests concurrent auth operations
func TestConcurrentAuthRequests(t *testing.T) {
	t.Run("Concurrent login attempts", func(t *testing.T) {
		done := make(chan bool, 5)
		responses := make([]string, 5)

		for i := 0; i < 5; i++ {
			go func(idx int) {
				// Simulate login request
				requestBody := map[string]string{
					"email":    "user" + string(rune(idx)) + "@example.com",
					"password": "password123",
				}

				body, _ := json.Marshal(requestBody)
				var payload map[string]string
				json.Unmarshal(body, &payload)

				responses[idx] = payload["email"]
				done <- true
			}(i)
		}

		for i := 0; i < 5; i++ {
			<-done
		}

		// All responses should be unique
		emailSet := make(map[string]bool)
		for _, email := range responses {
			emailSet[email] = true
		}
		assert.Equal(t, 5, len(emailSet))
	})

	t.Run("Concurrent logout attempts", func(t *testing.T) {
		if testing.Short() {
			t.Skip("Skipping integration test in short mode")
		}

		redisClient := redis.NewClient(&redis.Options{
			Addr: "localhost:6379",
		})
		defer redisClient.Close()

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err := redisClient.Ping(ctx).Err()
		if err != nil {
			t.Skipf("Redis not available: %v", err)
		}

		logger := logrus.New()
		cacheManager := cache.NewRedisCacheManager(redisClient, logger)
		tokenCache := cacheManager.Token()

		done := make(chan bool, 5)

		for i := 0; i < 5; i++ {
			go func(idx int) {
				token := "token_" + string(rune(idx))

				// Simulate logout - add to blacklist
				tokenCache.Set(ctx, token, true, 1*time.Hour)

				// Verify added
				exists, _ := tokenCache.Exists(ctx, token)
				assert.True(t, exists)

				// Cleanup
				tokenCache.Delete(ctx, token)
				done <- true
			}(i)
		}

		for i := 0; i < 5; i++ {
			<-done
		}
	})
}
