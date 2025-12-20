import { useState } from 'react'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import { usePostAuthResetPassword } from '@/shared/api/generated/api'
import { ApiError } from '@/shared/api/client'

const newPasswordSchema = z
  .object({
    password: z.string().min(8, 'Password must be at least 8 characters'),
    confirmPassword: z.string(),
  })
  .refine((data) => data.password === data.confirmPassword, {
    message: 'Passwords do not match',
    path: ['confirmPassword'],
  })

export type NewPasswordFormValues = z.infer<typeof newPasswordSchema>

export function useNewPasswordForm(token: string) {
  const [error, setError] = useState<string | null>(null)
  const [success, setSuccess] = useState(false)

  const resetPasswordMutation = usePostAuthResetPassword({
    mutation: {
      onSuccess: () => {
        setSuccess(true)
      },
      onError: (err) => {
        if (err instanceof ApiError) {
          const data = err.data as { error?: string }
          setError(data.error || 'Reset failed')
        } else {
          setError('An unexpected error occurred')
        }
      },
    },
  })

  const form = useForm<NewPasswordFormValues>({
    resolver: zodResolver(newPasswordSchema),
    defaultValues: {
      password: '',
      confirmPassword: '',
    },
  })

  const onSubmit = async (values: NewPasswordFormValues) => {
    setError(null)

    if (!token) {
      setError('Missing reset token')
      return
    }

    await resetPasswordMutation.mutateAsync({
      data: { token, newPassword: values.password },
    })
  }

  return {
    form,
    isLoading: resetPasswordMutation.isPending,
    error,
    success,
    onSubmit,
  }
}
