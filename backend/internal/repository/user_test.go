package repository

import (
	"testing"

	"github.com/ktc-plugin-hub/backend/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// 사용자 생성 및 ID 조회가 올바르게 동작하는지 확인
func TestShouldCreateUserAndFindByID(t *testing.T) {
	repo := NewMockUserRepository()
	user := &model.User{
		Email:        "test@example.com",
		PasswordHash: "hashed",
		Nickname:     "tester",
		Role:         "user",
	}

	err := repo.Create(user)
	require.NoError(t, err)
	assert.NotZero(t, user.ID)

	found, err := repo.FindByID(user.ID)
	require.NoError(t, err)
	assert.Equal(t, "test@example.com", found.Email)
}

// 이메일로 사용자를 조회할 수 있는지 확인
func TestShouldFindUserByEmail(t *testing.T) {
	repo := NewMockUserRepository()
	user := &model.User{
		Email:        "find@example.com",
		PasswordHash: "hashed",
		Nickname:     "finder",
		Role:         "user",
	}
	err := repo.Create(user)
	require.NoError(t, err)

	found, err := repo.FindByEmail("find@example.com")
	require.NoError(t, err)
	assert.Equal(t, "finder", found.Nickname)
}

// 존재하지 않는 사용자 조회 시 에러를 반환하는지 확인
func TestShouldReturnErrorWhenUserNotFound(t *testing.T) {
	repo := NewMockUserRepository()

	_, err := repo.FindByID(999)
	assert.Error(t, err)
	assert.Equal(t, ErrNotFound, err)

	_, err = repo.FindByEmail("nonexistent@example.com")
	assert.Error(t, err)
}

// 중복 이메일로 생성 시 에러를 반환하는지 확인
func TestShouldReturnErrorWhenDuplicateEmail(t *testing.T) {
	repo := NewMockUserRepository()
	user1 := &model.User{Email: "dup@example.com", PasswordHash: "h", Nickname: "a", Role: "user"}
	user2 := &model.User{Email: "dup@example.com", PasswordHash: "h", Nickname: "b", Role: "user"}

	err := repo.Create(user1)
	require.NoError(t, err)

	err = repo.Create(user2)
	assert.Error(t, err)
	assert.Equal(t, ErrDuplicateKey, err)
}

// MockUserRepository가 UserRepository 인터페이스를 구현하는지 확인
func TestMockUserRepositoryShouldImplementInterface(t *testing.T) {
	var _ UserRepository = NewMockUserRepository()
}
