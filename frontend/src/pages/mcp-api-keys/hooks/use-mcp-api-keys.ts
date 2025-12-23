import { useState } from 'react'
import { useQueryClient } from '@tanstack/react-query'
import { toast } from 'sonner'
import {
  useGetMcpKeys,
  usePostMcpKeys,
  useDeleteMcpKeysId,
  usePutMcpKeysId,
  getGetMcpKeysQueryKey,
  useGetDataSources,
} from '@/shared/api/generated/api'
import type { MCPAPIKey } from '@/shared/api/generated/models'

export function useMCPApiKeys() {
  const queryClient = useQueryClient()
  const [createDialogOpen, setCreateDialogOpen] = useState(false)
  const [editDialogOpen, setEditDialogOpen] = useState(false)
  const [deleteDialogOpen, setDeleteDialogOpen] = useState(false)
  const [keyToEdit, setKeyToEdit] = useState<MCPAPIKey | null>(null)
  const [keyToDelete, setKeyToDelete] = useState<MCPAPIKey | null>(null)
  const [apiKey, setApiKey] = useState<string | null>(null)

  // Fetch MCP API keys
  const { data: keysData, isLoading: keysLoading } = useGetMcpKeys()

  // Fetch data sources for the multi-select
  const { data: dsData, isLoading: dsLoading } = useGetDataSources()

  const createMutation = usePostMcpKeys({
    mutation: {
      onSuccess: (response) => {
        queryClient.invalidateQueries({ queryKey: getGetMcpKeysQueryKey() })
        setApiKey(response.apiKey ?? null)
        toast.success('MCP API key created successfully')
      },
      onError: (error: unknown) => {
        const message =
          (error as { response?: { data?: { error?: string } } })?.response
            ?.data?.error || 'Failed to create MCP API key'
        toast.error(message)
      },
    },
  })

  const updateMutation = usePutMcpKeysId({
    mutation: {
      onSuccess: () => {
        queryClient.invalidateQueries({ queryKey: getGetMcpKeysQueryKey() })
        setEditDialogOpen(false)
        setKeyToEdit(null)
        toast.success('MCP API key updated successfully')
      },
      onError: (error: unknown) => {
        const message =
          (error as { response?: { data?: { error?: string } } })?.response
            ?.data?.error || 'Failed to update MCP API key'
        toast.error(message)
      },
    },
  })

  const deleteMutation = useDeleteMcpKeysId({
    mutation: {
      onSuccess: () => {
        queryClient.invalidateQueries({ queryKey: getGetMcpKeysQueryKey() })
        setDeleteDialogOpen(false)
        setKeyToDelete(null)
        toast.success('MCP API key deleted')
      },
      onError: () => {
        toast.error('Failed to delete MCP API key')
      },
    },
  })

  const handleCreateKey = async (name: string, dataSourceIds: string[]) => {
    await createMutation.mutateAsync({
      data: { name, dataSourceIds },
    })
  }

  const handleEditKey = (key: MCPAPIKey) => {
    setKeyToEdit(key)
    setEditDialogOpen(true)
  }

  const handleUpdateKey = async (dataSourceIds: string[]) => {
    if (keyToEdit?.id) {
      await updateMutation.mutateAsync({
        id: keyToEdit.id,
        data: { dataSourceIds },
      })
    }
  }

  const handleDeleteKey = (key: MCPAPIKey) => {
    setKeyToDelete(key)
    setDeleteDialogOpen(true)
  }

  const confirmDeleteKey = async () => {
    if (keyToDelete?.id) {
      await deleteMutation.mutateAsync({ id: keyToDelete.id })
    }
  }

  const closeDialogs = () => {
    setCreateDialogOpen(false)
    setEditDialogOpen(false)
    setDeleteDialogOpen(false)
    setKeyToEdit(null)
    setKeyToDelete(null)
    setApiKey(null)
  }

  return {
    keys: keysData?.keys ?? [],
    dataSources: dsData?.dataSources ?? [],
    isLoading: keysLoading || dsLoading,
    createDialogOpen,
    setCreateDialogOpen,
    editDialogOpen,
    keyToEdit,
    deleteDialogOpen,
    keyToDelete,
    apiKey,
    isCreating: createMutation.isPending,
    isUpdating: updateMutation.isPending,
    isDeleting: deleteMutation.isPending,
    handleCreateKey,
    handleEditKey,
    handleUpdateKey,
    handleDeleteKey,
    confirmDeleteKey,
    closeDialogs,
  }
}
