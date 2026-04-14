package model

import "time"

// 플러그인 상태 상수
const (
	PluginStatusPending  = "pending"
	PluginStatusApproved = "approved"
	PluginStatusRejected = "rejected"
	PluginStatusHidden   = "hidden"
)

// Plugin 은 플러그인 모델
type Plugin struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	AuthorID      uint      `gorm:"index;not null" json:"author_id"`
	CategoryID    uint      `gorm:"index;not null" json:"category_id"`
	Name          string    `gorm:"uniqueIndex;size:100;not null" json:"name"`
	Description   string    `gorm:"type:text;not null" json:"description"`
	Status        string    `gorm:"size:20;not null;default:pending;index" json:"status"`
	IsOfficial    bool      `gorm:"not null;default:false" json:"is_official"`
	DownloadCount int       `gorm:"not null;default:0" json:"download_count"`
	AvgRating     float64   `gorm:"type:decimal(2,1);not null;default:0.0" json:"avg_rating"`
	ReviewCount   int       `gorm:"not null;default:0" json:"review_count"`
	CreatedAt     time.Time `gorm:"not null;autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time `gorm:"not null;autoUpdateTime" json:"updated_at"`

	// 관계
	Author   User     `gorm:"foreignKey:AuthorID" json:"author,omitempty"`
	Category Category `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
}

// TableName 은 테이블 이름을 반환한다
func (Plugin) TableName() string {
	return "plugins"
}
