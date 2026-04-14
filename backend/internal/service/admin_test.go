package service

import (
	"testing"

	"github.com/ktc-plugin-hub/backend/internal/model"
	"github.com/ktc-plugin-hub/backend/internal/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestAdminService() (AdminService, *repository.MockPluginRepository) {
	pluginRepo := repository.NewMockPluginRepository()
	return NewAdminService(pluginRepo), pluginRepo
}

func TestShouldGetPendingPlugins(t *testing.T) {
	svc, pluginRepo := newTestAdminService()
	pluginRepo.Create(&model.Plugin{Name: "pending-p", AuthorID: 1, CategoryID: 1, Description: "d", Status: model.PluginStatusPending})
	pluginRepo.Create(&model.Plugin{Name: "approved-p", AuthorID: 1, CategoryID: 1, Description: "d", Status: model.PluginStatusApproved})

	result, err := svc.GetPendingPlugins()
	require.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, "pending-p", result[0].Name)
}

func TestShouldApprovePlugin(t *testing.T) {
	svc, pluginRepo := newTestAdminService()
	pluginRepo.Create(&model.Plugin{Name: "p1", AuthorID: 1, CategoryID: 1, Description: "d", Status: model.PluginStatusPending})

	resp, err := svc.ApprovePlugin(1)
	require.NoError(t, err)
	assert.Equal(t, model.PluginStatusApproved, resp.Status)
}

func TestShouldRejectPlugin(t *testing.T) {
	svc, pluginRepo := newTestAdminService()
	pluginRepo.Create(&model.Plugin{Name: "p1", AuthorID: 1, CategoryID: 1, Description: "d", Status: model.PluginStatusPending})

	resp, err := svc.RejectPlugin(1, "부적절한 내용")
	require.NoError(t, err)
	assert.Equal(t, model.PluginStatusRejected, resp.Status)
}

func TestShouldHidePlugin(t *testing.T) {
	svc, pluginRepo := newTestAdminService()
	pluginRepo.Create(&model.Plugin{Name: "p1", AuthorID: 1, CategoryID: 1, Description: "d", Status: model.PluginStatusApproved})

	resp, err := svc.HidePlugin(1)
	require.NoError(t, err)
	assert.Equal(t, model.PluginStatusHidden, resp.Status)
}
