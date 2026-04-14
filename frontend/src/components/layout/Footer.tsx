'use client'

import React from 'react'

// 푸터 컴포넌트
export default function Footer() {
  return (
    <footer className="bg-gray-50 dark:bg-gray-900 border-t border-gray-200 dark:border-gray-700 mt-auto">
      <div className="max-w-7xl mx-auto px-4 py-8 sm:px-6 lg:px-8">
        <div className="text-center text-sm text-gray-500 dark:text-gray-400">
          <p>KTC Claude Plugin Hub</p>
          <p className="mt-1">클로드 코드 플러그인을 검색, 설치, 관리하는 통합 플랫폼</p>
        </div>
      </div>
    </footer>
  )
}
