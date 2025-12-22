import { useQueryClient } from '@tanstack/react-query'
import { toast } from 'sonner'
import {
  useGetDashboards,
  useGetDashboardsDefault,
  useGetDashboardsId,
  usePostDashboards,
  getGetDashboardsQueryKey,
} from '@/shared/api/generated/api'

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

  const createDashboard = async (name: string) => {
    await createDashboardMutation.mutateAsync({ data: { name } })
  }

  return {
    dashboards: dashboardsData?.dashboards ?? [],
    dashboard: dashboardData?.dashboard,
    isLoading,
    createDashboard,
  }
}
