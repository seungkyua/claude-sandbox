'use client'

import React from 'react'

// Card 컴포넌트 속성
interface CardProps {
  children: React.ReactNode
  className?: string
  onClick?: () => void
}

// 공통 카드 컴포넌트
export default function Card({ children, className = '', onClick }: CardProps) {
  return (
    <div
      className={`bg-white dark:bg-gray-800 rounded-lg shadow-md p-4 border border-gray-200 dark:border-gray-700 ${
        onClick ? 'cursor-pointer hover:shadow-lg transition-shadow' : ''
      } ${className}`}
      onClick={onClick}
    >
      {children}
    </div>
  )
}
