package service_impl

import (
	"context"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/rusgainew/tunduck-app/internal/models"
	"github.com/rusgainew/tunduck-app/internal/repository"
	"github.com/rusgainew/tunduck-app/internal/services"
	"github.com/rusgainew/tunduck-app/pkg/apperror"
	"github.com/rusgainew/tunduck-app/pkg/cache"
	"github.com/rusgainew/tunduck-app/pkg/entity"
	"github.com/rusgainew/tunduck-app/pkg/logger"
	"github.com/rusgainew/tunduck-app/pkg/transaction"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type userService struct {
	repo         repository.UserRepository
	db           *gorm.DB
	logger       *logger.Logger
	cacheManager cache.CacheManager
	cacheHelper  *cache.CacheHelper
}

// NewUserService создает новый user service с обязательными зависимостями
func NewUserService(repo repository.UserRepository, db *gorm.DB, log *logrus.Logger) services.UserService {
	if db == nil {
		log.Fatal("database connection is required for UserService")
		return nil
	}
	return &userService{
		repo:         repo,
		db:           db,
		logger:       logger.New(log),
		cacheManager: nil, // CacheManager будет установлен позже
		cacheHelper:  nil, // CacheHelper будет установлен позже
	}
}

// SetCacheManager устанавливает CacheManager для использования кеша в сервисе
func (s *userService) SetCacheManager(cacheManager cache.CacheManager) {
	s.cacheManager = cacheManager
	s.cacheHelper = cache.NewCacheHelper(cacheManager)
}

func (s *userService) Register(ctx context.Context, req *models.RegisterRequest) (*models.AuthResponse, error) {
	s.logger.Info(ctx, "Starting user registration", logrus.Fields{"username": req.Username, "email": req.Email})

	var response *models.AuthResponse

	// Используем транзакцию для регистрации
	err := transaction.Execute(ctx, s.db, &logrus.Logger{}, func(txDB *gorm.DB) error {
		// Проверяем, существует ли пользователь с таким username
		var existingUserCount int64
		if err := txDB.Model(&entity.User{}).Where("username = ?", req.Username).Count(&existingUserCount).Error; err != nil {
			s.logger.Error(ctx, "Failed to check username existence", err, logrus.Fields{"username": req.Username})
			return apperror.DatabaseError("checking username", err)
		}
		if existingUserCount > 0 {
			s.logger.Warn(ctx, "Registration failed: username already exists", logrus.Fields{"username": req.Username})
			return apperror.New(apperror.ErrUsernameExists, "username already exists")
		}

		// Проверяем email
		var existingEmailCount int64
		if err := txDB.Model(&entity.User{}).Where("email = ?", req.Email).Count(&existingEmailCount).Error; err != nil {
			s.logger.Error(ctx, "Failed to check email existence", err, logrus.Fields{"email": req.Email})
			return apperror.DatabaseError("checking email", err)
		}
		if existingEmailCount > 0 {
			s.logger.Warn(ctx, "Registration failed: email already exists", logrus.Fields{"email": req.Email})
			return apperror.New(apperror.ErrEmailExists, "email already exists")
		}

		// Хешируем пароль
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			s.logger.Error(ctx, "Failed to hash password", err)
			return apperror.New(apperror.ErrInternal, "password processing error")
		}

		// Создаём пользователя внутри транзакции
		user := &entity.User{
			ID:       uuid.New(),
			Username: req.Username,
			Email:    req.Email,
			FullName: req.FullName,
			Phone:    req.Phone,
			Password: string(hashedPassword),
			IsActive: true,
		}

		if err := txDB.Create(user).Error; err != nil {
			s.logger.Error(ctx, "Failed to create user", err, logrus.Fields{"user_id": user.ID})
			return apperror.DatabaseError("creating user", err)
		}

		s.logger.Info(ctx, "User registered successfully", logrus.Fields{"user_id": user.ID, "username": user.Username})

		// Генерируем JWT токен
		token, err := s.generateToken(ctx, user)
		if err != nil {
			return err
		}

		response = &models.AuthResponse{
			User: &models.UserInfo{
				ID:       user.ID,
				Username: user.Username,
				Email:    user.Email,
				FullName: user.FullName,
				Phone:    user.Phone,
				IsActive: user.IsActive,
			},
			Token: token,
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return response, nil
}

func (s *userService) Login(ctx context.Context, req *models.LoginRequest) (*models.AuthResponse, error) {
	s.logger.Info(ctx, "Starting user login", logrus.Fields{"username": req.Username})

	// Пытаемся получить из кеша
	var user *entity.User
	if s.cacheManager != nil {
		cacheKey := "username:" + req.Username
		cached, _ := s.cacheManager.User().Get(ctx, cacheKey)
		if cached != nil {
			if cachedUser, ok := cached.(*entity.User); ok {
				user = cachedUser
				s.logger.Debug(ctx, "User found in cache", logrus.Fields{"username": req.Username})
			}
		}
	}

	// Если не в кеше, ищем в БД
	if user == nil {
		var err error
		user, err = s.repo.GetByUsername(ctx, req.Username)
		if err != nil {
			s.logger.Error(ctx, "Failed to lookup user", err, logrus.Fields{"username": req.Username})
			return nil, apperror.DatabaseError("looking up user", err)
		}
	}

	if user == nil {
		s.logger.Warn(ctx, "Login failed: user not found", logrus.Fields{"username": req.Username})
		return nil, apperror.New(apperror.ErrInvalidCredentials, "invalid username or password")
	}

	// Проверяем активность
	if !user.IsActive {
		s.logger.Warn(ctx, "Login failed: account blocked", logrus.Fields{"user_id": user.ID, "username": user.Username})
		return nil, apperror.New(apperror.ErrAccountBlocked, "account is blocked")
	}

	// Проверяем пароль
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		s.logger.Warn(ctx, "Login failed: invalid password", logrus.Fields{"user_id": user.ID, "username": user.Username})
		return nil, apperror.New(apperror.ErrInvalidCredentials, "invalid username or password")
	}

	// Кешируем пользователя (1 час)
	if s.cacheManager != nil {
		cacheKey := "username:" + user.Username
		_ = s.cacheManager.User().Set(ctx, cacheKey, user, time.Hour)
		_ = s.cacheManager.User().Set(ctx, "email:"+user.Email, user, time.Hour)
		_ = s.cacheManager.User().Set(ctx, "id:"+user.ID.String(), user, time.Hour)
	}

	// Генерируем токен
	token, err := s.generateToken(ctx, user)
	if err != nil {
		return nil, err
	}

	s.logger.Info(ctx, "User logged in successfully", logrus.Fields{"user_id": user.ID, "username": user.Username})

	return &models.AuthResponse{
		Token: token,
		User: &models.UserInfo{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
			FullName: user.FullName,
			Phone:    user.Phone,
			IsActive: user.IsActive,
		},
	}, nil
}

func (s *userService) ValidateToken(tokenString string) (*models.UserInfo, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return nil, apperror.New(apperror.ErrConfigError, "JWT_SECRET is not configured")
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, apperror.New(apperror.ErrInvalidToken, "invalid token signing method")
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, apperror.New(apperror.ErrInvalidToken, "failed to parse token").WithError(err)
	}

	if !token.Valid {
		return nil, apperror.New(apperror.ErrInvalidToken, "token is invalid or expired")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, apperror.New(apperror.ErrInvalidToken, "invalid token claims")
	}

	userID, err := uuid.Parse(claims["user_id"].(string))
	if err != nil {
		return nil, apperror.New(apperror.ErrInvalidToken, "invalid user_id in token").WithError(err)
	}

	return &models.UserInfo{
		ID:       userID,
		Username: claims["username"].(string),
		Email:    claims["email"].(string),
		FullName: claims["full_name"].(string),
	}, nil
}

func (s *userService) generateToken(ctx context.Context, user *entity.User) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		s.logger.Error(ctx, "JWT_SECRET is not configured", nil)
		return "", apperror.New(apperror.ErrConfigError, "JWT_SECRET is not configured")
	}

	claims := jwt.MapClaims{
		"user_id":   user.ID.String(),
		"username":  user.Username,
		"email":     user.Email,
		"full_name": user.FullName,
		"exp":       time.Now().Add(time.Hour * 24 * 7).Unix(), // 7 дней
		"iat":       time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		s.logger.Error(ctx, "Failed to sign token", err, logrus.Fields{"user_id": user.ID})
		return "", apperror.New(apperror.ErrInternal, "failed to generate token").WithError(err)
	}

	s.logger.Debug(ctx, "Token generated", logrus.Fields{"user_id": user.ID, "expires_in": "7 days"})
	return tokenString, nil
}

