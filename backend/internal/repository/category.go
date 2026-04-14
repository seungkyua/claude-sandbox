package repository

import (
	"github.com/ktc-plugin-hub/backend/internal/model"
	"gorm.io/gorm"
)

// CategoryRepository 는 카테고리 데이터 접근 인터페이스
type CategoryRepository interface {
	FindAll() ([]model.Category, error)
	FindByID(id uint) (*model.Category, error)
	Create(category *model.Category) error
}

// categoryRepository 는 CategoryRepository의 GORM 구현체
type categoryRepository struct {
	db *gorm.DB
}

// NewCategoryRepository 는 CategoryRepository 인스턴스를 생성한다
func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &categoryRepository{db: db}
}

// FindAll 은 모든 카테고리를 정렬 순서대로 조회한다
func (r *categoryRepository) FindAll() ([]model.Category, error) {
	var categories []model.Category
	err := r.db.Order("sort_order ASC").Find(&categories).Error
	return categories, err
}

// FindByID 는 ID로 카테고리를 조회한다
func (r *categoryRepository) FindByID(id uint) (*model.Category, error) {
	var category model.Category
	err := r.db.First(&category, id).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

// Create 는 새 카테고리를 생성한다
func (r *categoryRepository) Create(category *model.Category) error {
	return r.db.Create(category).Error
}
