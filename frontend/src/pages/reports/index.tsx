import { Plus } from 'lucide-react'
import { AuthenticatedLayout } from '@/layouts/authenticated'
import { Button } from '@/shared/components/ui/button'
import { useAuth } from '@/shared/hooks/use-auth'
import { useReports } from './hooks/use-reports'
import { ReportList } from './ui/report-list'
import { CreateReportDialog } from './ui/report-dialog'

export function ReportsPage() {
  const { user } = useAuth()
  const canEdit = user?.role === 'admin' || user?.role === 'editor'

  const {
    reports,
    isLoading,
    createDialogOpen,
    setCreateDialogOpen,
    createReport,
    isCreating,
  } = useReports()

  return (
    <AuthenticatedLayout
      title="Reports"
      actions={
        canEdit ? (
          <Button onClick={() => setCreateDialogOpen(true)}>
            <Plus className="h-4 w-4" />
            New Report
          </Button>
        ) : undefined
      }
    >
      <ReportList reports={reports} isLoading={isLoading} />

      <CreateReportDialog
        open={createDialogOpen}
        onOpenChange={setCreateDialogOpen}
        onCreate={createReport}
        isLoading={isCreating}
      />
    </AuthenticatedLayout>
  )
}
