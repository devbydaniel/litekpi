import { useState } from 'react'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import { authApi } from '@/shared/api/auth'
import { ApiError } from '@/shared/api/client'

const resetPasswordSchema = z.object({
  email: z.string().email('Please enter a valid email'),
})

export type ResetPasswordFormValues = z.infer<typeof resetPasswordSchema>

export function useResetPasswordForm() {
  const [isLoading, setIsLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [success, setSuccess] = useState(false)

  const form = useForm<ResetPasswordFormValues>({
    resolver: zodResolver(resetPasswordSchema),
    defaultValues: {
      email: '',
    },
  })

  const onSubmit = async (values: ResetPasswordFormValues) => {
    setError(null)
    setIsLoading(true)

    try {
      await authApi.forgotPassword(values.email)
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
    form,
    isLoading,
    error,
    success,
    onSubmit,
  }
}
