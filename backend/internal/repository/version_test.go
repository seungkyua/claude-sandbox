package repository

import (
	"testing"

	"github.com/ktc-plugin-hub/backend/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestShouldCreateVersionAndFindByPluginID(t *testing.T) {
	repo := NewMockVersionRepository()
	v := &model.PluginVersion{PluginID: 1, Version: "1.0.0", FilePath: "/path", FileSize: 1024}

	err := repo.Create(v)
	require.NoError(t, err)

	versions, err := repo.FindByPluginID(1)
	require.NoError(t, err)
	assert.Len(t, versions, 1)
}

func TestShouldReturnErrorWhenDuplicateVersion(t *testing.T) {
	repo := NewMockVersionRepository()
	repo.Create(&model.PluginVersion{PluginID: 1, Version: "1.0.0", FilePath: "/p", FileSize: 1024})
	err := repo.Create(&model.PluginVersion{PluginID: 1, Version: "1.0.0", FilePath: "/p2", FileSize: 2048})
	assert.Error(t, err)
}

func TestShouldFindLatestVersion(t *testing.T) {
	repo := NewMockVersionRepository()
	repo.Create(&model.PluginVersion{PluginID: 1, Version: "1.0.0", FilePath: "/p", FileSize: 1024})
	repo.Create(&model.PluginVersion{PluginID: 1, Version: "2.0.0", FilePath: "/p2", FileSize: 2048})

	latest, err := repo.FindLatestByPluginID(1)
	require.NoError(t, err)
	assert.Equal(t, "2.0.0", latest.Version)
}

func TestMockVersionRepositoryShouldImplementInterface(t *testing.T) {
	var _ VersionRepository = NewMockVersionRepository()
}
