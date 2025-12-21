import { useState } from 'react'
import { ArrowLeft, MoreHorizontal, Pencil, Plus, Trash } from 'lucide-react'
import { AuthenticatedLayout } from '@/layouts/authenticated'
import { Button } from '@/shared/components/ui/button'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from '@/shared/components/ui/dropdown-menu'
import { useAuth } from '@/shared/hooks/use-auth'
import { useReports } from './hooks/use-reports'
import { ReportList } from './ui/report-list'
import { ReportDetail } from './ui/report-detail'
import {
  CreateReportDialog,
  EditReportDialog,
  DeleteReportDialog,
} from './ui/report-dialog'
import type { Report } from '@/shared/api/generated/models'

export function ReportsPage() {
  const { user } = useAuth()
  const canEdit = user?.role === 'admin' || user?.role === 'editor'

  const [selectedReport, setSelectedReport] = useState<Report | null>(null)
  const [addKpiOpen, setAddKpiOpen] = useState(false)

  const {
    reports,
    isLoading,
    createDialogOpen,
    setCreateDialogOpen,
    editingReport,
    setEditingReport,
    reportToDelete,
    setReportToDelete,
    createReport,
    updateReport,
    deleteReport,
    isCreating,
    isUpdating,
    isDeleting,
  } = useReports()

  // If viewing a specific report
  if (selectedReport) {
    return (
      <AuthenticatedLayout
        title={
          <div className="flex items-center gap-2">
            <Button
              variant="ghost"
              size="icon"
              className="h-8 w-8"
              onClick={() => setSelectedReport(null)}
            >
              <ArrowLeft className="h-4 w-4" />
            </Button>
            <span>{selectedReport.name}</span>
          </div>
        }
        actions={
          canEdit ? (
            <div className="flex items-center gap-2">
              <Button onClick={() => setAddKpiOpen(true)}>
                <Plus className="h-4 w-4" />
                Add KPI
              </Button>
              <DropdownMenu>
                <DropdownMenuTrigger asChild>
                  <Button variant="ghost" size="icon">
                    <MoreHorizontal className="h-4 w-4" />
                  </Button>
                </DropdownMenuTrigger>
                <DropdownMenuContent align="end">
                  <DropdownMenuItem onClick={() => setEditingReport(selectedReport)}>
                    <Pencil className="mr-2 h-4 w-4" />
                    Rename
                  </DropdownMenuItem>
                  <DropdownMenuItem
                    className="text-destructive focus:text-destructive"
                    onClick={() => setReportToDelete(selectedReport)}
                  >
                    <Trash className="mr-2 h-4 w-4" />
                    Delete
                  </DropdownMenuItem>
                </DropdownMenuContent>
              </DropdownMenu>
            </div>
          ) : undefined
        }
      >
        <ReportDetail
          report={selectedReport}
          canEdit={canEdit}
          addKpiOpen={addKpiOpen}
          onAddKpiOpenChange={setAddKpiOpen}
        />

        <EditReportDialog
          report={editingReport}
          open={!!editingReport}
          onOpenChange={(open) => {
            if (!open) setEditingReport(null)
          }}
          onUpdate={async (id, name) => {
            await updateReport(id, name)
            setSelectedReport({ ...selectedReport, name })
          }}
          isLoading={isUpdating}
        />

        <DeleteReportDialog
          report={reportToDelete}
          open={!!reportToDelete}
          onOpenChange={(open) => {
            if (!open) setReportToDelete(null)
          }}
          onDelete={async (id) => {
            await deleteReport(id)
            setSelectedReport(null)
          }}
          isLoading={isDeleting}
        />
      </AuthenticatedLayout>
    )
  }

  // List view
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
      <ReportList
        reports={reports}
        isLoading={isLoading}
        onSelect={setSelectedReport}
      />

      <CreateReportDialog
        open={createDialogOpen}
        onOpenChange={setCreateDialogOpen}
        onCreate={createReport}
        isLoading={isCreating}
      />
    </AuthenticatedLayout>
  )
}
