'use client'

import React, { useState } from 'react'
import Button from '../common/Button'
import Input from '../common/Input'

interface RegisterFormProps {
  onSubmit: (email: string, password: string, nickname: string) => Promise<void>
  error?: string
}

// 회원가입 폼 컴포넌트
export default function RegisterForm({ onSubmit, error }: RegisterFormProps) {
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [nickname, setNickname] = useState('')
  const [loading, setLoading] = useState(false)

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setLoading(true)
    try {
      await onSubmit(email, password, nickname)
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
        label="닉네임"
        type="text"
        value={nickname}
        onChange={(e) => setNickname(e.target.value)}
        placeholder="닉네임 (2-50자)"
        required
        minLength={2}
        maxLength={50}
      />
      <Input
        label="비밀번호"
        type="password"
        value={password}
        onChange={(e) => setPassword(e.target.value)}
        placeholder="비밀번호 (8자 이상)"
        required
        minLength={8}
      />
      {error && <p className="text-sm text-red-500">{error}</p>}
      <Button type="submit" loading={loading} className="w-full">
        회원가입
      </Button>
    </form>
  )
}
