import { useQueryClient } from '@tanstack/react-query'
import { toast } from 'sonner'
import {
  useGetDashboards,
  useGetDashboardsDefault,
  useGetDashboardsId,
  usePostDashboards,
  usePostDashboardsIdWidgets,
  useDeleteDashboardsIdWidgetsWidgetId,
  usePutDashboardsIdWidgetsWidgetId,
  getGetDashboardsDefaultQueryKey,
  getGetDashboardsIdQueryKey,
  getGetDashboardsQueryKey,
} from '@/shared/api/generated/api'
import type { CreateWidgetRequest, UpdateWidgetRequest } from '@/shared/api/generated/models'

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

  // Create widget mutation
  const createWidgetMutation = usePostDashboardsIdWidgets({
    mutation: {
      onSuccess: () => {
        if (dashboardId) {
          queryClient.invalidateQueries({ queryKey: getGetDashboardsIdQueryKey(dashboardId) })
        } else {
          queryClient.invalidateQueries({ queryKey: getGetDashboardsDefaultQueryKey() })
        }
        toast.success('Widget added')
      },
      onError: () => {
        toast.error('Failed to add widget')
      },
    },
  })

  // Update widget mutation
  const updateWidgetMutation = usePutDashboardsIdWidgetsWidgetId({
    mutation: {
      onSuccess: () => {
        if (dashboardId) {
          queryClient.invalidateQueries({ queryKey: getGetDashboardsIdQueryKey(dashboardId) })
        } else {
          queryClient.invalidateQueries({ queryKey: getGetDashboardsDefaultQueryKey() })
        }
        toast.success('Widget updated')
      },
      onError: () => {
        toast.error('Failed to update widget')
      },
    },
  })

  // Delete widget mutation
  const deleteWidgetMutation = useDeleteDashboardsIdWidgetsWidgetId({
    mutation: {
      onSuccess: () => {
        if (dashboardId) {
          queryClient.invalidateQueries({ queryKey: getGetDashboardsIdQueryKey(dashboardId) })
        } else {
          queryClient.invalidateQueries({ queryKey: getGetDashboardsDefaultQueryKey() })
        }
        toast.success('Widget removed')
      },
      onError: () => {
        toast.error('Failed to remove widget')
      },
    },
  })

  const createDashboard = async (name: string) => {
    await createDashboardMutation.mutateAsync({ data: { name } })
  }

  const addWidget = async (widget: CreateWidgetRequest) => {
    const id = dashboardData?.dashboard?.id
    if (!id) return
    await createWidgetMutation.mutateAsync({ id, data: widget })
  }

  const updateWidget = async (widgetId: string, widget: UpdateWidgetRequest) => {
    const id = dashboardData?.dashboard?.id
    if (!id) return
    await updateWidgetMutation.mutateAsync({ id, widgetId, data: widget })
  }

  const deleteWidget = async (widgetId: string) => {
    const id = dashboardData?.dashboard?.id
    if (!id) return
    await deleteWidgetMutation.mutateAsync({ id, widgetId })
  }

  return {
    dashboards: dashboardsData?.dashboards ?? [],
    dashboard: dashboardData?.dashboard,
    widgets: dashboardData?.widgets ?? [],
    isLoading,
    createDashboard,
    addWidget,
    updateWidget,
    deleteWidget,
    isAddingWidget: createWidgetMutation.isPending,
    isUpdatingWidget: updateWidgetMutation.isPending,
    isDeletingWidget: deleteWidgetMutation.isPending,
  }
}
