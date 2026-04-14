import { describe, it, expect, vi } from 'vitest'
import { render, screen } from '@testing-library/react'
import LoginForm from './LoginForm'
import RegisterForm from './RegisterForm'

describe('LoginForm', () => {
  it('shouldRenderEmailAndPasswordInputs', () => {
    render(<LoginForm onSubmit={vi.fn()} />)
    expect(screen.getByText('이메일')).toBeInTheDocument()
    expect(screen.getByText('비밀번호')).toBeInTheDocument()
    expect(screen.getByText('로그인')).toBeInTheDocument()
  })

  it('shouldRenderErrorMessage', () => {
    render(<LoginForm onSubmit={vi.fn()} error="로그인 실패" />)
    expect(screen.getByText('로그인 실패')).toBeInTheDocument()
  })
})

describe('RegisterForm', () => {
  it('shouldRenderAllInputFields', () => {
    render(<RegisterForm onSubmit={vi.fn()} />)
    expect(screen.getByText('이메일')).toBeInTheDocument()
    expect(screen.getByText('닉네임')).toBeInTheDocument()
    expect(screen.getByText('비밀번호')).toBeInTheDocument()
    expect(screen.getByText('회원가입')).toBeInTheDocument()
  })
})
