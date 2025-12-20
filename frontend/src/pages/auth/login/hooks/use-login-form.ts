import { useState } from 'react'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import { useNavigate } from '@tanstack/react-router'
import { usePostAuthLogin } from '@/shared/api/generated/api'
import { useAuthStore, type User } from '@/shared/stores/auth-store'
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
  const [error, setError] = useState<string | null>(initialError || null)

  const loginMutation = usePostAuthLogin({
    mutation: {
      onSuccess: (response) => {
        setAuth(response.user as User, response.token!)
        navigate({ to: '/' })
      },
      onError: (err) => {
        if (err instanceof ApiError) {
          const data = err.data as { error?: string }
          setError(data.error || 'Login failed')
        } else {
          setError('An unexpected error occurred')
        }
      },
    },
  })

  const form = useForm<LoginFormValues>({
    resolver: zodResolver(loginSchema),
    defaultValues: {
      email: '',
      password: '',
    },
  })

  const onSubmit = async (values: LoginFormValues) => {
    setError(null)
    await loginMutation.mutateAsync({ data: values })
  }

  return {
    form,
    isLoading: loginMutation.isPending,
    error,
    onSubmit,
  }
}
