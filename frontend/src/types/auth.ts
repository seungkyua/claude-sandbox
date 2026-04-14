// 회원가입 요청
export interface RegisterRequest {
  email: string
  password: string
  nickname: string
}

// 로그인 요청
export interface LoginRequest {
  email: string
  password: string
}

// 토큰 응답
export interface TokenResponse {
  access_token: string
  refresh_token: string
  expires_in: number
}

// 사용자 정보 응답
export interface UserResponse {
  id: number
  email: string
  nickname: string
  role: 'user' | 'admin'
  created_at: string
}

// 토큰 갱신 요청
export interface RefreshRequest {
  refresh_token: string
}
