import { describe, it, expect } from 'vitest'
import type { UserResponse, TokenResponse } from './auth'
import type { PluginResponse } from './plugin'
import type { ReviewResponse } from './review'
import type { PaginatedResponse, ErrorResponse } from './api'

describe('타입 정의', () => {
  it('shouldDefineUserResponseType', () => {
    // 타입이 컴파일되면 성공
    const user: UserResponse = {
      id: 1,
      email: 'test@example.com',
      nickname: 'tester',
      role: 'user',
      created_at: '2026-01-01T00:00:00Z',
    }
    expect(user.id).toBe(1)
  })

  it('shouldDefinePluginResponseType', () => {
    const plugin: PluginResponse = {
      id: 1,
      name: 'test',
      description: 'desc',
      author: { id: 1, nickname: 'author' },
      category: { id: 1, name: 'dev' },
      status: 'approved',
      is_official: false,
      download_count: 0,
      avg_rating: 0,
      review_count: 0,
      created_at: '',
      updated_at: '',
    }
    expect(plugin.name).toBe('test')
  })

  it('shouldDefinePaginatedResponseType', () => {
    const resp: PaginatedResponse<PluginResponse> = {
      data: [],
      total: 0,
      page: 1,
      size: 20,
    }
    expect(resp.total).toBe(0)
  })

  it('shouldDefineErrorResponseType', () => {
    const err: ErrorResponse = {
      type: 'INVALID_INPUT',
      title: 'error',
      status: 400,
      detail: 'details',
    }
    expect(err.status).toBe(400)
  })
})
