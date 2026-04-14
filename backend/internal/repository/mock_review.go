package repository

import (
	"github.com/ktc-plugin-hub/backend/internal/model"
)

type MockReviewRepository struct {
	reviews map[uint]*model.Review
	nextID  uint
}

func NewMockReviewRepository() *MockReviewRepository {
	return &MockReviewRepository{
		reviews: make(map[uint]*model.Review),
		nextID:  1,
	}
}

func (r *MockReviewRepository) Create(review *model.Review) error {
	// 1인 1리뷰 검증
	for _, rv := range r.reviews {
		if rv.UserID == review.UserID && rv.PluginID == review.PluginID {
			return ErrDuplicateKey
		}
	}
	review.ID = r.nextID
	r.nextID++
	r.reviews[review.ID] = review
	return nil
}

func (r *MockReviewRepository) FindByPluginID(pluginID uint, offset, limit int) ([]model.Review, int64, error) {
	var result []model.Review
	for _, rv := range r.reviews {
		if rv.PluginID == pluginID {
			result = append(result, *rv)
		}
	}
	total := int64(len(result))
	start := offset
	if start > len(result) {
		start = len(result)
	}
	end := start + limit
	if end > len(result) {
		end = len(result)
	}
	return result[start:end], total, nil
}

func (r *MockReviewRepository) FindByUserAndPlugin(userID uint, pluginID uint) (*model.Review, error) {
	for _, rv := range r.reviews {
		if rv.UserID == userID && rv.PluginID == pluginID {
			return rv, nil
		}
	}
	return nil, ErrNotFound
}

func (r *MockReviewRepository) Update(review *model.Review) error {
	if _, ok := r.reviews[review.ID]; !ok {
		return ErrNotFound
	}
	r.reviews[review.ID] = review
	return nil
}

func (r *MockReviewRepository) Delete(id uint) error {
	if _, ok := r.reviews[id]; !ok {
		return ErrNotFound
	}
	delete(r.reviews, id)
	return nil
}

func (r *MockReviewRepository) FindByID(id uint) (*model.Review, error) {
	rv, ok := r.reviews[id]
	if !ok {
		return nil, ErrNotFound
	}
	return rv, nil
}

func (r *MockReviewRepository) CalculateAvgRating(pluginID uint) (float64, int, error) {
	var total float64
	var count int
	for _, rv := range r.reviews {
		if rv.PluginID == pluginID {
			total += float64(rv.Rating)
			count++
		}
	}
	if count == 0 {
		return 0, 0, nil
	}
	return total / float64(count), count, nil
}
