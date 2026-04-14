package model

import "time"

// Category 는 플러그인 카테고리 모델
type Category struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"uniqueIndex;size:50;not null" json:"name"`
	Description string    `gorm:"size:255" json:"description"`
	SortOrder   int       `gorm:"not null;default:0" json:"sort_order"`
	CreatedAt   time.Time `gorm:"not null;autoCreateTime" json:"created_at"`
}

// TableName 은 테이블 이름을 반환한다
func (Category) TableName() string {
	return "categories"
}
