package cache

import "context"

// CacheHelper вспомогательный класс для работы с кешем в сервисах
type CacheHelper struct {
	cacheManager CacheManager
}

// NewCacheHelper создает новый helper для работы с кешем
func NewCacheHelper(cacheManager CacheManager) *CacheHelper {
	return &CacheHelper{
		cacheManager: cacheManager,
	}
}

// InvalidateUserCache инвалидирует кеш пользователя
func (h *CacheHelper) InvalidateUserCache(ctx context.Context, userID string) error {
	return h.cacheManager.User().Delete(ctx, userID)
}

// InvalidateUsersByEmailCache инвалидирует кеш поиска по email
func (h *CacheHelper) InvalidateUsersByEmailCache(ctx context.Context, email string) error {
	return h.cacheManager.User().Delete(ctx, "email:"+email)
}

// InvalidateUsersByUsernameCache инвалидирует кеш поиска по username
func (h *CacheHelper) InvalidateUsersByUsernameCache(ctx context.Context, username string) error {
	return h.cacheManager.User().Delete(ctx, "username:"+username)
}

// InvalidateAllUsersCache полностью очищает кеш пользователей
func (h *CacheHelper) InvalidateAllUsersCache(ctx context.Context) error {
	return h.cacheManager.User().Clear(ctx, "*")
}

// InvalidateOrgCache инвалидирует кеш организации
func (h *CacheHelper) InvalidateOrgCache(ctx context.Context, orgID string) error {
	return h.cacheManager.Organization().Delete(ctx, orgID)
}

// InvalidateAllOrgCache инвалидирует весь кеш организаций
func (h *CacheHelper) InvalidateAllOrgCache(ctx context.Context) error {
	return h.cacheManager.Organization().Clear(ctx, "*")
}

// InvalidateDocumentCache инвалидирует кеш документа
func (h *CacheHelper) InvalidateDocumentCache(ctx context.Context, docID string) error {
	return h.cacheManager.Document().Delete(ctx, docID)
}

// InvalidateAllDocumentCache инвалидирует весь кеш документов
func (h *CacheHelper) InvalidateAllDocumentCache(ctx context.Context) error {
	return h.cacheManager.Document().Clear(ctx, "*")
}

// InvalidateSessionCache инвалидирует сессию
func (h *CacheHelper) InvalidateSessionCache(ctx context.Context, sessionID string) error {
	return h.cacheManager.Session().Delete(ctx, sessionID)
}

// InvalidateTokenBlacklist добавляет токен в черный список
func (h *CacheHelper) InvalidateTokenBlacklist(ctx context.Context, token string, ttl int64) error {
	// Хранит токен в черном списке с TTL до истечения JWT
	return h.cacheManager.Token().Set(ctx, "blacklist:"+token, "revoked", 0)
}

// IsTokenBlacklisted проверяет находится ли токен в черном списке
func (h *CacheHelper) IsTokenBlacklisted(ctx context.Context, token string) (bool, error) {
	return h.cacheManager.Token().Exists(ctx, "blacklist:"+token)
}

// FlushAllCaches очищает все кеши
func (h *CacheHelper) FlushAllCaches(ctx context.Context) error {
	return h.cacheManager.Flush(ctx)
}
