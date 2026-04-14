package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/ktc-plugin-hub/backend/internal/model"
	"github.com/ktc-plugin-hub/backend/internal/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupCategoryRouter() *gin.Engine {
	catRepo := repository.NewMockCategoryRepository()
	catRepo.Create(&model.Category{Name: "개발 도구", Description: "개발 도구", SortOrder: 1})
	catRepo.Create(&model.Category{Name: "디자인", Description: "디자인 도구", SortOrder: 2})

	handler := NewCategoryHandler(catRepo)
	r := gin.New()
	r.GET("/api/v1/categories", handler.GetAll)
	return r
}

// GET /categories 가 카테고리 목록을 반환하는지 확인
func TestShouldReturnCategoryListWhenGetCategories(t *testing.T) {
	r := setupCategoryRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/categories", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	var resp map[string][]map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Len(t, resp["data"], 2)
}

// 카테고리가 없을 때 빈 목록을 반환하는지 확인
func TestShouldReturnEmptyListWhenNoCategories(t *testing.T) {
	catRepo := repository.NewMockCategoryRepository()
	handler := NewCategoryHandler(catRepo)
	r := gin.New()
	r.GET("/api/v1/categories", handler.GetAll)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/categories", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	// data가 nil이거나 빈 배열
	assert.Contains(t, w.Body.String(), `"data"`)
}
