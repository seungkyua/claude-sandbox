package repository

import (
	"sort"

	"github.com/ktc-plugin-hub/backend/internal/model"
)

// MockVersionRepository 는 테스트용 버전 리포지토리 모킹
type MockVersionRepository struct {
	versions map[uint]*model.PluginVersion
	nextID   uint
}

func NewMockVersionRepository() *MockVersionRepository {
	return &MockVersionRepository{
		versions: make(map[uint]*model.PluginVersion),
		nextID:   1,
	}
}

func (r *MockVersionRepository) Create(version *model.PluginVersion) error {
	if r.IsDuplicateVersion(version.PluginID, version.Version) {
		return ErrDuplicateKey
	}
	version.ID = r.nextID
	r.nextID++
	r.versions[version.ID] = version
	return nil
}

func (r *MockVersionRepository) FindByPluginID(pluginID uint) ([]model.PluginVersion, error) {
	var result []model.PluginVersion
	for _, v := range r.versions {
		if v.PluginID == pluginID {
			result = append(result, *v)
		}
	}
	sort.Slice(result, func(i, j int) bool { return result[i].ID > result[j].ID })
	return result, nil
}

func (r *MockVersionRepository) FindByID(id uint) (*model.PluginVersion, error) {
	v, ok := r.versions[id]
	if !ok {
		return nil, ErrNotFound
	}
	return v, nil
}

func (r *MockVersionRepository) FindLatestByPluginID(pluginID uint) (*model.PluginVersion, error) {
	versions, _ := r.FindByPluginID(pluginID)
	if len(versions) == 0 {
		return nil, ErrNotFound
	}
	return &versions[0], nil
}

func (r *MockVersionRepository) IsDuplicateVersion(pluginID uint, version string) bool {
	for _, v := range r.versions {
		if v.PluginID == pluginID && v.Version == version {
			return true
		}
	}
	return false
}
