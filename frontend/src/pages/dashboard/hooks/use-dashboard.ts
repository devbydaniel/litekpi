import { useQueryClient } from '@tanstack/react-query'
import { toast } from 'sonner'
import {
  useGetDashboards,
  useGetDashboardsDefault,
  useGetDashboardsId,
  usePostDashboards,
  usePostDashboardsIdTimeSeries,
  useDeleteDashboardsIdTimeSeriesTimeSeriesId,
  usePutDashboardsIdTimeSeriesTimeSeriesId,
  getGetDashboardsDefaultQueryKey,
  getGetDashboardsIdQueryKey,
  getGetDashboardsQueryKey,
} from '@/shared/api/generated/api'
import type { CreateTimeSeriesRequest, UpdateTimeSeriesRequest } from '@/shared/api/generated/models'

export function useDashboard(dashboardId?: string) {
  const queryClient = useQueryClient()

  // Fetch all dashboards for the switcher
  const { data: dashboardsData } = useGetDashboards()

  // Fetch specific dashboard or default
  const { data: defaultDashboardData, isLoading: isLoadingDefault } = useGetDashboardsDefault({
    query: { enabled: !dashboardId },
  })

  const { data: specificDashboardData, isLoading: isLoadingSpecific } = useGetDashboardsId(
    dashboardId ?? '',
    { query: { enabled: !!dashboardId } }
  )

  const dashboardData = dashboardId ? specificDashboardData : defaultDashboardData
  const isLoading = dashboardId ? isLoadingSpecific : isLoadingDefault

  // Create dashboard mutation
  const createDashboardMutation = usePostDashboards({
    mutation: {
      onSuccess: () => {
        queryClient.invalidateQueries({ queryKey: getGetDashboardsQueryKey() })
        toast.success('Dashboard created')
      },
      onError: () => {
        toast.error('Failed to create dashboard')
      },
    },
  })

  // Create time series mutation
  const createTimeSeriesMutation = usePostDashboardsIdTimeSeries({
    mutation: {
      onSuccess: () => {
        if (dashboardId) {
          queryClient.invalidateQueries({ queryKey: getGetDashboardsIdQueryKey(dashboardId) })
        } else {
          queryClient.invalidateQueries({ queryKey: getGetDashboardsDefaultQueryKey() })
        }
        toast.success('Chart added')
      },
      onError: () => {
        toast.error('Failed to add chart')
      },
    },
  })

  // Update time series mutation
  const updateTimeSeriesMutation = usePutDashboardsIdTimeSeriesTimeSeriesId({
    mutation: {
      onSuccess: () => {
        if (dashboardId) {
          queryClient.invalidateQueries({ queryKey: getGetDashboardsIdQueryKey(dashboardId) })
        } else {
          queryClient.invalidateQueries({ queryKey: getGetDashboardsDefaultQueryKey() })
        }
        toast.success('Chart updated')
      },
      onError: () => {
        toast.error('Failed to update chart')
      },
    },
  })

  // Delete time series mutation
  const deleteTimeSeriesMutation = useDeleteDashboardsIdTimeSeriesTimeSeriesId({
    mutation: {
      onSuccess: () => {
        if (dashboardId) {
          queryClient.invalidateQueries({ queryKey: getGetDashboardsIdQueryKey(dashboardId) })
        } else {
          queryClient.invalidateQueries({ queryKey: getGetDashboardsDefaultQueryKey() })
        }
        toast.success('Chart removed')
      },
      onError: () => {
        toast.error('Failed to remove chart')
      },
    },
  })

  const createDashboard = async (name: string) => {
    await createDashboardMutation.mutateAsync({ data: { name } })
  }

  const addTimeSeries = async (timeSeries: CreateTimeSeriesRequest) => {
    const id = dashboardData?.dashboard?.id
    if (!id) return
    await createTimeSeriesMutation.mutateAsync({ id, data: timeSeries })
  }

  const updateTimeSeries = async (timeSeriesId: string, timeSeries: UpdateTimeSeriesRequest) => {
    const id = dashboardData?.dashboard?.id
    if (!id) return
    await updateTimeSeriesMutation.mutateAsync({ id, timeSeriesId, data: timeSeries })
  }

  const deleteTimeSeries = async (timeSeriesId: string) => {
    const id = dashboardData?.dashboard?.id
    if (!id) return
    await deleteTimeSeriesMutation.mutateAsync({ id, timeSeriesId })
  }

  return {
    dashboards: dashboardsData?.dashboards ?? [],
    dashboard: dashboardData?.dashboard,
    timeSeries: dashboardData?.timeSeries ?? [],
    scalarMetrics: dashboardData?.scalarMetrics ?? [],
    isLoading,
    createDashboard,
    addTimeSeries,
    updateTimeSeries,
    deleteTimeSeries,
    isAddingTimeSeries: createTimeSeriesMutation.isPending,
    isUpdatingTimeSeries: updateTimeSeriesMutation.isPending,
    isDeletingTimeSeries: deleteTimeSeriesMutation.isPending,
  }
}
