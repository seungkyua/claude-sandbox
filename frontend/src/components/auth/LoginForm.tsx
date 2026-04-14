'use client'

import React, { useState } from 'react'
import Button from '../common/Button'
import Input from '../common/Input'

interface LoginFormProps {
  onSubmit: (email: string, password: string) => Promise<void>
  error?: string
}

// 로그인 폼 컴포넌트
export default function LoginForm({ onSubmit, error }: LoginFormProps) {
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [loading, setLoading] = useState(false)

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setLoading(true)
    try {
      await onSubmit(email, password)
    } finally {
      setLoading(false)
    }
  }

  return (
    <form onSubmit={handleSubmit} className="space-y-4">
      <Input
        label="이메일"
        type="email"
        value={email}
        onChange={(e) => setEmail(e.target.value)}
        placeholder="이메일을 입력하세요"
        required
      />
      <Input
        label="비밀번호"
        type="password"
        value={password}
        onChange={(e) => setPassword(e.target.value)}
        placeholder="비밀번호를 입력하세요"
        required
      />
      {error && <p className="text-sm text-red-500">{error}</p>}
      <Button type="submit" loading={loading} className="w-full">
        로그인
      </Button>
    </form>
  )
}
