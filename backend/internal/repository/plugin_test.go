package repository

import (
	"testing"

	"github.com/ktc-plugin-hub/backend/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createTestPlugin(name string, authorID uint) *model.Plugin {
	return &model.Plugin{
		AuthorID:    authorID,
		CategoryID:  1,
		Name:        name,
		Description: name + " 설명",
		Status:      model.PluginStatusApproved,
	}
}

// 플러그인 생성 및 조회 확인
func TestShouldCreatePluginAndFindByID(t *testing.T) {
	repo := NewMockPluginRepository()
	plugin := createTestPlugin("test-plugin", 1)

	err := repo.Create(plugin)
	require.NoError(t, err)
	assert.NotZero(t, plugin.ID)

	found, err := repo.FindByID(plugin.ID)
	require.NoError(t, err)
	assert.Equal(t, "test-plugin", found.Name)
}

// 이름으로 플러그인 조회 확인
func TestShouldFindPluginByName(t *testing.T) {
	repo := NewMockPluginRepository()
	repo.Create(createTestPlugin("my-plugin", 1))

	found, err := repo.FindByName("my-plugin")
	require.NoError(t, err)
	assert.Equal(t, "my-plugin", found.Name)
}

// 중복 이름으로 생성 시 에러
func TestShouldReturnErrorWhenDuplicatePluginName(t *testing.T) {
	repo := NewMockPluginRepository()
	repo.Create(createTestPlugin("dup-plugin", 1))
	err := repo.Create(createTestPlugin("dup-plugin", 2))
	assert.Error(t, err)
}

// 필터링 조회 (상태별)
func TestShouldFilterPluginsByStatus(t *testing.T) {
	repo := NewMockPluginRepository()
	p1 := createTestPlugin("approved-plugin", 1)
	p1.Status = model.PluginStatusApproved
	repo.Create(p1)

	p2 := createTestPlugin("pending-plugin", 1)
	p2.Status = model.PluginStatusPending
	repo.Create(p2)

	plugins, total, err := repo.FindAll(PluginFilter{
		Status: model.PluginStatusApproved, Offset: 0, Limit: 20,
	})
	require.NoError(t, err)
	assert.Equal(t, int64(1), total)
	assert.Len(t, plugins, 1)
	assert.Equal(t, "approved-plugin", plugins[0].Name)
}

// 키워드 검색 확인
func TestShouldFilterPluginsByKeyword(t *testing.T) {
	repo := NewMockPluginRepository()
	repo.Create(createTestPlugin("code-formatter", 1))
	repo.Create(createTestPlugin("image-viewer", 1))

	plugins, total, err := repo.FindAll(PluginFilter{
		Keyword: "code", Offset: 0, Limit: 20,
	})
	require.NoError(t, err)
	assert.Equal(t, int64(1), total)
	assert.Equal(t, "code-formatter", plugins[0].Name)
}

// 플러그인 업데이트 확인
func TestShouldUpdatePlugin(t *testing.T) {
	repo := NewMockPluginRepository()
	plugin := createTestPlugin("update-me", 1)
	repo.Create(plugin)

	plugin.Description = "updated description"
	err := repo.Update(plugin)
	require.NoError(t, err)

	found, _ := repo.FindByID(plugin.ID)
	assert.Equal(t, "updated description", found.Description)
}

// 플러그인 삭제 확인
func TestShouldDeletePlugin(t *testing.T) {
	repo := NewMockPluginRepository()
	plugin := createTestPlugin("delete-me", 1)
	repo.Create(plugin)

	err := repo.Delete(plugin.ID)
	require.NoError(t, err)

	_, err = repo.FindByID(plugin.ID)
	assert.Error(t, err)
}

// 다운로드 카운트 증가 확인
func TestShouldIncrementDownloadCount(t *testing.T) {
	repo := NewMockPluginRepository()
	plugin := createTestPlugin("dl-plugin", 1)
	repo.Create(plugin)

	repo.IncrementDownloadCount(plugin.ID)
	found, _ := repo.FindByID(plugin.ID)
	assert.Equal(t, 1, found.DownloadCount)
}

// MockPluginRepository가 인터페이스 구현 확인
func TestMockPluginRepositoryShouldImplementInterface(t *testing.T) {
	var _ PluginRepository = NewMockPluginRepository()
}
