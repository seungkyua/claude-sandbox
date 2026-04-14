package service

import (
	"github.com/ktc-plugin-hub/backend/internal/dto"
	"github.com/ktc-plugin-hub/backend/internal/model"
	"github.com/ktc-plugin-hub/backend/internal/repository"
)

// AdminService 는 관리자 비즈니스 로직 인터페이스
type AdminService interface {
	GetPendingPlugins() ([]dto.PluginResponse, error)
	ApprovePlugin(id uint) (*dto.PluginResponse, error)
	RejectPlugin(id uint, reason string) (*dto.PluginResponse, error)
	HidePlugin(id uint) (*dto.PluginResponse, error)
}

type adminService struct {
	pluginRepo repository.PluginRepository
}

func NewAdminService(pluginRepo repository.PluginRepository) AdminService {
	return &adminService{pluginRepo: pluginRepo}
}

func (s *adminService) GetPendingPlugins() ([]dto.PluginResponse, error) {
	plugins, _, err := s.pluginRepo.FindAll(repository.PluginFilter{
		Status: model.PluginStatusPending,
		Offset: 0,
		Limit:  100,
	})
	if err != nil {
		return nil, err
	}

	var result []dto.PluginResponse
	for _, p := range plugins {
		result = append(result, dto.PluginResponse{
			ID: p.ID, Name: p.Name, Description: p.Description,
			Author:   dto.AuthorResponse{ID: p.AuthorID},
			Category: dto.CategoryResponse{ID: p.CategoryID},
			Status:   p.Status, IsOfficial: p.IsOfficial,
			CreatedAt: p.CreatedAt,
		})
	}
	return result, nil
}

func (s *adminService) ApprovePlugin(id uint) (*dto.PluginResponse, error) {
	plugin, err := s.pluginRepo.FindByID(id)
	if err != nil {
		return nil, ErrNotFound
	}
	plugin.Status = model.PluginStatusApproved
	if err := s.pluginRepo.Update(plugin); err != nil {
		return nil, err
	}
	return s.toResponse(plugin), nil
}

func (s *adminService) RejectPlugin(id uint, reason string) (*dto.PluginResponse, error) {
	plugin, err := s.pluginRepo.FindByID(id)
	if err != nil {
		return nil, ErrNotFound
	}
	plugin.Status = model.PluginStatusRejected
	if err := s.pluginRepo.Update(plugin); err != nil {
		return nil, err
	}
	return s.toResponse(plugin), nil
}

func (s *adminService) HidePlugin(id uint) (*dto.PluginResponse, error) {
	plugin, err := s.pluginRepo.FindByID(id)
	if err != nil {
		return nil, ErrNotFound
	}
	plugin.Status = model.PluginStatusHidden
	if err := s.pluginRepo.Update(plugin); err != nil {
		return nil, err
	}
	return s.toResponse(plugin), nil
}

func (s *adminService) toResponse(p *model.Plugin) *dto.PluginResponse {
	return &dto.PluginResponse{
		ID: p.ID, Name: p.Name, Description: p.Description,
		Author:   dto.AuthorResponse{ID: p.AuthorID},
		Category: dto.CategoryResponse{ID: p.CategoryID},
		Status:   p.Status, IsOfficial: p.IsOfficial,
		CreatedAt: p.CreatedAt, UpdatedAt: p.UpdatedAt,
	}
}
