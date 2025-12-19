import { useState } from 'react'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import { useNavigate } from '@tanstack/react-router'
import { authApi } from '@/shared/api/auth'
import { ApiError } from '@/shared/api/client'
import { useAuthStore } from '@/shared/stores/auth-store'

const completeSetupSchema = z.object({
  name: z.string().min(1, 'Name is required').max(255, 'Name is too long'),
  organizationName: z.string().min(1, 'Organization name is required').max(255, 'Organization name is too long'),
})

export type CompleteSetupFormValues = z.infer<typeof completeSetupSchema>

interface UseCompleteSetupFormOptions {
  token: string
  email: string
  initialName?: string
}

export function useCompleteSetupForm({ token, email, initialName }: UseCompleteSetupFormOptions) {
  const navigate = useNavigate()
  const setAuth = useAuthStore((state) => state.setAuth)
  const [isLoading, setIsLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  const form = useForm<CompleteSetupFormValues>({
    resolver: zodResolver(completeSetupSchema),
    defaultValues: {
      name: initialName || '',
      organizationName: '',
    },
  })

  const onSubmit = async (values: CompleteSetupFormValues) => {
    if (!token) {
      setError('Invalid setup token. Please try signing in again.')
      return
    }

    setError(null)
    setIsLoading(true)

    try {
      const response = await authApi.completeOAuthSetup({
        token,
        name: values.name,
        organizationName: values.organizationName,
      })

      setAuth(response.user, response.token)
      navigate({ to: '/' })
    } catch (err) {
      if (err instanceof ApiError) {
        const data = err.data as { error?: string }
        setError(data.error || 'Failed to complete setup')
      } else {
        setError('An unexpected error occurred')
      }
    } finally {
      setIsLoading(false)
    }
  }

  return {
    form,
    isLoading,
    error,
    email,
    onSubmit,
  }
}
