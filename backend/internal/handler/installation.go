package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ktc-plugin-hub/backend/internal/dto"
	mw "github.com/ktc-plugin-hub/backend/internal/middleware"
	"github.com/ktc-plugin-hub/backend/internal/service"
)

// InstallationHandler 는 설치 관련 HTTP 핸들러
type InstallationHandler struct {
	installService service.InstallationService
}

func NewInstallationHandler(installService service.InstallationService) *InstallationHandler {
	return &InstallationHandler{installService: installService}
}

// Install 은 플러그인을 설치한다
// POST /api/v1/plugins/:id/install
func (h *InstallationHandler) Install(c *gin.Context) {
	pluginID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse("INVALID_INPUT", "잘못된 ID", http.StatusBadRequest, "유효한 플러그인 ID를 입력하세요"))
		return
	}

	var req dto.InstallRequest
	c.ShouldBindJSON(&req) // 선택적 파라미터이므로 에러 무시

	userID, _ := mw.GetUserIDFromContext(c)

	resp, err := h.installService.Install(userID, uint(pluginID), req.VersionID)
	if err != nil {
		if err == service.ErrAlreadyInstalled {
			c.JSON(http.StatusConflict, dto.NewErrorResponse("ALREADY_INSTALLED", "이미 설치됨", http.StatusConflict, err.Error()))
			return
		}
		c.JSON(http.StatusInternalServerError, dto.NewErrorResponse("INTERNAL_ERROR", "서버 오류", http.StatusInternalServerError, err.Error()))
		return
	}

	c.JSON(http.StatusCreated, resp)
}

// Uninstall 은 플러그인을 삭제한다
// DELETE /api/v1/plugins/:id/install
func (h *InstallationHandler) Uninstall(c *gin.Context) {
	pluginID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse("INVALID_INPUT", "잘못된 ID", http.StatusBadRequest, "유효한 플러그인 ID를 입력하세요"))
		return
	}

	userID, _ := mw.GetUserIDFromContext(c)

	if err := h.installService.Uninstall(userID, uint(pluginID)); err != nil {
		c.JSON(http.StatusNotFound, dto.NewErrorResponse("NOT_FOUND", "설치 정보 없음", http.StatusNotFound, err.Error()))
		return
	}

	c.Status(http.StatusNoContent)
}

// ToggleActive 는 활성화/비활성화를 토글한다
// PATCH /api/v1/plugins/:id/install
func (h *InstallationHandler) ToggleActive(c *gin.Context) {
	pluginID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse("INVALID_INPUT", "잘못된 ID", http.StatusBadRequest, "유효한 플러그인 ID를 입력하세요"))
		return
	}

	var req dto.ToggleActiveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse("INVALID_INPUT", "입력 오류", http.StatusBadRequest, err.Error()))
		return
	}

	userID, _ := mw.GetUserIDFromContext(c)

	resp, err := h.installService.ToggleActive(userID, uint(pluginID), req.IsActive)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.NewErrorResponse("NOT_FOUND", "설치 정보 없음", http.StatusNotFound, err.Error()))
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetMyInstallations 은 내 설치 목록을 반환한다
// GET /api/v1/me/installations
func (h *InstallationHandler) GetMyInstallations(c *gin.Context) {
	userID, _ := mw.GetUserIDFromContext(c)

	result, err := h.installService.GetMyInstallations(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewErrorResponse("INTERNAL_ERROR", "서버 오류", http.StatusInternalServerError, err.Error()))
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": result})
}
