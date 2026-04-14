package repository

import (
	"github.com/ktc-plugin-hub/backend/internal/model"
	"gorm.io/gorm"
)

// InstallationRepository 는 설치 데이터 접근 인터페이스
type InstallationRepository interface {
	Create(installation *model.Installation) error
	Delete(userID uint, pluginID uint) error
	FindByUserID(userID uint) ([]model.Installation, error)
	FindByUserAndPlugin(userID uint, pluginID uint) (*model.Installation, error)
	UpdateActive(userID uint, pluginID uint, isActive bool) error
}

type installationRepository struct {
	db *gorm.DB
}

func NewInstallationRepository(db *gorm.DB) InstallationRepository {
	return &installationRepository{db: db}
}

func (r *installationRepository) Create(installation *model.Installation) error {
	return r.db.Create(installation).Error
}

func (r *installationRepository) Delete(userID uint, pluginID uint) error {
	return r.db.Where("user_id = ? AND plugin_id = ?", userID, pluginID).Delete(&model.Installation{}).Error
}

func (r *installationRepository) FindByUserID(userID uint) ([]model.Installation, error) {
	var installations []model.Installation
	err := r.db.Where("user_id = ?", userID).Preload("Plugin").Preload("Version").Find(&installations).Error
	return installations, err
}

func (r *installationRepository) FindByUserAndPlugin(userID uint, pluginID uint) (*model.Installation, error) {
	var installation model.Installation
	err := r.db.Where("user_id = ? AND plugin_id = ?", userID, pluginID).First(&installation).Error
	if err != nil {
		return nil, err
	}
	return &installation, nil
}

func (r *installationRepository) UpdateActive(userID uint, pluginID uint, isActive bool) error {
	return r.db.Model(&model.Installation{}).Where("user_id = ? AND plugin_id = ?", userID, pluginID).
		Update("is_active", isActive).Error
}
