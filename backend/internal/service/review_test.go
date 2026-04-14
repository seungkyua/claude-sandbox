package service

import (
	"testing"

	"github.com/ktc-plugin-hub/backend/internal/dto"
	"github.com/ktc-plugin-hub/backend/internal/model"
	"github.com/ktc-plugin-hub/backend/internal/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestReviewService() (ReviewService, *repository.MockPluginRepository) {
	reviewRepo := repository.NewMockReviewRepository()
	pluginRepo := repository.NewMockPluginRepository()
	pluginRepo.Create(&model.Plugin{Name: "p1", AuthorID: 1, CategoryID: 1, Description: "d", Status: "approved"})
	return NewReviewService(reviewRepo, pluginRepo), pluginRepo
}

func TestShouldCreateReviewSuccessfully(t *testing.T) {
	svc, _ := newTestReviewService()
	resp, err := svc.CreateReview(1, 2, &dto.CreateReviewRequest{Rating: 5, Content: "great!"})
	require.NoError(t, err)
	assert.Equal(t, 5, resp.Rating)
}

func TestShouldReturnErrorWhenSelfReview(t *testing.T) {
	svc, _ := newTestReviewService()
	_, err := svc.CreateReview(1, 1, &dto.CreateReviewRequest{Rating: 5, Content: "self"})
	assert.Equal(t, ErrSelfReview, err)
}

func TestShouldReturnErrorWhenDuplicateReview(t *testing.T) {
	svc, _ := newTestReviewService()
	svc.CreateReview(1, 2, &dto.CreateReviewRequest{Rating: 4, Content: "first"})
	_, err := svc.CreateReview(1, 2, &dto.CreateReviewRequest{Rating: 5, Content: "second"})
	assert.Equal(t, ErrDuplicateReview, err)
}

func TestShouldRecalculateRatingAfterReview(t *testing.T) {
	svc, pluginRepo := newTestReviewService()
	svc.CreateReview(1, 2, &dto.CreateReviewRequest{Rating: 4, Content: "good"})
	svc.CreateReview(1, 3, &dto.CreateReviewRequest{Rating: 2, Content: "meh"})

	plugin, _ := pluginRepo.FindByID(1)
	assert.Equal(t, 3.0, plugin.AvgRating)
	assert.Equal(t, 2, plugin.ReviewCount)
}

func TestShouldUpdateReviewByOwner(t *testing.T) {
	svc, _ := newTestReviewService()
	created, _ := svc.CreateReview(1, 2, &dto.CreateReviewRequest{Rating: 3, Content: "ok"})

	newRating := 5
	updated, err := svc.UpdateReview(1, created.ID, 2, &dto.UpdateReviewRequest{Rating: &newRating})
	require.NoError(t, err)
	assert.Equal(t, 5, updated.Rating)
}

func TestShouldDeleteReviewByOwner(t *testing.T) {
	svc, _ := newTestReviewService()
	created, _ := svc.CreateReview(1, 2, &dto.CreateReviewRequest{Rating: 3, Content: "ok"})

	err := svc.DeleteReview(1, created.ID, 2, "user")
	require.NoError(t, err)
}

func TestShouldReturnForbiddenWhenNonOwnerDeletesReview(t *testing.T) {
	svc, _ := newTestReviewService()
	created, _ := svc.CreateReview(1, 2, &dto.CreateReviewRequest{Rating: 3, Content: "ok"})

	err := svc.DeleteReview(1, created.ID, 99, "user")
	assert.Equal(t, ErrForbidden, err)
}
