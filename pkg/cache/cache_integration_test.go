package cache

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/rusgainew/tunduck-app/pkg/entity"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// TestCacheIntegration tests cache hits and misses with real Redis
func TestCacheIntegration_UserService(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Setup Redis client
	redisClient := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	defer redisClient.Close()

	// Verify Redis connectivity
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := redisClient.Ping(ctx).Err()
	if err != nil {
		t.Skipf("Redis not available: %v", err)
	}

	// Clean up test keys before test
	cleanupTestKeys(t, redisClient)

	logger := logrus.New()
	cacheManager := NewRedisCacheManager(redisClient, logger)

	t.Run("Cache Hit on GetByUsername", func(t *testing.T) {
		testKey := "test_user_username_johndoe"
		testUser := &entity.User{
			ID:       uuid.New(),
			Username: "johndoe",
			Email:    "john@example.com",
			Password: "hashed_pwd",
		}

		// Set value in cache using Generic cache
		err := cacheManager.Generic().Set(ctx, testKey, testUser, 1*time.Hour)
		require.NoError(t, err)

		// Retrieve from cache - JSON десериализуется в map[string]interface{}
		val, err := cacheManager.Generic().Get(ctx, testKey)
		require.NoError(t, err)
		require.NotNil(t, val)

		// Проверяем что получили map (так работает JSON unmarshal в interface{})
		cachedData, ok := val.(map[string]interface{})
		assert.True(t, ok, "Value should be deserialized as map[string]interface{}")
		if ok {
			assert.Equal(t, testUser.Username, cachedData["username"])
			assert.Equal(t, testUser.Email, cachedData["email"])
		}
	})

	t.Run("Cache Miss returns nil", func(t *testing.T) {
		testKey := "non_existent_key_" + time.Now().Format("20060102150405")

		val, err := cacheManager.Generic().Get(ctx, testKey)
		require.NoError(t, err)
		assert.Nil(t, val)
	})

	t.Run("Cache Expiration after TTL", func(t *testing.T) {
		testKey := "test_expiration_key"
		testValue := "test_value"

		// Set with 1 second TTL
		err := cacheManager.Generic().Set(ctx, testKey, testValue, 1*time.Second)
		require.NoError(t, err)

		// Verify it exists
		val, err := cacheManager.Generic().Get(ctx, testKey)
		require.NoError(t, err)
		assert.NotNil(t, val)

		// Wait for expiration
		time.Sleep(1500 * time.Millisecond)

		// Should be expired
		val, err = cacheManager.Generic().Get(ctx, testKey)
		require.NoError(t, err)
		assert.Nil(t, val)
	})

	t.Run("Cache Invalidation with Delete", func(t *testing.T) {
		testKey := "test_invalidate_key"
		testValue := "test_value"

		// Set value
		err := cacheManager.Generic().Set(ctx, testKey, testValue, 1*time.Hour)
		require.NoError(t, err)

		// Verify it exists
		val, err := cacheManager.Generic().Get(ctx, testKey)
		require.NoError(t, err)
		require.NotNil(t, val)

		// Delete from cache
		err = cacheManager.Generic().Delete(ctx, testKey)
		require.NoError(t, err)

		// Should be gone
		val, err = cacheManager.Generic().Get(ctx, testKey)
		require.NoError(t, err)
		assert.Nil(t, val)
	})

	t.Run("Cache Batch Operations", func(t *testing.T) {
		keys := []string{"batch_key_1", "batch_key_2", "batch_key_3"}
		values := []string{"value_1", "value_2", "value_3"}

		// Set multiple values
		for i, key := range keys {
			err := cacheManager.Generic().Set(ctx, key, values[i], 1*time.Hour)
			require.NoError(t, err)
		}

		// Retrieve all
		for i, key := range keys {
			val, err := cacheManager.Generic().Get(ctx, key)
			require.NoError(t, err)
			require.NotNil(t, val)
			assert.Equal(t, values[i], val)
		}

		// Delete all
		for _, key := range keys {
			err := cacheManager.Generic().Delete(ctx, key)
			require.NoError(t, err)
		}

		// Verify all deleted
		for _, key := range keys {
			val, err := cacheManager.Generic().Get(ctx, key)
			require.NoError(t, err)
			assert.Nil(t, val)
		}
	})

	t.Run("Cache Exists Check", func(t *testing.T) {
		testKey := "test_exists_key"

		// Should not exist
		exists, err := cacheManager.Generic().Exists(ctx, testKey)
		require.NoError(t, err)
		assert.False(t, exists)

		// Set value
		err = cacheManager.Generic().Set(ctx, testKey, "value", 1*time.Hour)
		require.NoError(t, err)

		// Should exist
		exists, err = cacheManager.Generic().Exists(ctx, testKey)
		require.NoError(t, err)
		assert.True(t, exists)
	})

	// Cleanup
	cleanupTestKeys(t, redisClient)
}

