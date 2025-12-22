import { useQueryClient } from '@tanstack/react-query'
import { toast } from 'sonner'
import {
  useGetDashboardsIdMetrics,
  useGetDashboardsIdMetricsCompute,
  usePostDashboardsIdMetrics,
  usePutDashboardsIdMetricsMetricId,
  useDeleteDashboardsIdMetricsMetricId,
  getGetDashboardsIdMetricsQueryKey,
  getGetDashboardsIdMetricsComputeQueryKey,
} from '@/shared/api/generated/api'
import type { CreateMetricRequest, UpdateMetricRequest } from '@/shared/api/generated/models'

export function useMetrics(dashboardId: string | undefined) {
  const queryClient = useQueryClient()

  // Fetch metrics for the dashboard
  const { data: metricsData, isLoading: isLoadingMetrics } = useGetDashboardsIdMetrics(
    dashboardId ?? '',
    { query: { enabled: !!dashboardId } }
  )

  // Fetch computed metrics (with actual values)
  const { data: computedMetricsData, isLoading: isComputingMetrics } = useGetDashboardsIdMetricsCompute(
    dashboardId ?? '',
    { query: { enabled: !!dashboardId } }
  )

  // Invalidate helper
  const invalidateMetrics = () => {
    if (dashboardId) {
      queryClient.invalidateQueries({ queryKey: getGetDashboardsIdMetricsQueryKey(dashboardId) })
      queryClient.invalidateQueries({ queryKey: getGetDashboardsIdMetricsComputeQueryKey(dashboardId) })
    }
  }

  // Create metric mutation
  const createMetricMutation = usePostDashboardsIdMetrics({
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

  // Update metric mutation
  const updateMetricMutation = usePutDashboardsIdMetricsMetricId({
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

  // Delete metric mutation
  const deleteMetricMutation = useDeleteDashboardsIdMetricsMetricId({
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

  const addMetric = async (metric: CreateMetricRequest) => {
    if (!dashboardId) return
    await createMetricMutation.mutateAsync({ id: dashboardId, data: metric })
  }

  const updateMetric = async (metricId: string, metric: UpdateMetricRequest) => {
    if (!dashboardId) return
    await updateMetricMutation.mutateAsync({ id: dashboardId, metricId, data: metric })
  }

  const deleteMetric = async (metricId: string) => {
    if (!dashboardId) return
    await deleteMetricMutation.mutateAsync({ id: dashboardId, metricId })
  }

  return {
    metrics: metricsData?.metrics ?? [],
    computedMetrics: computedMetricsData?.metrics ?? [],
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
