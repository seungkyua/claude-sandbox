package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ktc-plugin-hub/backend/internal/dto"
	mw "github.com/ktc-plugin-hub/backend/internal/middleware"
	"github.com/ktc-plugin-hub/backend/internal/service"
)

// AuthHandler 는 인증 관련 HTTP 핸들러
type AuthHandler struct {
	authService service.AuthService
}

// NewAuthHandler 는 AuthHandler를 생성한다
func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// Register 는 회원가입을 처리한다
// POST /api/v1/auth/register
func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse(
			"INVALID_INPUT", "입력 오류", http.StatusBadRequest, err.Error(),
		))
		return
	}

	resp, err := h.authService.Register(&req)
	if err != nil {
		if err == service.ErrDuplicateEmail {
			c.JSON(http.StatusConflict, dto.NewErrorResponse(
				"DUPLICATE_EMAIL", "이메일 중복", http.StatusConflict, err.Error(),
			))
			return
		}
		c.JSON(http.StatusInternalServerError, dto.NewErrorResponse(
			"INTERNAL_ERROR", "서버 오류", http.StatusInternalServerError, err.Error(),
		))
		return
	}

	c.JSON(http.StatusCreated, resp)
}

// Login 은 로그인을 처리한다
// POST /api/v1/auth/login
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse(
			"INVALID_INPUT", "입력 오류", http.StatusBadRequest, err.Error(),
		))
		return
	}

	resp, err := h.authService.Login(&req)
	if err != nil {
		if err == service.ErrInvalidCredentials {
			c.JSON(http.StatusUnauthorized, dto.NewErrorResponse(
				"INVALID_CREDENTIALS", "인증 실패", http.StatusUnauthorized, err.Error(),
			))
			return
		}
		c.JSON(http.StatusInternalServerError, dto.NewErrorResponse(
			"INTERNAL_ERROR", "서버 오류", http.StatusInternalServerError, err.Error(),
		))
		return
	}

	c.JSON(http.StatusOK, resp)
}

// Refresh 는 토큰을 갱신한다
// POST /api/v1/auth/refresh
func (h *AuthHandler) Refresh(c *gin.Context) {
	var req dto.RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse(
			"INVALID_INPUT", "입력 오류", http.StatusBadRequest, err.Error(),
		))
		return
	}

	resp, err := h.authService.RefreshToken(req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.NewErrorResponse(
			"INVALID_TOKEN", "토큰 오류", http.StatusUnauthorized, err.Error(),
		))
		return
	}

	c.JSON(http.StatusOK, resp)
}

// Me 는 현재 로그인한 사용자 정보를 반환한다
// GET /api/v1/me
func (h *AuthHandler) Me(c *gin.Context) {
	userID, ok := mw.GetUserIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, dto.NewErrorResponse(
			"UNAUTHORIZED", "인증 필요", http.StatusUnauthorized, "로그인이 필요합니다",
		))
		return
	}

	resp, err := h.authService.GetUserByID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.NewErrorResponse(
			"NOT_FOUND", "사용자 없음", http.StatusNotFound, "사용자를 찾을 수 없습니다",
		))
		return
	}

	c.JSON(http.StatusOK, resp)
}
