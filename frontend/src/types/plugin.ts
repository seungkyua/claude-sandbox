// 플러그인 작성자 정보
export interface Author {
  id: number
  nickname: string
}

// 카테고리 정보
export interface Category {
  id: number
  name: string
}

// 플러그인 응답
export interface PluginResponse {
  id: number
  name: string
  description: string
  author: Author
  category: Category
  status: 'pending' | 'approved' | 'rejected' | 'hidden'
  is_official: boolean
  download_count: number
  avg_rating: number
  review_count: number
  latest_version?: string
  created_at: string
  updated_at: string
}

// 플러그인 등록 요청
export interface CreatePluginRequest {
  name: string
  description: string
  category_id: number
  version: string
  changelog?: string
}

// 플러그인 수정 요청
export interface UpdatePluginRequest {
  name?: string
  description?: string
  category_id?: number
}

// 버전 정보
export interface VersionResponse {
  id: number
  plugin_id: number
  version: string
  changelog: string
  file_size: number
  created_at: string
}

// 설치 정보
export interface InstallationResponse {
  id: number
  plugin_id: number
  plugin?: { id: number; name: string; is_official: boolean }
  version_id: number
  version?: { id: number; version: string }
  is_active: boolean
  installed_at: string
}
