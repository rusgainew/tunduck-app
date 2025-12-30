package cache

import (
	"context"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Note: These tests require Redis to be running on localhost:6379
// Run with: go test -v ./pkg/cache/

func setupTestRedis(t *testing.T) (*redis.Client, *logrus.Logger) {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	// Test connection
	if err := client.Ping(context.Background()).Err(); err != nil {
		t.Skip("Redis not running, skipping tests")
	}

	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)

	return client, logger
}

func TestRedisCache_Set_Get(t *testing.T) {
	client, logger := setupTestRedis(t)
	defer client.Close()

	cache := NewRedisCache(client, logger, "test")
	ctx := context.Background()

	// Test data
	testData := map[string]interface{}{
		"id":     "123",
		"name":   "John",
		"email":  "john@example.com",
		"active": true,
		"score":  95.5,
	}

	// Set value
	err := cache.Set(ctx, "user:123", testData, 1*time.Hour)
	require.NoError(t, err)

	// Get value
	result, err := cache.Get(ctx, "user:123")
	require.NoError(t, err)
	require.NotNil(t, result)

	t.Logf("Retrieved value: %v", result)
}

func TestRedisCache_Delete(t *testing.T) {
	client, logger := setupTestRedis(t)
	defer client.Close()

	cache := NewRedisCache(client, logger, "test")
	ctx := context.Background()

	// Set value
	err := cache.Set(ctx, "temp:key", "value", 1*time.Hour)
	require.NoError(t, err)

	// Verify it exists
	exists, err := cache.Exists(ctx, "temp:key")
	require.NoError(t, err)
	assert.True(t, exists)

	// Delete
	err = cache.Delete(ctx, "temp:key")
	require.NoError(t, err)

	// Verify it's gone
	exists, err = cache.Exists(ctx, "temp:key")
	require.NoError(t, err)
	assert.False(t, exists)
}

func TestRedisCache_GetMultiple_SetMultiple(t *testing.T) {
	client, logger := setupTestRedis(t)
	defer client.Close()

	cache := NewRedisCache(client, logger, "test")
	ctx := context.Background()

	// Prepare batch data
	batchData := map[string]interface{}{
		"user:1": map[string]string{"name": "Alice", "role": "admin"},
		"user:2": map[string]string{"name": "Bob", "role": "user"},
		"user:3": map[string]string{"name": "Carol", "role": "user"},
	}

	// Set multiple
	err := cache.SetMultiple(ctx, batchData, 1*time.Hour)
	require.NoError(t, err)

	// Get multiple
	result, err := cache.GetMultiple(ctx, []string{"user:1", "user:2", "user:3"})
	require.NoError(t, err)
	assert.Equal(t, 3, len(result))

	for key := range batchData {
		assert.Contains(t, result, key)
	}
}

func TestRedisCache_TTL_Expiration(t *testing.T) {
	client, logger := setupTestRedis(t)
	defer client.Close()

	cache := NewRedisCache(client, logger, "test")
	ctx := context.Background()

	// Set value with short TTL
	err := cache.Set(ctx, "temp:ttl", "expires soon", 1*time.Second)
	require.NoError(t, err)

	// Should exist immediately
	exists, err := cache.Exists(ctx, "temp:ttl")
	require.NoError(t, err)
	assert.True(t, exists)

	// Wait for expiration
	time.Sleep(2 * time.Second)

	// Should not exist anymore
	exists, err = cache.Exists(ctx, "temp:ttl")
	require.NoError(t, err)
	assert.False(t, exists)
}

func TestRedisCache_Clear_Pattern(t *testing.T) {
	client, logger := setupTestRedis(t)
	defer client.Close()

	cache := NewRedisCache(client, logger, "test")
	ctx := context.Background()

	// Set multiple keys with pattern
	for i := 1; i <= 5; i++ {
		key := "pattern:item:" + string(rune('0'+i))
		err := cache.Set(ctx, key, "data", 1*time.Hour)
		require.NoError(t, err)
	}

	// Clear by pattern
	err := cache.Clear(ctx, "pattern:*")
	require.NoError(t, err)

	// Verify all are gone
	result, err := cache.GetMultiple(ctx, []string{
		"pattern:item:1",
		"pattern:item:2",
		"pattern:item:3",
	})
	require.NoError(t, err)
	assert.Equal(t, 0, len(result))
}

func TestRedisCacheManager_Namespaces(t *testing.T) {
	client, logger := setupTestRedis(t)
	defer client.Close()

	manager := NewRedisCacheManager(client, logger)
	ctx := context.Background()

	// Test different namespaces
	testCases := []struct {
		name      string
		namespace Cache
		prefix    string
	}{
		{"User", manager.User(), "user"},
		{"Organization", manager.Organization(), "org"},
		{"Document", manager.Document(), "doc"},
		{"Session", manager.Session(), "session"},
		{"Token", manager.Token(), "token"},
		{"Generic", manager.Generic(), "cache"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.namespace.Set(ctx, "key", "value", 1*time.Hour)
			require.NoError(t, err)

			result, err := tc.namespace.Get(ctx, "key")
			require.NoError(t, err)
			require.NotNil(t, result)

			// Cleanup
			_ = tc.namespace.Delete(ctx, "key")
		})
	}
}

func TestRedisCacheManager_Flush(t *testing.T) {
	client, logger := setupTestRedis(t)
	defer client.Close()

	manager := NewRedisCacheManager(client, logger)
	ctx := context.Background()

	// Set data in different caches
	_ = manager.User().Set(ctx, "user:1", "data", 1*time.Hour)
	_ = manager.Organization().Set(ctx, "org:1", "data", 1*time.Hour)
	_ = manager.Document().Set(ctx, "doc:1", "data", 1*time.Hour)

	// Flush all
	err := manager.Flush(ctx)
	require.NoError(t, err)

	// Verify all are cleared
	exists, err := manager.User().Exists(ctx, "user:1")
	require.NoError(t, err)
	assert.False(t, exists)

	exists, err = manager.Organization().Exists(ctx, "org:1")
	require.NoError(t, err)
	assert.False(t, exists)
}

func TestCacheHelper_TokenBlacklist(t *testing.T) {
	client, logger := setupTestRedis(t)
	defer client.Close()

	manager := NewRedisCacheManager(client, logger)
	helper := NewCacheHelper(manager)
	ctx := context.Background()

	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."

	// Blacklist token
	err := helper.InvalidateTokenBlacklist(ctx, token, 3600)
	require.NoError(t, err)

	// Check if blacklisted
	isBlacklisted, err := helper.IsTokenBlacklisted(ctx, token)
	require.NoError(t, err)
	assert.True(t, isBlacklisted)
}

func TestCacheHelper_InvalidateUserCache(t *testing.T) {
	client, logger := setupTestRedis(t)
	defer client.Close()

	manager := NewRedisCacheManager(client, logger)
	helper := NewCacheHelper(manager)
	ctx := context.Background()

	userID := "123"

	// Set user cache
	_ = manager.User().Set(ctx, userID, "user data", 1*time.Hour)

	// Verify it exists
	exists, err := manager.User().Exists(ctx, userID)
	require.NoError(t, err)
	assert.True(t, exists)

	// Invalidate
	err = helper.InvalidateUserCache(ctx, userID)
	require.NoError(t, err)

	// Verify it's gone
	exists, err = manager.User().Exists(ctx, userID)
	require.NoError(t, err)
	assert.False(t, exists)
}
