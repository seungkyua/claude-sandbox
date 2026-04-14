'use client'

import React from 'react'
import Link from 'next/link'
import { useAuthStore } from '../../stores/authStore'
import { useThemeStore } from '../../stores/themeStore'
import Button from '../common/Button'

// 헤더 컴포넌트 (로고, 검색바, 로그인/프로필, 다크모드 토글)
export default function Header() {
  const { user, isAuthenticated, logout } = useAuthStore()
  const { isDark, toggle } = useThemeStore()

  return (
    <header className="sticky top-0 z-40 bg-white dark:bg-gray-900 border-b border-gray-200 dark:border-gray-700">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="flex items-center justify-between h-16">
          {/* 로고 */}
          <Link href="/" className="text-xl font-bold text-blue-600 dark:text-blue-400">
            Plugin Hub
          </Link>

          {/* 검색바 */}
          <div className="hidden md:flex flex-1 max-w-md mx-8">
            <input
              type="text"
              placeholder="플러그인 검색..."
              className="w-full px-4 py-2 rounded-lg border border-gray-300 dark:border-gray-600 dark:bg-gray-800 dark:text-white focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
          </div>

          {/* 우측 메뉴 */}
          <div className="flex items-center space-x-4">
            {/* 다크모드 토글 */}
            <button
              onClick={toggle}
              className="p-2 rounded-lg hover:bg-gray-100 dark:hover:bg-gray-800"
              aria-label="다크모드 토글"
            >
              {isDark ? '☀' : '☾'}
            </button>

            {isAuthenticated && user ? (
              <div className="flex items-center space-x-3">
                <Link href="/dashboard" className="text-sm text-gray-600 dark:text-gray-300 hover:text-blue-600">
                  대시보드
                </Link>
                {user.role === 'admin' && (
                  <Link href="/admin" className="text-sm text-gray-600 dark:text-gray-300 hover:text-blue-600">
                    관리자
                  </Link>
                )}
                <span className="text-sm text-gray-600 dark:text-gray-300">{user.nickname}</span>
                <Button variant="ghost" size="sm" onClick={logout}>
                  로그아웃
                </Button>
              </div>
            ) : (
              <div className="flex items-center space-x-2">
                <Link href="/login">
                  <Button variant="ghost" size="sm">로그인</Button>
                </Link>
                <Link href="/register">
                  <Button size="sm">회원가입</Button>
                </Link>
              </div>
            )}
          </div>
        </div>
      </div>
    </header>
  )
}
