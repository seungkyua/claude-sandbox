import { create } from 'zustand'

// 테마 상태 인터페이스
interface ThemeState {
  isDark: boolean
  toggle: () => void
  loadFromStorage: () => void
}

// 테마 스토어 (Zustand)
export const useThemeStore = create<ThemeState>((set) => ({
  isDark: false,

  // 다크/라이트 모드 토글
  toggle: () => {
    set((state) => {
      const newIsDark = !state.isDark
      if (typeof window !== 'undefined') {
        localStorage.setItem('theme', newIsDark ? 'dark' : 'light')
        document.documentElement.classList.toggle('dark', newIsDark)
      }
      return { isDark: newIsDark }
    })
  },

  // localStorage에서 테마 복원
  loadFromStorage: () => {
    if (typeof window !== 'undefined') {
      const theme = localStorage.getItem('theme')
      const isDark = theme === 'dark'
      document.documentElement.classList.toggle('dark', isDark)
      set({ isDark })
    }
  },
}))
