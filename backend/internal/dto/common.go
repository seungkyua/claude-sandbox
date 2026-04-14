package dto

// PaginationRequest 는 페이지네이션 요청 파라미터
type PaginationRequest struct {
	Page int `form:"page" json:"page"`
	Size int `form:"size" json:"size"`
}

// Normalize 는 페이지네이션 값을 정규화한다
func (p *PaginationRequest) Normalize() {
	if p.Page < 1 {
		p.Page = 1
	}
	if p.Size < 1 {
		p.Size = 20
	}
	if p.Size > 50 {
		p.Size = 50
	}
}

// Offset 은 데이터베이스 쿼리용 오프셋을 반환한다
func (p *PaginationRequest) Offset() int {
	return (p.Page - 1) * p.Size
}

// PaginatedResponse 는 페이지네이션된 응답 구조체
type PaginatedResponse struct {
	Data  interface{} `json:"data"`
	Total int64       `json:"total"`
	Page  int         `json:"page"`
	Size  int         `json:"size"`
}

// ErrorResponse 는 RFC 7807 에러 응답 포맷
type ErrorResponse struct {
	Type   string `json:"type"`
	Title  string `json:"title"`
	Status int    `json:"status"`
	Detail string `json:"detail"`
}

// NewErrorResponse 는 에러 응답을 생성한다
func NewErrorResponse(errType string, title string, status int, detail string) *ErrorResponse {
	return &ErrorResponse{
		Type:   errType,
		Title:  title,
		Status: status,
		Detail: detail,
	}
}
