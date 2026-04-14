package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/ktc-plugin-hub/backend/internal/dto"
)

// 컨텍스트에 저장되는 키
const (
	ContextUserID = "userID"
	ContextEmail  = "email"
	ContextRole   = "role"
)

// AuthMiddleware 는 Bearer Token을 검증하고 사용자 정보를 컨텍스트에 주입하는 미들웨어
func AuthMiddleware(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, dto.NewErrorResponse(
				"UNAUTHORIZED", "인증 필요", http.StatusUnauthorized, "Authorization 헤더가 없습니다",
			))
			c.Abort()
			return
		}

		// Bearer 토큰 추출
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, dto.NewErrorResponse(
				"UNAUTHORIZED", "인증 형식 오류", http.StatusUnauthorized, "Bearer 토큰 형식이 아닙니다",
			))
			c.Abort()
			return
		}

		// 토큰 검증
		claims, err := ValidateToken(parts[1], jwtSecret)
		if err != nil {
			c.JSON(http.StatusUnauthorized, dto.NewErrorResponse(
				"UNAUTHORIZED", "토큰 검증 실패", http.StatusUnauthorized, "유효하지 않거나 만료된 토큰입니다",
			))
			c.Abort()
			return
		}

		// Access Token만 허용
		if claims.TokenType != "access" {
			c.JSON(http.StatusUnauthorized, dto.NewErrorResponse(
				"UNAUTHORIZED", "토큰 타입 오류", http.StatusUnauthorized, "Access Token이 아닙니다",
			))
			c.Abort()
			return
		}

		// 사용자 정보를 컨텍스트에 주입
		c.Set(ContextUserID, claims.UserID)
		c.Set(ContextEmail, claims.Email)
		c.Set(ContextRole, claims.Role)

		c.Next()
	}
}

// AdminMiddleware 는 관리자 권한을 검증하는 미들웨어 (AuthMiddleware 이후에 사용)
func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get(ContextRole)
		if !exists || role != "admin" {
			c.JSON(http.StatusForbidden, dto.NewErrorResponse(
				"FORBIDDEN", "권한 부족", http.StatusForbidden, "관리자 권한이 필요합니다",
			))
			c.Abort()
			return
		}
		c.Next()
	}
}

// GetUserIDFromContext 는 컨텍스트에서 사용자 ID를 추출하는 헬퍼 함수
func GetUserIDFromContext(c *gin.Context) (uint, bool) {
	userID, exists := c.Get(ContextUserID)
	if !exists {
		return 0, false
	}
	id, ok := userID.(uint)
	return id, ok
}

// GetUserRoleFromContext 는 컨텍스트에서 사용자 역할을 추출하는 헬퍼 함수
func GetUserRoleFromContext(c *gin.Context) (string, bool) {
	role, exists := c.Get(ContextRole)
	if !exists {
		return "", false
	}
	r, ok := role.(string)
	return r, ok
}
