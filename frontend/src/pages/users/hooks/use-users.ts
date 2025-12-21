import { useState } from 'react'
import { useQueryClient } from '@tanstack/react-query'
import { toast } from 'sonner'
import {
  useGetAuthUsers,
  useDeleteAuthUsersId,
  usePatchAuthUsersIdRole,
  getGetAuthUsersQueryKey,
} from '@/shared/api/generated/api'
import type { User, Role } from '@/shared/api/generated/models'

export function useUsers() {
  const queryClient = useQueryClient()
  const [changeRoleDialogOpen, setChangeRoleDialogOpen] = useState(false)
  const [removeUserDialogOpen, setRemoveUserDialogOpen] = useState(false)
  const [selectedUser, setSelectedUser] = useState<User | null>(null)

  const { data, isLoading } = useGetAuthUsers()

  const updateRoleMutation = usePatchAuthUsersIdRole({
    mutation: {
      onSuccess: () => {
        queryClient.invalidateQueries({ queryKey: getGetAuthUsersQueryKey() })
        setChangeRoleDialogOpen(false)
        setSelectedUser(null)
        toast.success('User role updated')
      },
      onError: (error) => {
        const message = error?.error || 'Failed to update user role'
        toast.error(message)
      },
    },
  })

  const removeUserMutation = useDeleteAuthUsersId({
    mutation: {
      onSuccess: () => {
        queryClient.invalidateQueries({ queryKey: getGetAuthUsersQueryKey() })
        setRemoveUserDialogOpen(false)
        setSelectedUser(null)
        toast.success('User removed')
      },
      onError: (error) => {
        const message = error?.error || 'Failed to remove user'
        toast.error(message)
      },
    },
  })

  const handleChangeRole = (user: User) => {
    setSelectedUser(user)
    setChangeRoleDialogOpen(true)
  }

  const confirmChangeRole = async (role: Role) => {
    if (selectedUser?.id) {
      await updateRoleMutation.mutateAsync({
        id: selectedUser.id,
        data: { role },
      })
    }
  }

  const handleRemoveUser = (user: User) => {
    setSelectedUser(user)
    setRemoveUserDialogOpen(true)
  }

  const confirmRemoveUser = async () => {
    if (selectedUser?.id) {
      await removeUserMutation.mutateAsync({ id: selectedUser.id })
    }
  }

  const closeDialogs = () => {
    setChangeRoleDialogOpen(false)
    setRemoveUserDialogOpen(false)
    setSelectedUser(null)
  }

  return {
    users: data?.users ?? [],
    isLoading,
    changeRoleDialogOpen,
    removeUserDialogOpen,
    selectedUser,
    isUpdatingRole: updateRoleMutation.isPending,
    isRemovingUser: removeUserMutation.isPending,
    handleChangeRole,
    confirmChangeRole,
    handleRemoveUser,
    confirmRemoveUser,
    closeDialogs,
  }
}
