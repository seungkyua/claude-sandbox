package repository

import (
	"github.com/ktc-plugin-hub/backend/internal/model"
)

type MockInstallationRepository struct {
	installations map[uint]*model.Installation
	nextID        uint
}

func NewMockInstallationRepository() *MockInstallationRepository {
	return &MockInstallationRepository{
		installations: make(map[uint]*model.Installation),
		nextID:        1,
	}
}

func (r *MockInstallationRepository) Create(installation *model.Installation) error {
	// 중복 검증 (user_id + plugin_id)
	for _, i := range r.installations {
		if i.UserID == installation.UserID && i.PluginID == installation.PluginID {
			return ErrDuplicateKey
		}
	}
	installation.ID = r.nextID
	r.nextID++
	r.installations[installation.ID] = installation
	return nil
}

func (r *MockInstallationRepository) Delete(userID uint, pluginID uint) error {
	for id, i := range r.installations {
		if i.UserID == userID && i.PluginID == pluginID {
			delete(r.installations, id)
			return nil
		}
	}
	return ErrNotFound
}

func (r *MockInstallationRepository) FindByUserID(userID uint) ([]model.Installation, error) {
	var result []model.Installation
	for _, i := range r.installations {
		if i.UserID == userID {
			result = append(result, *i)
		}
	}
	return result, nil
}

func (r *MockInstallationRepository) FindByUserAndPlugin(userID uint, pluginID uint) (*model.Installation, error) {
	for _, i := range r.installations {
		if i.UserID == userID && i.PluginID == pluginID {
			return i, nil
		}
	}
	return nil, ErrNotFound
}

func (r *MockInstallationRepository) UpdateActive(userID uint, pluginID uint, isActive bool) error {
	for _, i := range r.installations {
		if i.UserID == userID && i.PluginID == pluginID {
			i.IsActive = isActive
			return nil
		}
	}
	return ErrNotFound
}
