package entity

import (
	"fmt"
	"regexp"
	"time"

	"github.com/google/uuid"
	"github.com/rusgainew/tunduck-app/pkg/rbac"
	"gorm.io/gorm"
)

// User представляет пользователя системы
type User struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	Username  string         `gorm:"uniqueIndex;not null" json:"username"`
	Email     string         `gorm:"uniqueIndex;not null" json:"email"`
	FullName  string         `gorm:"not null" json:"fullName"`
	Phone     string         `gorm:"not null" json:"phone"`
	Password  string         `gorm:"not null" json:"-"`                   // Хешированный пароль, не возвращается в JSON
	Role      rbac.Role      `gorm:"default:'user';not null" json:"role"` // Роль пользователя
	IsActive  bool           `gorm:"default:true" json:"isActive"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName возвращает имя таблицы для GORM
func (User) TableName() string {
	return "users"
}

// Validate проверяет валидность данных пользователя
func (u *User) Validate() error {
	if u.Username == "" {
		return fmt.Errorf("username cannot be empty")
	}
	if len(u.Username) < 3 {
		return fmt.Errorf("username must be at least 3 characters long")
	}
	if len(u.Username) > 50 {
		return fmt.Errorf("username must be at most 50 characters long")
	}

	if !isValidEmail(u.Email) {
		return fmt.Errorf("invalid email format")
	}

	if u.FullName == "" {
		return fmt.Errorf("full name cannot be empty")
	}
	if len(u.FullName) < 2 {
		return fmt.Errorf("full name must be at least 2 characters long")
	}
	if len(u.FullName) > 100 {
		return fmt.Errorf("full name must be at most 100 characters long")
	}

	if u.Phone == "" {
		return fmt.Errorf("phone cannot be empty")
	}
	if len(u.Phone) < 10 || len(u.Phone) > 20 {
		return fmt.Errorf("phone must be between 10 and 20 characters long")
	}

	if u.Password == "" {
		return fmt.Errorf("password cannot be empty")
	}

	if !u.Role.IsValid() {
		return fmt.Errorf("invalid role: %s", u.Role)
	}

	return nil
}

// isValidEmail проверяет корректность формата email
func isValidEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}
