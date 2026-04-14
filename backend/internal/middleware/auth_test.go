package middleware

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/ktc-plugin-hub/backend/internal/dto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// 유효한 토큰으로 요청하면 사용자 정보가 컨텍스트에 주입되는지 확인
func TestShouldInjectUserInfoWhenValidTokenProvided(t *testing.T) {
	token, err := GenerateAccessToken(1, "test@example.com", "user", testSecret, 3600)
	require.NoError(t, err)

	r := gin.New()
	r.Use(AuthMiddleware(testSecret))
	r.GET("/test", func(c *gin.Context) {
		userID, _ := GetUserIDFromContext(c)
		role, _ := GetUserRoleFromContext(c)
		c.JSON(200, gin.H{"user_id": userID, "role": role})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), `"user_id":1`)
	assert.Contains(t, w.Body.String(), `"role":"user"`)
}

// Authorization 헤더가 없으면 401을 반환하는지 확인
func TestShouldReturn401WhenNoAuthorizationHeader(t *testing.T) {
	r := gin.New()
	r.Use(AuthMiddleware(testSecret))
	r.GET("/test", func(c *gin.Context) {
		c.String(200, "ok")
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 401, w.Code)

	var errResp dto.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &errResp)
	assert.NoError(t, err)
	assert.Equal(t, "UNAUTHORIZED", errResp.Type)
}

// 관리자 미들웨어가 일반 사용자를 차단하는지 확인
func TestShouldReturn403WhenUserIsNotAdmin(t *testing.T) {
	token, err := GenerateAccessToken(1, "user@example.com", "user", testSecret, 3600)
	require.NoError(t, err)

	r := gin.New()
	r.Use(AuthMiddleware(testSecret))
	r.Use(AdminMiddleware())
	r.GET("/admin", func(c *gin.Context) {
		c.String(200, "admin only")
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/admin", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	r.ServeHTTP(w, req)

	assert.Equal(t, 403, w.Code)
}

// 관리자 미들웨어가 관리자를 허용하는지 확인
func TestShouldAllowAdminWhenUserIsAdmin(t *testing.T) {
	token, err := GenerateAccessToken(1, "admin@example.com", "admin", testSecret, 3600)
	require.NoError(t, err)

	r := gin.New()
	r.Use(AuthMiddleware(testSecret))
	r.Use(AdminMiddleware())
	r.GET("/admin", func(c *gin.Context) {
		c.String(200, "admin only")
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/admin", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}

// Refresh Token으로 인증하면 거부되는지 확인
func TestShouldRejectRefreshTokenAsAccessToken(t *testing.T) {
	token, err := GenerateRefreshToken(1, testSecret, 604800)
	require.NoError(t, err)

	r := gin.New()
	r.Use(AuthMiddleware(testSecret))
	r.GET("/test", func(c *gin.Context) {
		c.String(200, "ok")
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	r.ServeHTTP(w, req)

	assert.Equal(t, 401, w.Code)
}
