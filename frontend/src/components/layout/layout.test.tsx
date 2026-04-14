import { describe, it, expect, vi } from 'vitest'
import { render, screen } from '@testing-library/react'
import Footer from './Footer'

// Next.js의 Link를 모킹
vi.mock('next/link', () => ({
  default: ({ children, href }: { children: React.ReactNode; href: string }) => (
    <a href={href}>{children}</a>
  ),
}))

describe('Footer', () => {
  it('shouldRenderFooterText', () => {
    render(<Footer />)
    expect(screen.getByText('KTC Claude Plugin Hub')).toBeInTheDocument()
  })
})
