package dto

import "time"

// CreateReviewRequest 는 리뷰 작성 요청 DTO
type CreateReviewRequest struct {
	Rating  int    `json:"rating" binding:"required,min=1,max=5"`
	Content string `json:"content" binding:"required,min=1,max=1000"`
}

// UpdateReviewRequest 는 리뷰 수정 요청 DTO
type UpdateReviewRequest struct {
	Rating  *int    `json:"rating" binding:"omitempty,min=1,max=5"`
	Content *string `json:"content" binding:"omitempty,min=1,max=1000"`
}

// ReviewResponse 는 리뷰 응답 DTO
type ReviewResponse struct {
	ID        uint           `json:"id"`
	User      AuthorResponse `json:"user"`
	Rating    int            `json:"rating"`
	Content   string         `json:"content"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}

// RejectRequest 는 플러그인 반려 요청 DTO
type RejectRequest struct {
	Reason string `json:"reason" binding:"required"`
}
