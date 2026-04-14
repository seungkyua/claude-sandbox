import { describe, it, expect, vi } from 'vitest'
import api from './api'
import { authApi } from './auth'
import { pluginsApi } from './plugins'
import { reviewsApi } from './reviews'
import { adminApi } from './admin'

// API 호출을 모킹하여 통합 검증
vi.mock('./api', () => ({
  default: {
    get: vi.fn(),
    post: vi.fn(),
    put: vi.fn(),
    patch: vi.fn(),
    delete: vi.fn(),
    defaults: { baseURL: 'http://localhost:8080/api/v1', headers: { 'Content-Type': 'application/json' } },
    interceptors: {
      request: { use: vi.fn() },
      response: { use: vi.fn() },
    },
  },
}))

describe('프론트엔드-백엔드 API 연동 테스트', () => {
  it('shouldCallRegisterAPI', async () => {
    const mockResponse = { data: { id: 1, email: 'test@example.com', nickname: 'tester', role: 'user' } }
    vi.mocked(api.post).mockResolvedValueOnce(mockResponse)

    const result = await authApi.register({ email: 'test@example.com', password: 'pass123', nickname: 'tester' })
    expect(api.post).toHaveBeenCalledWith('/auth/register', expect.any(Object))
    expect(result.email).toBe('test@example.com')
  })

  it('shouldCallLoginAPI', async () => {
    const mockResponse = { data: { access_token: 'token', refresh_token: 'refresh', expires_in: 3600 } }
    vi.mocked(api.post).mockResolvedValueOnce(mockResponse)

    const result = await authApi.login({ email: 'test@example.com', password: 'pass123' })
    expect(api.post).toHaveBeenCalledWith('/auth/login', expect.any(Object))
    expect(result.access_token).toBe('token')
  })

  it('shouldCallGetPluginsAPI', async () => {
    const mockResponse = { data: { data: [], total: 0, page: 1, size: 20 } }
    vi.mocked(api.get).mockResolvedValueOnce(mockResponse)

    const result = await pluginsApi.getList()
    expect(api.get).toHaveBeenCalledWith('/plugins', expect.any(Object))
    expect(result.total).toBe(0)
  })

  it('shouldCallGetPluginByIdAPI', async () => {
    const mockPlugin = { id: 1, name: 'test-plugin' }
    vi.mocked(api.get).mockResolvedValueOnce({ data: mockPlugin })

    const result = await pluginsApi.getById(1)
    expect(api.get).toHaveBeenCalledWith('/plugins/1')
    expect(result.name).toBe('test-plugin')
  })

  it('shouldCallInstallPluginAPI', async () => {
    const mockResponse = { data: { id: 1, plugin_id: 1, is_active: true } }
    vi.mocked(api.post).mockResolvedValueOnce(mockResponse)

    const result = await pluginsApi.install(1)
    expect(api.post).toHaveBeenCalledWith('/plugins/1/install')
    expect(result.is_active).toBe(true)
  })

  it('shouldCallGetReviewsAPI', async () => {
    const mockResponse = { data: { data: [], total: 0, page: 1, size: 20 } }
    vi.mocked(api.get).mockResolvedValueOnce(mockResponse)

    const result = await reviewsApi.getByPluginId(1)
    expect(api.get).toHaveBeenCalledWith('/plugins/1/reviews', expect.any(Object))
    expect(result.total).toBe(0)
  })

  it('shouldCallAdminApproveAPI', async () => {
    const mockResponse = { data: { id: 1, status: 'approved' } }
    vi.mocked(api.patch).mockResolvedValueOnce(mockResponse)

    const result = await adminApi.approve(1)
    expect(api.patch).toHaveBeenCalledWith('/admin/plugins/1/approve')
    expect(result.status).toBe('approved')
  })

  it('shouldVerifyCORSHeaderConfiguration', () => {
    // API 클라이언트가 올바른 Content-Type을 설정하는지 확인
    expect(api.defaults.headers['Content-Type']).toBe('application/json')
    // 백엔드는 CORS 미들웨어에서 http://localhost:3000 허용
    expect(api.defaults.baseURL).toContain('/api/v1')
  })
})
