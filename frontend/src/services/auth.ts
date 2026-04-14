import api from './api'
import type { RegisterRequest, LoginRequest, TokenResponse, UserResponse } from '../types/auth'

// 인증 관련 API 서비스
export const authApi = {
  // 회원가입
  register: async (data: RegisterRequest): Promise<UserResponse> => {
    const response = await api.post('/auth/register', data)
    return response.data
  },

  // 로그인
  login: async (data: LoginRequest): Promise<TokenResponse> => {
    const response = await api.post('/auth/login', data)
    return response.data
  },

  // 토큰 갱신
  refresh: async (refreshToken: string): Promise<TokenResponse> => {
    const response = await api.post('/auth/refresh', { refresh_token: refreshToken })
    return response.data
  },

  // 내 정보 조회
  getMe: async (): Promise<UserResponse> => {
    const response = await api.get('/me')
    return response.data
  },
}
