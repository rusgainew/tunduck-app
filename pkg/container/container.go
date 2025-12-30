package container

import (
	"github.com/go-playground/validator/v10"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"github.com/rusgainew/tunduck-app/internal/repository"
	repositorypostgres "github.com/rusgainew/tunduck-app/internal/repository/repository_postgres"
	"github.com/rusgainew/tunduck-app/internal/services"
	"github.com/rusgainew/tunduck-app/internal/services/service_impl"
	"github.com/rusgainew/tunduck-app/pkg/cache"
	"github.com/rusgainew/tunduck-app/pkg/logger"
	"github.com/rusgainew/tunduck-app/pkg/ratelimit"
)

// Container управляет всеми зависимостями приложения
type Container struct {
	// Logger
	logger *logger.Logger
	logrus *logrus.Logger

	// Database
	db *gorm.DB

	// Redis
	redisClient *redis.Client

	// Cache
	cacheManager cache.CacheManager

	// Rate Limiter
	rateLimiter *ratelimit.RateLimiter

	// Repositories
	userRepository repository.UserRepository
	docRepository  repository.EsfDocumentRepository

	// Services
	userService     services.UserService
	documentService services.EsfDocumentService

	// Validators
	validator *validator.Validate
}

// NewContainer создает и инициализирует контейнер зависимостей
func NewContainer(db *gorm.DB, log *logrus.Logger, redisClient *redis.Client) *Container {
	c := &Container{
		db:           db,
		logrus:       log,
		logger:       logger.New(log),
		validator:    validator.New(),
		redisClient:  redisClient,
		cacheManager: cache.NewRedisCacheManager(redisClient, log),
		rateLimiter:  ratelimit.NewRateLimiter(redisClient),
	}

	// Инициализируем repositories
	c.initRepositories()

	// Инициализируем services
	c.initServices()

	return c
}

// initRepositories инициализирует все repositories
func (c *Container) initRepositories() {
	c.userRepository = repositorypostgres.NewUserRepositoryPostgres(c.db, c.logrus)
	c.docRepository = repositorypostgres.NewEsfDocumentRepositoryPostgres(c.db, c.logrus)
}

// initServices инициализирует все services
func (c *Container) initServices() {
	// Передаем DB через конструктор
	c.userService = service_impl.NewUserService(c.userRepository, c.db, c.logrus)
	c.documentService = service_impl.NewEsfDocumentService(c.docRepository, c.db, c.logrus)

	// Установляем CacheManager в сервисы
	if c.cacheManager != nil {
		c.userService.SetCacheManager(c.cacheManager)
		c.documentService.SetCacheManager(c.cacheManager)
		// OrgService would be initialized here if added to container
		// orgService.SetCacheManager(c.cacheManager)
	}
}

// Getters для repositories
func (c *Container) GetUserRepository() repository.UserRepository {
	return c.userRepository
}

func (c *Container) GetEsfDocumentRepository() repository.EsfDocumentRepository {
	return c.docRepository
}

// Getters для services
func (c *Container) GetUserService() services.UserService {
	return c.userService
}

func (c *Container) GetEsfDocumentService() services.EsfDocumentService {
	return c.documentService
}

// Getters для других компонентов
func (c *Container) GetLogger() *logger.Logger {
	return c.logger
}

func (c *Container) GetLogrus() *logrus.Logger {
	return c.logrus
}

func (c *Container) GetDatabase() *gorm.DB {
	return c.db
}

func (c *Container) GetValidator() *validator.Validate {
	return c.validator
}

func (c *Container) GetCacheManager() cache.CacheManager {
	return c.cacheManager
}

func (c *Container) GetRateLimiter() *ratelimit.RateLimiter {
	return c.rateLimiter
}

func (c *Container) GetRedisClient() *redis.Client {
	return c.redisClient
}
