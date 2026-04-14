package repository

import (
	"github.com/ktc-plugin-hub/backend/internal/model"
	"gorm.io/gorm"
)

// VersionRepository 는 플러그인 버전 데이터 접근 인터페이스
type VersionRepository interface {
	Create(version *model.PluginVersion) error
	FindByPluginID(pluginID uint) ([]model.PluginVersion, error)
	FindByID(id uint) (*model.PluginVersion, error)
	FindLatestByPluginID(pluginID uint) (*model.PluginVersion, error)
	IsDuplicateVersion(pluginID uint, version string) bool
}

type versionRepository struct {
	db *gorm.DB
}

func NewVersionRepository(db *gorm.DB) VersionRepository {
	return &versionRepository{db: db}
}

func (r *versionRepository) Create(version *model.PluginVersion) error {
	return r.db.Create(version).Error
}

func (r *versionRepository) FindByPluginID(pluginID uint) ([]model.PluginVersion, error) {
	var versions []model.PluginVersion
	err := r.db.Where("plugin_id = ?", pluginID).Order("created_at DESC").Find(&versions).Error
	return versions, err
}

func (r *versionRepository) FindByID(id uint) (*model.PluginVersion, error) {
	var version model.PluginVersion
	err := r.db.First(&version, id).Error
	if err != nil {
		return nil, err
	}
	return &version, nil
}

func (r *versionRepository) FindLatestByPluginID(pluginID uint) (*model.PluginVersion, error) {
	var version model.PluginVersion
	err := r.db.Where("plugin_id = ?", pluginID).Order("created_at DESC").First(&version).Error
	if err != nil {
		return nil, err
	}
	return &version, nil
}

func (r *versionRepository) IsDuplicateVersion(pluginID uint, version string) bool {
	var count int64
	r.db.Model(&model.PluginVersion{}).Where("plugin_id = ? AND version = ?", pluginID, version).Count(&count)
	return count > 0
}
