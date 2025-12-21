import { useState } from 'react'
import { useQueryClient } from '@tanstack/react-query'
import { toast } from 'sonner'
import {
  useGetReports,
  usePostReports,
  usePutReportsId,
  useDeleteReportsId,
  getGetReportsQueryKey,
} from '@/shared/api/generated/api'
import type { Report } from '@/shared/api/generated/models'

export function useReports() {
  const queryClient = useQueryClient()
  const [createDialogOpen, setCreateDialogOpen] = useState(false)
  const [editingReport, setEditingReport] = useState<Report | null>(null)
  const [reportToDelete, setReportToDelete] = useState<Report | null>(null)

  // Fetch all reports
  const { data: reportsData, isLoading } = useGetReports()
  const reports = reportsData?.reports ?? []

  // Create report mutation
  const createMutation = usePostReports({
    mutation: {
      onSuccess: () => {
        queryClient.invalidateQueries({ queryKey: getGetReportsQueryKey() })
        toast.success('Report created')
        setCreateDialogOpen(false)
      },
      onError: () => {
        toast.error('Failed to create report')
      },
    },
  })

  // Update report mutation
  const updateMutation = usePutReportsId({
    mutation: {
      onSuccess: () => {
        queryClient.invalidateQueries({ queryKey: getGetReportsQueryKey() })
        toast.success('Report updated')
        setEditingReport(null)
      },
      onError: () => {
        toast.error('Failed to update report')
      },
    },
  })

  // Delete report mutation
  const deleteMutation = useDeleteReportsId({
    mutation: {
      onSuccess: () => {
        queryClient.invalidateQueries({ queryKey: getGetReportsQueryKey() })
        toast.success('Report deleted')
        setReportToDelete(null)
      },
      onError: () => {
        toast.error('Failed to delete report')
      },
    },
  })

  const createReport = async (name: string) => {
    await createMutation.mutateAsync({ data: { name } })
  }

  const updateReport = async (id: string, name: string) => {
    await updateMutation.mutateAsync({ id, data: { name } })
  }

  const deleteReport = async (id: string) => {
    await deleteMutation.mutateAsync({ id })
  }

  return {
    reports,
    isLoading,
    createDialogOpen,
    setCreateDialogOpen,
    editingReport,
    setEditingReport,
    reportToDelete,
    setReportToDelete,
    createReport,
    updateReport,
    deleteReport,
    isCreating: createMutation.isPending,
    isUpdating: updateMutation.isPending,
    isDeleting: deleteMutation.isPending,
  }
}
