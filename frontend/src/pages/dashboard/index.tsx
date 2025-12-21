import { useState } from 'react'
import { Plus, BarChart3 } from 'lucide-react'
import { AuthenticatedLayout } from '@/layouts/authenticated'
import { Button } from '@/shared/components/ui/button'
import { Skeleton } from '@/shared/components/ui/skeleton'
import { EmptyState } from '@/shared/components/ui/empty-state'
import { useAuth } from '@/shared/hooks/use-auth'
import { useDashboard } from './hooks/use-dashboard'
import { WidgetCard } from './ui/widget-card'
import { AddWidgetDialog } from './ui/add-widget-dialog'
import { DashboardSwitcher } from './ui/dashboard-switcher'

interface DashboardPageProps {
  dashboardId?: string
}

export function DashboardPage({ dashboardId }: DashboardPageProps) {
  const { user } = useAuth()
  const canEdit = user?.role === 'admin' || user?.role === 'editor'

  const [addWidgetOpen, setAddWidgetOpen] = useState(false)

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

  const title = dashboards.length > 1 ? (
    <DashboardSwitcher dashboards={dashboards} currentDashboard={dashboard} />
  ) : (
    dashboard?.name ?? 'Dashboard'
  )

  return (
    <AuthenticatedLayout
      title={title}
      actions={
        canEdit ? (
          <Button onClick={() => setAddWidgetOpen(true)}>
            <Plus className="h-4 w-4" />
            Add Widget
          </Button>
        ) : undefined
      }
    >
      {isLoading ? (
        <div className="space-y-6">
          <Skeleton className="h-[400px] w-full" />
          <Skeleton className="h-[400px] w-full" />
        </div>
      ) : widgets.length === 0 ? (
        <EmptyState
          icon={BarChart3}
          title="No widgets yet"
          description={
            canEdit
              ? 'Add your first widget to start tracking metrics.'
              : 'No widgets have been added to this dashboard.'
          }
        />
      ) : (
        <div className="space-y-6">
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
    </AuthenticatedLayout>
  )
}
