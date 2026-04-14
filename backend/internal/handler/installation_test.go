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

func setupInstallationRouter() *gin.Engine {
	installRepo := repository.NewMockInstallationRepository()
	pluginRepo := repository.NewMockPluginRepository()
	versionRepo := repository.NewMockVersionRepository()

	pluginRepo.Create(&model.Plugin{Name: "test-plugin", AuthorID: 1, CategoryID: 1, Description: "d", Status: "approved"})
	versionRepo.Create(&model.PluginVersion{PluginID: 1, Version: "1.0.0", FilePath: "/p", FileSize: 1024})

	installSvc := service.NewInstallationService(installRepo, pluginRepo, versionRepo)
	handler := NewInstallationHandler(installSvc)

	r := gin.New()
	auth := r.Group("/api/v1", mw.AuthMiddleware(testJWTConfig.Secret))
	{
		auth.POST("/plugins/:id/install", handler.Install)
		auth.DELETE("/plugins/:id/install", handler.Uninstall)
		auth.PATCH("/plugins/:id/install", handler.ToggleActive)
		auth.GET("/me/installations", handler.GetMyInstallations)
	}
	return r
}

func TestShouldReturn201WhenInstallPlugin(t *testing.T) {
	r := setupInstallationRouter()
	token, _ := mw.GenerateAccessToken(2, "user@test.com", "user", testJWTConfig.Secret, 3600)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/plugins/1/install", bytes.NewBufferString("{}"))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	r.ServeHTTP(w, req)

	assert.Equal(t, 201, w.Code)
}

func TestShouldReturn204WhenUninstallPlugin(t *testing.T) {
	r := setupInstallationRouter()
	token, _ := mw.GenerateAccessToken(2, "user@test.com", "user", testJWTConfig.Secret, 3600)

	// 설치
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/plugins/1/install", bytes.NewBufferString("{}"))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	r.ServeHTTP(w, req)

	// 삭제
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("DELETE", "/api/v1/plugins/1/install", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	r.ServeHTTP(w, req)

	assert.Equal(t, 204, w.Code)
}

func TestShouldToggleActiveSuccessfullyViaAPI(t *testing.T) {
	r := setupInstallationRouter()
	token, _ := mw.GenerateAccessToken(2, "user@test.com", "user", testJWTConfig.Secret, 3600)

	// 설치
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/plugins/1/install", bytes.NewBufferString("{}"))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	r.ServeHTTP(w, req)

	// 비활성화
	body, _ := json.Marshal(dto.ToggleActiveRequest{IsActive: false})
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("PATCH", "/api/v1/plugins/1/install", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}

func TestShouldReturnInstallationListWhenGetMyInstallations(t *testing.T) {
	r := setupInstallationRouter()
	token, _ := mw.GenerateAccessToken(2, "user@test.com", "user", testJWTConfig.Secret, 3600)

	// 설치
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/plugins/1/install", bytes.NewBufferString("{}"))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	r.ServeHTTP(w, req)

	// 조회
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/v1/me/installations", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), `"data"`)
}
