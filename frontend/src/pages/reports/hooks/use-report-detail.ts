import { useQueryClient } from '@tanstack/react-query'
import { toast } from 'sonner'
import {
  useGetReportsId,
  useGetReportsIdCompute,
  usePostReportsIdKpis,
  usePutReportsIdKpisKpiId,
  useDeleteReportsIdKpisKpiId,
  getGetReportsIdQueryKey,
  getGetReportsIdComputeQueryKey,
} from '@/shared/api/generated/api'
import type { CreateKPIRequest, UpdateKPIRequest } from '@/shared/api/generated/models'

export function useReportDetail(reportId: string | undefined) {
  const queryClient = useQueryClient()

  // Fetch report with KPIs
  const { data: reportData, isLoading: isLoadingReport } = useGetReportsId(
    reportId ?? '',
    { query: { enabled: !!reportId } }
  )

  // Fetch computed report with KPI values
  const { data: computedData, isLoading: isComputingKpis } = useGetReportsIdCompute(
    reportId ?? '',
    { query: { enabled: !!reportId } }
  )

  // Invalidate helper
  const invalidateReport = () => {
    if (reportId) {
      queryClient.invalidateQueries({ queryKey: getGetReportsIdQueryKey(reportId) })
      queryClient.invalidateQueries({ queryKey: getGetReportsIdComputeQueryKey(reportId) })
    }
  }

  // Create KPI mutation
  const createKpiMutation = usePostReportsIdKpis({
    mutation: {
      onSuccess: () => {
        invalidateReport()
        toast.success('KPI added')
      },
      onError: () => {
        toast.error('Failed to add KPI')
      },
    },
  })

  // Update KPI mutation
  const updateKpiMutation = usePutReportsIdKpisKpiId({
    mutation: {
      onSuccess: () => {
        invalidateReport()
        toast.success('KPI updated')
      },
      onError: () => {
        toast.error('Failed to update KPI')
      },
    },
  })

  // Delete KPI mutation
  const deleteKpiMutation = useDeleteReportsIdKpisKpiId({
    mutation: {
      onSuccess: () => {
        invalidateReport()
        toast.success('KPI removed')
      },
      onError: () => {
        toast.error('Failed to remove KPI')
      },
    },
  })

  const addKpi = async (kpi: CreateKPIRequest) => {
    if (!reportId) return
    await createKpiMutation.mutateAsync({ id: reportId, data: kpi })
  }

  const updateKpi = async (kpiId: string, kpi: UpdateKPIRequest) => {
    if (!reportId) return
    await updateKpiMutation.mutateAsync({ id: reportId, kpiId, data: kpi })
  }

  const deleteKpi = async (kpiId: string) => {
    if (!reportId) return
    await deleteKpiMutation.mutateAsync({ id: reportId, kpiId })
  }

  return {
    report: reportData,
    computedReport: computedData,
    kpis: reportData?.kpis ?? [],
    computedKpis: computedData?.kpis ?? [],
    isLoadingReport,
    isComputingKpis,
    addKpi,
    updateKpi,
    deleteKpi,
    isAddingKpi: createKpiMutation.isPending,
    isUpdatingKpi: updateKpiMutation.isPending,
    isDeletingKpi: deleteKpiMutation.isPending,
  }
}
