import { useQueryClient } from '@tanstack/react-query'
import { toast } from 'sonner'
import {
  useGetDashboardsIdScalarMetrics,
  useGetDashboardsIdScalarMetricsCompute,
  usePostDashboardsIdScalarMetrics,
  usePutDashboardsIdScalarMetricsScalarMetricId,
  useDeleteDashboardsIdScalarMetricsScalarMetricId,
  getGetDashboardsIdScalarMetricsQueryKey,
  getGetDashboardsIdScalarMetricsComputeQueryKey,
} from '@/shared/api/generated/api'
import type { CreateScalarMetricRequest, UpdateScalarMetricRequest } from '@/shared/api/generated/models'

export function useScalarMetrics(dashboardId: string | undefined) {
  const queryClient = useQueryClient()

  // Fetch scalar metrics for the dashboard
  const { data: metricsData, isLoading: isLoadingMetrics } = useGetDashboardsIdScalarMetrics(
    dashboardId ?? '',
    { query: { enabled: !!dashboardId } }
  )

  // Fetch computed scalar metrics (with actual values)
  const { data: computedMetricsData, isLoading: isComputingMetrics } = useGetDashboardsIdScalarMetricsCompute(
    dashboardId ?? '',
    { query: { enabled: !!dashboardId } }
  )

  // Invalidate helper
  const invalidateMetrics = () => {
    if (dashboardId) {
      queryClient.invalidateQueries({ queryKey: getGetDashboardsIdScalarMetricsQueryKey(dashboardId) })
      queryClient.invalidateQueries({ queryKey: getGetDashboardsIdScalarMetricsComputeQueryKey(dashboardId) })
    }
  }

  // Create scalar metric mutation
  const createMetricMutation = usePostDashboardsIdScalarMetrics({
    mutation: {
      onSuccess: () => {
        invalidateMetrics()
        toast.success('Metric added')
      },
      onError: () => {
        toast.error('Failed to add metric')
      },
    },
  })

  // Update scalar metric mutation
  const updateMetricMutation = usePutDashboardsIdScalarMetricsScalarMetricId({
    mutation: {
      onSuccess: () => {
        invalidateMetrics()
        toast.success('Metric updated')
      },
      onError: () => {
        toast.error('Failed to update metric')
      },
    },
  })

  // Delete scalar metric mutation
  const deleteMetricMutation = useDeleteDashboardsIdScalarMetricsScalarMetricId({
    mutation: {
      onSuccess: () => {
        invalidateMetrics()
        toast.success('Metric removed')
      },
      onError: () => {
        toast.error('Failed to remove metric')
      },
    },
  })

  const addMetric = async (metric: CreateScalarMetricRequest) => {
    if (!dashboardId) return
    await createMetricMutation.mutateAsync({ id: dashboardId, data: metric })
  }

  const updateMetric = async (scalarMetricId: string, metric: UpdateScalarMetricRequest) => {
    if (!dashboardId) return
    await updateMetricMutation.mutateAsync({ id: dashboardId, scalarMetricId, data: metric })
  }

  const deleteMetric = async (scalarMetricId: string) => {
    if (!dashboardId) return
    await deleteMetricMutation.mutateAsync({ id: dashboardId, scalarMetricId })
  }

  return {
    scalarMetrics: metricsData?.scalarMetrics ?? [],
    computedMetrics: computedMetricsData?.scalarMetrics ?? [],
    isLoadingMetrics,
    isComputingMetrics,
    addMetric,
    updateMetric,
    deleteMetric,
    isAddingMetric: createMetricMutation.isPending,
    isUpdatingMetric: updateMetricMutation.isPending,
    isDeletingMetric: deleteMetricMutation.isPending,
  }
}
