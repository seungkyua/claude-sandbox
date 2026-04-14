package service

import (
	"errors"

	"github.com/ktc-plugin-hub/backend/internal/dto"
	"github.com/ktc-plugin-hub/backend/internal/model"
	"github.com/ktc-plugin-hub/backend/internal/repository"
)

var (
	ErrDuplicateVersion = errors.New("이미 존재하는 버전 번호입니다")
)

// VersionService 는 버전 관련 비즈니스 로직 인터페이스
type VersionService interface {
	CreateVersion(pluginID uint, req *dto.CreateVersionRequest, filePath string, fileSize int64, userID uint, role string) (*dto.VersionResponse, error)
	GetVersionsByPluginID(pluginID uint) ([]dto.VersionResponse, error)
	DownloadVersion(pluginID uint, versionID uint) (*model.PluginVersion, error)
}

type versionService struct {
	versionRepo repository.VersionRepository
	pluginRepo  repository.PluginRepository
}

func NewVersionService(versionRepo repository.VersionRepository, pluginRepo repository.PluginRepository) VersionService {
	return &versionService{versionRepo: versionRepo, pluginRepo: pluginRepo}
}

func (s *versionService) CreateVersion(pluginID uint, req *dto.CreateVersionRequest, filePath string, fileSize int64, userID uint, role string) (*dto.VersionResponse, error) {
	// 플러그인 존재 확인
	plugin, err := s.pluginRepo.FindByID(pluginID)
	if err != nil {
		return nil, ErrNotFound
	}

	// 권한 검증
	if plugin.AuthorID != userID && role != "admin" {
		return nil, ErrForbidden
	}

	// 중복 버전 검증
	if s.versionRepo.IsDuplicateVersion(pluginID, req.Version) {
		return nil, ErrDuplicateVersion
	}

	version := &model.PluginVersion{
		PluginID:  pluginID,
		Version:   req.Version,
		Changelog: req.Changelog,
		FilePath:  filePath,
		FileSize:  fileSize,
	}

	if err := s.versionRepo.Create(version); err != nil {
		return nil, err
	}

	return &dto.VersionResponse{
		ID:        version.ID,
		PluginID:  version.PluginID,
		Version:   version.Version,
		Changelog: version.Changelog,
		FileSize:  version.FileSize,
		CreatedAt: version.CreatedAt,
	}, nil
}

func (s *versionService) GetVersionsByPluginID(pluginID uint) ([]dto.VersionResponse, error) {
	versions, err := s.versionRepo.FindByPluginID(pluginID)
	if err != nil {
		return nil, err
	}

	var result []dto.VersionResponse
	for _, v := range versions {
		result = append(result, dto.VersionResponse{
			ID:        v.ID,
			PluginID:  v.PluginID,
			Version:   v.Version,
			Changelog: v.Changelog,
			FileSize:  v.FileSize,
			CreatedAt: v.CreatedAt,
		})
	}
	return result, nil
}

func (s *versionService) DownloadVersion(pluginID uint, versionID uint) (*model.PluginVersion, error) {
	version, err := s.versionRepo.FindByID(versionID)
	if err != nil {
		return nil, ErrNotFound
	}

	if version.PluginID != pluginID {
		return nil, ErrNotFound
	}

	// 다운로드 횟수 증가
	s.pluginRepo.IncrementDownloadCount(pluginID)

	return version, nil
}
