package service

import (
	"testing"

	"github.com/ktc-plugin-hub/backend/internal/model"
	"github.com/ktc-plugin-hub/backend/internal/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestInstallationService() InstallationService {
	installRepo := repository.NewMockInstallationRepository()
	pluginRepo := repository.NewMockPluginRepository()
	versionRepo := repository.NewMockVersionRepository()

	pluginRepo.Create(&model.Plugin{Name: "test-plugin", AuthorID: 1, CategoryID: 1, Description: "d", Status: "approved"})
	versionRepo.Create(&model.PluginVersion{PluginID: 1, Version: "1.0.0", FilePath: "/p", FileSize: 1024})

	return NewInstallationService(installRepo, pluginRepo, versionRepo)
}

func TestShouldInstallPluginSuccessfully(t *testing.T) {
	svc := newTestInstallationService()

	resp, err := svc.Install(1, 1, nil)
	require.NoError(t, err)
	assert.True(t, resp.IsActive)
	assert.Equal(t, uint(1), resp.PluginID)
}

func TestShouldReturnErrorWhenAlreadyInstalled(t *testing.T) {
	svc := newTestInstallationService()
	svc.Install(1, 1, nil)

	_, err := svc.Install(1, 1, nil)
	assert.Equal(t, ErrAlreadyInstalled, err)
}

func TestShouldUninstallSuccessfully(t *testing.T) {
	svc := newTestInstallationService()
	svc.Install(1, 1, nil)

	err := svc.Uninstall(1, 1)
	require.NoError(t, err)
}

func TestShouldToggleActiveSuccessfully(t *testing.T) {
	svc := newTestInstallationService()
	svc.Install(1, 1, nil)

	resp, err := svc.ToggleActive(1, 1, false)
	require.NoError(t, err)
	assert.False(t, resp.IsActive)
}

func TestShouldGetMyInstallations(t *testing.T) {
	svc := newTestInstallationService()
	svc.Install(1, 1, nil)

	result, err := svc.GetMyInstallations(1)
	require.NoError(t, err)
	assert.Len(t, result, 1)
}
