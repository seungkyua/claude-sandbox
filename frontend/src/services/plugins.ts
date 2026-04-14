import api from './api'
import type { PluginResponse, CreatePluginRequest, UpdatePluginRequest, InstallationResponse } from '../types/plugin'
import type { PaginatedResponse } from '../types/api'

// 플러그인 관련 API 서비스
export const pluginsApi = {
  // 플러그인 목록 조회
  getList: async (params?: Record<string, string | number>): Promise<PaginatedResponse<PluginResponse>> => {
    const response = await api.get('/plugins', { params })
    return response.data
  },

  // 플러그인 상세 조회
  getById: async (id: number): Promise<PluginResponse> => {
    const response = await api.get(`/plugins/${id}`)
    return response.data
  },

  // 플러그인 등록
  create: async (data: CreatePluginRequest): Promise<PluginResponse> => {
    const response = await api.post('/plugins', data)
    return response.data
  },

  // 플러그인 수정
  update: async (id: number, data: UpdatePluginRequest): Promise<PluginResponse> => {
    const response = await api.put(`/plugins/${id}`, data)
    return response.data
  },

  // 플러그인 삭제
  delete: async (id: number): Promise<void> => {
    await api.delete(`/plugins/${id}`)
  },

  // 플러그인 설치
  install: async (id: number): Promise<InstallationResponse> => {
    const response = await api.post(`/plugins/${id}/install`)
    return response.data
  },

  // 플러그인 삭제 (내 설치에서)
  uninstall: async (id: number): Promise<void> => {
    await api.delete(`/plugins/${id}/install`)
  },

  // 활성화/비활성화 토글
  toggleActive: async (id: number, isActive: boolean): Promise<InstallationResponse> => {
    const response = await api.patch(`/plugins/${id}/install`, { is_active: isActive })
    return response.data
  },

  // 내 설치 목록
  getMyInstallations: async (): Promise<{ data: InstallationResponse[] }> => {
    const response = await api.get('/me/installations')
    return response.data
  },
}
