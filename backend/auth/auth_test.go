package auth

import (
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testJWTSecret = "test-jwt-secret-for-testing-only"

func TestHashPassword(t *testing.T) {
	password := "testpassword123"

	hash, err := HashPassword(password)
	require.NoError(t, err)
	assert.NotEmpty(t, hash)
	assert.NotEqual(t, password, hash)

	// Hash should be different each time
	hash2, err := HashPassword(password)
	require.NoError(t, err)
	assert.NotEqual(t, hash, hash2)
}

func TestHashPassword_EmptyPassword(t *testing.T) {
	hash, err := HashPassword("")
	// bcrypt can handle empty passwords, but they're still hashed
	require.NoError(t, err)
	assert.NotEmpty(t, hash)
}

func TestVerifyPassword(t *testing.T) {
	password := "testpassword123"
	hash, err := HashPassword(password)
	require.NoError(t, err)

	// Correct password should match
	result := VerifyPassword(password, hash)
	assert.True(t, result)

	// Wrong password should not match
	result = VerifyPassword("wrongpassword", hash)
	assert.False(t, result)
}

func TestVerifyPassword_InvalidHash(t *testing.T) {
	result := VerifyPassword("password", "invalid-hash")
	assert.False(t, result)
}

func TestCustomClaims_Creation(t *testing.T) {
	userID := uuid.New().String()

	claims := &CustomClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	assert.Equal(t, userID, claims.UserID)
	assert.False(t, claims.RegisteredClaims.ExpiresAt.Before(time.Now()))
}

func TestJWTTokenGeneration(t *testing.T) {
	userID := uuid.New().String()
	secret := []byte(testJWTSecret)

	claims := &CustomClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(secret)

	require.NoError(t, err)
	assert.NotEmpty(t, tokenString)

	// Token should contain three parts separated by dots
	parts := strings.Split(tokenString, ".")
	assert.Equal(t, 3, len(parts))
}

func TestJWTTokenValidation(t *testing.T) {
	userID := uuid.New().String()
	secret := []byte(testJWTSecret)

	// Create token
	claims := &CustomClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(secret)
	require.NoError(t, err)

	// Parse and validate token
	parsedToken, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})

	require.NoError(t, err)
	assert.True(t, parsedToken.Valid)

	// Check claims
	if parsedClaims, ok := parsedToken.Claims.(*CustomClaims); ok {
		assert.Equal(t, userID, parsedClaims.UserID)
	} else {
		t.Fatal("Failed to parse claims")
	}
}

func TestJWTTokenValidation_InvalidToken(t *testing.T) {
	secret := []byte(testJWTSecret)

	_, err := jwt.ParseWithClaims("invalid.token.here", &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})

	assert.Error(t, err)
}

func TestJWTTokenValidation_WrongSecret(t *testing.T) {
	userID := uuid.New().String()
	correctSecret := []byte(testJWTSecret)
	wrongSecret := []byte("wrong-secret")

	// Create token with correct secret
	claims := &CustomClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(correctSecret)
	require.NoError(t, err)

	// Try to validate with wrong secret
	_, err = jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return wrongSecret, nil
	})

	assert.Error(t, err)
}

func TestJWTTokenExpiration(t *testing.T) {
	userID := uuid.New().String()
	secret := []byte(testJWTSecret)

	// Create expired token
	claims := &CustomClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-1 * time.Hour)), // Expired 1 hour ago
			IssuedAt:  jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(secret)
	require.NoError(t, err)

	// Try to parse expired token
	_, err = jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})

	assert.Error(t, err)
}

func TestPasswordHashingConsistency(t *testing.T) {
	password := "consistencytest123"

	// Hash the same password multiple times
	for i := 0; i < 5; i++ {
		hash, err := HashPassword(password)
		require.NoError(t, err)

		// Each hash should verify the original password
		result := VerifyPassword(password, hash)
		assert.True(t, result)

		// But should not verify wrong passwords
		result = VerifyPassword("wrongpassword", hash)
		assert.False(t, result)
	}
}

// Benchmark tests for performance monitoring
func BenchmarkHashPassword(b *testing.B) {
	password := "benchmarkpassword123"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := HashPassword(password)
		if err != nil {
			b.Fatalf("Failed to hash password: %v", err)
		}
	}
}

func BenchmarkVerifyPassword(b *testing.B) {
	password := "benchmarkpassword123"
	hash, err := HashPassword(password)
	if err != nil {
		b.Fatalf("Failed to hash password: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := VerifyPassword(password, hash)
		if !result {
			b.Fatalf("Failed to verify password")
		}
	}
}

func BenchmarkJWTGeneration(b *testing.B) {
	userID := uuid.New().String()
	secret := []byte(testJWTSecret)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		claims := &CustomClaims{
			UserID: userID,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
			},
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		_, err := token.SignedString(secret)
		if err != nil {
			b.Fatalf("Failed to generate JWT: %v", err)
		}
	}
}

func BenchmarkJWTValidation(b *testing.B) {
	userID := uuid.New().String()
	secret := []byte(testJWTSecret)

	// Pre-generate token
	claims := &CustomClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(secret)
	if err != nil {
		b.Fatalf("Failed to generate JWT: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
			return secret, nil
		})
		if err != nil {
			b.Fatalf("Failed to validate JWT: %v", err)
		}
	}
}

func BenchmarkJWTRoundTrip(b *testing.B) {
	userID := uuid.New().String()
	secret := []byte(testJWTSecret)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Generate token
		claims := &CustomClaims{
			UserID: userID,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
			},
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString(secret)
		if err != nil {
			b.Fatalf("Failed to generate JWT: %v", err)
		}

		// Validate token
		_, err = jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
			return secret, nil
		})
		if err != nil {
			b.Fatalf("Failed to validate JWT: %v", err)
		}
	}
}