// GetByUsername получает пользователя по username с кешированием
func (s *userService) GetByUsername(ctx context.Context, username string) (*entity.User, error) {
	// Проверяем кеш
	if s.cacheManager != nil {
		cacheKey := "username:" + username
		cached, _ := s.cacheManager.User().Get(ctx, cacheKey)
		if cached != nil {
			user := cached.(*entity.User)
			s.logger.Debug(ctx, "User retrieved from cache", logrus.Fields{"username": username})
			return user, nil
		}
	}

	// Кеш-мисс, получаем из БД
	user, err := s.repo.GetByUsername(ctx, username)
	if err != nil {
		return nil, err
	}

	// Кешируем результат на 1 час
	if s.cacheManager != nil && user != nil {
		cacheKey := "username:" + username
		_ = s.cacheManager.User().Set(ctx, cacheKey, user, time.Hour)
	}

	return user, nil
}

// GetByEmail получает пользователя по email с кешированием
func (s *userService) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	// Проверяем кеш
	if s.cacheManager != nil {
		cacheKey := "email:" + email
		cached, _ := s.cacheManager.User().Get(ctx, cacheKey)
		if cached != nil {
			user := cached.(*entity.User)
			s.logger.Debug(ctx, "User retrieved from cache", logrus.Fields{"email": email})
			return user, nil
		}
	}

	// Кеш-мисс, получаем из БД
	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	// Кешируем результат на 1 час
	if s.cacheManager != nil && user != nil {
		cacheKey := "email:" + email
		_ = s.cacheManager.User().Set(ctx, cacheKey, user, time.Hour)
	}

	return user, nil
}

// CacheWarmUsers предварительно загружает пользователей в кеш
func (s *userService) CacheWarmUsers(ctx context.Context, limit int) error {
	if s.cacheManager == nil {
		return nil // Кеш не настроен, пропускаем
	}

	s.logger.Info(ctx, "Starting cache warming for users", logrus.Fields{"limit": limit})

	// Получаем пользователей из БД
	users, err := s.repo.GetAll(ctx, limit)
	if err != nil {
		s.logger.Error(ctx, "Failed to get users for cache warming", err)
		return err
	}

	// Подготавливаем данные для пакетного кеширования
	batchData := make(map[string]interface{})
	for _, user := range users {
		batchData["id:"+user.ID.String()] = user
		batchData["username:"+user.Username] = user
		batchData["email:"+user.Email] = user
	}

	// Кешируем все сразу (более эффективно)
	if err := s.cacheManager.User().SetMultiple(ctx, batchData, time.Hour); err != nil {
		s.logger.Error(ctx, "Failed to warm cache", err)
		return err
	}

	s.logger.Info(ctx, "Cache warming completed", logrus.Fields{"users_cached": len(users)})
	return nil
}
