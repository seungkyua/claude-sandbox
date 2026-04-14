package service

import (
	"errors"
	"math"

	"github.com/ktc-plugin-hub/backend/internal/dto"
	"github.com/ktc-plugin-hub/backend/internal/model"
	"github.com/ktc-plugin-hub/backend/internal/repository"
)

var (
	ErrSelfReview      = errors.New("본인 플러그인에 리뷰를 작성할 수 없습니다")
	ErrDuplicateReview = errors.New("이미 리뷰를 작성했습니다")
)

// ReviewService 는 리뷰 비즈니스 로직 인터페이스
type ReviewService interface {
	CreateReview(pluginID uint, userID uint, req *dto.CreateReviewRequest) (*dto.ReviewResponse, error)
	GetReviewsByPluginID(pluginID uint, page, size int) (*dto.PaginatedResponse, error)
	UpdateReview(pluginID uint, reviewID uint, userID uint, req *dto.UpdateReviewRequest) (*dto.ReviewResponse, error)
	DeleteReview(pluginID uint, reviewID uint, userID uint, role string) error
}

type reviewService struct {
	reviewRepo repository.ReviewRepository
	pluginRepo repository.PluginRepository
}

func NewReviewService(reviewRepo repository.ReviewRepository, pluginRepo repository.PluginRepository) ReviewService {
	return &reviewService{reviewRepo: reviewRepo, pluginRepo: pluginRepo}
}

func (s *reviewService) CreateReview(pluginID uint, userID uint, req *dto.CreateReviewRequest) (*dto.ReviewResponse, error) {
	// 플러그인 존재 확인
	plugin, err := s.pluginRepo.FindByID(pluginID)
	if err != nil {
		return nil, ErrNotFound
	}

	// 본인 플러그인 리뷰 불가
	if plugin.AuthorID == userID {
		return nil, ErrSelfReview
	}

	// 중복 리뷰 확인
	_, err = s.reviewRepo.FindByUserAndPlugin(userID, pluginID)
	if err == nil {
		return nil, ErrDuplicateReview
	}

	review := &model.Review{
		PluginID: pluginID,
		UserID:   userID,
		Rating:   req.Rating,
		Content:  req.Content,
	}

	if err := s.reviewRepo.Create(review); err != nil {
		return nil, err
	}

	// 평균 평점 재계산
	s.recalculateRating(pluginID)

	return &dto.ReviewResponse{
		ID:        review.ID,
		User:      dto.AuthorResponse{ID: userID},
		Rating:    review.Rating,
		Content:   review.Content,
		CreatedAt: review.CreatedAt,
		UpdatedAt: review.UpdatedAt,
	}, nil
}

func (s *reviewService) GetReviewsByPluginID(pluginID uint, page, size int) (*dto.PaginatedResponse, error) {
	if page < 1 {
		page = 1
	}
	if size < 1 || size > 50 {
		size = 20
	}
	offset := (page - 1) * size

	reviews, total, err := s.reviewRepo.FindByPluginID(pluginID, offset, size)
	if err != nil {
		return nil, err
	}

	var data []dto.ReviewResponse
	for _, rv := range reviews {
		data = append(data, dto.ReviewResponse{
			ID:        rv.ID,
			User:      dto.AuthorResponse{ID: rv.UserID, Nickname: rv.User.Nickname},
			Rating:    rv.Rating,
			Content:   rv.Content,
			CreatedAt: rv.CreatedAt,
			UpdatedAt: rv.UpdatedAt,
		})
	}

	return &dto.PaginatedResponse{Data: data, Total: total, Page: page, Size: size}, nil
}

func (s *reviewService) UpdateReview(pluginID uint, reviewID uint, userID uint, req *dto.UpdateReviewRequest) (*dto.ReviewResponse, error) {
	review, err := s.reviewRepo.FindByID(reviewID)
	if err != nil {
		return nil, ErrNotFound
	}

	if review.UserID != userID {
		return nil, ErrForbidden
	}

	if req.Rating != nil {
		review.Rating = *req.Rating
	}
	if req.Content != nil {
		review.Content = *req.Content
	}

	if err := s.reviewRepo.Update(review); err != nil {
		return nil, err
	}

	s.recalculateRating(pluginID)

	return &dto.ReviewResponse{
		ID: review.ID, User: dto.AuthorResponse{ID: userID},
		Rating: review.Rating, Content: review.Content,
		CreatedAt: review.CreatedAt, UpdatedAt: review.UpdatedAt,
	}, nil
}

func (s *reviewService) DeleteReview(pluginID uint, reviewID uint, userID uint, role string) error {
	review, err := s.reviewRepo.FindByID(reviewID)
	if err != nil {
		return ErrNotFound
	}

	if review.UserID != userID && role != "admin" {
		return ErrForbidden
	}

	if err := s.reviewRepo.Delete(reviewID); err != nil {
		return err
	}

	s.recalculateRating(pluginID)
	return nil
}

func (s *reviewService) recalculateRating(pluginID uint) {
	avg, count, err := s.reviewRepo.CalculateAvgRating(pluginID)
	if err != nil {
		return
	}
	// 소수점 1자리 반올림
	avg = math.Round(avg*10) / 10
	s.pluginRepo.UpdateRating(pluginID, avg, count)
}
