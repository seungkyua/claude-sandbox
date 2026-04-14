import { describe, it, expect } from 'vitest'
import { render, screen } from '@testing-library/react'
import Button from './Button'
import Input from './Input'
import Card from './Card'
import Badge from './Badge'
import Pagination from './Pagination'
import Modal from './Modal'

describe('Button 컴포넌트', () => {
  it('shouldRenderButtonWithText', () => {
    render(<Button>클릭</Button>)
    expect(screen.getByText('클릭')).toBeInTheDocument()
  })

  it('shouldBeDisabledWhenLoadingIsTrue', () => {
    render(<Button loading>로딩중</Button>)
    expect(screen.getByText('로딩중').closest('button')).toBeDisabled()
  })
})

describe('Input 컴포넌트', () => {
  it('shouldRenderWithLabel', () => {
    render(<Input label="이메일" placeholder="입력" />)
    expect(screen.getByText('이메일')).toBeInTheDocument()
  })

  it('shouldRenderErrorMessage', () => {
    render(<Input error="필수 항목입니다" />)
    expect(screen.getByText('필수 항목입니다')).toBeInTheDocument()
  })
})

describe('Card 컴포넌트', () => {
  it('shouldRenderChildren', () => {
    render(<Card>카드 내용</Card>)
    expect(screen.getByText('카드 내용')).toBeInTheDocument()
  })
})

describe('Badge 컴포넌트', () => {
  it('shouldRenderBadgeText', () => {
    render(<Badge variant="official">공식</Badge>)
    expect(screen.getByText('공식')).toBeInTheDocument()
  })
})

describe('Pagination 컴포넌트', () => {
  it('shouldRenderPageNumbers', () => {
    render(<Pagination currentPage={1} totalPages={5} onPageChange={() => {}} />)
    expect(screen.getByText('1')).toBeInTheDocument()
    expect(screen.getByText('5')).toBeInTheDocument()
  })

  it('shouldNotRenderWhenSinglePage', () => {
    const { container } = render(<Pagination currentPage={1} totalPages={1} onPageChange={() => {}} />)
    expect(container.innerHTML).toBe('')
  })
})

describe('Modal 컴포넌트', () => {
  it('shouldRenderWhenOpen', () => {
    render(<Modal isOpen={true} onClose={() => {}} title="테스트 모달">내용</Modal>)
    expect(screen.getByText('테스트 모달')).toBeInTheDocument()
    expect(screen.getByText('내용')).toBeInTheDocument()
  })

  it('shouldNotRenderWhenClosed', () => {
    const { container } = render(<Modal isOpen={false} onClose={() => {}} title="숨김">내용</Modal>)
    expect(container.innerHTML).toBe('')
  })
})