// TestCachePerformance tests cache performance with multiple operations
func TestCachePerformance(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping performance test in short mode")
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	defer redisClient.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := redisClient.Ping(ctx).Err()
	if err != nil {
		t.Skipf("Redis not available: %v", err)
	}

	logger := logrus.New()
	cacheManager := NewRedisCacheManager(redisClient, logger)

	t.Run("High throughput writes", func(t *testing.T) {
		start := time.Now()
		iterations := 1000

		for i := 0; i < iterations; i++ {
			key := "perf_key_" + time.Now().Format("20060102150405") + "_" + string(rune(i))
			err := cacheManager.Generic().Set(ctx, key, "value", 1*time.Hour)
			require.NoError(t, err)
		}

		elapsed := time.Since(start)
		opsPerSec := float64(iterations) / elapsed.Seconds()
		t.Logf("Wrote %d items in %v (%.0f ops/sec)", iterations, elapsed, opsPerSec)
		assert.Greater(t, opsPerSec, 100.0, "Should support at least 100 writes/sec")
	})

	t.Run("High throughput reads", func(t *testing.T) {
		// Setup: write some keys
		for i := 0; i < 100; i++ {
			key := "read_perf_key_" + string(rune(i))
			cacheManager.Generic().Set(ctx, key, "value", 1*time.Hour)
		}

		start := time.Now()
		iterations := 1000

		for i := 0; i < iterations; i++ {
			key := "read_perf_key_" + string(rune(i%100))
			cacheManager.Generic().Get(ctx, key)
		}

		elapsed := time.Since(start)
		opsPerSec := float64(iterations) / elapsed.Seconds()
		t.Logf("Read %d items in %v (%.0f ops/sec)", iterations, elapsed, opsPerSec)
		assert.Greater(t, opsPerSec, 1000.0, "Should support at least 1000 reads/sec")
	})
}

// TestCacheWithDatabase tests cache with actual database fallback
func TestCacheWithDatabase(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Check if PostgreSQL is available
	dsn := "host=localhost user=postgres password=postgres dbname=tunduc_test port=5432 sslmode=disable"
	_, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Skipf("PostgreSQL not available: %v", err)
	}

	// Check if Redis is available
	redisClient := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = redisClient.Ping(ctx).Err()
	if err != nil {
		t.Skipf("Redis not available: %v", err)
	}

	logger := logrus.New()
	cacheManager := NewRedisCacheManager(redisClient, logger)

	t.Run("Cache miss falls back to continue", func(t *testing.T) {
		testKey := "cache_fallback_test"

		// Ensure key doesn't exist in cache
		cacheManager.Generic().Delete(ctx, testKey)

		// Try to get non-existent key
		val, err := cacheManager.Generic().Get(ctx, testKey)
		require.NoError(t, err)
		assert.Nil(t, val)

		// Application would then query database and populate cache
		dbResult := "value_from_db"
		err = cacheManager.Generic().Set(ctx, testKey, dbResult, 1*time.Hour)
		require.NoError(t, err)

		// Now should get from cache
		val, err = cacheManager.Generic().Get(ctx, testKey)
		require.NoError(t, err)
		require.NotNil(t, val)
		assert.Equal(t, dbResult, val)
	})

	cleanupTestKeys(t, redisClient)
}

// Helper function to clean up test keys
func cleanupTestKeys(_ *testing.T, redisClient *redis.Client) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Очищаем с префиксом "cache:" так как используется Generic cache
	patterns := []string{
		"cache:test_*",
		"cache:batch_key_*",
		"cache:perf_key_*",
		"cache:read_perf_key_*",
		"test_*", // для совместимости
		"batch_key_*",
		"perf_key_*",
		"read_perf_key_*",
	}

	for _, pattern := range patterns {
		iter := redisClient.Scan(ctx, 0, pattern, 0).Iterator()
		for iter.Next(ctx) {
			redisClient.Del(ctx, iter.Val())
		}
	}

	// Очистка специальных тестовых ключей
	redisClient.Del(ctx, "cache_fallback_test")
}
