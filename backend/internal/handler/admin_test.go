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

func setupAdminRouter() *gin.Engine {
	pluginRepo := repository.NewMockPluginRepository()
	pluginRepo.Create(&model.Plugin{Name: "pending-p", AuthorID: 2, CategoryID: 1, Description: "d", Status: model.PluginStatusPending})

	adminSvc := service.NewAdminService(pluginRepo)
	handler := NewAdminHandler(adminSvc)

	r := gin.New()
	admin := r.Group("/api/v1/admin", mw.AuthMiddleware(testJWTConfig.Secret), mw.AdminMiddleware())
	{
		admin.GET("/plugins/pending", handler.GetPendingPlugins)
		admin.PATCH("/plugins/:id/approve", handler.ApprovePlugin)
		admin.PATCH("/plugins/:id/reject", handler.RejectPlugin)
		admin.PATCH("/plugins/:id/hide", handler.HidePlugin)
	}
	return r
}

func TestShouldReturnPendingPluginsWhenAdmin(t *testing.T) {
	r := setupAdminRouter()
	token, _ := mw.GenerateAccessToken(1, "admin@test.com", "admin", testJWTConfig.Secret, 3600)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/admin/plugins/pending", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "pending-p")
}

func TestShouldApprovePluginWhenAdmin(t *testing.T) {
	r := setupAdminRouter()
	token, _ := mw.GenerateAccessToken(1, "admin@test.com", "admin", testJWTConfig.Secret, 3600)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PATCH", "/api/v1/admin/plugins/1/approve", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "approved")
}

func TestShouldRejectPluginWhenAdmin(t *testing.T) {
	r := setupAdminRouter()
	token, _ := mw.GenerateAccessToken(1, "admin@test.com", "admin", testJWTConfig.Secret, 3600)

	body, _ := json.Marshal(dto.RejectRequest{Reason: "부적절한 내용"})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PATCH", "/api/v1/admin/plugins/1/reject", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "rejected")
}

func TestShouldReturn403WhenNonAdminAccessesAdminAPI(t *testing.T) {
	r := setupAdminRouter()
	token, _ := mw.GenerateAccessToken(1, "user@test.com", "user", testJWTConfig.Secret, 3600)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/admin/plugins/pending", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	r.ServeHTTP(w, req)

	assert.Equal(t, 403, w.Code)
}
