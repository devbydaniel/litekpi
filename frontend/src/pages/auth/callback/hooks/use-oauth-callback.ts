import { useEffect, useState } from 'react'
import { useNavigate } from '@tanstack/react-router'
import { useAuthStore } from '@/shared/stores/auth-store'
import type { User } from '@/shared/types'

interface UseOAuthCallbackOptions {
  token: string
  userEncoded: string
  initialError?: string
}

export function useOAuthCallback({ token, userEncoded, initialError }: UseOAuthCallbackOptions) {
  const navigate = useNavigate()
  const setAuth = useAuthStore((state) => state.setAuth)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    if (initialError) {
      setError(initialError)
      return
    }

    if (!token || !userEncoded) {
      setError('Invalid callback parameters')
      return
    }

    try {
      // Decode user from base64
      const userJson = atob(userEncoded)
      const user: User = JSON.parse(userJson)

      // Set auth state
      setAuth(user, token)

      // Redirect to home
      navigate({ to: '/' })
    } catch {
      setError('Failed to process authentication')
    }
  }, [token, userEncoded, initialError, setAuth, navigate])

  return { error, navigate }
}
