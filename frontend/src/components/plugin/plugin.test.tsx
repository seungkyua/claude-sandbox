import { describe, it, expect, vi } from 'vitest'
import { render, screen } from '@testing-library/react'
import PluginCard from './PluginCard'
import PluginList from './PluginList'
import ReviewItem from './ReviewItem'
import type { PluginResponse } from '../../types/plugin'
import type { ReviewResponse } from '../../types/review'

const mockPlugin: PluginResponse = {
  id: 1,
  name: 'test-plugin',
  description: 'A test plugin description',
  author: { id: 1, nickname: 'tester' },
  category: { id: 1, name: '개발 도구' },
  status: 'approved',
  is_official: true,
  download_count: 1234,
  avg_rating: 4.5,
  review_count: 10,
  created_at: '2026-01-01',
  updated_at: '2026-01-01',
}

const mockReview: ReviewResponse = {
  id: 1,
  user: { id: 1, nickname: 'reviewer' },
  rating: 5,
  content: 'Great plugin!',
  created_at: '2026-01-01T00:00:00Z',
  updated_at: '2026-01-01T00:00:00Z',
}

describe('PluginCard', () => {
  it('shouldRenderPluginName', () => {
    render(<PluginCard plugin={mockPlugin} />)
    expect(screen.getByText('test-plugin')).toBeInTheDocument()
  })

  it('shouldRenderOfficialBadge', () => {
    render(<PluginCard plugin={mockPlugin} />)
    expect(screen.getByText('공식')).toBeInTheDocument()
  })

  it('shouldRenderDownloadCount', () => {
    render(<PluginCard plugin={mockPlugin} />)
    expect(screen.getByText('1,234 다운로드')).toBeInTheDocument()
  })
})

describe('PluginList', () => {
  it('shouldRenderPluginCards', () => {
    render(<PluginList plugins={[mockPlugin]} />)
    expect(screen.getByText('test-plugin')).toBeInTheDocument()
  })

  it('shouldRenderEmptyMessageWhenNoPlugins', () => {
    render(<PluginList plugins={[]} emptyMessage="결과 없음" />)
    expect(screen.getByText('결과 없음')).toBeInTheDocument()
  })
})

describe('ReviewItem', () => {
  it('shouldRenderReviewContent', () => {
    render(<ReviewItem review={mockReview} />)
    expect(screen.getByText('Great plugin!')).toBeInTheDocument()
    expect(screen.getByText('reviewer')).toBeInTheDocument()
  })
})
