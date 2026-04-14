import { create } from 'zustand'
import type { UserResponse } from '../types/auth'

// 인증 상태 인터페이스
interface AuthState {
  user: UserResponse | null
  accessToken: string | null
  isAuthenticated: boolean
  setAuth: (user: UserResponse, accessToken: string, refreshToken: string) => void
  logout: () => void
  loadFromStorage: () => void
}

// 인증 스토어 (Zustand)
export const useAuthStore = create<AuthState>((set) => ({
  user: null,
  accessToken: null,
  isAuthenticated: false,

  // 로그인 성공 시 상태 업데이트 및 토큰 저장
  setAuth: (user, accessToken, refreshToken) => {
    if (typeof window !== 'undefined') {
      localStorage.setItem('access_token', accessToken)
      localStorage.setItem('refresh_token', refreshToken)
      localStorage.setItem('user', JSON.stringify(user))
    }
    set({ user, accessToken, isAuthenticated: true })
  },

  // 로그아웃 시 상태 초기화 및 토큰 제거
  logout: () => {
    if (typeof window !== 'undefined') {
      localStorage.removeItem('access_token')
      localStorage.removeItem('refresh_token')
      localStorage.removeItem('user')
    }
    set({ user: null, accessToken: null, isAuthenticated: false })
  },

  // 새로고침 시 localStorage에서 복원
  loadFromStorage: () => {
    if (typeof window !== 'undefined') {
      const token = localStorage.getItem('access_token')
      const userStr = localStorage.getItem('user')
      if (token && userStr) {
        try {
          const user = JSON.parse(userStr) as UserResponse
          set({ user, accessToken: token, isAuthenticated: true })
        } catch {
          // 파싱 실패 시 무시
        }
      }
    }
  },
}))
