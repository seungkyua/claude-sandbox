import { describe, it, expect, beforeEach } from 'vitest'
import { useAuthStore } from './authStore'
import { useThemeStore } from './themeStore'

describe('AuthStore', () => {
  beforeEach(() => {
    useAuthStore.setState({
      user: null,
      accessToken: null,
      isAuthenticated: false,
    })
  })

  it('shouldSetAuthWhenLogin', () => {
    const store = useAuthStore.getState()
    store.setAuth(
      { id: 1, email: 'test@example.com', nickname: 'tester', role: 'user', created_at: '' },
      'access-token',
      'refresh-token'
    )

    const state = useAuthStore.getState()
    expect(state.isAuthenticated).toBe(true)
    expect(state.user?.email).toBe('test@example.com')
    expect(state.accessToken).toBe('access-token')
  })

  it('shouldClearAuthWhenLogout', () => {
    const store = useAuthStore.getState()
    store.setAuth(
      { id: 1, email: 'test@example.com', nickname: 'tester', role: 'user', created_at: '' },
      'token',
      'refresh'
    )
    store.logout()

    const state = useAuthStore.getState()
    expect(state.isAuthenticated).toBe(false)
    expect(state.user).toBeNull()
    expect(state.accessToken).toBeNull()
  })
})

describe('ThemeStore', () => {
  beforeEach(() => {
    useThemeStore.setState({ isDark: false })
  })

  it('shouldToggleTheme', () => {
    const store = useThemeStore.getState()
    expect(store.isDark).toBe(false)

    store.toggle()
    expect(useThemeStore.getState().isDark).toBe(true)

    store.toggle()
    expect(useThemeStore.getState().isDark).toBe(false)
  })
})
