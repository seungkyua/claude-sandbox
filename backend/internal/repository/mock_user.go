package repository

import (
	"errors"

	"github.com/ktc-plugin-hub/backend/internal/model"
)

// 공통 에러 정의
var (
	ErrNotFound      = errors.New("레코드를 찾을 수 없습니다")
	ErrDuplicateKey  = errors.New("중복된 키입니다")
)

// MockUserRepository 는 테스트용 사용자 리포지토리 모킹
type MockUserRepository struct {
	users  map[uint]*model.User
	nextID uint
}

// NewMockUserRepository 는 MockUserRepository를 생성한다
func NewMockUserRepository() *MockUserRepository {
	return &MockUserRepository{
		users:  make(map[uint]*model.User),
		nextID: 1,
	}
}

// Create 는 사용자를 생성한다 (이메일 중복 검증 포함)
func (r *MockUserRepository) Create(user *model.User) error {
	for _, u := range r.users {
		if u.Email == user.Email {
			return ErrDuplicateKey
		}
	}
	user.ID = r.nextID
	r.nextID++
	r.users[user.ID] = user
	return nil
}

// FindByEmail 은 이메일로 사용자를 조회한다
func (r *MockUserRepository) FindByEmail(email string) (*model.User, error) {
	for _, u := range r.users {
		if u.Email == email {
			return u, nil
		}
	}
	return nil, ErrNotFound
}

// FindByID 는 ID로 사용자를 조회한다
func (r *MockUserRepository) FindByID(id uint) (*model.User, error) {
	u, ok := r.users[id]
	if !ok {
		return nil, ErrNotFound
	}
	return u, nil
}
