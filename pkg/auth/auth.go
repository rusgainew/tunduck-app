package auth

import (
	"crypto/sha256"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// HashPassword hashes a password using bcrypt
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(hash), nil
}

// VerifyPassword verifies a password against its hash
func VerifyPassword(hash string, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// HashTokenForBlacklist creates a SHA256 hash of JWT token for blacklist storage
// This prevents storing actual tokens in Redis for security reasons
func HashTokenForBlacklist(token string) string {
	hash := sha256.Sum256([]byte(token))
	return fmt.Sprintf("%x", hash)
}

// GenerateToken generates a JWT token with proper claims and signature
func GenerateToken(userID string, email string, secretKey string, duration interface{}) (string, error) {
	if userID == "" && email == "" {
		return "", fmt.Errorf("userID or email must be provided")
	}

	if secretKey == "" {
		return "", fmt.Errorf("secretKey is required")
	}

	// Преобразуем duration в time.Duration
	var tokenDuration time.Duration
	switch d := duration.(type) {
	case time.Duration:
		tokenDuration = d
	case int:
		tokenDuration = time.Duration(d) * time.Hour
	case int64:
		tokenDuration = time.Duration(d) * time.Hour
	default:
		tokenDuration = 24 * time.Hour // по умолчанию 24 часа
	}

	// Создаем claims с необходимой информацией
	claims := jwt.MapClaims{
		"user_id": userID,
		"email":   email,
		"iat":     time.Now().Unix(),
		"exp":     time.Now().Add(tokenDuration).Unix(),
	}

	// Создаем токен с алгоритмом HS256
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Подписываем токен секретным ключом
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}
