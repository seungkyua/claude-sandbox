package repository

import (
	"sort"

	"github.com/ktc-plugin-hub/backend/internal/model"
)

// MockCategoryRepository 는 테스트용 카테고리 리포지토리 모킹
type MockCategoryRepository struct {
	categories map[uint]*model.Category
	nextID     uint
}

// NewMockCategoryRepository 는 MockCategoryRepository를 생성한다
func NewMockCategoryRepository() *MockCategoryRepository {
	return &MockCategoryRepository{
		categories: make(map[uint]*model.Category),
		nextID:     1,
	}
}

// FindAll 은 모든 카테고리를 정렬 순서대로 반환한다
func (r *MockCategoryRepository) FindAll() ([]model.Category, error) {
	var categories []model.Category
	for _, c := range r.categories {
		categories = append(categories, *c)
	}
	sort.Slice(categories, func(i, j int) bool {
		return categories[i].SortOrder < categories[j].SortOrder
	})
	return categories, nil
}

// FindByID 는 ID로 카테고리를 조회한다
func (r *MockCategoryRepository) FindByID(id uint) (*model.Category, error) {
	c, ok := r.categories[id]
	if !ok {
		return nil, ErrNotFound
	}
	return c, nil
}

// Create 는 새 카테고리를 생성한다
func (r *MockCategoryRepository) Create(category *model.Category) error {
	for _, c := range r.categories {
		if c.Name == category.Name {
			return ErrDuplicateKey
		}
	}
	category.ID = r.nextID
	r.nextID++
	r.categories[category.ID] = category
	return nil
}
