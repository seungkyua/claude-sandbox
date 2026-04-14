package model

import "time"

// Review 는 리뷰 모델
type Review struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	PluginID  uint      `gorm:"index;not null" json:"plugin_id"`
	UserID    uint      `gorm:"not null" json:"user_id"`
	Rating    int       `gorm:"not null" json:"rating"` // 1~5
	Content   string    `gorm:"type:text;not null" json:"content"`
	CreatedAt time.Time `gorm:"not null;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"not null;autoUpdateTime" json:"updated_at"`

	// 관계
	Plugin Plugin `gorm:"foreignKey:PluginID" json:"-"`
	User   User   `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// TableName 은 테이블 이름을 반환한다
func (Review) TableName() string {
	return "reviews"
}
