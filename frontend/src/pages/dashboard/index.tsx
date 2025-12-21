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
import { useDashboardKpis } from './hooks/use-dashboard-kpis'
import { WidgetCard } from './ui/widget-card'
import { AddWidgetDialog } from './ui/add-widget-dialog'
import { AddKpiDialog } from './ui/add-kpi-dialog'
import { KpiGrid } from './ui/kpi-grid'
import { DashboardSwitcher } from './ui/dashboard-switcher'

interface DashboardPageProps {
  dashboardId?: string
}

export function DashboardPage({ dashboardId }: DashboardPageProps) {
  const { user } = useAuth()
  const canEdit = user?.role === 'admin' || user?.role === 'editor'

  const [addWidgetOpen, setAddWidgetOpen] = useState(false)
  const [addKpiOpen, setAddKpiOpen] = useState(false)

  const {
    dashboards,
    dashboard,
    widgets,
    isLoading,
    addWidget,
    updateWidget,
    deleteWidget,
    isAddingWidget,
  } = useDashboard(dashboardId)

  const {
    kpis,
    computedKpis,
    isComputingKpis,
    addKpi,
    updateKpi,
    deleteKpi,
    isAddingKpi,
    isUpdatingKpi,
  } = useDashboardKpis(dashboard?.id)

  const title = dashboards.length > 1 ? (
    <DashboardSwitcher dashboards={dashboards} currentDashboard={dashboard} />
  ) : (
    dashboard?.name ?? 'Dashboard'
  )

  const hasContent = widgets.length > 0 || computedKpis.length > 0

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
              <DropdownMenuItem onClick={() => setAddKpiOpen(true)}>
                Add KPI
              </DropdownMenuItem>
              <DropdownMenuItem onClick={() => setAddWidgetOpen(true)}>
                Add Widget
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
              ? 'Add KPIs or widgets to start tracking metrics.'
              : 'No content has been added to this dashboard.'
          }
        />
      ) : (
        <div className="space-y-6">
          <KpiGrid
            kpis={kpis}
            computedKpis={computedKpis}
            isLoading={isComputingKpis}
            canEdit={canEdit}
            onUpdate={updateKpi}
            onDelete={deleteKpi}
            isUpdating={isUpdatingKpi}
          />
          {widgets.map((widget) => (
            <WidgetCard
              key={widget.id}
              widget={widget}
              canEdit={canEdit}
              onDelete={() => widget.id && deleteWidget(widget.id)}
              onUpdate={updateWidget}
            />
          ))}
        </div>
      )}

      <AddWidgetDialog
        open={addWidgetOpen}
        onOpenChange={setAddWidgetOpen}
        onAdd={addWidget}
        isLoading={isAddingWidget}
      />

      <AddKpiDialog
        open={addKpiOpen}
        onOpenChange={setAddKpiOpen}
        onAdd={addKpi}
        isLoading={isAddingKpi}
      />
    </AuthenticatedLayout>
  )
}
