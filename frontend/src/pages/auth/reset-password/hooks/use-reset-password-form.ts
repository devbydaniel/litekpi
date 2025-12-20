import { useState } from 'react'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import { usePostAuthForgotPassword } from '@/shared/api/generated/api'
import { ApiError } from '@/shared/api/client'

const resetPasswordSchema = z.object({
  email: z.string().email('Please enter a valid email'),
})

export type ResetPasswordFormValues = z.infer<typeof resetPasswordSchema>

export function useResetPasswordForm() {
  const [error, setError] = useState<string | null>(null)
  const [success, setSuccess] = useState(false)

  const forgotPasswordMutation = usePostAuthForgotPassword({
    mutation: {
      onSuccess: () => {
        setSuccess(true)
      },
      onError: (err) => {
        if (err instanceof ApiError) {
          const data = err.data as { error?: string }
          setError(data.error || 'Request failed')
        } else {
          setError('An unexpected error occurred')
        }
      },
    },
  })

  const form = useForm<ResetPasswordFormValues>({
    resolver: zodResolver(resetPasswordSchema),
    defaultValues: {
      email: '',
    },
  })

  const onSubmit = async (values: ResetPasswordFormValues) => {
    setError(null)
    await forgotPasswordMutation.mutateAsync({ data: { email: values.email } })
  }

  return {
    form,
    isLoading: forgotPasswordMutation.isPending,
    error,
    success,
    onSubmit,
  }
}
