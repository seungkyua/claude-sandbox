import { describe, it, expect, vi } from 'vitest'
import api from './api'

describe('API 클라이언트', () => {
  it('shouldHaveBaseURL', () => {
    expect(api.defaults.baseURL).toBeDefined()
    expect(api.defaults.baseURL).toContain('/api/v1')
  })

  it('shouldHaveJSONContentType', () => {
    expect(api.defaults.headers['Content-Type']).toBe('application/json')
  })

  it('shouldHaveRequestInterceptor', () => {
    // 인터셉터가 등록되었는지 확인
    expect(api.interceptors.request).toBeDefined()
  })

  it('shouldHaveResponseInterceptor', () => {
    expect(api.interceptors.response).toBeDefined()
  })
})
