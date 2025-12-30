package service_impl

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/rusgainew/tunduck-app/internal/models"
	repositorypostgres "github.com/rusgainew/tunduck-app/internal/repository/repository_postgres"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Integration test для UserService с реальной БД
// Примечание: Требует запущенного PostgreSQL для полноценного тестирования
// Используйте testcontainers-go для автоматического управления контейнерами в CI/CD

func setupTestDB(t *testing.T) *gorm.DB {
	// TODO: Использовать testcontainers-go для создания контейнера PostgreSQL
	// Пример конфигурации:
	// req := testcontainers.ContainerRequest{
	//     Image:        "postgres:15",
	//     ExposedPorts: []string{"5432/tcp"},
	//     Env: map[string]string{
	//         "POSTGRES_USER":     "test",
	//         "POSTGRES_PASSWORD": "test",
	//         "POSTGRES_DB":       "testdb",
	//     },
	// }

	// На данный момент, создаем in-memory database (SQLite)
	// или используем тестовую БД если она настроена
	dsn := "user=postgres password=postgres dbname=tunduct_test host=localhost port=5432 sslmode=disable"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Skipf("Could not connect to test database: %v. Skipping integration tests.", err)
	}

	return db
}

// TestUserServiceRegisterIntegration тестирует регистрацию пользователя
func TestUserServiceRegisterIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	db := setupTestDB(t)
	log := logrus.New()

	repo := repositorypostgres.NewUserRepositoryPostgres(db, log)
	service := NewUserService(repo, db, log)

	ctx := context.Background()

	// Подготавливаем тестовые данные
	req := &models.RegisterRequest{
		Username:        "testuser" + uuid.New().String()[:8],
		Email:           "test" + uuid.New().String()[:8] + "@example.com",
		FullName:        "Test User",
		Phone:           "+1234567890",
		Password:        "SecurePassword123!",
		ConfirmPassword: "SecurePassword123!",
	}

	// Выполняем регистрацию
	resp, err := service.Register(ctx, req)

	// Проверяем результаты
	assert.NoError(t, err, "Registration should not error")
	assert.NotNil(t, resp, "Response should not be nil")
	assert.NotEmpty(t, resp.Token, "Token should not be empty")
	assert.Equal(t, req.Username, resp.User.Username, "Username should match")
	assert.Equal(t, req.Email, resp.User.Email, "Email should match")
}

// TestUserServiceLoginIntegration тестирует вход пользователя
func TestUserServiceLoginIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	db := setupTestDB(t)
	log := logrus.New()

	repo := repositorypostgres.NewUserRepositoryPostgres(db, log)
	service := NewUserService(repo, db, log)

	ctx := context.Background()

	// Регистрируем пользователя
	username := "testuser" + uuid.New().String()[:8]
	email := "test" + uuid.New().String()[:8] + "@example.com"
	password := "SecurePassword123!"

	registerReq := &models.RegisterRequest{
		Username:        username,
		Email:           email,
		FullName:        "Test User",
		Phone:           "+1234567890",
		Password:        password,
		ConfirmPassword: password,
	}

	_, err := service.Register(ctx, registerReq)
	require.NoError(t, err, "Registration should succeed")

	// Пытаемся войти
	loginReq := &models.LoginRequest{
		Username: username,
		Password: password,
	}

	resp, err := service.Login(ctx, loginReq)

	// Проверяем результаты
	assert.NoError(t, err, "Login should not error")
	assert.NotNil(t, resp, "Response should not be nil")
	assert.NotEmpty(t, resp.Token, "Token should not be empty")
	assert.Equal(t, username, resp.User.Username, "Username should match")
}

// TestUserServiceValidation тестирует валидацию пользователя
func TestUserServiceValidation(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	db := setupTestDB(t)
	log := logrus.New()

	repo := repositorypostgres.NewUserRepositoryPostgres(db, log)
	service := NewUserService(repo, db, log)

	ctx := context.Background()

	tests := []struct {
		name    string
		req     *models.RegisterRequest
		wantErr bool
	}{
		{
			name: "Invalid email format",
			req: &models.RegisterRequest{
				Username:        "testuser",
				Email:           "invalid-email",
				FullName:        "Test User",
				Phone:           "+1234567890",
				Password:        "SecurePassword123!",
				ConfirmPassword: "SecurePassword123!",
			},
			wantErr: true,
		},
		{
			name: "Password mismatch",
			req: &models.RegisterRequest{
				Username:        "testuser",
				Email:           "test@example.com",
				FullName:        "Test User",
				Phone:           "+1234567890",
				Password:        "SecurePassword123!",
				ConfirmPassword: "DifferentPassword!",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := service.Register(ctx, tt.req)
			if tt.wantErr {
				assert.Error(t, err, "Expected error for: "+tt.name)
			} else {
				assert.NoError(t, err, "Unexpected error for: "+tt.name)
			}
		})
	}
}
