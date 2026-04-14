package repository

import (
	"github.com/ktc-plugin-hub/backend/internal/model"
	"gorm.io/gorm"
)

// ReviewRepository 는 리뷰 데이터 접근 인터페이스
type ReviewRepository interface {
	Create(review *model.Review) error
	FindByPluginID(pluginID uint, offset, limit int) ([]model.Review, int64, error)
	FindByUserAndPlugin(userID uint, pluginID uint) (*model.Review, error)
	Update(review *model.Review) error
	Delete(id uint) error
	FindByID(id uint) (*model.Review, error)
	CalculateAvgRating(pluginID uint) (float64, int, error)
}

type reviewRepository struct {
	db *gorm.DB
}

func NewReviewRepository(db *gorm.DB) ReviewRepository {
	return &reviewRepository{db: db}
}

func (r *reviewRepository) Create(review *model.Review) error {
	return r.db.Create(review).Error
}

func (r *reviewRepository) FindByPluginID(pluginID uint, offset, limit int) ([]model.Review, int64, error) {
	var reviews []model.Review
	var total int64
	r.db.Model(&model.Review{}).Where("plugin_id = ?", pluginID).Count(&total)
	err := r.db.Where("plugin_id = ?", pluginID).Preload("User").
		Order("created_at DESC").Offset(offset).Limit(limit).Find(&reviews).Error
	return reviews, total, err
}

func (r *reviewRepository) FindByUserAndPlugin(userID uint, pluginID uint) (*model.Review, error) {
	var review model.Review
	err := r.db.Where("user_id = ? AND plugin_id = ?", userID, pluginID).First(&review).Error
	if err != nil {
		return nil, err
	}
	return &review, nil
}

func (r *reviewRepository) Update(review *model.Review) error {
	return r.db.Save(review).Error
}

func (r *reviewRepository) Delete(id uint) error {
	return r.db.Delete(&model.Review{}, id).Error
}

func (r *reviewRepository) FindByID(id uint) (*model.Review, error) {
	var review model.Review
	err := r.db.Preload("User").First(&review, id).Error
	if err != nil {
		return nil, err
	}
	return &review, nil
}

func (r *reviewRepository) CalculateAvgRating(pluginID uint) (float64, int, error) {
	var result struct {
		Avg   float64
		Count int
	}
	err := r.db.Model(&model.Review{}).Where("plugin_id = ?", pluginID).
		Select("COALESCE(AVG(rating), 0) as avg, COUNT(*) as count").Scan(&result).Error
	return result.Avg, result.Count, err
}
