package repository

import (
	"sort"
	"strings"

	"github.com/ktc-plugin-hub/backend/internal/model"
)

// MockPluginRepository 는 테스트용 플러그인 리포지토리 모킹
type MockPluginRepository struct {
	plugins map[uint]*model.Plugin
	nextID  uint
}

// NewMockPluginRepository 는 MockPluginRepository를 생성한다
func NewMockPluginRepository() *MockPluginRepository {
	return &MockPluginRepository{
		plugins: make(map[uint]*model.Plugin),
		nextID:  1,
	}
}

func (r *MockPluginRepository) Create(plugin *model.Plugin) error {
	for _, p := range r.plugins {
		if p.Name == plugin.Name {
			return ErrDuplicateKey
		}
	}
	plugin.ID = r.nextID
	r.nextID++
	r.plugins[plugin.ID] = plugin
	return nil
}

func (r *MockPluginRepository) FindByID(id uint) (*model.Plugin, error) {
	p, ok := r.plugins[id]
	if !ok {
		return nil, ErrNotFound
	}
	return p, nil
}

func (r *MockPluginRepository) FindAll(filter PluginFilter) ([]model.Plugin, int64, error) {
	var result []model.Plugin
	for _, p := range r.plugins {
		if filter.Status != "" && p.Status != filter.Status {
			continue
		}
		if filter.CategoryID != nil && p.CategoryID != *filter.CategoryID {
			continue
		}
		if filter.Keyword != "" {
			kw := strings.ToLower(filter.Keyword)
			if !strings.Contains(strings.ToLower(p.Name), kw) &&
				!strings.Contains(strings.ToLower(p.Description), kw) {
				continue
			}
		}
		if filter.AuthorID != nil && p.AuthorID != *filter.AuthorID {
			continue
		}
		result = append(result, *p)
	}

	total := int64(len(result))

	// 정렬
	switch filter.Sort {
	case "popular":
		sort.Slice(result, func(i, j int) bool { return result[i].DownloadCount > result[j].DownloadCount })
	case "rating":
		sort.Slice(result, func(i, j int) bool { return result[i].AvgRating > result[j].AvgRating })
	default:
		sort.Slice(result, func(i, j int) bool { return result[i].ID > result[j].ID })
	}

	// 페이지네이션
	start := filter.Offset
	if start > len(result) {
		start = len(result)
	}
	end := start + filter.Limit
	if end > len(result) {
		end = len(result)
	}

	return result[start:end], total, nil
}

func (r *MockPluginRepository) Update(plugin *model.Plugin) error {
	if _, ok := r.plugins[plugin.ID]; !ok {
		return ErrNotFound
	}
	r.plugins[plugin.ID] = plugin
	return nil
}

func (r *MockPluginRepository) Delete(id uint) error {
	if _, ok := r.plugins[id]; !ok {
		return ErrNotFound
	}
	delete(r.plugins, id)
	return nil
}

func (r *MockPluginRepository) FindByName(name string) (*model.Plugin, error) {
	for _, p := range r.plugins {
		if p.Name == name {
			return p, nil
		}
	}
	return nil, ErrNotFound
}

func (r *MockPluginRepository) IncrementDownloadCount(id uint) error {
	p, ok := r.plugins[id]
	if !ok {
		return ErrNotFound
	}
	p.DownloadCount++
	return nil
}

func (r *MockPluginRepository) UpdateRating(id uint, avgRating float64, reviewCount int) error {
	p, ok := r.plugins[id]
	if !ok {
		return ErrNotFound
	}
	p.AvgRating = avgRating
	p.ReviewCount = reviewCount
	return nil
}
