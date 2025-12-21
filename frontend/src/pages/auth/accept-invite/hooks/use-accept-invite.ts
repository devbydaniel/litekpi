import { useState } from 'react'
import { useNavigate } from '@tanstack/react-router'
import { toast } from 'sonner'
import {
  useGetAuthInvitesValidate,
  usePostAuthInvitesAccept,
} from '@/shared/api/generated/api'

export function useAcceptInvite(token: string | undefined) {
  const navigate = useNavigate()
  const [success, setSuccess] = useState(false)

  const { data: validationData, isLoading: isValidating, error: validationError } = useGetAuthInvitesValidate(
    { token: token ?? '' },
    {
      query: {
        enabled: !!token,
      },
    }
  )

  const acceptMutation = usePostAuthInvitesAccept({
    mutation: {
      onSuccess: () => {
        setSuccess(true)
        toast.success('Account created successfully')
      },
      onError: (error) => {
        const message = error?.error || 'Failed to accept invite'
        toast.error(message)
      },
    },
  })

  const handleAccept = async (name: string, password: string) => {
    if (!token) return
    await acceptMutation.mutateAsync({
      data: { token, name, password },
    })
  }

  const goToLogin = () => {
    navigate({ to: '/login' })
  }

  const isValid = validationData?.valid ?? false
  const inviteInfo = isValid ? validationData : null

  return {
    isValidating,
    validationError: validationError?.error,
    isValid,
    inviteInfo,
    success,
    isAccepting: acceptMutation.isPending,
    handleAccept,
    goToLogin,
  }
}
