package service

import (
	"testing"

	"github.com/ktc-plugin-hub/backend/internal/dto"
	"github.com/ktc-plugin-hub/backend/internal/model"
	"github.com/ktc-plugin-hub/backend/internal/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestVersionService() (VersionService, *repository.MockPluginRepository) {
	pluginRepo := repository.NewMockPluginRepository()
	versionRepo := repository.NewMockVersionRepository()
	return NewVersionService(versionRepo, pluginRepo), pluginRepo
}

func TestShouldCreateVersionSuccessfully(t *testing.T) {
	svc, pluginRepo := newTestVersionService()
	pluginRepo.Create(&model.Plugin{Name: "p1", AuthorID: 1, CategoryID: 1, Description: "d", Status: "approved"})

	resp, err := svc.CreateVersion(1, &dto.CreateVersionRequest{
		Version: "1.0.0", Changelog: "initial",
	}, "/path/file.zip", 1024, 1, "user")

	require.NoError(t, err)
	assert.Equal(t, "1.0.0", resp.Version)
}

func TestShouldReturnErrorWhenDuplicateVersionNumber(t *testing.T) {
	svc, pluginRepo := newTestVersionService()
	pluginRepo.Create(&model.Plugin{Name: "p1", AuthorID: 1, CategoryID: 1, Description: "d", Status: "approved"})

	svc.CreateVersion(1, &dto.CreateVersionRequest{Version: "1.0.0"}, "/p", 1024, 1, "user")
	_, err := svc.CreateVersion(1, &dto.CreateVersionRequest{Version: "1.0.0"}, "/p2", 2048, 1, "user")
	assert.Equal(t, ErrDuplicateVersion, err)
}

func TestShouldReturnForbiddenWhenNonOwnerCreatesVersion(t *testing.T) {
	svc, pluginRepo := newTestVersionService()
	pluginRepo.Create(&model.Plugin{Name: "p1", AuthorID: 1, CategoryID: 1, Description: "d", Status: "approved"})

	_, err := svc.CreateVersion(1, &dto.CreateVersionRequest{Version: "2.0.0"}, "/p", 1024, 99, "user")
	assert.Equal(t, ErrForbidden, err)
}

func TestShouldIncrementDownloadCountWhenDownload(t *testing.T) {
	svc, pluginRepo := newTestVersionService()
	pluginRepo.Create(&model.Plugin{Name: "p1", AuthorID: 1, CategoryID: 1, Description: "d", Status: "approved"})
	svc.CreateVersion(1, &dto.CreateVersionRequest{Version: "1.0.0"}, "/path", 1024, 1, "user")

	_, err := svc.DownloadVersion(1, 1)
	require.NoError(t, err)

	plugin, _ := pluginRepo.FindByID(1)
	assert.Equal(t, 1, plugin.DownloadCount)
}
