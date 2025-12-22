import { useState } from 'react'
import { Plus, BarChart3, ChevronDown } from 'lucide-react'
import { AuthenticatedLayout } from '@/layouts/authenticated'
import { Button } from '@/shared/components/ui/button'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from '@/shared/components/ui/dropdown-menu'
import { Skeleton } from '@/shared/components/ui/skeleton'
import { EmptyState } from '@/shared/components/ui/empty-state'
import { useAuth } from '@/shared/hooks/use-auth'
import { useDashboard } from './hooks/use-dashboard'
import { useScalarMetrics } from './hooks/use-scalar-metrics'
import { TimeSeriesCard } from './ui/time-series-card'
import { AddTimeSeriesDialog } from './ui/add-time-series-dialog'
import { AddScalarMetricDialog } from './ui/add-scalar-metric-dialog'
import { ScalarMetricGrid } from './ui/scalar-metric-grid'
import { DashboardSwitcher } from './ui/dashboard-switcher'

interface DashboardPageProps {
  dashboardId?: string
}

export function DashboardPage({ dashboardId }: DashboardPageProps) {
  const { user } = useAuth()
  const canEdit = user?.role === 'admin' || user?.role === 'editor'

  const [addTimeSeriesOpen, setAddTimeSeriesOpen] = useState(false)
  const [addMetricOpen, setAddMetricOpen] = useState(false)

  const {
    dashboards,
    dashboard,
    timeSeries,
    isLoading,
    addTimeSeries,
    updateTimeSeries,
    deleteTimeSeries,
    isAddingTimeSeries,
  } = useDashboard(dashboardId)

  const {
    scalarMetrics,
    computedMetrics,
    isComputingMetrics,
    addMetric,
    updateMetric,
    deleteMetric,
    isAddingMetric,
    isUpdatingMetric,
  } = useScalarMetrics(dashboard?.id)

  const title = dashboards.length > 1 ? (
    <DashboardSwitcher dashboards={dashboards} currentDashboard={dashboard} />
  ) : (
    dashboard?.name ?? 'Dashboard'
  )

  const hasContent = timeSeries.length > 0 || computedMetrics.length > 0

  return (
    <AuthenticatedLayout
      title={title}
      actions={
        canEdit ? (
          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <Button>
                <Plus className="h-4 w-4" />
                Add
                <ChevronDown className="ml-1 h-4 w-4" />
              </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent align="end">
              <DropdownMenuItem onClick={() => setAddMetricOpen(true)}>
                Add Metric
              </DropdownMenuItem>
              <DropdownMenuItem onClick={() => setAddTimeSeriesOpen(true)}>
                Add Chart
              </DropdownMenuItem>
            </DropdownMenuContent>
          </DropdownMenu>
        ) : undefined
      }
    >
      {isLoading ? (
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
              ? 'Add metrics or charts to start tracking data.'
              : 'No content has been added to this dashboard.'
          }
        />
      ) : (
        <div className="space-y-6">
          <ScalarMetricGrid
            metrics={scalarMetrics}
            computedMetrics={computedMetrics}
            isLoading={isComputingMetrics}
            canEdit={canEdit}
            onUpdate={updateMetric}
            onDelete={deleteMetric}
            isUpdating={isUpdatingMetric}
          />
          {timeSeries.map((ts) => (
            <TimeSeriesCard
              key={ts.id}
              timeSeries={ts}
              canEdit={canEdit}
              onDelete={() => ts.id && deleteTimeSeries(ts.id)}
              onUpdate={updateTimeSeries}
            />
          ))}
        </div>
      )}

      <AddTimeSeriesDialog
        open={addTimeSeriesOpen}
        onOpenChange={setAddTimeSeriesOpen}
        onAdd={addTimeSeries}
        isLoading={isAddingTimeSeries}
      />

      <AddScalarMetricDialog
        open={addMetricOpen}
        onOpenChange={setAddMetricOpen}
        onAdd={addMetric}
        isLoading={isAddingMetric}
      />
    </AuthenticatedLayout>
  )
}
