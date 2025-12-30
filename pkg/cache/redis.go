package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rusgainew/tunduck-app/pkg/apperror"
	"github.com/sirupsen/logrus"
)

// RedisCache реализация Cache с использованием Redis
type RedisCache struct {
	client *redis.Client
	logger *logrus.Logger
	prefix string
}

// NewRedisCache создает новый Redis кеш
func NewRedisCache(client *redis.Client, logger *logrus.Logger, prefix string) *RedisCache {
	return &RedisCache{
		client: client,
		logger: logger,
		prefix: prefix,
	}
}

// Get получает значение из кеша
func (r *RedisCache) Get(ctx context.Context, key string) (interface{}, error) {
	fullKey := r.getFullKey(key)

	val, err := r.client.Get(ctx, fullKey).Result()
	if err != nil {
		if err == redis.Nil {
			r.logger.WithFields(logrus.Fields{
				"key":    key,
				"prefix": r.prefix,
			}).Debug("Cache miss")
			return nil, nil
		}
		r.logger.WithError(err).WithField("key", fullKey).Error("Failed to get from cache")
		return nil, apperror.New(apperror.ErrInternal, fmt.Sprintf("cache get error: %v", err))
	}

	r.logger.WithFields(logrus.Fields{
		"key":    key,
		"prefix": r.prefix,
	}).Debug("Cache hit")

	// Пытаемся распарсить JSON
	var data interface{}
	if err := json.Unmarshal([]byte(val), &data); err != nil {
		// Если не JSON, возвращаем строку
		return val, nil
	}

	return data, nil
}

// Set устанавливает значение в кеш с TTL
func (r *RedisCache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	fullKey := r.getFullKey(key)

	// Сериализуем значение в JSON
	jsonValue, err := json.Marshal(value)
	if err != nil {
		r.logger.WithError(err).WithField("key", fullKey).Error("Failed to marshal cache value")
		return apperror.New(apperror.ErrInternal, "cache marshal error")
	}

	if err := r.client.Set(ctx, fullKey, jsonValue, ttl).Err(); err != nil {
		r.logger.WithError(err).WithField("key", fullKey).Error("Failed to set cache")
		return apperror.New(apperror.ErrInternal, fmt.Sprintf("cache set error: %v", err))
	}

	r.logger.WithFields(logrus.Fields{
		"key":    key,
		"prefix": r.prefix,
		"ttl":    ttl,
	}).Debug("Value set to cache")

	return nil
}

// Delete удаляет значение из кеша
func (r *RedisCache) Delete(ctx context.Context, key string) error {
	fullKey := r.getFullKey(key)

	if err := r.client.Del(ctx, fullKey).Err(); err != nil {
		r.logger.WithError(err).WithField("key", fullKey).Error("Failed to delete from cache")
		return apperror.New(apperror.ErrInternal, "cache delete error")
	}

	r.logger.WithField("key", key).Debug("Value deleted from cache")

	return nil
}

// Exists проверяет наличие ключа в кеше
func (r *RedisCache) Exists(ctx context.Context, key string) (bool, error) {
	fullKey := r.getFullKey(key)

	exists, err := r.client.Exists(ctx, fullKey).Result()
	if err != nil {
		r.logger.WithError(err).WithField("key", fullKey).Error("Failed to check cache key existence")
		return false, apperror.New(apperror.ErrInternal, "cache exists check error")
	}

	return exists > 0, nil
}

// Clear удаляет все значения по паттерну
func (r *RedisCache) Clear(ctx context.Context, pattern string) error {
	fullPattern := r.getFullKey(pattern)

	// Получаем все ключи по паттерну
	keys, err := r.client.Keys(ctx, fullPattern).Result()
	if err != nil {
		r.logger.WithError(err).WithField("pattern", fullPattern).Error("Failed to get cache keys by pattern")
		return apperror.New(apperror.ErrInternal, "cache pattern search error")
	}

	if len(keys) > 0 {
		if err := r.client.Del(ctx, keys...).Err(); err != nil {
			r.logger.WithError(err).WithField("pattern", fullPattern).Error("Failed to clear cache by pattern")
			return apperror.New(apperror.ErrInternal, "cache clear error")
		}
	}

	r.logger.WithFields(logrus.Fields{
		"pattern": pattern,
		"deleted": len(keys),
	}).Debug("Cache cleared by pattern")

	return nil
}

