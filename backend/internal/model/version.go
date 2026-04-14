package model

import "time"

// PluginVersion 은 플러그인 버전 모델
type PluginVersion struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	PluginID  uint      `gorm:"index;not null" json:"plugin_id"`
	Version   string    `gorm:"size:20;not null" json:"version"`
	Changelog string    `gorm:"type:text" json:"changelog"`
	FilePath  string    `gorm:"size:500;not null" json:"file_path"`
	FileSize  int64     `gorm:"not null" json:"file_size"`
	CreatedAt time.Time `gorm:"not null;autoCreateTime" json:"created_at"`

	// 관계
	Plugin Plugin `gorm:"foreignKey:PluginID" json:"-"`
}

// TableName 은 테이블 이름을 반환한다
func (PluginVersion) TableName() string {
	return "plugin_versions"
}
