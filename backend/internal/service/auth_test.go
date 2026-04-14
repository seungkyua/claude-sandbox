package service

import (
	"testing"

	"github.com/ktc-plugin-hub/backend/internal/config"
	"github.com/ktc-plugin-hub/backend/internal/dto"
	"github.com/ktc-plugin-hub/backend/internal/middleware"
	"github.com/ktc-plugin-hub/backend/internal/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var testJWTConfig = &config.JWTConfig{
	Secret:          "test-secret-key",
	AccessTokenTTL:  3600,
	RefreshTokenTTL: 604800,
}

func newTestAuthService() AuthService {
	return NewAuthService(repository.NewMockUserRepository(), testJWTConfig)
}

// 회원가입이 성공적으로 완료되는지 확인
func TestShouldRegisterUserSuccessfully(t *testing.T) {
	svc := newTestAuthService()

	resp, err := svc.Register(&dto.RegisterRequest{
		Email:    "new@example.com",
		Password: "password123",
		Nickname: "newuser",
	})

	require.NoError(t, err)
	assert.NotZero(t, resp.ID)
	assert.Equal(t, "new@example.com", resp.Email)
	assert.Equal(t, "newuser", resp.Nickname)
	assert.Equal(t, "user", resp.Role)
}

// 중복 이메일로 회원가입 시 에러가 발생하는지 확인
func TestShouldReturnErrorWhenRegisterWithDuplicateEmail(t *testing.T) {
	svc := newTestAuthService()

	_, err := svc.Register(&dto.RegisterRequest{
		Email: "dup@example.com", Password: "password123", Nickname: "user1",
	})
	require.NoError(t, err)

	_, err = svc.Register(&dto.RegisterRequest{
		Email: "dup@example.com", Password: "password456", Nickname: "user2",
	})
	assert.Error(t, err)
	assert.Equal(t, ErrDuplicateEmail, err)
}

// 로그인이 성공적으로 완료되는지 확인
func TestShouldLoginSuccessfully(t *testing.T) {
	svc := newTestAuthService()

	_, err := svc.Register(&dto.RegisterRequest{
		Email: "login@example.com", Password: "mypassword", Nickname: "loginer",
	})
	require.NoError(t, err)

	tokenResp, err := svc.Login(&dto.LoginRequest{
		Email: "login@example.com", Password: "mypassword",
	})
	require.NoError(t, err)
	assert.NotEmpty(t, tokenResp.AccessToken)
	assert.NotEmpty(t, tokenResp.RefreshToken)
	assert.Equal(t, 3600, tokenResp.ExpiresIn)
}

// 잘못된 비밀번호로 로그인 시 에러가 발생하는지 확인
func TestShouldReturnErrorWhenLoginWithWrongPassword(t *testing.T) {
	svc := newTestAuthService()

	_, _ = svc.Register(&dto.RegisterRequest{
		Email: "wrong@example.com", Password: "correctpass", Nickname: "wronger",
	})

	_, err := svc.Login(&dto.LoginRequest{
		Email: "wrong@example.com", Password: "wrongpass",
	})
	assert.Error(t, err)
	assert.Equal(t, ErrInvalidCredentials, err)
}

// 존재하지 않는 이메일로 로그인 시 에러가 발생하는지 확인
func TestShouldReturnErrorWhenLoginWithNonexistentEmail(t *testing.T) {
	svc := newTestAuthService()

	_, err := svc.Login(&dto.LoginRequest{
		Email: "nope@example.com", Password: "whatever",
	})
	assert.Error(t, err)
	assert.Equal(t, ErrInvalidCredentials, err)
}

// RefreshToken으로 새 토큰을 발급받을 수 있는지 확인
func TestShouldRefreshTokenSuccessfully(t *testing.T) {
	svc := newTestAuthService()

	_, _ = svc.Register(&dto.RegisterRequest{
		Email: "refresh@example.com", Password: "pass123", Nickname: "refresher",
	})

	tokenResp, err := svc.Login(&dto.LoginRequest{
		Email: "refresh@example.com", Password: "pass123",
	})
	require.NoError(t, err)

	newTokenResp, err := svc.RefreshToken(tokenResp.RefreshToken)
	require.NoError(t, err)
	assert.NotEmpty(t, newTokenResp.AccessToken)
	assert.NotEmpty(t, newTokenResp.RefreshToken)
}

// 잘못된 RefreshToken으로 갱신 시 에러가 발생하는지 확인
func TestShouldReturnErrorWhenRefreshWithInvalidToken(t *testing.T) {
	svc := newTestAuthService()

	_, err := svc.RefreshToken("invalid-token")
	assert.Error(t, err)
}

// Access Token으로 Refresh 시도하면 에러가 발생하는지 확인
func TestShouldReturnErrorWhenRefreshWithAccessToken(t *testing.T) {
	svc := newTestAuthService()

	accessToken, _ := middleware.GenerateAccessToken(1, "test@example.com", "user", testJWTConfig.Secret, 3600)

	_, err := svc.RefreshToken(accessToken)
	assert.Error(t, err)
	assert.Equal(t, ErrInvalidToken, err)
}

// GetUserByID가 올바르게 동작하는지 확인
func TestShouldGetUserByID(t *testing.T) {
	svc := newTestAuthService()

	registered, _ := svc.Register(&dto.RegisterRequest{
		Email: "getme@example.com", Password: "pass123", Nickname: "getmyinfo",
	})

	user, err := svc.GetUserByID(registered.ID)
	require.NoError(t, err)
	assert.Equal(t, "getme@example.com", user.Email)
}
