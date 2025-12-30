package docs

// EXAMPLE: How to use Redis caching in your services

/*
package service_impl

import (
	"context"
	"time"

	"github.com/rusgainew/tunduck-app/internal/models"
	"github.com/rusgainew/tunduck-app/internal/repository"
	"github.com/rusgainew/tunduck-app/internal/services"
	"github.com/rusgainew/tunduck-app/pkg/cache"
	"github.com/sirupsen/logrus"
)

type userServiceWithCache struct {
	repo          repository.UserRepository
	cacheManager  cache.CacheManager
	cacheHelper   *cache.CacheHelper
	logger        *logrus.Logger
}

// Example: GetByUsername with caching
func (s *userServiceWithCache) GetByUsername(ctx context.Context, username string) (*entity.User, error) {
	cacheKey := "username:" + username

	// Try to get from cache first
	cached, err := s.cacheManager.User().Get(ctx, cacheKey)
	if err == nil && cached != nil {
		s.logger.WithField("username", username).Debug("Cache hit")
		return cached.(*entity.User), nil
	}

	// Cache miss - query database
	user, err := s.repo.GetByUsername(ctx, username)
	if err != nil {
		return nil, err
	}

	// Cache the result for 1 hour
	if err := s.cacheManager.User().Set(ctx, cacheKey, user, time.Hour); err != nil {
		s.logger.WithError(err).Warn("Failed to cache user")
		// Continue anyway - cache is optional
	}

	return user, nil
}

// Example: Update with cache invalidation
func (s *userServiceWithCache) Update(ctx context.Context, user *entity.User) error {
	// Update in database
	if err := s.repo.Update(ctx, user); err != nil {
		return err
	}

	// Invalidate all related caches
	_ = s.cacheHelper.InvalidateUserCache(ctx, user.ID.String())
	_ = s.cacheHelper.InvalidateUsersByEmailCache(ctx, user.Email)
	_ = s.cacheHelper.InvalidateUsersByUsernameCache(ctx, user.Username)

	s.logger.WithField("user_id", user.ID.String()).Info("User caches invalidated")
	return nil
}

// Example: Register with transaction AND caching
func (s *userServiceWithCache) Register(ctx context.Context, req *models.RegisterRequest) (*models.AuthResponse, error) {
	// Register user (with transaction)
	response, err := s.registerWithTransaction(ctx, req)
	if err != nil {
		return nil, err
	}

	// Cache the new user
	cacheKey := "username:" + req.Username
	if err := s.cacheManager.User().Set(ctx, cacheKey, response.User, time.Hour); err != nil {
		s.logger.WithError(err).Warn("Failed to cache new user")
	}

	return response, nil
}

// Example: Logout with token blacklist
func (s *userServiceWithCache) Logout(ctx context.Context, token string, expiryTime time.Time) error {
	// Add token to blacklist (stored in Redis with TTL)
	if err := middleware.AddTokenToBlacklist(ctx, token, expiryTime, s.cacheManager); err != nil {
		s.logger.WithError(err).Error("Failed to blacklist token")
		return err
	}

	s.logger.Debug("Token added to blacklist")
	return nil
}

// Example: Batch operations for better performance
func (s *userServiceWithCache) CacheWarmUsers(ctx context.Context, limit int) error {
	// Get all users
	users, err := s.repo.GetAll(ctx, limit)
	if err != nil {
		return err
	}

	// Prepare batch data
	batchData := make(map[string]interface{})
	for _, user := range users {
		keyByID := "id:" + user.ID.String()
		keyByUsername := "username:" + user.Username
		keyByEmail := "email:" + user.Email

		batchData[keyByID] = user
		batchData[keyByUsername] = user
		batchData[keyByEmail] = user
	}

	// Set all at once (more efficient)
	if err := s.cacheManager.User().SetMultiple(ctx, batchData, time.Hour); err != nil {
		s.logger.WithError(err).Warn("Failed to warm cache")
		return err
	}

	s.logger.WithField("count", len(users)).Info("Cache warmed")
	return nil
}

// Example: Clear cache patterns for bulk invalidation
func (s *userServiceWithCache) ClearOrgUserCache(ctx context.Context, orgID string) error {
	// Clear all users in organization
	pattern := "org:" + orgID + ":users:*"
	if err := s.cacheManager.User().Clear(ctx, pattern); err != nil {
		s.logger.WithError(err).Error("Failed to clear org users cache")
		return err
	}

	s.logger.WithField("org_id", orgID).Info("Organization user cache cleared")
	return nil
}

// Example: Check if token is blacklisted (used in JWT middleware)
func (s *userServiceWithCache) IsTokenBlacklisted(ctx context.Context, token string) (bool, error) {
	return middleware.IsTokenBlacklisted(ctx, token, s.cacheManager)
}
*/
