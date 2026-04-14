package service

import (
	"errors"

	"github.com/ktc-plugin-hub/backend/internal/dto"
	"github.com/ktc-plugin-hub/backend/internal/model"
	"github.com/ktc-plugin-hub/backend/internal/repository"
)

// 플러그인 서비스 에러 정의
var (
	ErrDuplicateName = errors.New("이미 사용 중인 플러그인명입니다")
	ErrNotFound      = errors.New("플러그인을 찾을 수 없습니다")
	ErrForbidden     = errors.New("권한이 없습니다")
)

// PluginService 는 플러그인 비즈니스 로직 인터페이스
type PluginService interface {
	Create(req *dto.CreatePluginRequest, userID uint, role string) (*dto.PluginResponse, error)
	GetByID(id uint) (*dto.PluginResponse, error)
	GetList(req *dto.PluginListRequest) (*dto.PaginatedResponse, error)
	Update(id uint, req *dto.UpdatePluginRequest, userID uint, role string) (*dto.PluginResponse, error)
	Delete(id uint, userID uint, role string) error
	GetByAuthor(authorID uint) ([]dto.PluginResponse, error)
}

type pluginService struct {
	pluginRepo   repository.PluginRepository
	categoryRepo repository.CategoryRepository
}

// NewPluginService 는 PluginService 인스턴스를 생성한다
func NewPluginService(pluginRepo repository.PluginRepository, categoryRepo repository.CategoryRepository) PluginService {
	return &pluginService{pluginRepo: pluginRepo, categoryRepo: categoryRepo}
}

func (s *pluginService) Create(req *dto.CreatePluginRequest, userID uint, role string) (*dto.PluginResponse, error) {
	// 중복명 검증
	_, err := s.pluginRepo.FindByName(req.Name)
	if err == nil {
		return nil, ErrDuplicateName
	}

	// 관리자면 즉시 공개 + 공식, 일반 사용자는 pending
	status := model.PluginStatusPending
	isOfficial := false
	if role == "admin" {
		status = model.PluginStatusApproved
		isOfficial = true
	}

	plugin := &model.Plugin{
		AuthorID:    userID,
		CategoryID:  req.CategoryID,
		Name:        req.Name,
		Description: req.Description,
		Status:      status,
		IsOfficial:  isOfficial,
	}

	if err := s.pluginRepo.Create(plugin); err != nil {
		return nil, err
	}

	return s.toResponse(plugin), nil
}

func (s *pluginService) GetByID(id uint) (*dto.PluginResponse, error) {
	plugin, err := s.pluginRepo.FindByID(id)
	if err != nil {
		return nil, ErrNotFound
	}
	return s.toResponse(plugin), nil
}

func (s *pluginService) GetList(req *dto.PluginListRequest) (*dto.PaginatedResponse, error) {
	req.Normalize()

	filter := repository.PluginFilter{
		CategoryID: req.CategoryID,
		Keyword:    req.Keyword,
		Sort:       req.Sort,
		Status:     model.PluginStatusApproved, // 공개된 플러그인만
		Offset:     req.Offset(),
		Limit:      req.Size,
	}

	plugins, total, err := s.pluginRepo.FindAll(filter)
	if err != nil {
		return nil, err
	}

	var data []dto.PluginResponse
	for i := range plugins {
		data = append(data, *s.toResponse(&plugins[i]))
	}

	return &dto.PaginatedResponse{
		Data:  data,
		Total: total,
		Page:  req.Page,
		Size:  req.Size,
	}, nil
}

func (s *pluginService) Update(id uint, req *dto.UpdatePluginRequest, userID uint, role string) (*dto.PluginResponse, error) {
	plugin, err := s.pluginRepo.FindByID(id)
	if err != nil {
		return nil, ErrNotFound
	}

	// 권한 검증: 본인 또는 관리자만 수정 가능
	if plugin.AuthorID != userID && role != "admin" {
		return nil, ErrForbidden
	}

	if req.Name != nil {
		// 중복명 검증
		existing, err := s.pluginRepo.FindByName(*req.Name)
		if err == nil && existing.ID != id {
			return nil, ErrDuplicateName
		}
		plugin.Name = *req.Name
	}
	if req.Description != nil {
		plugin.Description = *req.Description
	}
	if req.CategoryID != nil {
		plugin.CategoryID = *req.CategoryID
	}

	if err := s.pluginRepo.Update(plugin); err != nil {
		return nil, err
	}

	return s.toResponse(plugin), nil
}

func (s *pluginService) Delete(id uint, userID uint, role string) error {
	plugin, err := s.pluginRepo.FindByID(id)
	if err != nil {
		return ErrNotFound
	}

	// 권한 검증: 본인 또는 관리자만 삭제 가능
	if plugin.AuthorID != userID && role != "admin" {
		return ErrForbidden
	}

	return s.pluginRepo.Delete(id)
}

func (s *pluginService) GetByAuthor(authorID uint) ([]dto.PluginResponse, error) {
	filter := repository.PluginFilter{
		AuthorID: &authorID,
		Offset:   0,
		Limit:    100,
	}
	plugins, _, err := s.pluginRepo.FindAll(filter)
	if err != nil {
		return nil, err
	}

	var result []dto.PluginResponse
	for i := range plugins {
		result = append(result, *s.toResponse(&plugins[i]))
	}
	return result, nil
}

func (s *pluginService) toResponse(p *model.Plugin) *dto.PluginResponse {
	return &dto.PluginResponse{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		Author: dto.AuthorResponse{
			ID:       p.AuthorID,
			Nickname: p.Author.Nickname,
		},
		Category: dto.CategoryResponse{
			ID:   p.CategoryID,
			Name: p.Category.Name,
		},
		Status:        p.Status,
		IsOfficial:    p.IsOfficial,
		DownloadCount: p.DownloadCount,
		AvgRating:     p.AvgRating,
		ReviewCount:   p.ReviewCount,
		CreatedAt:     p.CreatedAt,
		UpdatedAt:     p.UpdatedAt,
	}
}
