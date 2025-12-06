import { useEffect, useState } from 'react'
import { authApi } from '@/shared/api/auth'
import { ApiError } from '@/shared/api/client'

type VerificationStatus = 'loading' | 'success' | 'error'

export function useEmailVerification(token: string) {
  const [status, setStatus] = useState<VerificationStatus>('loading')
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    if (!token) {
      setStatus('error')
      setError('Missing verification token')
      return
    }

    const verifyEmail = async () => {
      try {
        await authApi.verifyEmail(token)
        setStatus('success')
      } catch (err) {
        setStatus('error')
        if (err instanceof ApiError) {
          const data = err.data as { error?: string }
          setError(data.error || 'Verification failed')
        } else {
          setError('An unexpected error occurred')
        }
      }
    }

    verifyEmail()
  }, [token])

  return { status, error }
}
