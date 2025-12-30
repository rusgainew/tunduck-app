package auth

import (
	"crypto/sha256"
	"fmt"

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

// GenerateToken generates a JWT token (placeholder for full JWT implementation)
// In production, this would use a proper JWT library like github.com/golang-jwt/jwt
func GenerateToken(userID string, email string, secretKey string, duration interface{}) (string, error) {
	// This is a simplified implementation
	// In production, would create a proper JWT with claims, signature, etc.
	if userID == "" && email == "" {
		return "", fmt.Errorf("userID or email must be provided")
	}

	// Placeholder: return a dummy token-like string for testing
	// Real implementation would use jwt.NewWithClaims()
	signature := fmt.Sprintf("%x", sha256.Sum256([]byte(userID+email+secretKey)))
	tokenString := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIiK3VzZXJJRCsiLCJlbWFpbCI6IiIrZW1haWwrIiJ9.signature_" + signature
	return tokenString, nil
}
