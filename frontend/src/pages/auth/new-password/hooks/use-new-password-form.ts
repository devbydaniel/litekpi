import { useState } from 'react'
import { authApi } from '@/shared/api/auth'
import { ApiError } from '@/shared/api/client'

export function useNewPasswordForm(token: string) {
  const [password, setPassword] = useState('')
  const [confirmPassword, setConfirmPassword] = useState('')
  const [isLoading, setIsLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [success, setSuccess] = useState(false)

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setError(null)

    if (!token) {
      setError('Missing reset token')
      return
    }

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
      await authApi.resetPassword({ token, newPassword: password })
      setSuccess(true)
    } catch (err) {
      if (err instanceof ApiError) {
        const data = err.data as { error?: string }
        setError(data.error || 'Reset failed')
      } else {
        setError('An unexpected error occurred')
      }
    } finally {
      setIsLoading(false)
    }
  }

  return {
    password,
    setPassword,
    confirmPassword,
    setConfirmPassword,
    isLoading,
    error,
    success,
    handleSubmit,
  }
}
