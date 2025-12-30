# Redis Caching Configuration

## Setup

### 1. Install Redis

```bash
# macOS
brew install redis
brew services start redis

# Linux (Ubuntu/Debian)
sudo apt-get install redis-server
sudo systemctl start redis-server

# Docker
docker run -d -p 6379:6379 redis:latest
```

### 2. Environment Variables

Add to your `.env` file:

```
REDIS_HOST=localhost
REDIS_PORT=6379
```

### 3. Verify Connection

```bash
redis-cli ping
# Should return: PONG
```

## Usage Examples

### Caching User Data

```go
// In user service
func (s *userService) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
    cacheKey := "email:" + email

    // Try cache first
    cached, _ := s.cacheManager.User().Get(ctx, cacheKey)
    if cached != nil {
        return cached.(*entity.User), nil
    }

    // Query database
    user, err := s.repo.GetByEmail(ctx, email)
    if err != nil {
        return nil, err
    }

    // Cache for 1 hour
    s.cacheManager.User().Set(ctx, cacheKey, user, time.Hour)
    return user, nil
}

// Invalidate on update
func (s *userService) Update(ctx context.Context, user *entity.User) error {
    if err := s.repo.Update(ctx, user); err != nil {
        return err
    }

    // Clear related caches
    s.cacheManager.User().Delete(ctx, user.ID.String())
    s.cacheManager.User().Delete(ctx, "email:" + user.Email)
    s.cacheManager.User().Delete(ctx, "username:" + user.Username)

    return nil
}
```

### JWT Token Blacklist (Logout)

```go
// Add logout endpoint
POST /api/auth/logout
Authorization: Bearer <token>

// Handler implementation
handler := middleware.CreateLogoutHandler(cacheManager, logger)

// Or manually:
err := middleware.AddTokenToBlacklist(
    ctx,
    token,
    tokenExpiryTime,
    cacheManager,
)
```

### Cache Warming

```go
func (s *userService) CacheWarmUsers(ctx context.Context, limit int) error {
    users, err := s.repo.GetAll(ctx, limit)
    if err != nil {
        return err
    }

    data := make(map[string]interface{})
    for _, user := range users {
        key := "id:" + user.ID.String()
        data[key] = user
    }

    return s.cacheManager.User().SetMultiple(ctx, data, time.Hour)
}
```

## Monitoring Redis

### Check Memory Usage

```bash
redis-cli info memory
```

### Monitor Commands

```bash
redis-cli monitor
```

### View All Keys

```bash
redis-cli keys "*"
```

### Clear Cache

```bash
redis-cli FLUSHDB      # Clear current database
redis-cli FLUSHALL     # Clear all databases
```

## Performance Tips

1. **Batch Operations** - Use GetMultiple/SetMultiple
2. **Appropriate TTL** - Don't cache too long (stale data) or too short (no benefit)
3. **Invalidate on Change** - Always clear cache when data is modified
4. **Monitor Hits** - Check DEBUG logs for cache effectiveness
5. **Use Patterns** - Clear related caches with pattern matching

## Troubleshooting

### Connection Refused

```
Error: Failed to connect to Redis
Solution:
- Check Redis is running: redis-cli ping
- Check REDIS_HOST and REDIS_PORT in .env
- Check firewall rules if using remote Redis
```

### Memory Issues

```
Check memory usage: redis-cli info memory
Clear old data: redis-cli FLUSHDB
Adjust TTL values if too high
Monitor with: redis-cli --stat
```

### Stale Data

```
Problem: Users see old data
Solution:
- Cache is working but not invalidated on updates
- Check that all Update/Delete operations clear cache
- Review cache invalidation logic in services
```

## Integration Checklist

- [ ] Redis installed and running
- [ ] REDIS_HOST and REDIS_PORT in .env
- [ ] Connection verified with redis-cli
- [ ] CacheManager in DI container
- [ ] Services updated with cache logic
- [ ] Cache invalidation on mutations
- [ ] JWT blacklist working for logout
- [ ] Monitoring setup for cache hits/misses
