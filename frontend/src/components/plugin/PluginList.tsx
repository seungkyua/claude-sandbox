'use client'

import React from 'react'
import PluginCard from './PluginCard'
import type { PluginResponse } from '../../types/plugin'

interface PluginListProps {
  plugins: PluginResponse[]
  onPluginClick?: (plugin: PluginResponse) => void
  emptyMessage?: string
}

// 플러그인 목록 컴포넌트
export default function PluginList({ plugins, onPluginClick, emptyMessage = '플러그인이 없습니다' }: PluginListProps) {
  if (plugins.length === 0) {
    return (
      <div className="text-center py-12 text-gray-500 dark:text-gray-400">
        {emptyMessage}
      </div>
    )
  }

  return (
    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
      {plugins.map((plugin) => (
        <PluginCard
          key={plugin.id}
          plugin={plugin}
          onClick={() => onPluginClick?.(plugin)}
        />
      ))}
    </div>
  )
}
