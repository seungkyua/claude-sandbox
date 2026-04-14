package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/ktc-plugin-hub/backend/internal/dto"
	mw "github.com/ktc-plugin-hub/backend/internal/middleware"
	"github.com/ktc-plugin-hub/backend/internal/repository"
	"github.com/ktc-plugin-hub/backend/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupPluginRouter() *gin.Engine {
	pluginRepo := repository.NewMockPluginRepository()
	catRepo := repository.NewMockCategoryRepository()
	pluginSvc := service.NewPluginService(pluginRepo, catRepo)
	handler := NewPluginHandler(pluginSvc)

	r := gin.New()
	// 공개 라우트
	r.GET("/api/v1/plugins", handler.GetList)
	r.GET("/api/v1/plugins/:id", handler.GetByID)

	// 인증 필요 라우트
	auth := r.Group("/api/v1", mw.AuthMiddleware(testJWTConfig.Secret))
	{
		auth.POST("/plugins", handler.Create)
		auth.PUT("/plugins/:id", handler.Update)
		auth.DELETE("/plugins/:id", handler.Delete)
	}

	return r
}

func createPluginViaAPI(r *gin.Engine, token string, name string) *httptest.ResponseRecorder {
	body, _ := json.Marshal(dto.CreatePluginRequest{
		Name: name, Description: "test desc", CategoryID: 1, Version: "1.0.0",
	})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/plugins", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	r.ServeHTTP(w, req)
	return w
}

// POST /plugins 가 201을 반환하는지 확인
func TestShouldReturn201WhenCreatePlugin(t *testing.T) {
	r := setupPluginRouter()
	token, _ := mw.GenerateAccessToken(1, "user@test.com", "user", testJWTConfig.Secret, 3600)

	w := createPluginViaAPI(r, token, "new-plugin")
	assert.Equal(t, 201, w.Code)

	var resp dto.PluginResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, "new-plugin", resp.Name)
	assert.Equal(t, "pending", resp.Status)
}

// GET /plugins 가 목록을 반환하는지 확인
func TestShouldReturnPluginListWhenGetPlugins(t *testing.T) {
	r := setupPluginRouter()
	token, _ := mw.GenerateAccessToken(1, "user@test.com", "admin", testJWTConfig.Secret, 3600)

	// approved 플러그인 생성 (관리자)
	createPluginViaAPI(r, token, "approved-plugin-1")
	createPluginViaAPI(r, token, "approved-plugin-2")

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/plugins", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}

// GET /plugins/:id 가 상세를 반환하는지 확인
func TestShouldReturnPluginDetailWhenGetByID(t *testing.T) {
	r := setupPluginRouter()
	token, _ := mw.GenerateAccessToken(1, "user@test.com", "user", testJWTConfig.Secret, 3600)

	createW := createPluginViaAPI(r, token, "detail-plugin")
	var created dto.PluginResponse
	json.Unmarshal(createW.Body.Bytes(), &created)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/plugins/1", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}

// DELETE /plugins/:id 가 204를 반환하는지 확인
func TestShouldReturn204WhenDeletePlugin(t *testing.T) {
	r := setupPluginRouter()
	token, _ := mw.GenerateAccessToken(1, "user@test.com", "user", testJWTConfig.Secret, 3600)

	createPluginViaAPI(r, token, "delete-plugin")

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/api/v1/plugins/1", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	r.ServeHTTP(w, req)

	assert.Equal(t, 204, w.Code)
}

// 타인 플러그인 삭제 시 403 확인
func TestShouldReturn403WhenNonOwnerDeletesPlugin(t *testing.T) {
	r := setupPluginRouter()
	ownerToken, _ := mw.GenerateAccessToken(1, "owner@test.com", "user", testJWTConfig.Secret, 3600)
	otherToken, _ := mw.GenerateAccessToken(2, "other@test.com", "user", testJWTConfig.Secret, 3600)

	createPluginViaAPI(r, ownerToken, "owners-plugin")

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/api/v1/plugins/1", nil)
	req.Header.Set("Authorization", "Bearer "+otherToken)
	r.ServeHTTP(w, req)

	assert.Equal(t, 403, w.Code)
}
