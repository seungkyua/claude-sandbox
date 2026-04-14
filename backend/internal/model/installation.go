package model

import "time"

// Installation 은 플러그인 설치 모델
type Installation struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	UserID      uint      `gorm:"index;not null" json:"user_id"`
	PluginID    uint      `gorm:"index;not null" json:"plugin_id"`
	VersionID   uint      `gorm:"not null" json:"version_id"`
	IsActive    bool      `gorm:"not null;default:true" json:"is_active"`
	InstalledAt time.Time `gorm:"not null;autoCreateTime" json:"installed_at"`
	UpdatedAt   time.Time `gorm:"not null;autoUpdateTime" json:"updated_at"`

	// 관계
	User    User          `gorm:"foreignKey:UserID" json:"-"`
	Plugin  Plugin        `gorm:"foreignKey:PluginID" json:"plugin,omitempty"`
	Version PluginVersion `gorm:"foreignKey:VersionID" json:"version,omitempty"`
}

// TableName 은 테이블 이름을 반환한다
func (Installation) TableName() string {
	return "installations"
}
