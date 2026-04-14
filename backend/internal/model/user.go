package model

import "time"

// User 는 사용자 모델
type User struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	Email        string    `gorm:"uniqueIndex;size:255;not null" json:"email"`
	PasswordHash string    `gorm:"size:255;not null" json:"-"` // JSON 응답에서 제외
	Nickname     string    `gorm:"size:50;not null" json:"nickname"`
	Role         string    `gorm:"size:20;not null;default:user" json:"role"`
	CreatedAt    time.Time `gorm:"not null;autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time `gorm:"not null;autoUpdateTime" json:"updated_at"`
}

// TableName 은 테이블 이름을 반환한다
func (User) TableName() string {
	return "users"
}

// IsAdmin 은 관리자 여부를 반환한다
func (u *User) IsAdmin() bool {
	return u.Role == "admin"
}
