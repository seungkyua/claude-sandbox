package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ktc-plugin-hub/backend/internal/dto"
	"github.com/ktc-plugin-hub/backend/internal/repository"
)

// CategoryHandler 는 카테고리 관련 HTTP 핸들러
type CategoryHandler struct {
	categoryRepo repository.CategoryRepository
}

// NewCategoryHandler 는 CategoryHandler를 생성한다
func NewCategoryHandler(categoryRepo repository.CategoryRepository) *CategoryHandler {
	return &CategoryHandler{categoryRepo: categoryRepo}
}

// GetAll 은 모든 카테고리를 반환한다
// GET /api/v1/categories
func (h *CategoryHandler) GetAll(c *gin.Context) {
	categories, err := h.categoryRepo.FindAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewErrorResponse(
			"INTERNAL_ERROR", "서버 오류", http.StatusInternalServerError, err.Error(),
		))
		return
	}

	// DTO 변환
	var result []dto.CategoryResponse
	for _, cat := range categories {
		result = append(result, dto.CategoryResponse{
			ID:   cat.ID,
			Name: cat.Name,
		})
	}

	c.JSON(http.StatusOK, gin.H{"data": result})
}
