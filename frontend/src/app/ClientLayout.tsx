'use client'

import React, { useEffect } from 'react'
import Header from '../components/layout/Header'
import Footer from '../components/layout/Footer'
import { useAuthStore } from '../stores/authStore'
import { useThemeStore } from '../stores/themeStore'

export default function ClientLayout({ children }: { children: React.ReactNode }) {
  const loadAuth = useAuthStore((s) => s.loadFromStorage)
  const loadTheme = useThemeStore((s) => s.loadFromStorage)

  useEffect(() => {
    loadAuth()
    loadTheme()
  }, [loadAuth, loadTheme])

  return (
    <>
      <Header />
      <main className="flex-1">{children}</main>
      <Footer />
    </>
  )
}