// GetMultiple получает несколько значений одновременно
func (r *RedisCache) GetMultiple(ctx context.Context, keys []string) (map[string]interface{}, error) {
	if len(keys) == 0 {
		return make(map[string]interface{}), nil
	}

	// Преобразуем ключи с префиксом
	fullKeys := make([]string, len(keys))
	for i, key := range keys {
		fullKeys[i] = r.getFullKey(key)
	}

	// Получаем значения
	vals, err := r.client.MGet(ctx, fullKeys...).Result()
	if err != nil {
		r.logger.WithError(err).WithField("count", len(keys)).Error("Failed to get multiple values from cache")
		return nil, apperror.New(apperror.ErrInternal, "cache mget error")
	}

	result := make(map[string]interface{})
	for i, val := range vals {
		if val != nil {
			result[keys[i]] = val
		}
	}

	r.logger.WithFields(logrus.Fields{
		"requested": len(keys),
		"found":     len(result),
	}).Debug("Multiple values retrieved from cache")

	return result, nil
}

// SetMultiple устанавливает несколько значений одновременно
func (r *RedisCache) SetMultiple(ctx context.Context, data map[string]interface{}, ttl time.Duration) error {
	if len(data) == 0 {
		return nil
	}

	pipe := r.client.Pipeline()
	for key, value := range data {
		fullKey := r.getFullKey(key)
		jsonValue, err := json.Marshal(value)
		if err != nil {
			r.logger.WithError(err).WithField("key", fullKey).Error("Failed to marshal cache value")
			continue
		}
		pipe.Set(ctx, fullKey, jsonValue, ttl)
	}

	_, err := pipe.Exec(ctx)
	if err != nil {
		r.logger.WithError(err).WithField("count", len(data)).Error("Failed to set multiple values in cache")
		return apperror.New(apperror.ErrInternal, "cache mset error")
	}

	r.logger.WithFields(logrus.Fields{
		"count": len(data),
		"ttl":   ttl,
	}).Debug("Multiple values set to cache")

	return nil
}

// getFullKey добавляет префикс к ключу
func (r *RedisCache) getFullKey(key string) string {
	if r.prefix == "" {
		return key
	}
	return fmt.Sprintf("%s:%s", r.prefix, key)
}

// RedisCacheManager реализация CacheManager с использованием Redis
type RedisCacheManager struct {
	client       *redis.Client
	logger       *logrus.Logger
	userCache    *RedisCache
	orgCache     *RedisCache
	docCache     *RedisCache
	sessionCache *RedisCache
	tokenCache   *RedisCache
	genericCache *RedisCache
}

// NewRedisCacheManager создает новый менеджер кешей
func NewRedisCacheManager(client *redis.Client, logger *logrus.Logger) *RedisCacheManager {
	return &RedisCacheManager{
		client:       client,
		logger:       logger,
		userCache:    NewRedisCache(client, logger, "user"),
		orgCache:     NewRedisCache(client, logger, "org"),
		docCache:     NewRedisCache(client, logger, "doc"),
		sessionCache: NewRedisCache(client, logger, "session"),
		tokenCache:   NewRedisCache(client, logger, "token"),
		genericCache: NewRedisCache(client, logger, "cache"),
	}
}

// User возвращает кеш пользователей
func (m *RedisCacheManager) User() Cache {
	return m.userCache
}

// Organization возвращает кеш организаций
func (m *RedisCacheManager) Organization() Cache {
	return m.orgCache
}

// Document возвращает кеш документов
func (m *RedisCacheManager) Document() Cache {
	return m.docCache
}

// Session возвращает кеш сессий
func (m *RedisCacheManager) Session() Cache {
	return m.sessionCache
}

// Token возвращает кеш токенов
func (m *RedisCacheManager) Token() Cache {
	return m.tokenCache
}

// Generic возвращает общий кеш
func (m *RedisCacheManager) Generic() Cache {
	return m.genericCache
}

// Flush очищает все кеши
func (m *RedisCacheManager) Flush(ctx context.Context) error {
	if err := m.client.FlushDB(ctx).Err(); err != nil {
		m.logger.WithError(err).Error("Failed to flush Redis cache")
		return apperror.New(apperror.ErrInternal, "cache flush error")
	}

	m.logger.Info("All caches flushed")
	return nil
}
