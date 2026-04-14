package service

import (
	"testing"

	"github.com/ktc-plugin-hub/backend/internal/dto"
	"github.com/ktc-plugin-hub/backend/internal/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestPluginService() PluginService {
	return NewPluginService(
		repository.NewMockPluginRepository(),
		repository.NewMockCategoryRepository(),
	)
}

// 일반 사용자가 플러그인 등록하면 pending 상태인지 확인
func TestShouldCreatePluginAsPendingWhenUserRegisters(t *testing.T) {
	svc := newTestPluginService()

	resp, err := svc.Create(&dto.CreatePluginRequest{
		Name: "my-plugin", Description: "desc", CategoryID: 1, Version: "1.0.0",
	}, 1, "user")

	require.NoError(t, err)
	assert.Equal(t, "pending", resp.Status)
	assert.False(t, resp.IsOfficial)
}

// 관리자가 플러그인 등록하면 즉시 approved + official인지 확인
func TestShouldCreatePluginAsApprovedWhenAdminRegisters(t *testing.T) {
	svc := newTestPluginService()

	resp, err := svc.Create(&dto.CreatePluginRequest{
		Name: "official-plugin", Description: "desc", CategoryID: 1, Version: "1.0.0",
	}, 1, "admin")

	require.NoError(t, err)
	assert.Equal(t, "approved", resp.Status)
	assert.True(t, resp.IsOfficial)
}

// 중복 플러그인명 등록 시 에러
func TestShouldReturnErrorWhenDuplicatePluginName(t *testing.T) {
	svc := newTestPluginService()
	svc.Create(&dto.CreatePluginRequest{
		Name: "dup-plugin", Description: "d", CategoryID: 1, Version: "1.0.0",
	}, 1, "user")

	_, err := svc.Create(&dto.CreatePluginRequest{
		Name: "dup-plugin", Description: "d2", CategoryID: 1, Version: "1.0.0",
	}, 2, "user")

	assert.Error(t, err)
	assert.Equal(t, ErrDuplicateName, err)
}

// 본인 플러그인 수정 확인
func TestShouldUpdatePluginWhenOwner(t *testing.T) {
	svc := newTestPluginService()
	created, _ := svc.Create(&dto.CreatePluginRequest{
		Name: "edit-me", Description: "original", CategoryID: 1, Version: "1.0.0",
	}, 1, "user")

	newDesc := "updated"
	updated, err := svc.Update(created.ID, &dto.UpdatePluginRequest{
		Description: &newDesc,
	}, 1, "user")

	require.NoError(t, err)
	assert.Equal(t, "updated", updated.Description)
}

// 타인 플러그인 수정 시 에러
func TestShouldReturnForbiddenWhenNonOwnerUpdates(t *testing.T) {
	svc := newTestPluginService()
	created, _ := svc.Create(&dto.CreatePluginRequest{
		Name: "not-mine", Description: "d", CategoryID: 1, Version: "1.0.0",
	}, 1, "user")

	newDesc := "hacked"
	_, err := svc.Update(created.ID, &dto.UpdatePluginRequest{
		Description: &newDesc,
	}, 99, "user")

	assert.Error(t, err)
	assert.Equal(t, ErrForbidden, err)
}

// 관리자는 타인 플러그인 수정 가능
func TestShouldAllowAdminToUpdateOthersPlugin(t *testing.T) {
	svc := newTestPluginService()
	created, _ := svc.Create(&dto.CreatePluginRequest{
		Name: "admin-edit", Description: "d", CategoryID: 1, Version: "1.0.0",
	}, 1, "user")

	newDesc := "admin edited"
	updated, err := svc.Update(created.ID, &dto.UpdatePluginRequest{
		Description: &newDesc,
	}, 999, "admin")

	require.NoError(t, err)
	assert.Equal(t, "admin edited", updated.Description)
}

// 삭제 확인
func TestShouldDeletePluginWhenOwner(t *testing.T) {
	svc := newTestPluginService()
	created, _ := svc.Create(&dto.CreatePluginRequest{
		Name: "delete-me", Description: "d", CategoryID: 1, Version: "1.0.0",
	}, 1, "user")

	err := svc.Delete(created.ID, 1, "user")
	require.NoError(t, err)

	_, err = svc.GetByID(created.ID)
	assert.Error(t, err)
}

// 타인 플러그인 삭제 시 에러
func TestShouldReturnForbiddenWhenNonOwnerDeletes(t *testing.T) {
	svc := newTestPluginService()
	created, _ := svc.Create(&dto.CreatePluginRequest{
		Name: "cant-delete", Description: "d", CategoryID: 1, Version: "1.0.0",
	}, 1, "user")

	err := svc.Delete(created.ID, 99, "user")
	assert.Equal(t, ErrForbidden, err)
}
