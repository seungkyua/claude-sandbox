package model

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

// User 모델의 테이블 이름이 올바른지 확인
func TestShouldReturnCorrectTableNameForUser(t *testing.T) {
	u := User{}
	assert.Equal(t, "users", u.TableName())
}

// User의 IsAdmin 메서드가 올바르게 동작하는지 확인
func TestShouldReturnTrueWhenUserIsAdmin(t *testing.T) {
	admin := &User{Role: "admin"}
	user := &User{Role: "user"}
	assert.True(t, admin.IsAdmin())
	assert.False(t, user.IsAdmin())
}

// User의 PasswordHash가 JSON 직렬화에서 제외되는지 확인
func TestShouldExcludePasswordHashFromJSON(t *testing.T) {
	u := User{
		ID:           1,
		Email:        "test@example.com",
		PasswordHash: "secret-hash",
		Nickname:     "tester",
		Role:         "user",
	}
	data, err := json.Marshal(u)
	assert.NoError(t, err)
	assert.NotContains(t, string(data), "secret-hash")
	assert.Contains(t, string(data), "test@example.com")
}

// 각 모델의 테이블 이름이 올바른지 확인
func TestShouldReturnCorrectTableNames(t *testing.T) {
	assert.Equal(t, "categories", Category{}.TableName())
	assert.Equal(t, "plugins", Plugin{}.TableName())
	assert.Equal(t, "plugin_versions", PluginVersion{}.TableName())
	assert.Equal(t, "reviews", Review{}.TableName())
	assert.Equal(t, "installations", Installation{}.TableName())
}

// 플러그인 상태 상수가 올바른지 확인
func TestShouldHaveCorrectPluginStatusConstants(t *testing.T) {
	assert.Equal(t, "pending", PluginStatusPending)
	assert.Equal(t, "approved", PluginStatusApproved)
	assert.Equal(t, "rejected", PluginStatusRejected)
	assert.Equal(t, "hidden", PluginStatusHidden)
}
