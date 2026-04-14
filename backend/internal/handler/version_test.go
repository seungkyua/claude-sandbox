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
	"github.com/ktc-plugin-hub/backend/internal/model"
	"github.com/ktc-plugin-hub/backend/internal/repository"
	"github.com/ktc-plugin-hub/backend/internal/service"
	"github.com/stretchr/testify/assert"
)

func setupVersionRouter() *gin.Engine {
	pluginRepo := repository.NewMockPluginRepository()
	versionRepo := repository.NewMockVersionRepository()
	pluginRepo.Create(&model.Plugin{Name: "test-plugin", AuthorID: 1, CategoryID: 1, Description: "d", Status: "approved"})

	versionSvc := service.NewVersionService(versionRepo, pluginRepo)
	handler := NewVersionHandler(versionSvc)

	r := gin.New()
	auth := r.Group("/api/v1", mw.AuthMiddleware(testJWTConfig.Secret))
	{
		auth.POST("/plugins/:id/versions", handler.CreateVersion)
		auth.GET("/plugins/:id/versions/:versionId/download", handler.Download)
	}
	return r
}

func TestShouldReturn201WhenCreateVersion(t *testing.T) {
	r := setupVersionRouter()
	token, _ := mw.GenerateAccessToken(1, "user@test.com", "user", testJWTConfig.Secret, 3600)

	body, _ := json.Marshal(dto.CreateVersionRequest{Version: "1.0.0", Changelog: "initial release"})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/plugins/1/versions", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	r.ServeHTTP(w, req)

	assert.Equal(t, 201, w.Code)
}

func TestShouldReturn409WhenDuplicateVersion(t *testing.T) {
	r := setupVersionRouter()
	token, _ := mw.GenerateAccessToken(1, "user@test.com", "user", testJWTConfig.Secret, 3600)

	body, _ := json.Marshal(dto.CreateVersionRequest{Version: "1.0.0"})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/plugins/1/versions", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	r.ServeHTTP(w, req)

	// 같은 버전 다시 생성
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/v1/plugins/1/versions", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	r.ServeHTTP(w, req)

	assert.Equal(t, 409, w.Code)
}

func TestShouldReturnFileInfoWhenDownload(t *testing.T) {
	r := setupVersionRouter()
	token, _ := mw.GenerateAccessToken(1, "user@test.com", "user", testJWTConfig.Secret, 3600)

	// 버전 생성
	body, _ := json.Marshal(dto.CreateVersionRequest{Version: "1.0.0"})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/plugins/1/versions", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	r.ServeHTTP(w, req)

	// 다운로드
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/v1/plugins/1/versions/1/download", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}
