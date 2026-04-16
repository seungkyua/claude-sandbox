'use client'

import React, { useEffect, useState } from 'react'
import { useRouter } from 'next/navigation'
import PluginList from '../components/plugin/PluginList'
import { pluginsApi } from '../services/plugins'
import type { PluginResponse } from '../types/plugin'

export default function Home() {
  const router = useRouter()
  const [officialPlugins, setOfficialPlugins] = useState<PluginResponse[]>([])
  const [communityPlugins, setCommunityPlugins] = useState<PluginResponse[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    async function fetchPlugins() {
      try {
        setLoading(true)
        const res = await pluginsApi.getList({ sort: 'popular', size: 50 })
        const plugins = res.data ?? []
        const approved = plugins.filter((p) => p.status === 'approved')
        setOfficialPlugins(approved.filter((p) => p.is_official))
        setCommunityPlugins(approved.filter((p) => !p.is_official))
      } catch {
        setError('플러그인 목록을 불러오는데 실패했습니다. 백엔드 서버가 실행 중인지 확인해주세요.')
      } finally {
        setLoading(false)
      }
    }
    fetchPlugins()
  }, [])

  const handlePluginClick = (plugin: PluginResponse) => {
    router.push(`/plugins/${plugin.id}`)
  }

  return (
    <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      {/* 히어로 섹션 */}
      <section className="text-center mb-12">
        <h1 className="text-4xl font-bold text-gray-900 dark:text-white mb-4">
          KTC Claude Plugin Hub
        </h1>
        <p className="text-lg text-gray-600 dark:text-gray-400 max-w-2xl mx-auto">
          클로드 코드 플러그인을 검색하고 설치하세요. 누구나 플러그인을 만들어 공유할 수 있습니다.
        </p>
      </section>

      {loading && (
        <div className="text-center py-20 text-gray-500 dark:text-gray-400">
          플러그인을 불러오는 중...
        </div>
      )}

      {error && (
        <div className="text-center py-20">
          <p className="text-red-500 dark:text-red-400 mb-4">{error}</p>
          <p className="text-sm text-gray-500 dark:text-gray-400">
            백엔드 서버: <code className="bg-gray-100 dark:bg-gray-800 px-2 py-1 rounded">go run ./cmd/server/</code>
          </p>
        </div>
      )}

      {!loading && !error && (
        <>
          {/* 공식 플러그인 섹션 */}
          <section className="mb-12">
            <div className="flex items-center mb-6">
              <h2 className="text-2xl font-bold text-gray-900 dark:text-white">
                공식 플러그인
              </h2>
              <span className="ml-3 px-3 py-1 text-xs font-medium bg-blue-100 text-blue-800 dark:bg-blue-900 dark:text-blue-200 rounded-full">
                Official
              </span>
            </div>
            <PluginList
              plugins={officialPlugins}
              onPluginClick={handlePluginClick}
              emptyMessage="아직 등록된 공식 플러그인이 없습니다"
            />
          </section>

          {/* 커뮤니티 플러그인 섹션 */}
          <section>
            <h2 className="text-2xl font-bold text-gray-900 dark:text-white mb-6">
              커뮤니티 플러그인
            </h2>
            <PluginList
              plugins={communityPlugins}
              onPluginClick={handlePluginClick}
              emptyMessage="아직 등록된 커뮤니티 플러그인이 없습니다"
            />
          </section>
        </>
      )}
    </div>
  )
}
