package repository

import (
	"github.com/ktc-plugin-hub/backend/internal/model"
	"gorm.io/gorm"
)

// PluginRepository 는 플러그인 데이터 접근 인터페이스
type PluginRepository interface {
	Create(plugin *model.Plugin) error
	FindByID(id uint) (*model.Plugin, error)
	FindAll(filter PluginFilter) ([]model.Plugin, int64, error)
	Update(plugin *model.Plugin) error
	Delete(id uint) error
	FindByName(name string) (*model.Plugin, error)
	IncrementDownloadCount(id uint) error
	UpdateRating(id uint, avgRating float64, reviewCount int) error
}

// PluginFilter 는 플러그인 목록 필터링 조건
type PluginFilter struct {
	CategoryID *uint
	Keyword    string
	Sort       string // popular, latest, rating
	Status     string
	AuthorID   *uint
	Offset     int
	Limit      int
}

// pluginRepository 는 PluginRepository의 GORM 구현체
type pluginRepository struct {
	db *gorm.DB
}

// NewPluginRepository 는 PluginRepository 인스턴스를 생성한다
func NewPluginRepository(db *gorm.DB) PluginRepository {
	return &pluginRepository{db: db}
}

func (r *pluginRepository) Create(plugin *model.Plugin) error {
	return r.db.Create(plugin).Error
}

func (r *pluginRepository) FindByID(id uint) (*model.Plugin, error) {
	var plugin model.Plugin
	err := r.db.Preload("Author").Preload("Category").First(&plugin, id).Error
	if err != nil {
		return nil, err
	}
	return &plugin, nil
}

func (r *pluginRepository) FindAll(filter PluginFilter) ([]model.Plugin, int64, error) {
	var plugins []model.Plugin
	var total int64

	query := r.db.Model(&model.Plugin{}).Preload("Author").Preload("Category")

	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}
	if filter.CategoryID != nil {
		query = query.Where("category_id = ?", *filter.CategoryID)
	}
	if filter.Keyword != "" {
		keyword := "%" + filter.Keyword + "%"
		query = query.Where("name ILIKE ? OR description ILIKE ?", keyword, keyword)
	}
	if filter.AuthorID != nil {
		query = query.Where("author_id = ?", *filter.AuthorID)
	}

	query.Count(&total)

	switch filter.Sort {
	case "popular":
		query = query.Order("download_count DESC")
	case "rating":
		query = query.Order("avg_rating DESC")
	default:
		query = query.Order("created_at DESC")
	}

	err := query.Offset(filter.Offset).Limit(filter.Limit).Find(&plugins).Error
	return plugins, total, err
}

func (r *pluginRepository) Update(plugin *model.Plugin) error {
	return r.db.Save(plugin).Error
}

func (r *pluginRepository) Delete(id uint) error {
	return r.db.Delete(&model.Plugin{}, id).Error
}

func (r *pluginRepository) FindByName(name string) (*model.Plugin, error) {
	var plugin model.Plugin
	err := r.db.Where("name = ?", name).First(&plugin).Error
	if err != nil {
		return nil, err
	}
	return &plugin, nil
}

func (r *pluginRepository) IncrementDownloadCount(id uint) error {
	return r.db.Model(&model.Plugin{}).Where("id = ?", id).
		UpdateColumn("download_count", gorm.Expr("download_count + 1")).Error
}

func (r *pluginRepository) UpdateRating(id uint, avgRating float64, reviewCount int) error {
	return r.db.Model(&model.Plugin{}).Where("id = ?", id).
		Updates(map[string]interface{}{
			"avg_rating":   avgRating,
			"review_count": reviewCount,
		}).Error
}
