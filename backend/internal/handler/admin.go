package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ktc-plugin-hub/backend/internal/dto"
	"github.com/ktc-plugin-hub/backend/internal/service"
)

// AdminHandler 는 관리자 관련 HTTP 핸들러
type AdminHandler struct {
	adminService service.AdminService
}

func NewAdminHandler(adminService service.AdminService) *AdminHandler {
	return &AdminHandler{adminService: adminService}
}

// GetPendingPlugins 은 심사 대기 플러그인 목록을 반환한다
// GET /api/v1/admin/plugins/pending
func (h *AdminHandler) GetPendingPlugins(c *gin.Context) {
	result, err := h.adminService.GetPendingPlugins()
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewErrorResponse("INTERNAL_ERROR", "서버 오류", http.StatusInternalServerError, err.Error()))
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": result})
}

// ApprovePlugin 은 플러그인을 승인한다
// PATCH /api/v1/admin/plugins/:id/approve
func (h *AdminHandler) ApprovePlugin(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	resp, err := h.adminService.ApprovePlugin(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, dto.NewErrorResponse("NOT_FOUND", "플러그인 없음", http.StatusNotFound, err.Error()))
		return
	}
	c.JSON(http.StatusOK, resp)
}

// RejectPlugin 은 플러그인을 반려한다
// PATCH /api/v1/admin/plugins/:id/reject
func (h *AdminHandler) RejectPlugin(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	var req dto.RejectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse("INVALID_INPUT", "입력 오류", http.StatusBadRequest, err.Error()))
		return
	}

	resp, err := h.adminService.RejectPlugin(uint(id), req.Reason)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.NewErrorResponse("NOT_FOUND", "플러그인 없음", http.StatusNotFound, err.Error()))
		return
	}
	c.JSON(http.StatusOK, resp)
}

// HidePlugin 은 플러그인을 비공개 처리한다
// PATCH /api/v1/admin/plugins/:id/hide
func (h *AdminHandler) HidePlugin(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	resp, err := h.adminService.HidePlugin(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, dto.NewErrorResponse("NOT_FOUND", "플러그인 없음", http.StatusNotFound, err.Error()))
		return
	}
	c.JSON(http.StatusOK, resp)
}
