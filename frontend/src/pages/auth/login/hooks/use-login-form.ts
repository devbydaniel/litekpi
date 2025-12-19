import { useState } from 'react'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import { useNavigate } from '@tanstack/react-router'
import { authApi } from '@/shared/api/auth'
import { useAuthStore } from '@/shared/stores/auth-store'
import { ApiError } from '@/shared/api/client'

const loginSchema = z.object({
  email: z.string().email('Please enter a valid email'),
  password: z.string().min(1, 'Password is required'),
})

export type LoginFormValues = z.infer<typeof loginSchema>

interface UseLoginFormOptions {
  initialError?: string
}

export function useLoginForm({ initialError }: UseLoginFormOptions = {}) {
  const navigate = useNavigate()
  const setAuth = useAuthStore((state) => state.setAuth)

  const [isLoading, setIsLoading] = useState(false)
  const [error, setError] = useState<string | null>(initialError || null)

  const form = useForm<LoginFormValues>({
    resolver: zodResolver(loginSchema),
    defaultValues: {
      email: '',
      password: '',
    },
  })

  const onSubmit = async (values: LoginFormValues) => {
    setError(null)
    setIsLoading(true)

    try {
      const response = await authApi.login(values)
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
    form,
    isLoading,
    error,
    onSubmit,
    handleOAuthLogin,
  }
}
