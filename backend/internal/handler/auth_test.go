package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/ktc-plugin-hub/backend/internal/config"
	"github.com/ktc-plugin-hub/backend/internal/dto"
	mw "github.com/ktc-plugin-hub/backend/internal/middleware"
	"github.com/ktc-plugin-hub/backend/internal/repository"
	"github.com/ktc-plugin-hub/backend/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var testJWTConfig = &config.JWTConfig{
	Secret:          "test-handler-secret",
	AccessTokenTTL:  3600,
	RefreshTokenTTL: 604800,
}

func init() {
	gin.SetMode(gin.TestMode)
}

func setupAuthRouter() (*gin.Engine, *AuthHandler) {
	userRepo := repository.NewMockUserRepository()
	authSvc := service.NewAuthService(userRepo, testJWTConfig)
	handler := NewAuthHandler(authSvc)

	r := gin.New()
	auth := r.Group("/api/v1/auth")
	{
		auth.POST("/register", handler.Register)
		auth.POST("/login", handler.Login)
		auth.POST("/refresh", handler.Refresh)
	}
	r.GET("/api/v1/me", mw.AuthMiddleware(testJWTConfig.Secret), handler.Me)

	return r, handler
}

// POST /auth/register 가 201을 반환하는지 확인
func TestShouldReturn201WhenRegisterSuccessfully(t *testing.T) {
	r, _ := setupAuthRouter()

	body, _ := json.Marshal(dto.RegisterRequest{
		Email: "new@example.com", Password: "password123", Nickname: "newuser",
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, 201, w.Code)

	var resp dto.UserResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, "new@example.com", resp.Email)
	assert.Equal(t, "user", resp.Role)
}

// 중복 이메일 등록 시 409를 반환하는지 확인
func TestShouldReturn409WhenRegisterWithDuplicateEmail(t *testing.T) {
	r, _ := setupAuthRouter()

	body, _ := json.Marshal(dto.RegisterRequest{
		Email: "dup@example.com", Password: "password123", Nickname: "user1",
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Equal(t, 201, w.Code)

	// 같은 이메일로 다시 등록
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Equal(t, 409, w.Code)
}

// POST /auth/login 이 토큰을 반환하는지 확인
func TestShouldReturnTokensWhenLoginSuccessfully(t *testing.T) {
	r, _ := setupAuthRouter()

	// 회원가입
	regBody, _ := json.Marshal(dto.RegisterRequest{
		Email: "login@example.com", Password: "mypassword", Nickname: "loginer",
	})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer(regBody))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	// 로그인
	loginBody, _ := json.Marshal(dto.LoginRequest{
		Email: "login@example.com", Password: "mypassword",
	})
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(loginBody))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	var resp dto.TokenResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.NotEmpty(t, resp.AccessToken)
	assert.NotEmpty(t, resp.RefreshToken)
}

// 잘못된 비밀번호로 로그인 시 401을 반환하는지 확인
func TestShouldReturn401WhenLoginWithWrongPassword(t *testing.T) {
	r, _ := setupAuthRouter()

	regBody, _ := json.Marshal(dto.RegisterRequest{
		Email: "wrong@example.com", Password: "correct", Nickname: "wrong",
	})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer(regBody))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	loginBody, _ := json.Marshal(dto.LoginRequest{
		Email: "wrong@example.com", Password: "incorrect",
	})
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(loginBody))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, 401, w.Code)
}

// GET /me 가 사용자 정보를 반환하는지 확인
func TestShouldReturnUserInfoWhenGetMe(t *testing.T) {
	r, _ := setupAuthRouter()

	// 회원가입 + 로그인
	regBody, _ := json.Marshal(dto.RegisterRequest{
		Email: "me@example.com", Password: "mypass123", Nickname: "meuser",
	})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer(regBody))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	loginBody, _ := json.Marshal(dto.LoginRequest{
		Email: "me@example.com", Password: "mypass123",
	})
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(loginBody))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	var tokenResp dto.TokenResponse
	json.Unmarshal(w.Body.Bytes(), &tokenResp)

	// GET /me
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/v1/me", nil)
	req.Header.Set("Authorization", "Bearer "+tokenResp.AccessToken)
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	var userResp dto.UserResponse
	err := json.Unmarshal(w.Body.Bytes(), &userResp)
	require.NoError(t, err)
	assert.Equal(t, "me@example.com", userResp.Email)
}

// POST /auth/refresh 가 새 토큰을 반환하는지 확인
func TestShouldReturnNewTokensWhenRefresh(t *testing.T) {
	r, _ := setupAuthRouter()

	// 회원가입 + 로그인
	regBody, _ := json.Marshal(dto.RegisterRequest{
		Email: "ref@example.com", Password: "mypass123", Nickname: "refuser",
	})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer(regBody))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	loginBody, _ := json.Marshal(dto.LoginRequest{
		Email: "ref@example.com", Password: "mypass123",
	})
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(loginBody))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	var tokenResp dto.TokenResponse
	json.Unmarshal(w.Body.Bytes(), &tokenResp)

	// Refresh
	refreshBody, _ := json.Marshal(dto.RefreshRequest{
		RefreshToken: tokenResp.RefreshToken,
	})
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/v1/auth/refresh", bytes.NewBuffer(refreshBody))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	var newTokenResp dto.TokenResponse
	err := json.Unmarshal(w.Body.Bytes(), &newTokenResp)
	require.NoError(t, err)
	assert.NotEmpty(t, newTokenResp.AccessToken)
}
