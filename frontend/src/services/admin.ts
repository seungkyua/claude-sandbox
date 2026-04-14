import api from './api'
import type { PluginResponse } from '../types/plugin'

// 관리자 관련 API 서비스
export const adminApi = {
  // 심사 대기 플러그인 목록
  getPendingPlugins: async (): Promise<{ data: PluginResponse[] }> => {
    const response = await api.get('/admin/plugins/pending')
    return response.data
  },

  // 플러그인 승인
  approve: async (id: number): Promise<PluginResponse> => {
    const response = await api.patch(`/admin/plugins/${id}/approve`)
    return response.data
  },

  // 플러그인 반려
  reject: async (id: number, reason: string): Promise<PluginResponse> => {
    const response = await api.patch(`/admin/plugins/${id}/reject`, { reason })
    return response.data
  },

  // 플러그인 비공개 처리
  hide: async (id: number): Promise<PluginResponse> => {
    const response = await api.patch(`/admin/plugins/${id}/hide`)
    return response.data
  },
}
