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

func setupReviewRouter() *gin.Engine {
	reviewRepo := repository.NewMockReviewRepository()
	pluginRepo := repository.NewMockPluginRepository()
	pluginRepo.Create(&model.Plugin{Name: "test-plugin", AuthorID: 1, CategoryID: 1, Description: "d", Status: "approved"})

	reviewSvc := service.NewReviewService(reviewRepo, pluginRepo)
	handler := NewReviewHandler(reviewSvc)

	r := gin.New()
	r.GET("/api/v1/plugins/:id/reviews", handler.GetReviews)

	auth := r.Group("/api/v1", mw.AuthMiddleware(testJWTConfig.Secret))
	{
		auth.POST("/plugins/:id/reviews", handler.CreateReview)
		auth.PUT("/plugins/:id/reviews/:reviewId", handler.UpdateReview)
		auth.DELETE("/plugins/:id/reviews/:reviewId", handler.DeleteReview)
	}
	return r
}

func TestShouldReturn201WhenCreateReview(t *testing.T) {
	r := setupReviewRouter()
	token, _ := mw.GenerateAccessToken(2, "user@test.com", "user", testJWTConfig.Secret, 3600)

	body, _ := json.Marshal(dto.CreateReviewRequest{Rating: 5, Content: "awesome"})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/plugins/1/reviews", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	r.ServeHTTP(w, req)

	assert.Equal(t, 201, w.Code)
}

func TestShouldReturn403WhenSelfReview(t *testing.T) {
	r := setupReviewRouter()
	// AuthorID=1인 플러그인에 userID=1이 리뷰
	token, _ := mw.GenerateAccessToken(1, "author@test.com", "user", testJWTConfig.Secret, 3600)

	body, _ := json.Marshal(dto.CreateReviewRequest{Rating: 5, Content: "self review"})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/plugins/1/reviews", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	r.ServeHTTP(w, req)

	assert.Equal(t, 403, w.Code)
}

func TestShouldReturn409WhenDuplicateReview(t *testing.T) {
	r := setupReviewRouter()
	token, _ := mw.GenerateAccessToken(2, "user@test.com", "user", testJWTConfig.Secret, 3600)

	body, _ := json.Marshal(dto.CreateReviewRequest{Rating: 4, Content: "first"})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/plugins/1/reviews", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	r.ServeHTTP(w, req)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/v1/plugins/1/reviews", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	r.ServeHTTP(w, req)

	assert.Equal(t, 409, w.Code)
}

func TestShouldReturnReviewListWhenGetReviews(t *testing.T) {
	r := setupReviewRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/plugins/1/reviews", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}
