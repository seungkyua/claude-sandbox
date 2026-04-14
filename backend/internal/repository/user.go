package repository

import (
	"github.com/ktc-plugin-hub/backend/internal/model"
	"gorm.io/gorm"
)

// UserRepository 는 사용자 데이터 접근 인터페이스
type UserRepository interface {
	Create(user *model.User) error
	FindByEmail(email string) (*model.User, error)
	FindByID(id uint) (*model.User, error)
}

// userRepository 는 UserRepository의 GORM 구현체
type userRepository struct {
	db *gorm.DB
}

// NewUserRepository 는 UserRepository 인스턴스를 생성한다
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

// Create 는 새 사용자를 생성한다
func (r *userRepository) Create(user *model.User) error {
	return r.db.Create(user).Error
}

// FindByEmail 은 이메일로 사용자를 조회한다
func (r *userRepository) FindByEmail(email string) (*model.User, error) {
	var user model.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByID 는 ID로 사용자를 조회한다
func (r *userRepository) FindByID(id uint) (*model.User, error) {
	var user model.User
	err := r.db.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
