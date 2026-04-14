package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ktc-plugin-hub/backend/internal/dto"
	mw "github.com/ktc-plugin-hub/backend/internal/middleware"
	"github.com/ktc-plugin-hub/backend/internal/service"
)

// PluginHandler 는 플러그인 관련 HTTP 핸들러
type PluginHandler struct {
	pluginService service.PluginService
}

// NewPluginHandler 는 PluginHandler를 생성한다
func NewPluginHandler(pluginService service.PluginService) *PluginHandler {
	return &PluginHandler{pluginService: pluginService}
}

// Create 는 플러그인을 등록한다
// POST /api/v1/plugins
func (h *PluginHandler) Create(c *gin.Context) {
	var req dto.CreatePluginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse(
			"INVALID_INPUT", "입력 오류", http.StatusBadRequest, err.Error(),
		))
		return
	}

	userID, _ := mw.GetUserIDFromContext(c)
	role, _ := mw.GetUserRoleFromContext(c)

	resp, err := h.pluginService.Create(&req, userID, role)
	if err != nil {
		if err == service.ErrDuplicateName {
			c.JSON(http.StatusConflict, dto.NewErrorResponse(
				"DUPLICATE_NAME", "플러그인명 중복", http.StatusConflict, err.Error(),
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

// GetList 는 플러그인 목록을 반환한다
// GET /api/v1/plugins
func (h *PluginHandler) GetList(c *gin.Context) {
	var req dto.PluginListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse(
			"INVALID_INPUT", "입력 오류", http.StatusBadRequest, err.Error(),
		))
		return
	}

	resp, err := h.pluginService.GetList(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewErrorResponse(
			"INTERNAL_ERROR", "서버 오류", http.StatusInternalServerError, err.Error(),
		))
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetByID 는 플러그인 상세를 반환한다
// GET /api/v1/plugins/:id
func (h *PluginHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse(
			"INVALID_INPUT", "잘못된 ID", http.StatusBadRequest, "유효한 플러그인 ID를 입력하세요",
		))
		return
	}

	resp, err := h.pluginService.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, dto.NewErrorResponse(
			"NOT_FOUND", "플러그인 없음", http.StatusNotFound, err.Error(),
		))
		return
	}

	c.JSON(http.StatusOK, resp)
}

// Update 는 플러그인을 수정한다
// PUT /api/v1/plugins/:id
func (h *PluginHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse(
			"INVALID_INPUT", "잘못된 ID", http.StatusBadRequest, "유효한 플러그인 ID를 입력하세요",
		))
		return
	}

	var req dto.UpdatePluginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse(
			"INVALID_INPUT", "입력 오류", http.StatusBadRequest, err.Error(),
		))
		return
	}

	userID, _ := mw.GetUserIDFromContext(c)
	role, _ := mw.GetUserRoleFromContext(c)

	resp, err := h.pluginService.Update(uint(id), &req, userID, role)
	if err != nil {
		switch err {
		case service.ErrNotFound:
			c.JSON(http.StatusNotFound, dto.NewErrorResponse("NOT_FOUND", "플러그인 없음", http.StatusNotFound, err.Error()))
		case service.ErrForbidden:
			c.JSON(http.StatusForbidden, dto.NewErrorResponse("FORBIDDEN", "권한 없음", http.StatusForbidden, err.Error()))
		case service.ErrDuplicateName:
			c.JSON(http.StatusConflict, dto.NewErrorResponse("DUPLICATE_NAME", "플러그인명 중복", http.StatusConflict, err.Error()))
		default:
			c.JSON(http.StatusInternalServerError, dto.NewErrorResponse("INTERNAL_ERROR", "서버 오류", http.StatusInternalServerError, err.Error()))
		}
		return
	}

	c.JSON(http.StatusOK, resp)
}

// Delete 는 플러그인을 삭제한다
// DELETE /api/v1/plugins/:id
func (h *PluginHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse(
			"INVALID_INPUT", "잘못된 ID", http.StatusBadRequest, "유효한 플러그인 ID를 입력하세요",
		))
		return
	}

	userID, _ := mw.GetUserIDFromContext(c)
	role, _ := mw.GetUserRoleFromContext(c)

	err = h.pluginService.Delete(uint(id), userID, role)
	if err != nil {
		switch err {
		case service.ErrNotFound:
			c.JSON(http.StatusNotFound, dto.NewErrorResponse("NOT_FOUND", "플러그인 없음", http.StatusNotFound, err.Error()))
		case service.ErrForbidden:
			c.JSON(http.StatusForbidden, dto.NewErrorResponse("FORBIDDEN", "권한 없음", http.StatusForbidden, err.Error()))
		default:
			c.JSON(http.StatusInternalServerError, dto.NewErrorResponse("INTERNAL_ERROR", "서버 오류", http.StatusInternalServerError, err.Error()))
		}
		return
	}

	c.Status(http.StatusNoContent)
}
