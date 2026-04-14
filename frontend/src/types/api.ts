// 페이지네이션 응답
export interface PaginatedResponse<T> {
  data: T[]
  total: number
  page: number
  size: number
}

// RFC 7807 에러 응답
export interface ErrorResponse {
  type: string
  title: string
  status: number
  detail: string
}
