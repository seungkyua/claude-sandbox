'use client'

import React from 'react'
import type { ReviewResponse } from '../../types/review'

interface ReviewItemProps {
  review: ReviewResponse
}

// 리뷰 아이템 컴포넌트
export default function ReviewItem({ review }: ReviewItemProps) {
  return (
    <div className="border-b border-gray-200 dark:border-gray-700 py-4 last:border-0">
      <div className="flex items-center justify-between">
        <div className="flex items-center space-x-2">
          <span className="font-medium text-gray-900 dark:text-white">
            {review.user.nickname}
          </span>
          <span className="text-yellow-500">
            {'★'.repeat(review.rating)}{'☆'.repeat(5 - review.rating)}
          </span>
        </div>
        <span className="text-sm text-gray-500">
          {new Date(review.created_at).toLocaleDateString('ko-KR')}
        </span>
      </div>
      <p className="mt-2 text-gray-700 dark:text-gray-300">{review.content}</p>
    </div>
  )
}
