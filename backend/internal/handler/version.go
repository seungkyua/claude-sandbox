package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ktc-plugin-hub/backend/internal/dto"
	mw "github.com/ktc-plugin-hub/backend/internal/middleware"
	"github.com/ktc-plugin-hub/backend/internal/service"
)

// VersionHandler 는 버전 관련 HTTP 핸들러
type VersionHandler struct {
	versionService service.VersionService
}

func NewVersionHandler(versionService service.VersionService) *VersionHandler {
	return &VersionHandler{versionService: versionService}
}

// CreateVersion 은 새 버전을 업로드한다
// POST /api/v1/plugins/:id/versions
func (h *VersionHandler) CreateVersion(c *gin.Context) {
	pluginID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse("INVALID_INPUT", "잘못된 ID", http.StatusBadRequest, "유효한 플러그인 ID를 입력하세요"))
		return
	}

	var req dto.CreateVersionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse("INVALID_INPUT", "입력 오류", http.StatusBadRequest, err.Error()))
		return
	}

	userID, _ := mw.GetUserIDFromContext(c)
	role, _ := mw.GetUserRoleFromContext(c)

	// 파일 업로드는 실제 환경에서 multipart로 처리, 여기서는 JSON으로 대체
	resp, err := h.versionService.CreateVersion(uint(pluginID), &req, "/uploads/placeholder", 0, userID, role)
	if err != nil {
		switch err {
		case service.ErrDuplicateVersion:
			c.JSON(http.StatusConflict, dto.NewErrorResponse("DUPLICATE_VERSION", "버전 중복", http.StatusConflict, err.Error()))
		case service.ErrNotFound:
			c.JSON(http.StatusNotFound, dto.NewErrorResponse("NOT_FOUND", "플러그인 없음", http.StatusNotFound, err.Error()))
		case service.ErrForbidden:
			c.JSON(http.StatusForbidden, dto.NewErrorResponse("FORBIDDEN", "권한 없음", http.StatusForbidden, err.Error()))
		default:
			c.JSON(http.StatusInternalServerError, dto.NewErrorResponse("INTERNAL_ERROR", "서버 오류", http.StatusInternalServerError, err.Error()))
		}
		return
	}

	c.JSON(http.StatusCreated, resp)
}

// Download 는 플러그인 파일을 다운로드한다 (다운로드 횟수 증가)
// GET /api/v1/plugins/:id/versions/:versionId/download
func (h *VersionHandler) Download(c *gin.Context) {
	pluginID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse("INVALID_INPUT", "잘못된 ID", http.StatusBadRequest, "유효한 플러그인 ID를 입력하세요"))
		return
	}
	versionID, err := strconv.ParseUint(c.Param("versionId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse("INVALID_INPUT", "잘못된 버전 ID", http.StatusBadRequest, "유효한 버전 ID를 입력하세요"))
		return
	}

	version, err := h.versionService.DownloadVersion(uint(pluginID), uint(versionID))
	if err != nil {
		c.JSON(http.StatusNotFound, dto.NewErrorResponse("NOT_FOUND", "버전 없음", http.StatusNotFound, err.Error()))
		return
	}

	// 실제로는 파일을 전송하지만, 테스트에서는 JSON으로 반환
	c.JSON(http.StatusOK, gin.H{
		"file_path": version.FilePath,
		"file_size": version.FileSize,
		"version":   version.Version,
	})
}
