import type { Author } from './plugin'

// 리뷰 응답
export interface ReviewResponse {
  id: number
  user: Author
  rating: number
  content: string
  created_at: string
  updated_at: string
}

// 리뷰 작성 요청
export interface CreateReviewRequest {
  rating: number
  content: string
}

// 리뷰 수정 요청
export interface UpdateReviewRequest {
  rating?: number
  content?: string
}
