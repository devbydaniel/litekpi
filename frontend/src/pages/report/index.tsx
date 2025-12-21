import { useState } from 'react'
import { Link, useNavigate } from '@tanstack/react-router'
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
import { useReportDetail } from './hooks/use-report-detail'
import { ReportDetail } from './ui/report-detail'
import { EditReportDialog, DeleteReportDialog } from './ui/report-dialog'

interface ReportPageProps {
  reportId: string
}

export function ReportPage({ reportId }: ReportPageProps) {
  const navigate = useNavigate()
  const { user } = useAuth()
  const canEdit = user?.role === 'admin' || user?.role === 'editor'

  const [addKpiOpen, setAddKpiOpen] = useState(false)
  const [editDialogOpen, setEditDialogOpen] = useState(false)
  const [deleteDialogOpen, setDeleteDialogOpen] = useState(false)

  const {
    report,
    isLoadingReport,
    updateReport,
    deleteReport,
    isUpdatingReport,
    isDeletingReport,
  } = useReportDetail(reportId)

  const handleDelete = async () => {
    await deleteReport()
    navigate({ to: '/reports' })
  }

  return (
    <AuthenticatedLayout
      title={
        <div className="flex items-center gap-2">
          <Button variant="ghost" size="icon" className="h-8 w-8" asChild>
            <Link to="/reports">
              <ArrowLeft className="h-4 w-4" />
            </Link>
          </Button>
          <span>{isLoadingReport ? 'Loading...' : report?.name}</span>
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
                <DropdownMenuItem onClick={() => setEditDialogOpen(true)}>
                  <Pencil className="mr-2 h-4 w-4" />
                  Rename
                </DropdownMenuItem>
                <DropdownMenuItem
                  className="text-destructive focus:text-destructive"
                  onClick={() => setDeleteDialogOpen(true)}
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
        reportId={reportId}
        canEdit={canEdit}
        addKpiOpen={addKpiOpen}
        onAddKpiOpenChange={setAddKpiOpen}
      />

      {report && (
        <>
          <EditReportDialog
            report={report}
            open={editDialogOpen}
            onOpenChange={setEditDialogOpen}
            onUpdate={async (_id, name) => {
              await updateReport(name)
            }}
            isLoading={isUpdatingReport}
          />

          <DeleteReportDialog
            report={report}
            open={deleteDialogOpen}
            onOpenChange={setDeleteDialogOpen}
            onDelete={handleDelete}
            isLoading={isDeletingReport}
          />
        </>
      )}
    </AuthenticatedLayout>
  )
}
