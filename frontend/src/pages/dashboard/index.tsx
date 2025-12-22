import { useState } from 'react'
import { Plus, BarChart3 } from 'lucide-react'
import { AuthenticatedLayout } from '@/layouts/authenticated'
import { Button } from '@/shared/components/ui/button'
import { Skeleton } from '@/shared/components/ui/skeleton'
import { EmptyState } from '@/shared/components/ui/empty-state'
import { useAuth } from '@/shared/hooks/use-auth'
import { useDashboard } from './hooks/use-dashboard'
import { useMetrics } from './hooks/use-metrics'
import { MetricGrid } from './ui/metric-grid'
import { AddMetricDialog } from './ui/add-metric-dialog'
import { DashboardSwitcher } from './ui/dashboard-switcher'

interface DashboardPageProps {
  dashboardId?: string
}

export function DashboardPage({ dashboardId }: DashboardPageProps) {
  const { user } = useAuth()
  const canEdit = user?.role === 'admin' || user?.role === 'editor'

  const [addMetricOpen, setAddMetricOpen] = useState(false)

  const { dashboards, dashboard, isLoading } = useDashboard(dashboardId)

  const {
    metrics,
    computedMetrics,
    isComputingMetrics,
    addMetric,
    updateMetric,
    deleteMetric,
    isAddingMetric,
    isUpdatingMetric,
  } = useMetrics(dashboard?.id)

  const title =
    dashboards.length > 1 ? (
      <DashboardSwitcher dashboards={dashboards} currentDashboard={dashboard} />
    ) : (
      dashboard?.name ?? 'Dashboard'
    )

  const hasContent = computedMetrics.length > 0

  return (
    <AuthenticatedLayout
      title={title}
      actions={
        canEdit ? (
          <Button onClick={() => setAddMetricOpen(true)}>
            <Plus className="h-4 w-4" />
            Add Metric
          </Button>
        ) : undefined
      }
    >
      {isLoading || isComputingMetrics ? (
        <div className="space-y-6">
          <Skeleton className="h-[400px] w-full" />
          <Skeleton className="h-[400px] w-full" />
        </div>
      ) : !hasContent ? (
        <EmptyState
          icon={BarChart3}
          title="No content yet"
          description={
            canEdit
              ? 'Add metrics to start tracking data.'
              : 'No content has been added to this dashboard.'
          }
        />
      ) : (
        <MetricGrid
          metrics={metrics}
          computedMetrics={computedMetrics}
          isLoading={isComputingMetrics}
          canEdit={canEdit}
          onUpdate={updateMetric}
          onDelete={deleteMetric}
          isUpdating={isUpdatingMetric}
        />
      )}

      <AddMetricDialog
        open={addMetricOpen}
        onOpenChange={setAddMetricOpen}
        onAdd={addMetric}
        isLoading={isAddingMetric}
      />
    </AuthenticatedLayout>
  )
}
