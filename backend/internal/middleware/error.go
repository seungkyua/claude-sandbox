package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ktc-plugin-hub/backend/internal/dto"
)

// JSONErrorHandler 는 모든 응답을 JSON 형식으로 반환하는 에러 핸들러
func JSONErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// 에러가 있으면 JSON 형식으로 응답
		if len(c.Errors) > 0 {
			err := c.Errors.Last()
			c.JSON(c.Writer.Status(), dto.NewErrorResponse(
				"INTERNAL_ERROR",
				"서버 내부 오류",
				http.StatusInternalServerError,
				err.Error(),
			))
		}
	}
}

// RespondError 는 RFC 7807 형식의 에러 응답을 보내는 헬퍼 함수
func RespondError(c *gin.Context, status int, errType string, title string, detail string) {
	c.JSON(status, dto.NewErrorResponse(errType, title, status, detail))
}
