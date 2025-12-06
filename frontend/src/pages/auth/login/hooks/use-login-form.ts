import { useState } from 'react'
import { useNavigate } from '@tanstack/react-router'
import { authApi } from '@/shared/api/auth'
import { useAuthStore } from '@/shared/stores/auth-store'
import { ApiError } from '@/shared/api/client'

interface UseLoginFormOptions {
  initialError?: string
}

export function useLoginForm({ initialError }: UseLoginFormOptions = {}) {
  const navigate = useNavigate()
  const setAuth = useAuthStore((state) => state.setAuth)

  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [isLoading, setIsLoading] = useState(false)
  const [error, setError] = useState<string | null>(initialError || null)

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setError(null)
    setIsLoading(true)

    try {
      const response = await authApi.login({ email, password })
      setAuth(response.user, response.token)
      navigate({ to: '/' })
    } catch (err) {
      if (err instanceof ApiError) {
        const data = err.data as { error?: string }
        setError(data.error || 'Login failed')
      } else {
        setError('An unexpected error occurred')
      }
    } finally {
      setIsLoading(false)
    }
  }

  const handleOAuthLogin = (provider: 'google' | 'github') => {
    const url = provider === 'google' ? authApi.getGoogleAuthUrl() : authApi.getGithubAuthUrl()
    window.location.href = url
  }

  return {
    email,
    setEmail,
    password,
    setPassword,
    isLoading,
    error,
    handleSubmit,
    handleOAuthLogin,
  }
}
