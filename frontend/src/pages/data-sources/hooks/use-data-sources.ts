import { useState } from 'react'
import { useQueryClient } from '@tanstack/react-query'
import { toast } from 'sonner'
import {
  useGetDataSources,
  usePostDataSources,
  useDeleteDataSourcesId,
  usePostDataSourcesIdRegenerateKey,
  usePostDataSourcesDemo,
  getGetDataSourcesQueryKey,
} from '@/shared/api/generated/api'
import type { DataSource } from '@/shared/api/generated/models'

export function useDataSources() {
  const queryClient = useQueryClient()
  const [createDialogOpen, setCreateDialogOpen] = useState(false)
  const [regenerateDialogOpen, setRegenerateDialogOpen] = useState(false)
  const [deleteDialogOpen, setDeleteDialogOpen] = useState(false)
  const [selectedDataSource, setSelectedDataSource] = useState<DataSource | null>(null)
  const [dataSourceToDelete, setDataSourceToDelete] = useState<DataSource | null>(null)
  const [apiKey, setApiKey] = useState<string | null>(null)

  const { data, isLoading } = useGetDataSources()

  const createMutation = usePostDataSources({
    mutation: {
      onSuccess: (response) => {
        queryClient.invalidateQueries({ queryKey: getGetDataSourcesQueryKey() })
        setApiKey(response.apiKey ?? null)
        toast.success('Data source created successfully')
      },
      onError: () => {
        toast.error('Failed to create data source')
      },
    },
  })

  const deleteMutation = useDeleteDataSourcesId({
    mutation: {
      onSuccess: () => {
        queryClient.invalidateQueries({ queryKey: getGetDataSourcesQueryKey() })
        setDeleteDialogOpen(false)
        setDataSourceToDelete(null)
        toast.success('Data source deleted')
      },
      onError: () => {
        toast.error('Failed to delete data source')
      },
    },
  })

  const regenerateMutation = usePostDataSourcesIdRegenerateKey({
    mutation: {
      onSuccess: (response) => {
        setApiKey(response.apiKey ?? null)
        toast.success('API key regenerated')
      },
      onError: () => {
        toast.error('Failed to regenerate API key')
      },
    },
  })

  const createDemoMutation = usePostDataSourcesDemo({
    mutation: {
      onSuccess: (response) => {
        queryClient.invalidateQueries({ queryKey: getGetDataSourcesQueryKey() })
        setApiKey(response.apiKey ?? null)
        setCreateDialogOpen(true)
        toast.success('Demo data source created successfully')
      },
      onError: () => {
        toast.error('Failed to create demo data source')
      },
    },
  })

  const handleCreateDataSource = async (name: string) => {
    await createMutation.mutateAsync({ data: { name } })
  }

  const handleDeleteDataSource = (dataSource: DataSource) => {
    setDataSourceToDelete(dataSource)
    setDeleteDialogOpen(true)
  }

  const confirmDeleteDataSource = async () => {
    if (dataSourceToDelete?.id) {
      await deleteMutation.mutateAsync({ id: dataSourceToDelete.id })
    }
  }

  const handleRegenerateKey = (dataSource: DataSource) => {
    setSelectedDataSource(dataSource)
    setRegenerateDialogOpen(true)
  }

  const confirmRegenerateKey = async () => {
    if (selectedDataSource?.id) {
      await regenerateMutation.mutateAsync({ id: selectedDataSource.id })
    }
  }

  const handleCreateDemo = async () => {
    await createDemoMutation.mutateAsync()
  }

  const closeDialogs = () => {
    setCreateDialogOpen(false)
    setRegenerateDialogOpen(false)
    setDeleteDialogOpen(false)
    setSelectedDataSource(null)
    setDataSourceToDelete(null)
    setApiKey(null)
  }

  return {
    dataSources: data?.dataSources ?? [],
    isLoading,
    createDialogOpen,
    setCreateDialogOpen,
    regenerateDialogOpen,
    deleteDialogOpen,
    selectedDataSource,
    dataSourceToDelete,
    apiKey,
    isCreating: createMutation.isPending,
    isCreatingDemo: createDemoMutation.isPending,
    isRegenerating: regenerateMutation.isPending,
    isDeleting: deleteMutation.isPending,
    handleCreateDataSource,
    handleCreateDemo,
    handleDeleteDataSource,
    confirmDeleteDataSource,
    handleRegenerateKey,
    confirmRegenerateKey,
    closeDialogs,
  }
}
