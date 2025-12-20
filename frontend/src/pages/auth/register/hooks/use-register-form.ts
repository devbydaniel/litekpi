import { useState } from 'react'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import { usePostAuthRegister } from '@/shared/api/generated/api'
import { ApiError } from '@/shared/api/client'

const registerSchema = z
  .object({
    name: z.string().min(1, 'Name is required').max(255, 'Name is too long'),
    organizationName: z.string().min(1, 'Organization name is required').max(255, 'Organization name is too long'),
    email: z.string().email('Please enter a valid email'),
    password: z.string().min(8, 'Password must be at least 8 characters'),
    confirmPassword: z.string(),
  })
  .refine((data) => data.password === data.confirmPassword, {
    message: 'Passwords do not match',
    path: ['confirmPassword'],
  })

export type RegisterFormValues = z.infer<typeof registerSchema>

export function useRegisterForm() {
  const [error, setError] = useState<string | null>(null)
  const [success, setSuccess] = useState(false)

  const registerMutation = usePostAuthRegister({
    mutation: {
      onSuccess: () => {
        setSuccess(true)
      },
      onError: (err) => {
        if (err instanceof ApiError) {
          const data = err.data as { error?: string }
          setError(data.error || 'Registration failed')
        } else {
          setError('An unexpected error occurred')
        }
      },
    },
  })

  const form = useForm<RegisterFormValues>({
    resolver: zodResolver(registerSchema),
    defaultValues: {
      name: '',
      organizationName: '',
      email: '',
      password: '',
      confirmPassword: '',
    },
  })

  const onSubmit = async (values: RegisterFormValues) => {
    setError(null)
    await registerMutation.mutateAsync({
      data: {
        email: values.email,
        password: values.password,
        name: values.name,
        organizationName: values.organizationName,
      },
    })
  }

  return {
    form,
    isLoading: registerMutation.isPending,
    error,
    success,
    email: form.watch('email'),
    onSubmit,
  }
}
