package models

import "github.com/google/uuid"

// RegisterRequest запрос на регистрацию
type RegisterRequest struct {
	Username        string `json:"username" validate:"required,min=3,max=50"`
	Email           string `json:"email" validate:"required,email"`
	FullName        string `json:"fullName" validate:"required,min=2,max=100"`
	Phone           string `json:"phone" validate:"required,min=10,max=20"`
	Password        string `json:"password" validate:"required,min=6,max=100"`
	ConfirmPassword string `json:"confirmPassword" validate:"required,eqfield=Password"`
}

// LoginRequest запрос на вход
type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// AuthResponse ответ при успешной аутентификации
type AuthResponse struct {
	Token string    `json:"token"`
	User  *UserInfo `json:"user"`
}

// UserInfo информация о пользователе
type UserInfo struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
	Email    string    `json:"email"`
	FullName string    `json:"fullName"`
	Phone    string    `json:"phone"`
	IsActive bool      `json:"isActive"`
}

// ErrorResponse ответ с ошибкой
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}
