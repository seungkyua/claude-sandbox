package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ktc-plugin-hub/backend/internal/dto"
	mw "github.com/ktc-plugin-hub/backend/internal/middleware"
	"github.com/ktc-plugin-hub/backend/internal/service"
)

// ReviewHandler 는 리뷰 관련 HTTP 핸들러
type ReviewHandler struct {
	reviewService service.ReviewService
}

func NewReviewHandler(reviewService service.ReviewService) *ReviewHandler {
	return &ReviewHandler{reviewService: reviewService}
}

// CreateReview 는 리뷰를 작성한다
// POST /api/v1/plugins/:id/reviews
func (h *ReviewHandler) CreateReview(c *gin.Context) {
	pluginID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse("INVALID_INPUT", "잘못된 ID", http.StatusBadRequest, "유효한 플러그인 ID를 입력하세요"))
		return
	}

	var req dto.CreateReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse("INVALID_INPUT", "입력 오류", http.StatusBadRequest, err.Error()))
		return
	}

	userID, _ := mw.GetUserIDFromContext(c)

	resp, err := h.reviewService.CreateReview(uint(pluginID), userID, &req)
	if err != nil {
		switch err {
		case service.ErrSelfReview:
			c.JSON(http.StatusForbidden, dto.NewErrorResponse("SELF_REVIEW", "본인 리뷰 불가", http.StatusForbidden, err.Error()))
		case service.ErrDuplicateReview:
			c.JSON(http.StatusConflict, dto.NewErrorResponse("DUPLICATE_REVIEW", "중복 리뷰", http.StatusConflict, err.Error()))
		default:
			c.JSON(http.StatusInternalServerError, dto.NewErrorResponse("INTERNAL_ERROR", "서버 오류", http.StatusInternalServerError, err.Error()))
		}
		return
	}

	c.JSON(http.StatusCreated, resp)
}

// GetReviews 는 리뷰 목록을 반환한다
// GET /api/v1/plugins/:id/reviews
func (h *ReviewHandler) GetReviews(c *gin.Context) {
	pluginID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse("INVALID_INPUT", "잘못된 ID", http.StatusBadRequest, "유효한 플러그인 ID를 입력하세요"))
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))

	resp, err := h.reviewService.GetReviewsByPluginID(uint(pluginID), page, size)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewErrorResponse("INTERNAL_ERROR", "서버 오류", http.StatusInternalServerError, err.Error()))
		return
	}

	c.JSON(http.StatusOK, resp)
}

// UpdateReview 는 리뷰를 수정한다
// PUT /api/v1/plugins/:id/reviews/:reviewId
func (h *ReviewHandler) UpdateReview(c *gin.Context) {
	pluginID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	reviewID, _ := strconv.ParseUint(c.Param("reviewId"), 10, 64)

	var req dto.UpdateReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse("INVALID_INPUT", "입력 오류", http.StatusBadRequest, err.Error()))
		return
	}

	userID, _ := mw.GetUserIDFromContext(c)
	resp, err := h.reviewService.UpdateReview(uint(pluginID), uint(reviewID), userID, &req)
	if err != nil {
		switch err {
		case service.ErrForbidden:
			c.JSON(http.StatusForbidden, dto.NewErrorResponse("FORBIDDEN", "권한 없음", http.StatusForbidden, err.Error()))
		default:
			c.JSON(http.StatusInternalServerError, dto.NewErrorResponse("INTERNAL_ERROR", "서버 오류", http.StatusInternalServerError, err.Error()))
		}
		return
	}

	c.JSON(http.StatusOK, resp)
}

// DeleteReview 는 리뷰를 삭제한다
// DELETE /api/v1/plugins/:id/reviews/:reviewId
func (h *ReviewHandler) DeleteReview(c *gin.Context) {
	pluginID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	reviewID, _ := strconv.ParseUint(c.Param("reviewId"), 10, 64)

	userID, _ := mw.GetUserIDFromContext(c)
	role, _ := mw.GetUserRoleFromContext(c)

	err := h.reviewService.DeleteReview(uint(pluginID), uint(reviewID), userID, role)
	if err != nil {
		switch err {
		case service.ErrForbidden:
			c.JSON(http.StatusForbidden, dto.NewErrorResponse("FORBIDDEN", "권한 없음", http.StatusForbidden, err.Error()))
		default:
			c.JSON(http.StatusNotFound, dto.NewErrorResponse("NOT_FOUND", "리뷰 없음", http.StatusNotFound, err.Error()))
		}
		return
	}

	c.Status(http.StatusNoContent)
}
