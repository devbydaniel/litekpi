import { useQueryClient } from '@tanstack/react-query'
import { toast } from 'sonner'
import {
  useGetDashboardsIdKpis,
  useGetDashboardsIdKpisCompute,
  usePostDashboardsIdKpis,
  usePutDashboardsIdKpisKpiId,
  useDeleteDashboardsIdKpisKpiId,
  getGetDashboardsIdKpisQueryKey,
  getGetDashboardsIdKpisComputeQueryKey,
} from '@/shared/api/generated/api'
import type { CreateKPIRequest, UpdateKPIRequest } from '@/shared/api/generated/models'

export function useDashboardKpis(dashboardId: string | undefined) {
  const queryClient = useQueryClient()

  // Fetch KPIs for the dashboard
  const { data: kpisData, isLoading: isLoadingKpis } = useGetDashboardsIdKpis(
    dashboardId ?? '',
    { query: { enabled: !!dashboardId } }
  )

  // Fetch computed KPIs (with actual values)
  const { data: computedKpisData, isLoading: isComputingKpis } = useGetDashboardsIdKpisCompute(
    dashboardId ?? '',
    { query: { enabled: !!dashboardId } }
  )

  // Invalidate helper
  const invalidateKpis = () => {
    if (dashboardId) {
      queryClient.invalidateQueries({ queryKey: getGetDashboardsIdKpisQueryKey(dashboardId) })
      queryClient.invalidateQueries({ queryKey: getGetDashboardsIdKpisComputeQueryKey(dashboardId) })
    }
  }

  // Create KPI mutation
  const createKpiMutation = usePostDashboardsIdKpis({
    mutation: {
      onSuccess: () => {
        invalidateKpis()
        toast.success('KPI added')
      },
      onError: () => {
        toast.error('Failed to add KPI')
      },
    },
  })

  // Update KPI mutation
  const updateKpiMutation = usePutDashboardsIdKpisKpiId({
    mutation: {
      onSuccess: () => {
        invalidateKpis()
        toast.success('KPI updated')
      },
      onError: () => {
        toast.error('Failed to update KPI')
      },
    },
  })

  // Delete KPI mutation
  const deleteKpiMutation = useDeleteDashboardsIdKpisKpiId({
    mutation: {
      onSuccess: () => {
        invalidateKpis()
        toast.success('KPI removed')
      },
      onError: () => {
        toast.error('Failed to remove KPI')
      },
    },
  })

  const addKpi = async (kpi: CreateKPIRequest) => {
    if (!dashboardId) return
    await createKpiMutation.mutateAsync({ id: dashboardId, data: kpi })
  }

  const updateKpi = async (kpiId: string, kpi: UpdateKPIRequest) => {
    if (!dashboardId) return
    await updateKpiMutation.mutateAsync({ id: dashboardId, kpiId, data: kpi })
  }

  const deleteKpi = async (kpiId: string) => {
    if (!dashboardId) return
    await deleteKpiMutation.mutateAsync({ id: dashboardId, kpiId })
  }

  return {
    kpis: kpisData?.kpis ?? [],
    computedKpis: computedKpisData?.kpis ?? [],
    isLoadingKpis,
    isComputingKpis,
    addKpi,
    updateKpi,
    deleteKpi,
    isAddingKpi: createKpiMutation.isPending,
    isUpdatingKpi: updateKpiMutation.isPending,
    isDeletingKpi: deleteKpiMutation.isPending,
  }
}
