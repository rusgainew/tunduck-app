package cache

import (
	"context"
	"time"
)

// Cache интерфейс для работы с кешем
type Cache interface {
	// Get получает значение из кеша
	Get(ctx context.Context, key string) (interface{}, error)

	// Set устанавливает значение в кеш с TTL
	Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error

	// Delete удаляет значение из кеша
	Delete(ctx context.Context, key string) error

	// Exists проверяет наличие ключа в кеше
	Exists(ctx context.Context, key string) (bool, error)

	// Clear удаляет все значения по паттерну
	Clear(ctx context.Context, pattern string) error

	// GetMultiple получает несколько значений одновременно
	GetMultiple(ctx context.Context, keys []string) (map[string]interface{}, error)

	// SetMultiple устанавливает несколько значений одновременно
	SetMultiple(ctx context.Context, data map[string]interface{}, ttl time.Duration) error
}

// CacheManager управляет различными кешами приложения
type CacheManager interface {
	// User кеш для пользователей
	User() Cache

	// Organization кеш для организаций
	Organization() Cache

	// Document кеш для документов
	Document() Cache

	// Session кеш для сессий
	Session() Cache

	// Token кеш для токенов (blacklist)
	Token() Cache

	// Generic общий кеш
	Generic() Cache

	// Flush очищает все кеши
	Flush(ctx context.Context) error
}
