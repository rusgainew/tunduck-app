package services

import (
	"context"

	"github.com/rusgainew/tunduck-app/internal/models"
	"github.com/rusgainew/tunduck-app/pkg/cache"
	"github.com/rusgainew/tunduck-app/pkg/entity"
)

// UserService интерфейс для работы с пользователями
type UserService interface {
	Register(ctx context.Context, req *models.RegisterRequest) (*models.AuthResponse, error)
	Login(ctx context.Context, req *models.LoginRequest) (*models.AuthResponse, error)
	ValidateToken(token string) (*models.UserInfo, error)
	GetByUsername(ctx context.Context, username string) (*entity.User, error)
	GetByEmail(ctx context.Context, email string) (*entity.User, error)
	CacheWarmUsers(ctx context.Context, limit int) error
	SetCacheManager(cacheManager cache.CacheManager)
}
