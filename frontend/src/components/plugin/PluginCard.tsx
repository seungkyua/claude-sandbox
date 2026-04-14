'use client'

import React from 'react'
import Card from '../common/Card'
import Badge from '../common/Badge'
import type { PluginResponse } from '../../types/plugin'

interface PluginCardProps {
  plugin: PluginResponse
  onClick?: () => void
}

// 플러그인 카드 컴포넌트
export default function PluginCard({ plugin, onClick }: PluginCardProps) {
  return (
    <Card onClick={onClick} className="hover:border-blue-300 dark:hover:border-blue-600">
      <div className="flex items-start justify-between">
        <div className="flex-1 min-w-0">
          <div className="flex items-center space-x-2">
            <h3 className="text-lg font-semibold text-gray-900 dark:text-white truncate">
              {plugin.name}
            </h3>
            {plugin.is_official && <Badge variant="official">공식</Badge>}
          </div>
          <p className="mt-1 text-sm text-gray-500 dark:text-gray-400 line-clamp-2">
            {plugin.description}
          </p>
        </div>
      </div>
      <div className="mt-3 flex items-center justify-between text-sm text-gray-500 dark:text-gray-400">
        <span>{plugin.author.nickname}</span>
        <div className="flex items-center space-x-3">
          <span>{'★'} {plugin.avg_rating.toFixed(1)}</span>
          <span>{plugin.download_count.toLocaleString()} 다운로드</span>
        </div>
      </div>
    </Card>
  )
}
