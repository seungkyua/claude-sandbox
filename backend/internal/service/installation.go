package service

import (
	"errors"

	"github.com/ktc-plugin-hub/backend/internal/dto"
	"github.com/ktc-plugin-hub/backend/internal/model"
	"github.com/ktc-plugin-hub/backend/internal/repository"
)

var ErrAlreadyInstalled = errors.New("이미 설치된 플러그인입니다")

// InstallationService 는 설치 관련 비즈니스 로직 인터페이스
type InstallationService interface {
	Install(userID uint, pluginID uint, versionID *uint) (*dto.InstallationResponse, error)
	Uninstall(userID uint, pluginID uint) error
	ToggleActive(userID uint, pluginID uint, isActive bool) (*dto.InstallationResponse, error)
	GetMyInstallations(userID uint) ([]dto.InstallationResponse, error)
}

type installationService struct {
	installRepo repository.InstallationRepository
	pluginRepo  repository.PluginRepository
	versionRepo repository.VersionRepository
}

func NewInstallationService(installRepo repository.InstallationRepository, pluginRepo repository.PluginRepository, versionRepo repository.VersionRepository) InstallationService {
	return &installationService{installRepo: installRepo, pluginRepo: pluginRepo, versionRepo: versionRepo}
}

func (s *installationService) Install(userID uint, pluginID uint, versionID *uint) (*dto.InstallationResponse, error) {
	// 이미 설치 확인
	_, err := s.installRepo.FindByUserAndPlugin(userID, pluginID)
	if err == nil {
		return nil, ErrAlreadyInstalled
	}

	// 버전 ID 결정 (미지정 시 최신 버전)
	var vid uint
	if versionID != nil {
		vid = *versionID
	} else {
		latest, err := s.versionRepo.FindLatestByPluginID(pluginID)
		if err != nil {
			return nil, ErrNotFound
		}
		vid = latest.ID
	}

	installation := &model.Installation{
		UserID:   userID,
		PluginID: pluginID,
		VersionID: vid,
		IsActive: true,
	}

	if err := s.installRepo.Create(installation); err != nil {
		return nil, err
	}

	// 다운로드 횟수 증가
	s.pluginRepo.IncrementDownloadCount(pluginID)

	return &dto.InstallationResponse{
		ID:        installation.ID,
		PluginID:  installation.PluginID,
		VersionID: installation.VersionID,
		IsActive:  installation.IsActive,
		InstalledAt: installation.InstalledAt,
	}, nil
}

func (s *installationService) Uninstall(userID uint, pluginID uint) error {
	return s.installRepo.Delete(userID, pluginID)
}

func (s *installationService) ToggleActive(userID uint, pluginID uint, isActive bool) (*dto.InstallationResponse, error) {
	err := s.installRepo.UpdateActive(userID, pluginID, isActive)
	if err != nil {
		return nil, ErrNotFound
	}

	installation, err := s.installRepo.FindByUserAndPlugin(userID, pluginID)
	if err != nil {
		return nil, ErrNotFound
	}

	return &dto.InstallationResponse{
		ID:        installation.ID,
		PluginID:  installation.PluginID,
		VersionID: installation.VersionID,
		IsActive:  installation.IsActive,
		InstalledAt: installation.InstalledAt,
	}, nil
}

func (s *installationService) GetMyInstallations(userID uint) ([]dto.InstallationResponse, error) {
	installations, err := s.installRepo.FindByUserID(userID)
	if err != nil {
		return nil, err
	}

	var result []dto.InstallationResponse
	for _, i := range installations {
		result = append(result, dto.InstallationResponse{
			ID:        i.ID,
			PluginID:  i.PluginID,
			VersionID: i.VersionID,
			IsActive:  i.IsActive,
			InstalledAt: i.InstalledAt,
		})
	}
	return result, nil
}
