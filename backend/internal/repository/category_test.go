package repository

import (
	"testing"

	"github.com/ktc-plugin-hub/backend/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// 카테고리 생성 및 조회가 올바르게 동작하는지 확인
func TestShouldCreateCategoryAndFindByID(t *testing.T) {
	repo := NewMockCategoryRepository()
	cat := &model.Category{Name: "개발 도구", Description: "개발 관련 플러그인", SortOrder: 1}

	err := repo.Create(cat)
	require.NoError(t, err)
	assert.NotZero(t, cat.ID)

	found, err := repo.FindByID(cat.ID)
	require.NoError(t, err)
	assert.Equal(t, "개발 도구", found.Name)
}

// 모든 카테고리를 정렬 순서대로 조회하는지 확인
func TestShouldFindAllCategoriesOrderBySortOrder(t *testing.T) {
	repo := NewMockCategoryRepository()
	repo.Create(&model.Category{Name: "유틸리티", SortOrder: 3})
	repo.Create(&model.Category{Name: "개발 도구", SortOrder: 1})
	repo.Create(&model.Category{Name: "디자인", SortOrder: 2})

	categories, err := repo.FindAll()
	require.NoError(t, err)
	assert.Len(t, categories, 3)
	assert.Equal(t, "개발 도구", categories[0].Name)
	assert.Equal(t, "디자인", categories[1].Name)
	assert.Equal(t, "유틸리티", categories[2].Name)
}

// 존재하지 않는 카테고리 조회 시 에러를 반환하는지 확인
func TestShouldReturnErrorWhenCategoryNotFound(t *testing.T) {
	repo := NewMockCategoryRepository()
	_, err := repo.FindByID(999)
	assert.Error(t, err)
}

// MockCategoryRepository가 인터페이스를 구현하는지 확인
func TestMockCategoryRepositoryShouldImplementInterface(t *testing.T) {
	var _ CategoryRepository = NewMockCategoryRepository()
}
