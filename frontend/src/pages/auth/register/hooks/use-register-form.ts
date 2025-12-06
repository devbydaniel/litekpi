import { useState } from 'react'
import { authApi } from '@/shared/api/auth'
import { ApiError } from '@/shared/api/client'

export function useRegisterForm() {
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [confirmPassword, setConfirmPassword] = useState('')
  const [isLoading, setIsLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [success, setSuccess] = useState(false)

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setError(null)

    if (password !== confirmPassword) {
      setError('Passwords do not match')
      return
    }

    if (password.length < 8) {
      setError('Password must be at least 8 characters')
      return
    }

    setIsLoading(true)

    try {
      await authApi.register({ email, password })
      setSuccess(true)
    } catch (err) {
      if (err instanceof ApiError) {
        const data = err.data as { error?: string }
        setError(data.error || 'Registration failed')
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
    confirmPassword,
    setConfirmPassword,
    isLoading,
    error,
    success,
    handleSubmit,
    handleOAuthLogin,
  }
}
