import api from './api'
import type { ReviewResponse, CreateReviewRequest, UpdateReviewRequest } from '../types/review'
import type { PaginatedResponse } from '../types/api'

// 리뷰 관련 API 서비스
export const reviewsApi = {
  // 리뷰 목록 조회
  getByPluginId: async (pluginId: number, page = 1, size = 20): Promise<PaginatedResponse<ReviewResponse>> => {
    const response = await api.get(`/plugins/${pluginId}/reviews`, { params: { page, size } })
    return response.data
  },

  // 리뷰 작성
  create: async (pluginId: number, data: CreateReviewRequest): Promise<ReviewResponse> => {
    const response = await api.post(`/plugins/${pluginId}/reviews`, data)
    return response.data
  },

  // 리뷰 수정
  update: async (pluginId: number, reviewId: number, data: UpdateReviewRequest): Promise<ReviewResponse> => {
    const response = await api.put(`/plugins/${pluginId}/reviews/${reviewId}`, data)
    return response.data
  },

  // 리뷰 삭제
  delete: async (pluginId: number, reviewId: number): Promise<void> => {
    await api.delete(`/plugins/${pluginId}/reviews/${reviewId}`)
  },
}
