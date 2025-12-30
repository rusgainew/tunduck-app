package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestHashPassword tests password hashing functionality
func TestHashPassword(t *testing.T) {
	t.Run("Hash password produces hash", func(t *testing.T) {
		password := "mySecurePassword123!"
		hash, err := HashPassword(password)
		require.NoError(t, err)
		assert.NotEmpty(t, hash)
		assert.NotEqual(t, password, hash)
	})

	t.Run("Different passwords produce different hashes", func(t *testing.T) {
		hash1, err := HashPassword("password1")
		require.NoError(t, err)

		hash2, err := HashPassword("password2")
		require.NoError(t, err)

		assert.NotEqual(t, hash1, hash2)
	})

	t.Run("Same password produces different hashes (salt)", func(t *testing.T) {
		password := "samePassword123!"
		hash1, err := HashPassword(password)
		require.NoError(t, err)

		hash2, err := HashPassword(password)
		require.NoError(t, err)

		assert.NotEqual(t, hash1, hash2)
	})
}

// TestVerifyPassword tests password verification
func TestVerifyPassword(t *testing.T) {
	t.Run("Verify correct password returns true", func(t *testing.T) {
		password := "correctPassword123!"
		hash, err := HashPassword(password)
		require.NoError(t, err)

		valid := VerifyPassword(hash, password)
		assert.True(t, valid)
	})

	t.Run("Verify incorrect password returns false", func(t *testing.T) {
		password := "correctPassword123!"
		hash, err := HashPassword(password)
		require.NoError(t, err)

		valid := VerifyPassword(hash, "wrongPassword123!")
		assert.False(t, valid)
	})

	t.Run("Verify against empty password fails", func(t *testing.T) {
		hash, err := HashPassword("password")
		require.NoError(t, err)

		valid := VerifyPassword(hash, "")
		assert.False(t, valid)
	})
}

// TestHashTokenForBlacklist tests token blacklist hashing
func TestHashTokenForBlacklist(t *testing.T) {
	t.Run("Token hash consistency", func(t *testing.T) {
		token := "test_jwt_token_12345"
		hash1 := HashTokenForBlacklist(token)
		hash2 := HashTokenForBlacklist(token)

		assert.Equal(t, hash1, hash2)
		assert.NotEmpty(t, hash1)
		assert.Greater(t, len(hash1), 30)
	})

	t.Run("Different tokens produce different hashes", func(t *testing.T) {
		hash1 := HashTokenForBlacklist("token1")
		hash2 := HashTokenForBlacklist("token2")

		assert.NotEqual(t, hash1, hash2)
	})

	t.Run("Token hash is hex string", func(t *testing.T) {
		token := "jwt_token_test"
		hash := HashTokenForBlacklist(token)

		// SHA256 hex should be 64 characters
		assert.Equal(t, 64, len(hash))

		// Should only contain hex characters
		for _, char := range hash {
			assert.True(t, (char >= '0' && char <= '9') ||
				(char >= 'a' && char <= 'f'),
				"Hash should only contain hex characters")
		}
	})
}

// TestPasswordEdgeCases tests edge cases in password handling
func TestPasswordEdgeCases(t *testing.T) {
	t.Run("Long password (within bcrypt limit)", func(t *testing.T) {
		// bcrypt has a 72-byte limit
		longPassword := ""
		for i := 0; i < 70; i++ {
			longPassword += "a"
		}

		hash, err := HashPassword(longPassword)
		require.NoError(t, err)
		assert.NotEmpty(t, hash)

		valid := VerifyPassword(hash, longPassword)
		assert.True(t, valid)
	})

	t.Run("Special characters in password", func(t *testing.T) {
		specialPassword := "P@ssw0rd!#$%^&*()_+-=[]{}|;:',.<>?/`~"
		hash, err := HashPassword(specialPassword)
		require.NoError(t, err)

		valid := VerifyPassword(hash, specialPassword)
		assert.True(t, valid)
	})

	t.Run("Unicode characters in password", func(t *testing.T) {
		unicodePassword := "Пароль123!密码"
		hash, err := HashPassword(unicodePassword)
		require.NoError(t, err)

		valid := VerifyPassword(hash, unicodePassword)
		assert.True(t, valid)
	})

	t.Run("Whitespace in password", func(t *testing.T) {
		passwordWithSpace := "pass word 123"
		hash, err := HashPassword(passwordWithSpace)
		require.NoError(t, err)

		valid := VerifyPassword(hash, passwordWithSpace)
		assert.True(t, valid)

		// Without space should fail
		valid = VerifyPassword(hash, "password123")
		assert.False(t, valid)
	})
}

// TestConcurrentPasswordOperations tests thread-safety
func TestConcurrentPasswordOperations(t *testing.T) {
	t.Run("Concurrent password hashing", func(t *testing.T) {
		done := make(chan bool, 10)
		hashes := make(map[string]bool)
		password := "testPassword123!"

		for i := 0; i < 10; i++ {
			go func() {
				hash, err := HashPassword(password)
				require.NoError(t, err)
				hashes[hash] = true
				done <- true
			}()
		}

		for i := 0; i < 10; i++ {
			<-done
		}

		// All hashes should be different (different salts)
		assert.Greater(t, len(hashes), 1)
	})

	t.Run("Concurrent password verification", func(t *testing.T) {
		password := "testPassword123!"
		hash, err := HashPassword(password)
		require.NoError(t, err)

		done := make(chan bool, 10)
		results := make([]bool, 10)

		for i := 0; i < 10; i++ {
			go func(idx int) {
				results[idx] = VerifyPassword(hash, password)
				done <- true
			}(i)
		}

		for i := 0; i < 10; i++ {
			<-done
		}

		// All verifications should succeed
		for _, result := range results {
			assert.True(t, result)
		}
	})
}
