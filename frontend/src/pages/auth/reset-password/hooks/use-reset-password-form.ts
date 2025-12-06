import { useState } from 'react'
import { authApi } from '@/shared/api/auth'
import { ApiError } from '@/shared/api/client'

export function useResetPasswordForm() {
  const [email, setEmail] = useState('')
  const [isLoading, setIsLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [success, setSuccess] = useState(false)

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setError(null)
    setIsLoading(true)

    try {
      await authApi.forgotPassword(email)
      setSuccess(true)
    } catch (err) {
      if (err instanceof ApiError) {
        const data = err.data as { error?: string }
        setError(data.error || 'Request failed')
      } else {
        setError('An unexpected error occurred')
      }
    } finally {
      setIsLoading(false)
    }
  }

  return {
    email,
    setEmail,
    isLoading,
    error,
    success,
    handleSubmit,
  }
}
