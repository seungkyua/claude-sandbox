package dto

import "time"

// CreatePluginRequest 는 플러그인 등록 요청 DTO
type CreatePluginRequest struct {
	Name        string `json:"name" binding:"required,min=1,max=100"`
	Description string `json:"description" binding:"required"`
	CategoryID  uint   `json:"category_id" binding:"required"`
	Version     string `json:"version" binding:"required"`
	Changelog   string `json:"changelog"`
}

// UpdatePluginRequest 는 플러그인 수정 요청 DTO
type UpdatePluginRequest struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
	CategoryID  *uint   `json:"category_id"`
}

// PluginListRequest 는 플러그인 목록 조회 요청 DTO
type PluginListRequest struct {
	PaginationRequest
	CategoryID *uint  `form:"category_id"`
	Keyword    string `form:"keyword"`
	Sort       string `form:"sort"` // popular, latest, rating
}

// PluginResponse 는 플러그인 응답 DTO
type PluginResponse struct {
	ID            uint             `json:"id"`
	Name          string           `json:"name"`
	Description   string           `json:"description"`
	Author        AuthorResponse   `json:"author"`
	Category      CategoryResponse `json:"category"`
	Status        string           `json:"status"`
	IsOfficial    bool             `json:"is_official"`
	DownloadCount int              `json:"download_count"`
	AvgRating     float64          `json:"avg_rating"`
	ReviewCount   int              `json:"review_count"`
	LatestVersion string           `json:"latest_version,omitempty"`
	CreatedAt     time.Time        `json:"created_at"`
	UpdatedAt     time.Time        `json:"updated_at"`
}

// AuthorResponse 는 작성자 정보 응답 DTO
type AuthorResponse struct {
	ID       uint   `json:"id"`
	Nickname string `json:"nickname"`
}

// CategoryResponse 는 카테고리 정보 응답 DTO
type CategoryResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

// VersionResponse 는 버전 응답 DTO
type VersionResponse struct {
	ID        uint      `json:"id"`
	PluginID  uint      `json:"plugin_id"`
	Version   string    `json:"version"`
	Changelog string    `json:"changelog"`
	FileSize  int64     `json:"file_size"`
	CreatedAt time.Time `json:"created_at"`
}

// CreateVersionRequest 는 버전 등록 요청 DTO
type CreateVersionRequest struct {
	Version   string `json:"version" binding:"required"`
	Changelog string `json:"changelog"`
}

// InstallRequest 는 플러그인 설치 요청 DTO
type InstallRequest struct {
	VersionID *uint `json:"version_id"`
}

// InstallationResponse 는 설치 정보 응답 DTO
type InstallationResponse struct {
	ID          uint           `json:"id"`
	PluginID    uint           `json:"plugin_id"`
	Plugin      *PluginBrief   `json:"plugin,omitempty"`
	VersionID   uint           `json:"version_id"`
	Version     *VersionBrief  `json:"version,omitempty"`
	IsActive    bool           `json:"is_active"`
	InstalledAt time.Time      `json:"installed_at"`
}

// PluginBrief 는 플러그인 간략 정보
type PluginBrief struct {
	ID         uint   `json:"id"`
	Name       string `json:"name"`
	IsOfficial bool   `json:"is_official"`
}

// VersionBrief 는 버전 간략 정보
type VersionBrief struct {
	ID      uint   `json:"id"`
	Version string `json:"version"`
}

// ToggleActiveRequest 는 활성화/비활성화 요청 DTO
type ToggleActiveRequest struct {
	IsActive bool `json:"is_active"`
}
