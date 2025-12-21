import { useState } from 'react'
import { useQueryClient } from '@tanstack/react-query'
import { toast } from 'sonner'
import {
  useGetAuthInvites,
  useGetAuthEmailConfig,
  usePostAuthInvites,
  useDeleteAuthInvitesId,
  getGetAuthInvitesQueryKey,
} from '@/shared/api/generated/api'
import type { InviteWithInviter, Role } from '@/shared/api/generated/models'

export function useInvites() {
  const queryClient = useQueryClient()
  const [inviteDialogOpen, setInviteDialogOpen] = useState(false)
  const [inviteLinkDialogOpen, setInviteLinkDialogOpen] = useState(false)
  const [inviteLink, setInviteLink] = useState<string | null>(null)
  const [cancelInviteDialogOpen, setCancelInviteDialogOpen] = useState(false)
  const [selectedInvite, setSelectedInvite] = useState<InviteWithInviter | null>(null)

  const { data: invitesData, isLoading: isLoadingInvites } = useGetAuthInvites()
  const { data: emailConfig } = useGetAuthEmailConfig()

  const createInviteMutation = usePostAuthInvites({
    mutation: {
      onSuccess: (response) => {
        queryClient.invalidateQueries({ queryKey: getGetAuthInvitesQueryKey() })
        setInviteDialogOpen(false)

        if (response.inviteUrl) {
          setInviteLink(response.inviteUrl)
          setInviteLinkDialogOpen(true)
        } else {
          toast.success('Invitation sent')
        }
      },
      onError: (error) => {
        const message = error?.error || 'Failed to create invite'
        toast.error(message)
      },
    },
  })

  const cancelInviteMutation = useDeleteAuthInvitesId({
    mutation: {
      onSuccess: () => {
        queryClient.invalidateQueries({ queryKey: getGetAuthInvitesQueryKey() })
        setCancelInviteDialogOpen(false)
        setSelectedInvite(null)
        toast.success('Invite cancelled')
      },
      onError: (error) => {
        const message = error?.error || 'Failed to cancel invite'
        toast.error(message)
      },
    },
  })

  const handleCreateInvite = async (email: string, role: Role) => {
    await createInviteMutation.mutateAsync({ data: { email, role } })
  }

  const handleCancelInvite = (invite: InviteWithInviter) => {
    setSelectedInvite(invite)
    setCancelInviteDialogOpen(true)
  }

  const confirmCancelInvite = async () => {
    if (selectedInvite?.id) {
      await cancelInviteMutation.mutateAsync({ id: selectedInvite.id })
    }
  }

  const closeDialogs = () => {
    setInviteDialogOpen(false)
    setInviteLinkDialogOpen(false)
    setCancelInviteDialogOpen(false)
    setInviteLink(null)
    setSelectedInvite(null)
  }

  return {
    invites: invitesData?.invites ?? [],
    isLoading: isLoadingInvites,
    emailEnabled: emailConfig?.enabled ?? false,
    inviteDialogOpen,
    setInviteDialogOpen,
    inviteLinkDialogOpen,
    inviteLink,
    cancelInviteDialogOpen,
    selectedInvite,
    isCreatingInvite: createInviteMutation.isPending,
    isCancellingInvite: cancelInviteMutation.isPending,
    handleCreateInvite,
    handleCancelInvite,
    confirmCancelInvite,
    closeDialogs,
  }
}
