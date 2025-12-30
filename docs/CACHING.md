# Redis Caching Architecture

## Overview

Redis caching layer provides high-performance data caching with automatic TTL management and multi-backend support through abstraction interfaces.

## Architecture Components

### 1. Cache Interface (`pkg/cache/interface.go`)

Defines the contract for all cache implementations:

```go
type Cache interface {
    Get(ctx context.Context, key string) (interface{}, error)
    Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error
    Delete(ctx context.Context, key string) error
    Exists(ctx context.Context, key string) (bool, error)
    Clear(ctx context.Context, pattern string) error
    GetMultiple(ctx context.Context, keys []string) (map[string]interface{}, error)
    SetMultiple(ctx context.Context, data map[string]interface{}, ttl time.Duration) error
}
```

### 2. Redis Implementation (`pkg/cache/redis.go`)

- **RedisCache**: Implements Cache interface with Redis backend
- **RedisCacheManager**: Manages 6 separate cache namespaces
  - User cache (prefix: `user`)
  - Organization cache (prefix: `org`)
  - Document cache (prefix: `doc`)
  - Session cache (prefix: `session`)
  - Token cache (prefix: `token`)
  - Generic cache (prefix: `cache`)

### 3. Cache Helper (`pkg/cache/helper.go`)

Convenience methods for service layer:

- `InvalidateUserCache(ctx, userID)` - Remove user from cache
- `InvalidateUsersByEmailCache(ctx, email)` - Remove user by email lookup
- `InvalidateUsersByUsernameCache(ctx, username)` - Remove user by username lookup
- `InvalidateAllUsersCache(ctx)` - Clear all user cache
- `InvalidateTokenBlacklist(ctx, token, ttl)` - Add token to blacklist
- `IsTokenBlacklisted(ctx, token)` - Check if token is blacklisted
- `FlushAllCaches(ctx)` - Clear all caches

## Integration Points

### 1. DI Container (`pkg/container/container.go`)

```go
// Redis client is passed to NewContainer
container.NewContainer(db, logger, redisClient)

// Access cache manager from container
cacheManager := container.GetCacheManager()
```

### 2. Application Bootstrap (`cmd/api/app.go`)

```go
// Initialize Redis connection
redisClient := redis.NewClient(&redis.Options{
    Addr: "localhost:6379",
})

// Verify connection
redisClient.Ping(context.Background())

// Pass to container
container.NewContainer(db, logger, redisClient)
```

### 3. Service Layer Integration

```go
type userService struct {
    repo            repository.UserRepository
    db              *gorm.DB
    logger          *logger.Logger
    cacheManager    cache.CacheManager
    cacheHelper     *cache.CacheHelper
}

// Usage in GetByUsername
func (s *userService) GetByUsername(ctx context.Context, username string) (*entity.User, error) {
    // Try cache first
    cacheKey := "username:" + username
    cached, _ := s.cacheManager.User().Get(ctx, cacheKey)
    if cached != nil {
        return cached.(*entity.User), nil
    }

    // Query database
    user, err := s.repo.GetByUsername(ctx, username)
    if err != nil {
        return nil, err
    }

    // Cache result (1 hour TTL)
    s.cacheManager.User().Set(ctx, cacheKey, user, time.Hour)
    return user, nil
}
```

## Cache Key Patterns

### User Cache

- `user:<userID>` - User by ID
- `user:email:<email>` - User by email
- `user:username:<username>` - User by username

### Organization Cache

- `org:<orgID>` - Organization by ID
- `org:users:<orgID>` - Users in organization

### Document Cache

- `doc:<docID>` - Document by ID
- `doc:org:<orgID>` - Documents in organization

### Session Cache

- `session:<sessionID>` - Active session

### Token Cache (Blacklist)

- `token:blacklist:<tokenHash>` - Revoked JWT token

## TTL Strategy

| Cache Type         | TTL              | Use Case                     |
| ------------------ | ---------------- | ---------------------------- |
| User (static data) | 1 hour           | Lookup optimization          |
| User (search)      | 30 minutes       | Frequently accessed searches |
| Organization       | 2 hours          | Less frequent changes        |
| Document           | 30 minutes       | Document metadata            |
| Session            | 24 hours         | Active user sessions         |
| Token blacklist    | Until JWT expiry | Revoked tokens               |

## Error Handling

All cache operations include proper error handling:

- Cache misses return `nil` (not an error)
- Connection failures return `apperror.ErrInternal`
- Graceful degradation: app continues without cache

## Configuration

### Environment Variables

```
REDIS_HOST=localhost      # Default: localhost
REDIS_PORT=6379          # Default: 6379
```

### Graceful Degradation

- If Redis is unavailable at startup, logs warning but app continues
- If cache operation fails, app continues with DB query fallback

## Best Practices

1. **Always use context** - Cache operations respect deadlines
2. **Consistent key patterns** - Use `<type>:<id>` or `<type>:<field>:<value>`
3. **Appropriate TTL** - Balance between freshness and cache hits
4. **Invalidation on mutation** - Clear cache when data changes
5. **Batch operations** - Use GetMultiple/SetMultiple for bulk data

## Monitoring and Debugging

### Cache Hit/Miss Logging

All operations are logged with Debug level:

```
level=debug msg="Cache hit" key=user:123 prefix=user
level=debug msg="Cache miss" key=user:456 prefix=user
```

### Performance Metrics

Track in services:

- Cache hit rate
- Average response time (cached vs uncached)
- Memory usage

## Future Enhancements

1. **Cache warming** - Pre-load frequently accessed data
2. **Cache tags** - Group related keys for bulk invalidation
3. **Cache statistics** - Built-in hit/miss counters
4. **Distributed caching** - Redis cluster support
5. **Cache compression** - For large objects
