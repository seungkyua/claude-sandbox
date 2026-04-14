package middleware

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testSecret = "test-jwt-secret-key"

// Access Token 생성 및 검증이 올바르게 동작하는지 확인
func TestShouldGenerateAndValidateAccessToken(t *testing.T) {
	userID := uint(1)
	email := "test@example.com"
	role := "user"
	ttl := 3600 // 1시간

	token, err := GenerateAccessToken(userID, email, role, testSecret, ttl)
	require.NoError(t, err)
	assert.NotEmpty(t, token)

	claims, err := ValidateToken(token, testSecret)
	require.NoError(t, err)
	assert.Equal(t, userID, claims.UserID)
	assert.Equal(t, email, claims.Email)
	assert.Equal(t, role, claims.Role)
	assert.Equal(t, "access", claims.TokenType)
}

// Refresh Token 생성 및 검증이 올바르게 동작하는지 확인
func TestShouldGenerateAndValidateRefreshToken(t *testing.T) {
	userID := uint(1)
	ttl := 604800 // 7일

	token, err := GenerateRefreshToken(userID, testSecret, ttl)
	require.NoError(t, err)
	assert.NotEmpty(t, token)

	claims, err := ValidateToken(token, testSecret)
	require.NoError(t, err)
	assert.Equal(t, userID, claims.UserID)
	assert.Equal(t, "refresh", claims.TokenType)
}

// 잘못된 시크릿으로 검증하면 에러가 발생하는지 확인
func TestShouldReturnErrorWhenTokenHasInvalidSecret(t *testing.T) {
	token, err := GenerateAccessToken(1, "test@example.com", "user", testSecret, 3600)
	require.NoError(t, err)

	_, err = ValidateToken(token, "wrong-secret")
	assert.Error(t, err)
}

// 만료된 토큰 검증 시 에러가 발생하는지 확인
func TestShouldReturnErrorWhenTokenIsExpired(t *testing.T) {
	token, err := GenerateAccessToken(1, "test@example.com", "user", testSecret, -1)
	require.NoError(t, err)

	// 약간 대기 후 검증
	time.Sleep(10 * time.Millisecond)
	_, err = ValidateToken(token, testSecret)
	assert.Error(t, err)
}
